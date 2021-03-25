package v1alpha2

import (
	"fmt"

	"github.com/emicklei/go-restful"
	"k8s.io/klog"

	"kubesphere.io/kubesphere/pkg/api"
	meteringv1alpha1 "kubesphere.io/kubesphere/pkg/api/metering/v1alpha1"
	"kubesphere.io/kubesphere/pkg/apiserver/request"
	monitoringv1alpha3 "kubesphere.io/kubesphere/pkg/kapis/monitoring/v1alpha3"
	"kubesphere.io/kubesphere/pkg/models/metering"
	"kubesphere.io/kubesphere/pkg/models/monitoring"
	monitoringclient "kubesphere.io/kubesphere/pkg/simple/client/monitoring"
)

func (h *tenantHandler) QueryMeterings(req *restful.Request, resp *restful.Response) {

	u, ok := request.UserFrom(req.Request.Context())
	if !ok {
		err := fmt.Errorf("cannot obtain user info")
		klog.Errorln(err)
		api.HandleForbidden(resp, req, err)
		return
	}

	q := meteringv1alpha1.ParseQueryParameter(req)

	res, err := h.tenant.Metering(u, q)
	if err != nil {
		api.HandleBadRequest(resp, nil, err)
		return
	}

	if q.Operation == monitoringv1alpha3.OperationExport {
		monitoringv1alpha3.ExportMetrics(resp, res)
		return
	}

	resp.WriteAsJson(res)
}

func (h *tenantHandler) QueryMeteringsHierarchy(req *restful.Request, resp *restful.Response) {
	u, ok := request.UserFrom(req.Request.Context())
	if !ok {
		err := fmt.Errorf("cannot obtain user info")
		klog.Errorln(err)
		api.HandleForbidden(resp, req, err)
		return
	}

	q := meteringv1alpha1.ParseQueryParameter(req)
	q.Level = monitoringclient.LevelPod

	resourceStats, err := h.tenant.MeteringHierarchy(u, q)
	if err != nil {
		api.HandleBadRequest(resp, nil, err)
		return
	}

	resp.WriteAsJson(resourceStats)
}

func (h *tenantHandler) HandlePriceInfoQuery(req *restful.Request, resp *restful.Response) {

	var priceResponse metering.PriceResponse
	priceResponse.Init()

	meterConfig, err := monitoring.LoadYaml()
	if err != nil {
		klog.Warning(err)
		resp.WriteAsJson(priceResponse)
		return
	}

	priceInfo := meterConfig.GetPriceInfo()
	priceResponse.RetentionDay = meterConfig.RetentionDay
	priceResponse.Currency = priceInfo.CurrencyUnit
	priceResponse.CpuPerCorePerHour = priceInfo.CpuPerCorePerHour
	priceResponse.MemPerGigabytesPerHour = priceInfo.MemPerGigabytesPerHour
	priceResponse.IngressNetworkTrafficPerMegabytesPerHour = priceInfo.IngressNetworkTrafficPerMegabytesPerHour
	priceResponse.EgressNetworkTrafficPerMegabytesPerHour = priceInfo.EgressNetworkTrafficPerMegabytesPerHour
	priceResponse.PvcPerGigabytesPerHour = priceInfo.PvcPerGigabytesPerHour

	resp.WriteAsJson(priceResponse)

	return
}
