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
	//sqlInsertOrUpdateResourceTag = `
	//	INSERT INTO resource_tag ( type, t_key, t_value, description, resource_id, weight, create_at)
	//	VALUES
	//		( ?,?,?,?,?,?,? )
	//		ON DUPLICATE KEY UPDATE description =
	//	IF
	//		( type != 1,?, description ),
	//		weight =
	//	IF
	//		( type != 1,?, weight );
	//`
)
