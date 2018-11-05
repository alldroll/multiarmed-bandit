CREATE TABLE `experiment` (
    `name` VARCHAR(255) NOT NULL,
    `gamma` FLOAT NOT NULL,
    `enabled` TINYINT(1) NOT NULL DEFAULT 0,
    PRIMARY KEY (`name`),
    KEY `enabled` (`enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8;

CREATE TABLE `variant` (
    `id` TINYINT(3) UNSIGNED NOT NULL,
    `experiment_name` VARCHAR(255) NOT NULL,
    `shows` INT(11) UNSIGNED NOT NULL DEFAULT 0,
    `rewards` INT(11) UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`experiment_name`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8;
