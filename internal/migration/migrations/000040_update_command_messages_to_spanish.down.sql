-- =====================================================
-- Rollback Migration 000040
-- =====================================================
-- Restore original English command messages

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"👋 *Bienvenido al Asistente del Instituto*\n\nSoy tu asistente virtual y puedo ayudarte con:\n\n🎓 *Información Académica*\n   • Programas y carreras\n   • Requisitos de admisión\n   • Proceso de matrícula\n   • Calendario académico\n\n📚 *Consultas Generales*\n   Solo escribe tu pregunta y te ayudaré a encontrar la información que necesitas.\n\n⚡ *Comandos Disponibles*\n   /help - Muestra esta ayuda\n   /horarios - Consulta horarios de clases\n   /comandos - Lista todos los comandos\n\n💬 También puedes hacer preguntas directamente, por ejemplo:\n   \"¿Cuál es el proceso de matrícula?\"\n \"¿Qué carreras ofrecen?\"\n\n¿En qué puedo ayudarte hoy?"'::jsonb
)
WHERE prm_code = 'MESSAGE_HELP';

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"⚡ *Comandos Disponibles*\n\n/help - Muestra ayuda general del bot\n/horarios - Consulta horarios de clases\n/comandos - Muestra esta lista de comandos\n/start - Reinicia la conversación\n\n💡 *Tip*: No necesitas usar comandos para hacer preguntas. ¡Solo escribe tu consulta!"'::jsonb
)
WHERE prm_code = 'MESSAGE_COMMANDS';

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"👋 ¡Hola! Soy el asistente virtual del Instituto.\n\nEstoy aquí para ayudarte con información sobre:\n • Programas académicos\n   • Admisiones y matrículas\n   • Horarios y calendarios\n • Y mucho más...\n\nEscribe /help para ver todo lo que puedo hacer, o simplemente hazme una pregunta.\n\n¿En qué puedo ayudarte?"'::jsonb
)
WHERE prm_code = 'MESSAGE_START';

UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{message}',
    '"❓ Comando no reconocido.\n\nEscribe /help para ver los comandos disponibles, o simplemente hazme tu pregunta directamente."'::jsonb
)
WHERE prm_code = 'MESSAGE_UNKNOWN_COMMAND';
