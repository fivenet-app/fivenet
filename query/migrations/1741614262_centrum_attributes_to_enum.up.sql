BEGIN;

UPDATE `fivenet_centrum_units` SET `attributes` = REPLACE(`attributes`, '"no_dispatch_auto_assign"', '"UNIT_ATTRIBUTE_NO_DISPATCH_AUTO_ASSIGN"') WHERE `attributes` LIKE '%"no_dispatch_auto_assign"%';
UPDATE `fivenet_centrum_units` SET `attributes` = REPLACE(`attributes`, '"static"', '"UNIT_ATTRIBUTE_STATIC"') WHERE `attributes` LIKE '%"static"%';

UPDATE `fivenet_centrum_dispatches` SET `attributes` = REPLACE(`attributes`, '"multiple"', '"DISPATCH_ATTRIBUTE_MULTIPLE"') WHERE `attributes` LIKE '%"multiple"%';
UPDATE `fivenet_centrum_dispatches` SET `attributes` = REPLACE(`attributes`, '"duplicate"', '"DISPATCH_ATTRIBUTE_DUPLICATE"') WHERE `attributes` LIKE '%"duplicate"%';
UPDATE `fivenet_centrum_dispatches` SET `attributes` = REPLACE(`attributes`, '"too_old"', '"DISPATCH_ATTRIBUTE_TOO_OLD"') WHERE `attributes` LIKE '%"too_old"%';

COMMIT;
