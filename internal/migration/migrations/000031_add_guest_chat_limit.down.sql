-- Remove guest chat limit parameters
DELETE FROM cht_parameters WHERE prm_code IN (
    'GUEST_CHAT_LIMIT',
    'MESSAGE_GUEST_LIMIT_REACHED',
    'MESSAGE_GUEST_LIMIT_WARNING'
);
