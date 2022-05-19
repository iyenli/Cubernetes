package httpserver

import (
	"Cubernetes/cmd/apiserver/httpserver/restful"
	"net/http"
)

var restfulList = []Handler{
	{http.MethodGet, "/health", restful.GetHealth},

	{http.MethodGet, "/apis/pod/:uid", restful.GetPod},
	{http.MethodGet, "/apis/pods", restful.GetPods},
	{http.MethodPost, "/apis/pod", restful.PostPod},
	{http.MethodPut, "/apis/pod/:uid", restful.PutPod},
	{http.MethodDelete, "/apis/pod/:uid", restful.DelPod},
	{http.MethodPost, "/apis/select/pods", restful.SelectPods},
	{http.MethodPut, "/apis/pod/status/:uid", restful.UpdatePodStatus},

	{http.MethodGet, "/apis/service/:uid", restful.GetService},
	{http.MethodGet, "/apis/services", restful.GetServices},
	{http.MethodPost, "/apis/service", restful.PostService},
	{http.MethodPut, "/apis/service/:uid", restful.PutService},
	{http.MethodDelete, "/apis/service/:uid", restful.DelService},
	{http.MethodPost, "/apis/select/services", restful.SelectServices},

	{http.MethodGet, "/apis/replicaSet/:uid", restful.GetReplicaSet},
	{http.MethodGet, "/apis/replicaSets", restful.GetReplicaSets},
	{http.MethodPost, "/apis/replicaSet", restful.PostReplicaSet},
	{http.MethodPut, "/apis/replicaSet/:uid", restful.PutReplicaSet},
	{http.MethodDelete, "/apis/replicaSet/:uid", restful.DelReplicaSet},
	{http.MethodPost, "/apis/select/replicaSets", restful.SelectReplicaSets},

	{http.MethodGet, "/apis/node/:uid", restful.GetNode},
	{http.MethodGet, "/apis/nodes", restful.GetNodes},
	{http.MethodPost, "/apis/node", restful.PostNode},
	{http.MethodPut, "/apis/node/:uid", restful.PutNode},
	{http.MethodDelete, "/apis/node/:uid", restful.DelNode},
	{http.MethodPost, "/apis/select/nodes", restful.SelectNodes},

	{http.MethodGet, "/apis/dns/:uid", restful.GetDns},
	{http.MethodGet, "/apis/dnses", restful.GetDnses},
	{http.MethodPost, "/apis/dns", restful.PostDns},
	{http.MethodPut, "/apis/dns/:uid", restful.PutDns},
	{http.MethodDelete, "/apis/dns/:uid", restful.DelDns},
	{http.MethodPost, "/apis/select/dnses", restful.SelectDnses},

	{http.MethodGet, "/apis/autoScaler/:uid", restful.GetAutoScaler},
	{http.MethodGet, "/apis/autoScalers", restful.GetAutoScalers},
	{http.MethodPost, "/apis/autoScaler", restful.PostAutoScaler},
	{http.MethodPut, "/apis/autoScaler/:uid", restful.PutAutoScaler},
	{http.MethodDelete, "/apis/autoScaler/:uid", restful.DelAutoScaler},
	{http.MethodPost, "/apis/select/autoScalers", restful.SelectAutoScalers},

	{http.MethodGet, "/apis/gpuJob/:uid", restful.GetGpuJob},
	{http.MethodGet, "/apis/gpuJobs", restful.GetGpuJobs},
	{http.MethodPost, "/apis/gpuJob", restful.PostGpuJob},
	{http.MethodPut, "/apis/gpuJob/:uid", restful.PutGpuJob},
	{http.MethodDelete, "/apis/gpuJob/:uid", restful.DelGpuJob},
	{http.MethodPost, "/apis/select/gpuJobs", restful.SelectGpuJobs},
}
