CREATE TABLE IF NOT EXISTS `books` (
  `id` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '对象Id',
  `create_at` bigint NOT NULL COMMENT '创建时间(13位时间戳)',
  `create_by` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '创建人',
  `update_at` bigint NOT NULL COMMENT '更新时间',
  `update_by` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '更新人',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '书名',
  `author` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '作者',
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`) USING BTREE COMMENT '用于书名搜索',
  KEY `idx_author` (`author`) USING BTREE COMMENT '用于作者搜索'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

# TABLE IF NOT 如果不存在就创建
CREATE TABLE IF NOT EXISTS `resource` (
                                          `id` char(64) NOT NULL COMMENT '全局唯一Id, 直接使用个云商自己的Id',
                                          `resource_type` tinyint(2) NOT NULL COMMENT '资源类型',
                                          `vendor` tinyint(1) NOT NULL COMMENT '厂商',
                                          `region` varchar(64) NOT NULL COMMENT '地域',
                                          `zone` varchar(64) NOT NULL COMMENT '区域',
                                          `create_at` bigint(13) NOT NULL COMMENT '创建时间',
                                          `expire_at` bigint(13) DEFAULT NULL COMMENT '过期时间',
                                          `category` varchar(64) NOT NULL COMMENT '种类',
                                          `type` varchar(120) NOT NULL COMMENT '规格',
                                          `name` varchar(255) NOT NULL COMMENT '名称',
                                          `description` varchar(255) DEFAULT NULL COMMENT '描述',
                                          `status` varchar(255) NOT NULL COMMENT '服务商中的状态',
                                          `update_at` bigint(13) DEFAULT NULL COMMENT '更新时间',
                                          `sync_at` bigint(13) DEFAULT NULL COMMENT '同步时间',
                                          `sync_accout` varchar(255) DEFAULT NULL COMMENT '同步账号',
                                          `public_ip` varchar(64) DEFAULT NULL COMMENT '公网IP',
                                          `private_ip` varchar(64) DEFAULT NULL COMMENT '内网IP',
                                          `pay_type` varchar(255) DEFAULT NULL COMMENT '实例付费方式',
                                          `describe_hash` varchar(255) NOT NULL COMMENT '描述数据Hash',
                                          `resource_hash` varchar(255) NOT NULL COMMENT '基础数据Hash',
                                          `secret_id` varchar(64) NOT NULL COMMENT '关联的同于同步的secret id',
                                          `domain` varchar(255) NOT NULL COMMENT '资源所属域',
                                          `namespace` varchar(255) NOT NULL COMMENT '资源所属空间',
                                          `env` varchar(255) NOT NULL COMMENT '资源所属环境',
                                          `usage_mode` tinyint(2) NOT NULL COMMENT '资源使用方式',
                                          PRIMARY KEY (`id`) USING BTREE,
                                          KEY `idx_name` (`name`) USING BTREE,
                                          KEY `idx_status` (`status`) USING BTREE,
                                          KEY `idx_private_ip` (`public_ip`) USING BTREE,
                                          KEY `idx_public_ip` (`public_ip`) USING BTREE,
                                          KEY `idx_domain` (`domain`) USING HASH,
                                          KEY `idx_namespace` (`namespace`) USING HASH,
                                          KEY `idx_env` (`env`) USING HASH
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资源基础信息,所有资源的公共信息, 用于全局解索';

CREATE TABLE IF NOT EXISTS `resource_tag` (
                                              `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '标签Id',
                                              `t_key` varchar(255) NOT NULL COMMENT '标签的名称',
                                              `t_value` varchar(255) NOT NULL COMMENT '标签的值',
                                              `description` varchar(255) NOT NULL COMMENT '值的描述信息',
                                              `resource_id` varchar(64)  NOT NULL COMMENT '标签关联的资源Id',
                                              `weight` int(11) NOT NULL COMMENT '标签权重',
                                              `type` tinyint(4) NOT NULL COMMENT '标签类型',
                                              `create_at` bigint(13) NOT NULL COMMENT '标签创建时间',
                                              PRIMARY KEY (`id`),
                                              UNIQUE KEY `idx_id` (`t_key`,`t_value`,`resource_id`) COMMENT '一个资源同一个key value只允许有一对',
                                              KEY `idx_key` (`t_key`) USING HASH,
                                              KEY `idx_value` (`t_value`) USING BTREE,
                                              KEY `idx_resource_id` (`resource_id`) USING HASH
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COMMENT='资源标签';

CREATE TABLE IF NOT EXISTS `resource_host` (
                                               `resource_id` varchar(64)  NOT NULL COMMENT '关联的资源Id',
                                               `cpu` tinyint(4) NOT NULL COMMENT 'cpu核数',
                                               `memory` int(13) NOT NULL COMMENT '内存大小',
                                               `gpu_amount` tinyint(4) DEFAULT NULL COMMENT 'gpu核数',
                                               `gpu_spec` varchar(255)  DEFAULT NULL COMMENT 'gpu规格',
                                               `os_type` varchar(255)  DEFAULT NULL COMMENT '操作系统类型',
                                               `os_name` varchar(255) DEFAULT NULL COMMENT '操作系统名称',
                                               `serial_number` varchar(120)  DEFAULT NULL COMMENT '系统序列号',
                                               `image_id` char(64)  DEFAULT NULL COMMENT '镜像Id',
                                               `internet_max_bandwidth_out` int(10) DEFAULT NULL COMMENT '外网最大出口带宽',
                                               `internet_max_bandwidth_in` int(10) DEFAULT NULL COMMENT '外网最大入口带宽',
                                               `key_pair_name` varchar(255)  DEFAULT NULL COMMENT 'ssh key关联Id',
                                               `security_groups` varchar(255)  DEFAULT NULL COMMENT '安全组Id列表',
                                               PRIMARY KEY (`resource_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务器主机信息';

CREATE TABLE IF NOT EXISTS `secret` (
                                        `id` varchar(64) NOT NULL COMMENT '凭证Id',
                                        `create_at` bigint(13) NOT NULL COMMENT '创建时间',
                                        `description` varchar(255) NOT NULL COMMENT '凭证描述',
                                        `vendor` tinyint(1) NOT NULL COMMENT '资源提供商',
                                        `address` varchar(255)  NOT NULL COMMENT '体验提供方访问地址',
                                        `allow_regions` text  NOT NULL COMMENT '允许同步的Region列表',
                                        `crendential_type` tinyint(1) NOT NULL COMMENT '凭证类型',
                                        `api_key` varchar(255) NOT NULL COMMENT '凭证key',
                                        `api_secret` text  NOT NULL COMMENT '凭证secret',
                                        `request_rate` int(11) NOT NULL COMMENT '请求速率',
                                        PRIMARY KEY (`id`),
                                        UNIQUE KEY `idx_key` (`api_key`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资源提供商同步凭证管理';

CREATE TABLE IF NOT EXISTS `task` (
                                      `id` varchar(64) NOT NULL COMMENT '任务Id',
                                      `region` varchar(64) NOT NULL COMMENT '资源所属Region',
                                      `resource_type` tinyint(1) NOT NULL COMMENT '资源类型',
                                      `secret_id` varchar(64) NOT NULL COMMENT '用于操作资源的凭证Id',
                                      `secret_desc` text NOT NULL COMMENT '凭证描述',
                                      `timeout` int(11) NOT NULL COMMENT '任务超时时间',
                                      `status` tinyint(1) NOT NULL COMMENT '任务当前状态',
                                      `message` text NOT NULL COMMENT '任务失败相关信息',
                                      `start_at` bigint(20) NOT NULL COMMENT '任务开始时间',
                                      `end_at` bigint(20) NOT NULL COMMENT '任务结束时间',
                                      `total_succeed` int(11) NOT NULL COMMENT '总共操作成功的资源数量',
                                      `total_failed` int(11) NOT NULL COMMENT '总共操作失败的资源数量',
                                      PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资源操作任务管理';