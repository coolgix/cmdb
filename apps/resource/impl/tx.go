package impl

//为其他表提供公共方法的模块

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/coolgix/cmdb/apps/resource"
	"time"
)

func SaveResource(ctx context.Context, tx *sql.Tx, base *resource.Base, info *resource.Information) error {
	//=======
	// 保存到resource表
	//=======

	// 避免SQL注入, 请使用Prepare
	//使用的就是插入的sql
	stmt, err := tx.PrepareContext(ctx, sqlInsertResource)
	if err != nil {
		return fmt.Errorf("prepare insert resource error, %s", err)
	}
	defer stmt.Close()

	// 保存资源数据，有IP [10.10.1.1, 10.20.2.2] ---> 10.10.1.1,10.20.2.2
	//ip的数组转为字符串。存到数据库的一个字段里面
	_, err = stmt.ExecContext(ctx,
		base.Id, base.ResourceType, base.Vendor, base.Region, base.Zone, base.CreateAt, info.ExpireAt, info.Category, info.Type,
		info.Name, info.Description, info.Status, info.UpdateAt, base.SyncAt, info.SyncAccount, info.PublicIPToString(),
		info.PrivateIPToString(), info.PayType, base.DescribeHash, base.ResourceHash, base.SecretId,
		base.Domain, base.Namespace, base.Env, base.UsageMode,
	)
	if err != nil {
		return fmt.Errorf("save host resource info error, %s", err)
	}

	//=======
	// 保存到resource_tag 表
	//=======
	//如果表里面有tag字段需要更新。如果information有tag，这个数据也需要写到数据库
	if err := updateResourceTag(ctx, tx, base.Id, info.Tags); err != nil {
		return err
	}

	return nil
}

//用于做tag表的更新
func updateResourceTag(ctx context.Context, tx *sql.Tx, resourceId string, tags []*resource.Tag) error {
	// 保存资源标签
	stmt, err := tx.PrepareContext(ctx, sqlInsertOrUpdateResourceTag)
	if err != nil {
		return fmt.Errorf("prepare update resource tag error, %s", err)
	}
	defer stmt.Close()

	for i := range tags {
		t := tags[i]
		if t.Weight == 0 {
			t.Weight = 1
		}
		_, err = stmt.ExecContext(ctx,
			t.Type, t.Key, t.Value, t.Describe, resourceId, t.Weight, time.Now().UnixMilli(),
			t.Describe, t.Weight,
		)
		if err != nil {
			return fmt.Errorf("save resource tag error, %s", err)
		}
	}

	return nil
}

func UpdateResource(ctx context.Context, tx *sql.Tx, base *resource.Base, info *resource.Information) error {
	// 避免SQL注入, 请使用Prepare
	stmt, err := tx.PrepareContext(ctx, sqlUpdateResource)
	if err != nil {
		return fmt.Errorf("prepare update reousrce sql error, %s", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		info.ExpireAt, info.Category, info.Type, info.Name, info.Description,
		info.Status, info.UpdateAt, base.SyncAt, info.SyncAccount,
		info.PublicIPToString(), info.PrivateIPToString(), info.PayType, base.DescribeHash, base.ResourceHash,
		base.SecretId, base.Namespace, base.Env, base.UsageMode,
		base.Id,
	)
	if err != nil {
		return fmt.Errorf("update resource base info error, %s", err)
	}

	if err := updateResourceTag(ctx, tx, base.Id, info.Tags); err != nil {
		return fmt.Errorf("update resource tag error, %s", err)
	}

	return nil
}

func DeleteResource(ctx context.Context, tx *sql.Tx, id string) error {
	// 删除Resource表里面的数据
	stmt, err := tx.PrepareContext(ctx, sqlDeleteResource)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	// 删除Resource_Tag表里面的数据
	stmt, err = tx.PrepareContext(ctx, sqlDeleteResourceTag)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
