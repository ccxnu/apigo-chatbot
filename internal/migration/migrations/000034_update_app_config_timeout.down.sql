-- Rollback: Restore APP_CONFIG contextTimeout to 2 seconds

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{contextTimeout}', '2')
WHERE prm_code = 'APP_CONFIG';
