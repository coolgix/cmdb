package api

import (
	"github.com/coolgix/cmdb/apps/resource"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/http/label"
	"github.com/infraboard/mcube/http/response"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

//暴露的http 接口
//把resource_grpc.pb.go 的三个接口暴露出去
//Search(context.Context, *SearchRequest) (*ResourceSet, error)
//	QueryTag(context.Context, *QueryTagRequest) (*TagSet, error)
//	UpdateTag(context.Context, *UpdateTagRequest) (*Resource, error)

//对象hanler
var (
	h = &handler{}
)

//定义结构体
//handler对依赖我们的resource.ServiceServer 实体类
type handler struct {
	service resource.ServiceServer
	log     logger.Logger
}

//作为一个restful的api需要实现什么？
//初始化他的依赖， logger以及他依赖的服务
func (h *handler) Config() error {
	h.log = zap.L().Named(resource.AppName)
	h.service = app.GetGrpcApp(resource.AppName).(resource.ServiceServer)
	return nil
}

//实现类的名字
func (h *handler) Name() string {
	return resource.AppName

}

//当前的版本是我们的apiversion
//restful api version
//命名规则 /cmdb/api/v1/resource 包含api version
//方便前端对api进行通配
//通过version函数对外进行定义
func (h *handler) Version() string {
	return "v1"
}

func (h *handler) Registry(ws *restful.WebService) {
	tags := []string{h.Name()}

	// RESTful API,   resource = cmdb_resource, action: list, auth: true
	//这个接口到达search路径处理
	ws.Route(ws.GET("/search").To(h.SearchResource).
		//下面都是装饰器
		//用于生成swagger 文档
		//这些label使用权限系统隔离使用
		Doc("get all resources").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.List.Value()).
		Metadata(label.Auth, label.Enable).
		Reads(resource.SearchRequest{}).
		Writes(response.NewData(resource.ResourceSet{})).
		Returns(200, "OK", resource.ResourceSet{}))

	// 资源标签管理
	// ws.Route(ws.POST("/").To(h.AddTag).
	// 	Doc("add resource tags").
	// 	Metadata(restfulspec.KeyOpenAPITags, tags).
	// 	Metadata(label.ResourceLableKey, "tags").
	// 	Metadata(label.ActionLableKey, label.Create.Value()).
	// 	Metadata(label.AuthLabelKey, label.Enable).
	// 	Reads([]*resource.Tag{}).
	// 	Writes(response.NewData(resource.Resource{})))
	// ws.Route(ws.DELETE("/").To(h.RemoveTag).
	// 	Doc("remove resource tags").
	// 	Metadata(restfulspec.KeyOpenAPITags, tags).
	// 	Metadata(label.ResourceLableKey, "tags").
	// 	Metadata(label.ActionLableKey, label.Delete.Value()).
	// 	Metadata(label.AuthLabelKey, label.Enable).
	// 	Reads([]*resource.Tag{}).
	// 	Writes(response.NewData(resource.Resource{})))

	// // 资源发现
	// ws.Route(ws.GET("/discovery/prometheus").To(h.DiscoveryPrometheus).
	// 	Doc("discovery resoruce for prometheus").
	// 	Metadata(restfulspec.KeyOpenAPITags, tags).
	// 	Metadata(label.ResourceLableKey, "prometheus_resource").
	// 	Metadata(label.ActionLableKey, label.List.Value()).
	// 	Reads(resource.SearchRequest{}).
	// 	Writes(response.NewData(resource.ResourceSet{})).
	// 	Returns(200, "OK", resource.ResourceSet{}))

}

func init() {
	app.RegistryRESTfulApp(h)
}
