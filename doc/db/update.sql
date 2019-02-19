
DROP TABLE IF EXISTS `task_info`;
CREATE TABLE `task_info` (
  `id`           INT(10) UNSIGNED NOT NULL AUTO_INCREMENT
  COMMENT '主键',
  `task_no`      INT(10)          NULL     DEFAULT NULL
  COMMENT '作业编号',
  `task_name`    VARCHAR(64)      NULL     DEFAULT NULL
  COMMENT '作业名称',
  `task_status`  VARCHAR(64)      NOT NULL
  COMMENT '作业状态 success||failure||unknown',
  `task_type`    VARCHAR(64)      NOT NULL
  COMMENT '作业类型 file||script',
  `task_channel` VARCHAR(64)      NOT NULL
  COMMENT '作业通道 shell||salt',
  `exec_user`    VARCHAR(64)      NOT NULL
  COMMENT '执行用户root或其他',
  `expired_time` INT(10)          NOT NULL
  COMMENT '超时时间',
  `extend`       JSON             NULL
  COMMENT '扩展信息，文件的目标路径，文件名称； 脚本类型shell|python，脚本内容',
  `create_at`    TIMESTAMP        NOT NULL
  COMMENT '创建时间',
  `creator`      VARCHAR(255)              DEFAULT NULL
  COMMENT '创建人',
  `is_active`    BOOLEAN COMMENT '是否有效',
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 13
  DEFAULT CHARSET = utf8;

DROP TABLE IF EXISTS `task_result`;
CREATE TABLE `task_result` (
  `id`         INT(10) UNSIGNED NOT NULL         AUTO_INCREMENT
  COMMENT '主键',
  `task_id`    INT(10)          NOT NULL
  COMMENT '关联作业主键',
  `task_no`    INT(10)          NOT NULL
  COMMENT '作业编号',
  `sn`         VARCHAR(64)      NOT NULL
  COMMENT '设备序列号',
  `hostname`   VARCHAR(64)      NOT NULL
  COMMENT '主机名称',
  `ip`         VARCHAR(64)      NOT NULL
  COMMENT '业务IP',
  `start_time` TIMESTAMP        NOT NULL
  COMMENT '开始时间',
  `end_time`   TIMESTAMP        NOT NULL
  COMMENT '结束时间',
  `total_time` INT(10)          NOT NULL
  COMMENT '执行耗时',
  `status`     VARCHAR(64)      NOT NULL
  COMMENT '执行状态',
  `content`    VARCHAR(255)     NULL
  COMMENT '执行结果',
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 13
  DEFAULT CHARSET = utf8;
