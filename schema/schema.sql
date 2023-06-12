CREATE DATABASE IF NOT EXISTS log_bot;

CREATE TABLE IF NOT EXISTS app_log_file (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `app` VARCHAR(25) NOT NULL COMMENT 'app name',
  `file` VARCHAR(255) COMMENT 'log file',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Application Log File';

CREATE TABLE IF NOT EXISTS error_log (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `app` VARCHAR(25) NOT NULL COMMENT 'app name',
  `err_msg` TEXT COMMENT 'error msg',
  `ctime` timestamp default current_timestamp,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Application Error Log';