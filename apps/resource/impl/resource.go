package impl

import (
	"context"
	"fmt"
	"strings"

	"github.com/coolgix/cmdb/apps/resource"
	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/sqlbuilder"
)

func (s *.service) Search(ctx context.Context, req *resource.SearchRequest) (
	*resource.ResourceSet, error) {
	//sql 是一个模版，到底是左链接还是右链接，取决于我们是否需要关联tag表
	// left join 是先扫描左表，right join 先扫描右表,当有tag 过滤，需要关联右表，可以以右表为准
	//  如果扫描tag的成本比扫描resource的表的成本低，我们就用right join
	//需要判断请求有没有传tag的条件给我们

	join := "LEFT"
	if req.Hastag(){
		join = "RIGHT"

	}
	//构建sql具体使用那种方法
	query := sqlbuilder.NewQuery(fmt.Sprintf(sqlQueryResource,join))

	// 构建过滤条件
	//关键字语句
	if req.Keywords !=  "" {
		//补充sql语句的条件
		query.Where("")

	}

	return nil, nil
}

func (s *.service) QueryTag(ctx context.Context, req *resource.QueryTagRequest) (
	*resource.TagSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryTag not implemented")
}

func (s *.service) UpdateTag(ctx context.Context, req *resource.UpdateTagRequest) (
	*resource.Resource, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTag not implemented")
}
