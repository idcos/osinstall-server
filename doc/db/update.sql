CREATE TABLE `task_info` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_no` varchar(64) DEFAULT NULL COMMENT '作业编号',
  `task_name` varchar(64) DEFAULT NULL COMMENT '作业名称',
  `task_status` varchar(64) NOT NULL COMMENT '作业状态 success||failure||unknown',
  `task_type` varchar(64) NOT NULL COMMENT '作业类型 file||script',
  `task_channel` varchar(64) NOT NULL COMMENT '作业通道 shell||salt',
  `run_as` varchar(64) NOT NULL COMMENT '执行用户root或其他',
  `timeout` int(10) NOT NULL COMMENT '超时时间',
  `extend` varchar(255) DEFAULT NULL COMMENT '扩展信息，文件的目标路径，文件名称； 脚本类型shell|python，脚本内容',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
  `creator` varchar(255) DEFAULT NULL COMMENT '创建人',
  `is_active` tinyint(1) DEFAULT NULL COMMENT '是否有效',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=113 DEFAULT CHARSET=utf8;


CREATE TABLE `task_result` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(10) NOT NULL COMMENT '关联作业主键',
  `task_no` varchar(64) NOT NULL COMMENT '作业编号',
  `sn` varchar(64) NOT NULL COMMENT '设备序列号',
  `hostname` varchar(64) NOT NULL COMMENT '主机名称',
  `ip` varchar(64) NOT NULL COMMENT '业务IP',
  `start_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '开始时间',
  `end_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '结束时间',
  `total_time` int(10) NOT NULL COMMENT '执行耗时',
  `status` varchar(64) NOT NULL COMMENT '执行状态',
  `content` varchar(255) DEFAULT NULL COMMENT '执行结果',
  `created_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=73 DEFAULT CHARSET=utf8
