-- create company, user, vendor, vendor_bank_account and invoice tables at the initial stage
CREATE TABLE `companies` (
  `id` CHAR(36) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `representative_name` VARCHAR(255) NOT NULL,
  `phone_number` VARCHAR(20) NOT NULL,
  `zip_code` VARCHAR(10) NOT NULL,
  `address` VARCHAR(500) NOT NULL,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `deleted_at` DATETIME(6) NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = 'Companies';
CREATE TABLE `users` (
  `id` CHAR(36) NOT NULL,
  `company_id` CHAR(36) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) NOT NULL UNIQUE,
  `password` BINARY(60) NOT NULL COMMENT 'Hashed password using bcrypt',
  `role` TINYINT UNSIGNED NOT NULL COMMENT '{"0": "NORMAL", "1": "ADMIN"}',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `deleted_at` DATETIME(6) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_email_unique` (`email`),
  CONSTRAINT `fk_company_id_on_users` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = 'Users';
CREATE TABLE `vendors` (
  `id` CHAR(36) NOT NULL,
  `company_id` CHAR(36) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `representative_name` VARCHAR(255) NOT NULL,
  `phone_number` VARCHAR(20) NOT NULL,
  `zip_code` VARCHAR(10) NOT NULL,
  `address` VARCHAR(500) NOT NULL,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `deleted_at` DATETIME(6) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_company_id_on_vendors` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = 'Vendors';
CREATE TABLE `vendor_bank_accounts` (
  `id` CHAR(36) NOT NULL,
  `vendor_id` CHAR(36) NOT NULL,
  `bank_name` VARCHAR(255) NOT NULL,
  `branch_name` VARCHAR(255) NOT NULL,
  `account_number` VARCHAR(50) NOT NULL,
  `account_name` VARCHAR(255) NOT NULL,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `deleted_at` DATETIME(6) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_vendor_id_on_vendor_bank_accounts` FOREIGN KEY (`vendor_id`) REFERENCES `vendors` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = 'Vendor Bank Accounts';
CREATE TABLE `invoices` (
  `id` CHAR(36) NOT NULL,
  `company_id` CHAR(36) NOT NULL,
  `user_id` CHAR(36) NOT NULL,
  `vendor_id` CHAR(36) NOT NULL,
  `vendor_bank_account_id` CHAR(36) NOT NULL,
  `payment_amount` DECIMAL(10, 2) NOT NULL,
  `status` TINYINT UNSIGNED NOT NULL COMMENT '{"0": "UNPROCESSED", "1": "PROCESSING", "2": "PAID", "3": "ERROR"}',
  `fee` DECIMAL(10, 2) NOT NULL,
  `fee_rate` DECIMAL(5, 4) NOT NULL,
  `consumption_tax` DECIMAL(10, 2) NOT NULL,
  `consumption_tax_rate` DECIMAL(5, 4) NOT NULL,
  `billed_amount` DECIMAL(10, 2) NOT NULL,
  `issue_date` DATE NOT NULL,
  `due_date` DATE NOT NULL,
  `created_by` CHAR(36) NOT NULL,
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_by` CHAR(36) NOT NULL,
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `deleted_by` CHAR(36) NULL DEFAULT NULL,
  `deleted_at` DATETIME(6) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_due_date` (`due_date`),
  CONSTRAINT `fk_company_id_on_invoices` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_user_id_on_invoices` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_vendor_id_on_invoices` FOREIGN KEY (`vendor_id`) REFERENCES `vendors` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_vendor_bank_account_id_on_invoices` FOREIGN KEY (`vendor_bank_account_id`) REFERENCES `vendor_bank_accounts` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = 'Invoices';