package impl

import (
	"database/sql"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/coolgix/cmdb/apps/resource"

	"github.com/coolgix/cmdb/conf"

	"github.com/infraboard/mcube/app"

	"google.golang.org/grpc"

)

//接口实现的类

var {
	// service 服务实例

	svr = & service{}
}

//依赖于一个db，logger，
type service struct {
	db *sql.DB
	log logger.Logger

	resource.UnimplementedServiceServer //作为一个grpc的服务实现嵌入生成的一个UnimplementedServiceServer
}

// Config 作为grpc的实现需要有配置方法实现
func (s *service) Config() error{
	db, err := conf.C().MySQL.GetDB() //数据库配置从全局配置文档中获取
	if err != nil {
		return err
	}

	s.log = zap.L().Named(s.Name()) //初始化一个子logger
	s.db = db
	return nil
}

// Name 作为grpc的实现需要有对象的名称
func (s *service) Name() string {
	return resource.AppName

}

// Registry 作为grpc的实现需要有服务的注册方法
func (s *service) Registry(server grpc.Server) {
	resouce.RegisterServiceServer(server.svr)
}


func init(){
	app.RegistryGrpcApp(svr)

}

