package impl

import (
	"context"
	"fmt"
	"strings"

	"github.com/coolgix/cmdb/apps/resource"
	"github.com/infraboard/mcube/sqlbuilder"
)

func (s *service) Search(ctx context.Context, req *resource.SearchRequest) (
	*resource.ResourceSet, error) {
	//sql 是一个模版，到底是左链接还是右链接，取决于我们是否需要关联tag表
	// left join 是先扫描左表，right join 先扫描右表,当有tag 过滤，需要关联右表，可以以右表为准
	//  如果扫描tag的成本比扫描resource的表的成本低，我们就用right join
	//需要判断请求有没有传tag的条件给我们

	join := "LEFT"
	if req.Hastag() {
		join = "RIGHT"
	}

	// 构建过滤条件
	//构建sql具体使用那种方法
	builder := sqlbuilder.NewQuery(fmt.Sprintf(sqlQueryResource, join))
	s.buildQuery(builder, req)

	// =========
	// 计数统计: COUNT语句
	// =========
	set := resource.NewResourceSet()

	//获取total select count(*) FROMT t where ....
	countSQL, args := builderCountWith("COUNT(DISTINCT r.id)")
	countStmt, err := s.db.Prepare(countSQL)
	if err != nil {
		s.log.Debugf("count sql,%s,%v", countSQL, args)
		return nil, exception.NewInternalServerError("prepare count sql erro, %s", err)
	}
	defer countStmt.Close()
	err = countStmt.QueryRow(args...).Scan(&set.Total)

	return set, nil
}
func (s *service) buildQuery(builder *sqlbuilder.Builder, req *resource.SearchRequest) {
	// 构建过滤条件
	//关键字语句，动态过滤参数
	//参数里面有模糊匹配与关键字匹配
	if req.Keywords != "" {
		if req.ExactMatch {
			//补充sql语句的条件
			//精确匹配
			builder.Where("r.name = ? OR r.id = ? OR r.private_ip = ? OR r.public_ip = ?",
				req.Keywords,
				req.Keywords,
				req.Keywords,
				req.Keywords,
			)
		} else {
			// 模糊匹配
			builder.Where("r.name LIKE ? OR r.id = ? OR r.private_ip LIKE ? OR r.public_ip LIKE ?",
				"%"+req.Keywords+"%", //name 是前后通配
				"%"+req.Keywords+"%", //description 前后通配
				req.Keywords+"%",     //private_ip 需要一个后缀的通配
				req.Keywords+"%",     //public_ip 需要一个后缀的通配
			)
		}

	}

	// 按照资源属性过滤
	//补充where语句条件
	if req.Domain != "" {
		builder.Where("r.domain = ?", req.Domain)
	}
	if req.Namespace != "" {
		builder.Where("r.namespace = ?", req.Namespace)
	}
	if req.Env != "" {
		builder.Where("r.env = ?", req.Env)
	}
	if req.UsageMode != nil {
		builder.Where("r.usage_mode = ?", req.UsageMode)
	}
	if req.Vendor != nil {
		builder.Where("r.vendor = ?", req.Vendor)
	}
	if req.SyncAccount != "" {
		builder.Where("r.sync_accout = ?", req.SyncAccount)
	}
	if req.Type != nil {
		builder.Where("r.resource_type = ?", req.Type)
	}
	if req.Status != "" {
		builder.Where("r.status = ?", req.Status)
	}

	//tag过滤
	//在resource.tag表进行关联查询
	//通过tag key 和tag value 进行连表查询，配上where条件
	//我们允许输入多个tag来对资源进行解锁，多个tag之间的关系，到底是and还是or app=v1,product=p2
	//我们实现的策略是基于and 实现
	for i := range req.Tags {
		//取出selector 做拼接，如果selector的key为空，就不继续操作
		selector := req.Tags[i]

		//tag：=v1 ，作为tag查询，tag的key是必须的
		if selector.Key == "" {
			continue
		}
		//添加key为过滤条件，
		// like 默认是匹配所有，所以我们就把* 替换为%
		//.* 转为%号的操作
		//tag_key="xxxx", .* ,定制key如何统配
		query.Where("t.t_key LIKE ?", strings.ReplaceAll(selector.Key, ".*", "%"))

		//场景1：添加value过滤条件
		//定制value如何统配,app=["app1","app2","app3"] 或的关系，
		//tag value 是数组
		//tag_value=? or tag_value=?.有几个tag value

		//场景2：用户给予的tag是一个带有比较符号的条件
		//app_count>1
		//在app定义这个场景的问题

		//tag_value LIKe ? or tag_value LIKe ?
		var condtions []string
		var args []any
		for _, v := range selector.Values {
			//=,!=,=~,!~ 四种操作的统配
			//t.t_value [=,!=,=~,!~] value //这样一种表达式
			//where 条件语句
			condtions = append(condtions, fmt.Sprintf("t.t_value % ?", selector.Operator))
			//条件参数args
			args = append(args, string.ReplaceAll(v, ".*", "%"))

			//args=append(args,v) 这种也没有有问题，上面的写法是吧tag_value .* 变为 % 占位符 做的特殊处理，为了匹配正则里面的
			//.* 专门做的处理 如果app=product1.*匹配出来，就可以用%代替.%,
			//使用原声sql写是这样的app=product1.% ，用户传递不可能用%传递

		}
		//如果tag的value是有condtions多个条件做成的
		//app=~app1,app2, 根据符号来决定我们这个value之间的关系
		if len(condtions) > 0 {
			vmwhere := fmt.Sprintf("(%s)", strings.Join(condtions, selector.RelationShip()))
			builder.where(vmwhere, args...)
		}
	}
}

func (s *service) QueryTag(ctx context.Context, req *resource.QueryTagRequest) (
	*resource.TagSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryTag not implemented")
}

func (s *service) UpdateTag(ctx context.Context, req *resource.UpdateTagRequest) (
	*resource.Resource, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTag not implemented")
}
