-- Update APP_CONFIG contextTimeout from 2 to 10 seconds

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{contextTimeout}', '10')
WHERE prm_code = 'APP_CONFIG';
