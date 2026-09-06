CREATE DATABASE IF NOT EXISTS `compane_profile`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE `compane_profile`;

CREATE TABLE IF NOT EXISTS `company` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `intro_title` VARCHAR(255) NOT NULL DEFAULT '',
  `intro_body` TEXT NOT NULL,
  `design_values` TEXT NOT NULL,
  `years_label` VARCHAR(128) NOT NULL DEFAULT '',
  `address` VARCHAR(255) NOT NULL DEFAULT '',
  `logo_original_url` VARCHAR(512) NOT NULL DEFAULT '',
  `logo_thumb_url` VARCHAR(512) NOT NULL DEFAULT '',
  `logo_hor_original_url` VARCHAR(512) NOT NULL DEFAULT '',
  `logo_hor_thumb_url` VARCHAR(512) NOT NULL DEFAULT '',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `portfolio` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `slug` VARCHAR(64) NOT NULL,
  `sort_order` INT NOT NULL DEFAULT 0,
  `address` VARCHAR(255) NOT NULL DEFAULT '',
  `area` VARCHAR(64) NOT NULL DEFAULT '',
  `style` VARCHAR(128) NOT NULL DEFAULT '',
  `heart_flow` TEXT NOT NULL,
  `cover_original_url` VARCHAR(512) NOT NULL DEFAULT '',
  `cover_thumb_url` VARCHAR(512) NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_slug` (`slug`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `portfolio_image` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `portfolio_id` BIGINT UNSIGNED NOT NULL,
  `kind` VARCHAR(16) NOT NULL COMMENT 'render|real',
  `sort_order` INT NOT NULL DEFAULT 0,
  `original_url` VARCHAR(512) NOT NULL DEFAULT '',
  `thumb_url` VARCHAR(512) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `idx_portfolio_kind_sort` (`portfolio_id`, `kind`, `sort_order`),
  CONSTRAINT `fk_portfolio_image_portfolio`
    FOREIGN KEY (`portfolio_id`) REFERENCES `portfolio` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `pricing` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `original_url` VARCHAR(512) NOT NULL DEFAULT '',
  `thumb_url` VARCHAR(512) NOT NULL DEFAULT '',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
