package impl

const (
	//在表里插入数据
	insertHostSQL = `
	INSERT INTO resource_host  (
		resource_id,cpu,memory,gpu_amount,gpu_spec,os_type,os_name,
		serial_number,image_id,internet_max_bandwidth_out,
		internet_max_bandwidth_in,key_pair_name,security_groups
	) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?);
	`
	//根据resource_id跟新表里的数据
	updateHostSQL = `
	UPDATE resource_host  SET 
		cpu=?,memory=?,gpu_amount=?,gpu_spec=?,os_type=?,os_name=?,
		image_id=?,internet_max_bandwidth_out=?,
		internet_max_bandwidth_in=?,key_pair_name=?,security_groups=?
	WHERE resource_id = ?
	`

	//链接两张表查询
	queryHostSQL = `
	SELECT
	r.*,
	h.*
	FROM
	resource AS r
	LEFT JOIN resource_host  h ON r.id = h.resource_id
	LEFT JOIN resource_tag t ON r.id = t.resource_id`

	countHostSQL = `SELECT
	COUNT(DISTINCT r.id)
	FROM
	resource AS r
	LEFT JOIN resource_host  h ON r.id = h.resource_id
	LEFT JOIN resource_tag t ON r.id = t.resource_id
	`

	deleteHostSQL = `DELETE FROM resource_host  WHERE resource_id = ?;`
)
