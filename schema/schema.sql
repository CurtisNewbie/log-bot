CREATE DATABASE IF NOT EXISTS logbot;

CREATE TABLE IF NOT EXISTS logbot.error_log (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `node` VARCHAR(25) NOT NULL COMMENT 'node name',
  `app` VARCHAR(25) NOT NULL COMMENT 'app name',
  `func` varchar(50) NOT NULL COMMENT 'func name',
  `trace_id` varchar(25) NOT NULL DEFAULT '' COMMENT 'trace id',
  `span_id` varchar(25) NOT NULL DEFAULT '' COMMENT 'trace id',
  `err_msg` TEXT COMMENT 'error msg',
  `ctime` timestamp default current_timestamp,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Application Error Log';