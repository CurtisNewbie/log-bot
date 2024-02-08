CREATE DATABASE IF NOT EXISTS logbot;

CREATE TABLE IF NOT EXISTS logbot.error_log (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `node` VARCHAR(25) NOT NULL COMMENT 'node name',
  `app` VARCHAR(25) NOT NULL COMMENT 'app name',
  `caller` varchar(255) NOT NULL COMMENT 'caller name',
  `trace_id` varchar(25) NOT NULL DEFAULT '' COMMENT 'trace id',
  `span_id` varchar(25) NOT NULL DEFAULT '' COMMENT 'trace id',
  `err_msg` TEXT COMMENT 'error msg',
  `rtime` timestamp default current_timestamp COMMENT 'report time',
  `ctime` timestamp default current_timestamp COMMENT 'create time',
  PRIMARY KEY (`id`),
  INDEX idx_rtime (rtime)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Application Error Log';