-- Add guest/unregistered user chat limit parameter
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'GUEST_CHAT_LIMIT') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES (
            'CHAT_CONFIGURATION',
            'GUEST_CHAT_LIMIT',
            '{"value": 5}'::jsonb,
            'Maximum number of messages unregistered users can send per day'
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'MESSAGE_GUEST_LIMIT_REACHED') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES (
            'CHAT_CONFIGURATION',
            'MESSAGE_GUEST_LIMIT_REACHED',
            '{
                "message": "📊 Has alcanzado el límite de mensajes para usuarios no registrados.\n\n✅ Para continuar chateando sin límites, regístrate usando el comando:\n\n/register\n\nEl registro es rápido y te permite:\n🎓 Acceso ilimitado al asistente\n📚 Respuestas personalizadas según tu perfil\n⚡ Mejor experiencia de uso"
            }'::jsonb,
            'Message shown when unregistered user reaches chat limit'
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'MESSAGE_GUEST_LIMIT_WARNING') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES (
            'CHAT_CONFIGURATION',
            'MESSAGE_GUEST_LIMIT_WARNING',
            '{
                "template": "⚠️ Te quedan %d mensajes disponibles hoy.\n\n💡 Regístrate con /register para chat ilimitado."
            }'::jsonb,
            'Warning message showing remaining messages for unregistered users'
        );
    END IF;
END $$;
