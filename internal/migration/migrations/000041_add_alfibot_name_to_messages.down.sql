-- =====================================================
-- Rollback Migration 000041
-- =====================================================
-- Remove "Alfibot" name from messages

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"👋 *Bienvenido al Asistente del Instituto*\n\nSoy tu asistente virtual y puedo ayudarte con:\n\n🎓 *Información Académica*\n   • Programas y carreras\n   • Requisitos de admisión\n   • Proceso de matrícula\n   • Calendario académico\n\n📚 *Consultas Generales*\n   Solo escribe tu pregunta y te ayudaré a encontrar la información que necesitas.\n\n⚡ *Comandos Disponibles*\n   /ayuda - Muestra esta ayuda\n   /inicio - Mensaje de bienvenida\n   /horarios - Consulta horarios disponibles\n   /registrar - Registrarse en el sistema\n   /cancelar - Cancelar registro en curso\n   /comandos - Lista todos los comandos\n\n💬 También puedes hacer preguntas directamente, por ejemplo:\n   \"¿Cuál es el proceso de matrícula?\"\n   \"¿Qué carreras ofrecen?\"\n\n¿En qué puedo ayudarte hoy?"'::jsonb
)
WHERE prm_code = 'MESSAGE_HELP';

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"👋 ¡Hola! Soy el asistente virtual del Instituto.\n\nEstoy aquí para ayudarte con información sobre:\n   • Programas académicos\n   • Admisiones y matrículas\n   • Horarios y calendarios\n   • Y mucho más...\n\nEscribe /ayuda para ver todo lo que puedo hacer, o simplemente hazme una pregunta.\n\n¿En qué puedo ayudarte?"'::jsonb
)
WHERE prm_code = 'MESSAGE_START';

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"👋 ¡Hola! Bienvenido al asistente virtual del Instituto.\n\nPara poder ayudarte, necesito que te registres primero.\n\nPor favor, envíame tu número de cédula (10 dígitos).\n\nEjemplo: 1234567890"'::jsonb
)
WHERE prm_code = 'MESSAGE_REQUEST_CEDULA';

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"📊 Has alcanzado el límite de mensajes para usuarios no registrados.\n\n✅ Para continuar chateando sin límites, regístrate usando:\n\n/registrar"'::jsonb
)
WHERE prm_code = 'MESSAGE_GUEST_LIMIT_REACHED';

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"⚠️ Te queda %d mensaje disponible hoy.\n\n💡 Regístrate con /registrar para chat ilimitado."'::jsonb
)
WHERE prm_code = 'MESSAGE_GUEST_LIMIT_WARNING';

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"👋 ¡Hola %s! Bienvenido.\n\n¿En qué puedo ayudarte hoy?"'::jsonb
)
WHERE prm_code = 'MESSAGE_WELCOME_REGISTERED';

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"⏳ Estoy esperando tu número de cédula para continuar con el registro.\n\nPor favor envía tu cédula de 10 dígitos.\n\nEjemplo: 1234567890"'::jsonb
)
WHERE prm_code = 'REG_MSG_WAITING_CEDULA';

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"⏳ Necesito tu nombre completo y correo electrónico.\n\nFormato: *Nombre Completo / correo@email.com*\n\nEjemplo:\nJuan Pérez / juan.perez@gmail.com"'::jsonb
)
WHERE prm_code = 'REG_MSG_WAITING_EMAIL_NAME';
