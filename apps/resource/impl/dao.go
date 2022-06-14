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

func QueryTag(ctx context.Context, db *sql.DB, resourceIds []string) (
	tags []*resource.Tag, err error) {
	if len(resourceIds) == 0 {
		return
	}

	// 通过SQL 拼凑一个 IN (?,?,?,...)
	// args 具体参数, pos 代表占位符
	args, pos := []any{}, []string{}
	for _, id := range resourceIds {
		args = append(args, id)
		pos = append(pos, "?")
	}

	// query build
	query := sqlbuilder.NewQuery(sqlQueryResource)
	// 拼凑一个 IN (?,?,?,...)
	inWhere := fmt.Sprintf("resource_id IN (%s)", strings.Join(pos, ","))
	query.Where(inWhere, args...)
	querySQL, args := query.BuildQuery()
	zap.L().Debugf("sql: %s", querySQL)

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
