package api

import (
	"github.com/coolgix/cmdb/apps/resource"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/http/response"
)

//暴露的handler
//SearchResource resource有一套restful的接口
//restful.Request，restful.Response
func (h *handler) SearchResource(r *restful.Request, w *restful.Response) {
	// NewSearchRequestFromHTTP 从http里面解析参数
	req, err := resource.NewSearchRequestFromHTTP(r.Request)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	set, err := h.service.Search(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	response.Success(w.ResponseWriter, set)
}
