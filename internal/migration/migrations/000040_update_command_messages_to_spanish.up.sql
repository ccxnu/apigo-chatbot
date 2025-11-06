-- =====================================================
-- Migration 000040: Update command messages to Spanish
-- =====================================================
-- Update MESSAGE_HELP, MESSAGE_COMMANDS, MESSAGE_START, MESSAGE_UNKNOWN_COMMAND
-- to use Spanish commands as primary

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"👋 *Bienvenido al Asistente del Instituto*\n\nSoy tu asistente virtual y puedo ayudarte con:\n\n🎓 *Información Académica*\n   • Programas y carreras\n   • Requisitos de admisión\n   • Proceso de matrícula\n   • Calendario académico\n\n📚 *Consultas Generales*\n   Solo escribe tu pregunta y te ayudaré a encontrar la información que necesitas.\n\n⚡ *Comandos Disponibles*\n   /ayuda - Muestra esta ayuda\n   /inicio - Mensaje de bienvenida\n   /horarios - Consulta horarios disponibles\n   /registrar - Registrarse en el sistema\n   /cancelar - Cancelar registro en curso\n   /comandos - Lista todos los comandos\n\n💬 También puedes hacer preguntas directamente, por ejemplo:\n   \"¿Cuál es el proceso de matrícula?\"\n   \"¿Qué carreras ofrecen?\"\n\n¿En qué puedo ayudarte hoy?"'::jsonb
)
WHERE prm_code = 'MESSAGE_HELP';

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"⚡ *Comandos Disponibles*\n\n/ayuda - Ayuda general del bot\n/inicio - Mensaje de bienvenida\n/horarios - Consulta horarios disponibles\n/registrar - Registrarse en el sistema\n/cancelar - Cancelar registro en curso\n/comandos - Ver esta lista de comandos\n\n💡 _También puedes escribir tus preguntas directamente sin usar comandos._"'::jsonb
)
WHERE prm_code = 'MESSAGE_COMMANDS';

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"👋 ¡Hola! Soy el asistente virtual del Instituto.\n\nEstoy aquí para ayudarte con información sobre:\n   • Programas académicos\n   • Admisiones y matrículas\n   • Horarios y calendarios\n   • Y mucho más...\n\nEscribe /ayuda para ver todo lo que puedo hacer, o simplemente hazme una pregunta.\n\n¿En qué puedo ayudarte?"'::jsonb
)
WHERE prm_code = 'MESSAGE_START';

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"❓ Comando no reconocido.\n\nUsa /comandos para ver los comandos disponibles."'::jsonb
)
WHERE prm_code = 'MESSAGE_UNKNOWN_COMMAND';

COMMENT ON TABLE cht_parameters IS 'Updated command messages to use Spanish as primary language';
