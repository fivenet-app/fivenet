BEGIN;

-- Inactive wanted and zero fine values are represented by an absent property.
-- This lets empty_value_display control whether they appear in profile views.
DELETE state
FROM `fivenet_user_property_state` state
INNER JOIN `fivenet_user_property_values` value
  ON value.`user_id` = state.`user_id` AND value.`property_key` = state.`property_key`
WHERE value.`property_key` = 'user.wanted'
  AND value.`value_bool` = FALSE;

DELETE FROM `fivenet_user_property_values`
WHERE (`property_key` = 'user.wanted' AND `value_bool` = FALSE)
   OR (`property_key` = 'user.open_fines' AND `value_integer` <= 0);

COMMIT;
