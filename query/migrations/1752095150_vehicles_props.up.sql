BEGIN;

-- TABLE: `fivenet_vehicles_props` - Vehicle Properties
CREATE TABLE IF NOT EXISTS `fivenet_vehicles_props` (
    `plate` VARCHAR(12) NOT NULL,
    `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
    `wanted` tinyint(1) NOT NULL DEFAULT 0,
    `wanted_reason` VARCHAR(255) DEFAULT NULL,
    PRIMARY KEY (`plate`),
    KEY `idx_fivenet_vehicles_props_wanted` (`wanted`),
    CONSTRAINT `fk_fivenet_vehicles_props_plate` FOREIGN KEY (`plate`) REFERENCES `{{- if .ESXCompat }}owned_vehicles{{ else }}fivenet_owned_vehicles{{ end -}}` (`plate`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
