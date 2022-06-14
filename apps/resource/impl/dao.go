package impl

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/coolgix/cmdb/apps/resource"
	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/infraboard/mcube/sqlbuilder"
	"strings"
)

// 标准的查询逻辑
//传递一个*sql.DB sql链接过来，直接做一个tag的查询
// 相当于写了一个 select * from resource_tag where resource_id in (?,?,?);
//通过resource_id 把所有的tag找出来
func QueryTag(ctx context.Context, db *sql.DB, resourceIds []string) (
	tags []*resource.Tag, err error) {

	//传进来一堆id 需要检查这个id是否存在
	//这个是tag 为空的情况
	if len(resourceIds) == 0 {
		return
	}

	//正常的range逻辑
	// 把id通过SQL 拼凑一个 IN (?,?,?,...)
	//可能有多个条件，如果有10个就有10个占位符，到底有多个占位符号通过resourceIds []string来确定
	//所以做了一个range，后会生成两个参数，这个参数都是interface
	//args 具体参数, pos 代表占位符
	//有一个参数就放到args里面，any类型
	//pos 就是string类型 给予占位符
	args, pos := []any{}, []string{}
	for _, id := range resourceIds {
		args = append(args, id)
		pos = append(pos, "?")
	}

	// 动态生成一个query build
	//获取到sqlQueryResource 语句
	query := sqlbuilder.NewQuery(sqlQueryResource)
	// 加入一个inwhere语句，拼凑一个 IN (?,?,?,...)
	//通过逗号，把占位符放进来
	inWhere := fmt.Sprintf("resource_id IN (%s)", strings.Join(pos, ","))
	//把这个where语句加入到这个builder里面进行一个构造
	query.Where(inWhere, args...)
	querySQL, args := query.BuildQuery()
	//构造完成后打印这个语句
	zap.L().Debugf("sql: %s", querySQL)

	//最后我们PrepareContext.把值传入进行一个执行
	queryStmt, err := db.PrepareContext(ctx, querySQL)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare query resource tag error, %s", err.Error())
	}
	defer queryStmt.Close()

	rows, err := queryStmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	//执行完成后返回这个tag
	//NewDefaultTag 需要构造这个方法
	for rows.Next() {
		ins := resource.NewDefaultTag()
		err := rows.Scan(
			&ins.Key, &ins.Value, &ins.Describe, &ins.ResourceId, &ins.Weight, &ins.Type,
		)
		if err != nil {
			return nil, exception.NewInternalServerError("query resource tag error, %s", err.Error())
		}
		tags = append(tags, ins)
	}

	return
}
