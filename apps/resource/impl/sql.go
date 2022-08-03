package impl

//sqlInsertResource插入数据使用
// (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) 代表占位符
//

const (
	sqlInsertResource = `INSERT INTO resource (
		id,resource_type,vendor,region,zone,create_at,expire_at,category,type,
		name,description,status,update_at,sync_at,sync_accout,public_ip,
		private_ip,pay_type,describe_hash,resource_hash,secret_id,domain,
		namespace,env,usage_mode
	) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`

	//sqlUpdateResource为定义的用于变更information结构体属性
	sqlUpdateResource = `UPDATE resource SET 
		expire_at=?,category=?,type=?,name=?,description=?,
		status=?,update_at=?,sync_at=?,sync_accout=?,
		public_ip=?,private_ip=?,pay_type=?,describe_hash=?,resource_hash=?,
		secret_id=?,namespace=?,env=?,usage_mode=?
	WHERE id = ?`

	//sqlUpdateResource删除信息，通过id进行删除
	sqlDeleteResource = `DELETE FROM resource WHERE id = ?;`

	//sqlQueryResource，联表查询，以左边的表为准进行查询
	// resource_tag t 为别名
	//join 需要条件使用on链接，通过t.resource_id 这个字段进行链接
	sqlQueryResource = `SELECT r.* FROM resource r %s JOIN resource_tag t ON r.id = t.resource_id`

	//使用count 会重复计数
	//对字段去重使用DISTINCT，统计出当前资源查出来的个数
	//用于分页使用获取总数
	sqlCountResource = `SELECT COUNT(DISTINCT r.id) FROM resource r %s JOIN resource_tag t ON r.id = t.resource_id`

	sqlDeleteResourceTag = `
		DELETE 
		FROM
			resource_tag 
		WHERE
			resource_id =? 
			AND t_key =? 
			AND t_value =?;
	`

	//操作tag的sql
	//通过resource_id，把resource_tag表里的tag查询出来
	sqlQueryResourceTag = `SELECT t_key,t_value,description,resource_id,weight,type FROM resource_tag`

	//sqlDeleteResourceTag = `
	//	DELETE
	//	FROM
	//		resource_tag
	//	WHERE
	//		resource_id =?
	//		AND t_key =?
	//		AND t_value =?;
	//`

	//同时包括insert和uodate 的逻辑
	//如果存在第三方就不更新表
	//更新description和weight是有条件的
	//当type不等于1 的时候就使用传进来的值
	//ON DUPLICATE KEY UPDATE description = 如果出现一个主键冲突报错的话，用on关键字 当产生了主键冲突我们就执行
	//update操作 只更新description，这个update我们又不能完全update，需要根据具类型进行update
	//resource的type类型不等于第三方就允许更新，
	sqlInsertOrUpdateResourceTag = `
		INSERT INTO resource_tag ( type, t_key, t_value, description, resource_id, weight, create_at)
		VALUES
			( ?,?,?,?,?,?,? )
			ON DUPLICATE KEY UPDATE description =
		IF
			( type != 1,?, description ),
			weight =
		IF
			( type != 1,?, weight );
	`
)
