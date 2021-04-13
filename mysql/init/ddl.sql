CREATE SCHEMA IF NOT EXISTS `database` DEFAULT CHARACTER SET utf8mb4;
USE `database`;

SET CHARSET utf8mb4;

CREATE TABLE IF NOT EXISTS `database`.`users` (
  `id` VARCHAR(128) NOT NULL COMMENT 'ユーザID',
  `name` VARCHAR(128) NOT NULL COMMENT 'ユーザ名',
  `high_score` INT UNSIGNED NOT NULL COMMENT 'ハイスコア',
  PRIMARY KEY (`id`),
  INDEX `idx_high_score` (`high_score` DESC))
ENGINE = InnoDB
COMMENT = 'ユーザ';