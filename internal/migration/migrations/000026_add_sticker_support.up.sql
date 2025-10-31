-- Add sticker configuration parameter
INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
VALUES (
    'STICKER_CONFIG',
    'BOT_STICKERS',
    '{
        "enabled": true,
        "stickers": {
            "welcome": "https://example.com/stickers/welcome.webp",
            "thanks": "https://example.com/stickers/thanks.webp",
            "thinking": "https://example.com/stickers/thinking.webp",
            "happy": "https://example.com/stickers/happy.webp",
            "confused": "https://example.com/stickers/confused.webp"
        },
        "usage": "El bot puede enviar stickers para expresar emociones. Configura las URLs de tus stickers en formato WebP."
    }'::jsonb,
    'Sticker URLs for bot expressions. Update URLs to your own hosted stickers.'
)
ON CONFLICT (prm_code) DO UPDATE
SET prm_data = EXCLUDED.prm_data,
    prm_description = EXCLUDED.prm_description;

COMMENT ON COLUMN cht_conversation_messages.cvm_message_type IS 'Message type: text, image, document, audio, video, or sticker';
