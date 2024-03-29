package impl

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/coolgix/cmdb/apps/host"
	"github.com/coolgix/cmdb/apps/resource/impl"
	"time"
)

//SyncAt是在保存到数据库补充的，补充了一个毫秒时间UnixMilli
//
func (s *service) save(ctx context.Context, h *host.Host) error {
	if h.Base.SyncAt != 0 {
		h.Base.SyncAt = time.Now().UnixMilli()
	}

	var (
		stmt *sql.Stmt
		err  error
	)

	// 开启一个事物
	// 文档请参考: http://cngolib.com/database-sql.html#db-begintx
	// 关于事物级别可以参考文章: https://zhuanlan.zhihu.com/p/117476959
	// wiki: https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// 执行结果提交或者回滚事务
	// 当使用sql.Tx的操作方式操作数据后，需要我们使用sql.Tx的Commit()方法显式地提交事务，
	// 如果出错，则可以使用sql.Tx中的Rollback()方法回滚事务，保持数据的一致性
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				s.log.Errorf("rollback error, %s", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				s.log.Errorf("commit error, %s", err)
			}
		}
	}()

	// 生成描写信息的Hash，分别计算 Resource Inforamtion的Hash和Descirbe(host表独有属性)的Hash
	// 有专门把一个对象--> hash库 ,没选择这种做
	// 把对象--> json(string) --> hash(string)---> 对象的hash
	if err := h.GenHash(); err != nil {
		return err
	}

	// 保存资源基础信息(公共信息)
	// 保存在resource 表里面。
	err = impl.SaveResource(ctx, tx, h.Base, h.Information)
	if err != nil {
		return err
	}

	// 避免SQL注入, 请使用Prepare
	stmt, err = tx.PrepareContext(ctx, insertHostSQL)
	if err != nil {
		return fmt.Errorf("prepare insert host sql error, %s", err)
	}
	defer stmt.Close()

	desc := h.Describe
	_, err = stmt.ExecContext(ctx,
		h.Base.Id, desc.Cpu, desc.Memory, desc.GpuAmount, desc.GpuSpec, desc.OsType, desc.OsName,
		desc.SerialNumber, desc.ImageId, desc.InternetMaxBandwidthOut,
		desc.InternetMaxBandwidthIn, desc.KeyPairNameToString(), desc.SecurityGroupsToString(),
	)
	if err != nil {
		return fmt.Errorf("save host resource describe error, %s", err)
	}

	return err
}

func (s *service) update(ctx context.Context, ins *host.Host) error {
	var (
		stmt *sql.Stmt
		err  error
	)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start tx error, %s", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		} else {
			if err := tx.Commit(); err != nil {
				s.log.Errorf("commit error, %s", err)
			}
		}
	}()

	// 更新资源基本信息
	// 开启一个事务根据这个值有没有变化，根据外层host.go,放在dao之前就需要判断值是否有变化
	//Resource有变化就跟新resource 表
	if ins.Base.ResourceHashChanged {
		if err := impl.UpdateResource(ctx, tx, ins.Base, ins.Information); err != nil {
			return err
		}
	} else {
		s.log.Debug("resource data hash not changed, needn't update")
	}

	// 更新实例信息
	// describe表有信息变化，就跟新resource-host表，提交commit
	if ins.Base.DescribeHashChanged {
		stmt, err = tx.PrepareContext(ctx, updateHostSQL)
		if err != nil {
			return fmt.Errorf("prepare update host sql error, %s", err)
		}
		defer stmt.Close()

		base := ins.Base
		desc := ins.Describe
		_, err = stmt.ExecContext(ctx,
			desc.Cpu, desc.Memory, desc.GpuAmount, desc.GpuSpec, desc.OsType, desc.OsName,
			desc.ImageId, desc.InternetMaxBandwidthOut,
			desc.InternetMaxBandwidthIn, desc.KeyPairNameToString(), desc.SecurityGroupsToString(),
			base.Id,
		)
		if err != nil {
			return err
		}
	} else {
		s.log.Debug("describe data hash not changed, needn't update")
	}

	return err
}

//删除逻辑也需要开启事务，应为要开启两张表
//
func (s *service) delete(ctx context.Context, req *host.ReleaseHostRequest) error {
	var (
		stmt *sql.Stmt
		err  error
	)

	// 开启一个事物
	// 文档请参考: http://cngolib.com/database-sql.html#db-begintx
	// 关于事物级别可以参考文章: https://zhuanlan.zhihu.com/p/117476959
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// 执行结果提交或者回滚事务
	// 当使用sql.Tx的操作方式操作数据后，需要我们使用sql.Tx的Commit()方法显式地提交事务，
	// 如果出错，则可以使用sql.Tx中的Rollback()方法回滚事务，保持数据的一致性
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		} else {
			//没有报错就提交事务
			if err := tx.Commit(); err != nil {
				s.log.Errorf("commit error, %s", err)
			}
		}
	}()

	stmt, err = tx.PrepareContext(ctx, deleteHostSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, req.Id)
	if err != nil {
		return err
	}

	if err := impl.DeleteResource(ctx, tx, req.Id); err != nil {
		return err
	}

	return err
}
