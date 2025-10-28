#set document(
  title: "Manual de Usuario - Chatbot ISTS",
  author: "Pablo Moisés Cuenca Cuenca",
  date: datetime.today(),
)

#set page(
  paper: "us-letter",
  margin: (x: 1.18in, y: 1.18in),
  numbering: "1",
)

#set text(
  font: "Times New Roman",
  size: 12pt,
  lang: "es",
)

#show figure: set align(start)
#show figure.caption: set text(10pt)

#show heading.where(level: 1): set text(size: 1em/1.4)
#show heading.where(level: 2): set text(size: 1em/1.2)

#set heading(numbering: "1.1")

// Title page
#align(center)[
  *Instituto Superior Tecnológico Sudamericano*

  #image("assets/logo.png")

  #v(3em)

  #upper[
  Implementación de un asistente virtual basado en RAG y LLM como solución
  de negocio para la atención automatizada de servicios educativos
  ]

  #v(3em)

  #upper[*Manual de Usuario*]

  #v(3em)

  *Autor:*

  Pablo Moisés Cuenca

  #v(3em)

  *Director:*

  Ing. Cristhian Javier Villamarin Gaona

  #v(3em)

  Loja - Ecuador \
  2025
]

#pagebreak()

// Table of contents
#outline(
  title: [Índices],
  indent: auto,
)

#pagebreak()

= Introducción

== Bienvenida

Bienvenido al *Manual de Usuario del Sistema de Chatbot ISTS*.
Este documento está diseñado para ayudarle a comprender y utilizar
eficientemente el chatbot institucional a través de WhatsApp.

El Chatbot ISTS es un asistente virtual inteligente que proporciona
información sobre la institución, carreras, procesos administrativos,
y responde preguntas frecuentes de manera automática las 24 horas del día,
los 7 días de la semana.

#figure(
  image("assets/whatsapp_chat.png", width: 75%),
  caption: [Conversación con el chatbot]
)

== Objetivos del Manual

Este manual tiene los siguientes objetivos:

- Explicar cómo interactuar con el chatbot de manera efectiva
- Describir las funcionalidades disponibles para estudiantes y público general
- Proporcionar guías paso a paso para realizar consultas comunes
- Resolver problemas frecuentes de uso

== A Quién Está Dirigido

Este manual está dirigido a:

- *Estudiantes actuales* del ISTS
- *Aspirantes* interesados en estudiar en el instituto
- *Padres de familia* que buscan información
- *Público general* con consultas sobre la institución

== Requisitos Previos

Para utilizar el chatbot ISTS, necesita:

- Un celular inteligente con *WhatsApp* instalado
- Conexión a *Internet* (datos móviles o WiFi)
- El número de WhatsApp del chatbot institucional

#pagebreak()

= Conceptos Básicos

== ¿Qué es el Chatbot ISTS?

El Chatbot ISTS es un *asistente virtual automatizado* que utiliza inteligencia artificial
para responder preguntas sobre el instituto. Funciona a través de WhatsApp,
la aplicación de mensajería más popular, lo que facilita el acceso desde cualquier dispositivo móvil.

#figure(
  image("assets/welcome_chat.png", width: 70%),
  caption: [Mensaje de bienvenida del chatbot]
)

=== Características Principales

El chatbot ofrece:

- *Respuestas instantáneas* 24/7 sin necesidad de esperar atención humana
- *Búsqueda inteligente* en la base de conocimiento institucional
- *Información actualizada* sobre carreras, requisitos, fechas y procesos
- *Interfaz conversacional* natural en español
- *Consultas ilimitadas* sin costo adicional

=== Tecnología Utilizada

El sistema utiliza:

- *RAG (Retrieval Augmented Generation)*: Búsqueda semántica + generación de respuestas con IA
- *Procesamiento de Lenguaje Natural*: Comprende preguntas en lenguaje cotidiano
- *Base de Conocimiento Vectorial*: Encuentra información relevante rápidamente
- *Integración WhatsApp*: Comunicación a través de la plataforma familiar

== Ventajas de Usar el Chatbot

=== Disponibilidad 24/7

El chatbot está disponible en todo momento, incluso fuera del horario administrativo, fines de semana y feriados.

=== Respuestas Inmediatas

No necesita esperar en línea o agendar citas. Las respuestas son instantáneas.

=== Información Actualizada

El chatbot busca información en documentos oficiales de la institución, garantizando precisión.

=== Historial de Conversación

Puede revisar conversaciones anteriores en WhatsApp cuando necesite recordar información.

= Primeros Pasos

== Agregar el Chatbot a WhatsApp

=== Paso 1: Obtener el Número

Solicite el número oficial del Chatbot ISTS a través de:

- La página web institucional: #link("https://tecnologicosudamericano.edu.ec")
- Las redes sociales oficiales del instituto
- La oficina de información y atención al estudiante

#figure(
  image("assets/suda_page.png", width: 90%),
  caption: [Página oficial del Instituto Superior Tecnológico Sudamericano]
)

=== Paso 2: Agregar el Contacto


+ Abra WhatsApp en su teléfono
+ Toque el botón *"+"* o *"Nuevo chat"*
+ Seleccione *"Nuevo contacto"*
  #figure(
    image("assets/whatsapp_add_contact.jpeg", width: 60%),
    caption: [Proceso de agregar contacto en WhatsApp]
  )
+ Ingrese el número del chatbot
+ Guarde el contacto como *"Chatbot ISTS"* o similar
  #figure(
    image("assets/whatsapp_name_contact.jpeg", width: 60%),
    caption: [Ingresar datos del chatbot]
  )


=== Paso 3: Iniciar la Conversación

+ Abra el chat con el contacto guardado
+ Envíe un mensaje de saludo como *"Hola"* o *"Buenos días"*
+ El chatbot responderá automáticamente con un mensaje de registro o bienvenida

#figure(
  image("assets/whatsapp_first_message.jpeg", width: 45%),
  caption: [Primera conversación con el chatbot]
)

== Primer Mensaje de Bienvenida

Cuando inicie la conversación, el chatbot le enviará un mensaje de bienvenida que incluye:

- Presentación del servicio
- Instrucciones básicas de uso
- Ejemplos de preguntas que puede hacer
- Horarios de atención humana (si necesita soporte personalizado)

#pagebreak()

= Registro

= Uso del Chatbot

== Cómo Hacer Preguntas

=== Lenguaje Natural

El chatbot comprende preguntas en lenguaje natural. No necesita usar comandos especiales o palabras clave exactas.

*Ejemplos correctos:*

```
¿Cuándo son las inscripciones?
Quiero información sobre la carrera de software
¿Qué documentos necesito para matricularme?
Horarios de atención
```

// TODO: Add screenshot showing example questions
// #figure(
//   image("images/example_questions.png", width: 70%),
//   caption: [Ejemplos de preguntas al chatbot]
// )

=== Preguntas Claras y Específicas

Para obtener mejores respuestas:

- *Sea específico*: "¿Cuál es el costo de la matrícula en Desarrollo de Software?" es mejor que "¿Cuánto cuesta?"
- *Una pregunta a la vez*: Evite preguntas múltiples en un solo mensaje
- *Use ortografía correcta*: Aunque el bot tolera errores, la claridad mejora las respuestas

=== Reformular Preguntas

Si no obtiene la respuesta esperada:

- Reformule la pregunta con otras palabras
- Divida preguntas complejas en varias simples
- Use sinónimos o términos alternativos

== Tipos de Consultas Disponibles

=== Información sobre Carreras

Puede consultar sobre:

- Carreras disponibles y modalidades
- Mallas curriculares y asignaturas
- Duración de las carreras
- Perfil de egreso y campo laboral
- Requisitos de ingreso

*Ejemplo:*
```
Usuario: ¿Qué carreras técnicas tienen?
Bot: El ISTS ofrece las siguientes carreras técnicas:
     1. Desarrollo de Software
     2. Diseño Gráfico
     3. Contabilidad
     4. Administración de Empresas
     ...
```

// TODO: Add screenshot of career information response
// #figure(
//   image("images/career_query_example.png", width: 70%),
//   caption: [Consulta sobre carreras disponibles]
// )

=== Procesos Administrativos

Consulte sobre:

- Inscripciones y matriculación
- Documentos requeridos
- Costos y aranceles
- Becas y ayuda financiera
- Trámites académicos

*Ejemplo:*
```
Usuario: ¿Qué necesito para inscribirme?
Bot: Para inscribirte necesitas:
     1. Cédula de identidad (original y copia)
     2. Título de bachiller (original y copia)
     3. Certificado de votación
     ...
```

// TODO: Add screenshot of enrollment requirements response
// #figure(
//   image("images/enrollment_requirements.png", width: 70%),
//   caption: [Información sobre requisitos de inscripción]
// )

=== Fechas y Horarios

Información sobre:

- Calendario académico
- Fechas de inscripción
- Horarios de atención
- Eventos institucionales
- Período de exámenes

=== Ubicación y Contacto

Consulte:

- Dirección de la institución
- Números de teléfono
- Correos electrónicos
- Redes sociales oficiales
- Cómo llegar al campus

== Entender las Respuestas

=== Respuestas Directas

El chatbot proporciona respuestas directas y concisas cuando la información es específica.

=== Respuestas con Opciones

Para temas amplios, el bot puede ofrecer opciones para que elija:

```
Bot: Tenemos información sobre:
     1. Requisitos de admisión
     2. Proceso de inscripción
     3. Documentos necesarios

     ¿Sobre cuál te gustaría saber más?
```

// TODO: Add screenshot of bot offering multiple options
// #figure(
//   image("images/multiple_options_response.png", width: 70%),
//   caption: [Respuesta con múltiples opciones]
// )

=== Referencias a Documentos

El bot puede mencionar documentos oficiales como fuente de información.

=== Derivación a Atención Humana

Si su consulta es muy específica o requiere atención personalizada, el bot le indicará cómo contactar con personal administrativo.

= Casos de Uso Comunes

== Consulta sobre Inscripciones

*Escenario:* Desea inscribirse y necesita saber fechas y requisitos.

// TODO: Add screenshot of complete enrollment conversation flow
// #figure(
//   image("images/enrollment_conversation_flow.png", width: 80%),
//   caption: [Flujo de conversación sobre inscripciones]
// )

*Pasos:*

+ Pregunte: _"¿Cuándo son las inscripciones?"_
+ El bot le proporcionará las fechas del período actual
+ Pregunte: _"¿Qué documentos necesito?"_
+ Revise la lista proporcionada
+ Si tiene dudas adicionales, formule preguntas específicas

*Resultado esperado:* Información completa sobre fechas, requisitos y proceso.

== Información sobre una Carrera Específica

*Escenario:* Está interesado en la carrera de Desarrollo de Software.

*Pasos:*

+ Pregunte: _"Información sobre la carrera de Desarrollo de Software"_
+ El bot proporcionará un resumen general
+ Para más detalles, pregunte específicamente:
  - _"¿Cuánto dura la carrera?"_
  - _"¿Qué materias se estudian?"_
  - _"¿Cuál es el costo?"_

*Resultado esperado:* Información detallada sobre la carrera de su interés.

== Consulta de Costos y Aranceles

*Escenario:* Necesita saber cuánto cuesta estudiar en el ISTS.

*Pasos:*

+ Pregunte: _"¿Cuánto cuesta la matrícula?"_
+ El bot indicará el costo de matrícula
+ Pregunte: _"¿Y la pensión mensual?"_
+ Para becas, pregunte: _"¿Hay becas disponibles?"_

*Resultado esperado:* Información clara sobre costos y opciones de financiamiento.

== Ubicación y Contacto

*Escenario:* Necesita visitar el instituto o contactar a alguien.

*Pasos:*

+ Pregunte: _"¿Dónde está ubicado el instituto?"_
+ El bot proporcionará la dirección
+ Para contacto, pregunte: _"¿Cuál es el teléfono?"_
+ Para horarios: _"¿Cuál es el horario de atención?"_

*Resultado esperado:* Información completa de ubicación, contacto y horarios.

= Solución de Problemas

== El Chatbot No Responde

*Posibles causas y soluciones:*

// TODO: Add screenshot showing connection issues (one check mark)
// #figure(
//   image("images/connection_issue.png", width: 60%),
//   caption: [Mensaje sin entregar - un solo check]
// )

=== Problemas de Conexión

*Síntoma:* El mensaje no se envía o aparece con un solo check.

*Solución:*
- Verifique su conexión a Internet
- Intente conectarse a WiFi si usa datos móviles
- Reinicie WhatsApp
- Espere unos minutos y reintente

=== Número Incorrecto

*Síntoma:* WhatsApp indica que el número no existe.

*Solución:*
- Verifique que ingresó el número correctamente
- Consulte el número oficial en la página web del instituto
- Asegúrese de incluir el código de país

=== Servidor Ocupado

*Síntoma:* El mensaje se entrega pero no hay respuesta inmediata.

*Solución:*
- Espere unos segundos, el bot puede estar procesando múltiples consultas
- Si no responde en 1 minuto, reenvíe el mensaje
- En horario pico puede haber mayor latencia

== Respuestas Incorrectas o Confusas

*Situación:* El bot no entendió su pregunta o dio información incorrecta.

*Soluciones:*

+ *Reformule la pregunta*: Use palabras diferentes o más simples
+ *Sea más específico*: Agregue contexto o detalles adicionales
+ *Divida la consulta*: Haga una pregunta a la vez
+ *Reporte el problema*: Mencione al bot que la respuesta no fue útil

*Ejemplo:*
```
Usuario: La respuesta no fue clara
Bot: Disculpa, déjame intentar explicarte mejor.
     ¿Podrías reformular tu pregunta?
```

== No Encuentra la Información

*Situación:* El bot indica que no tiene información sobre su consulta.

*Acciones:*

+ *Use términos alternativos*: Pruebe sinónimos o formas diferentes de preguntar
+ *Pregunte de forma más general*: Amplíe la consulta
+ *Contacte atención humana*: Para temas muy específicos o nuevos, solicite contacto directo

*Ejemplo:*
```
Usuario: ¿Tienen descuento para hermanos?
Bot: No encuentro información específica sobre eso.
     Por favor, contacta a admisiones al 02-XXXXXXX
     para consultas sobre descuentos especiales.
```

== Bloqueo o Suspensión del Servicio

*Situación:* Ya no puede enviar mensajes al chatbot.

*Causas posibles:*
- Uso indebido o spam
- Problemas técnicos temporales
- Mantenimiento del sistema

*Solución:*
- Espere 24 horas e intente nuevamente
- Contacte a soporte técnico institucional
- Verifique anuncios oficiales sobre mantenimiento

= Buenas Prácticas

== Uso Responsable

=== Evite el Spam

- No envíe mensajes repetitivos innecesariamente
- Espere la respuesta antes de reenviar
- Una sola pregunta clara es mejor que múltiples mensajes

=== Respete el Propósito

El chatbot está diseñado para:
- Información institucional
- Consultas académicas
- Procesos administrativos

*No use el chatbot para:*
- Conversaciones personales sin relación con el instituto
- Mensajes ofensivos o inapropiados
- Publicidad o promociones

=== Verifique Información Crítica

Para decisiones importantes:
- Confirme información crucial con personal administrativo
- Consulte documentos oficiales cuando esté disponible
- Use el chatbot como primera fuente, pero verifique datos críticos

== Maximizar la Efectividad

=== Prepare Sus Preguntas

Antes de consultar:
- Anote sus dudas específicas
- Identifique qué información exactamente necesita
- Tenga a mano documentos si necesita hacer referencia

=== Use Horario Apropiado

Aunque el bot está disponible 24/7:
- Consultas complejas son mejor respondidas en horario laboral
- Si necesita derivación a humano, contacte en horario de oficina
- Mantenimientos programados suelen ser en madrugada

=== Guarde Información Importante

- Tome capturas de pantalla de información relevante
- Guarde respuestas importantes en notas de WhatsApp
- Copie números de contacto y fechas importantes

= Preguntas Frecuentes (FAQ)

== General

=== ¿El chatbot cobra por las respuestas?

No, el servicio del chatbot es completamente *gratuito*. Solo consume sus datos de Internet o WiFi de WhatsApp.

=== ¿Puedo usar el chatbot desde otro país?

Sí, siempre que tenga acceso a WhatsApp y conexión a Internet.

=== ¿El chatbot guarda mi información personal?

El chatbot registra las conversaciones para mejorar el servicio, pero respeta la privacidad. No comparte información personal con terceros.

=== ¿Hay límite de preguntas?

No hay límite. Puede hacer todas las consultas que necesite.

== Funcionalidad

=== ¿El chatbot entiende errores de ortografía?

Sí, el sistema tolera errores ortográficos menores, aunque la claridad mejora las respuestas.

=== ¿Puedo enviar imágenes o documentos?

Actualmente, el chatbot solo procesa texto. No puede analizar imágenes o archivos adjuntos.

=== ¿El chatbot recuerda conversaciones anteriores?

El bot tiene contexto de la conversación actual, pero cada nueva sesión es independiente.

=== ¿Qué idiomas soporta?

El chatbot está configurado para *español*. Puede entender algunas frases en inglés, pero responde en español.

== Soporte

=== ¿Qué hago si encuentro un error?

Reporte el problema:
- Mencione al bot que la respuesta fue incorrecta
- Contacte al soporte técnico institucional
- Envíe un correo a #link("soporte@ists.edu.ec")

=== ¿Cómo contacto con una persona real?

Solicite al bot: _"Quiero hablar con un humano"_ o _"Necesito atención personalizada"_. El bot le proporcionará los contactos correspondientes.

=== ¿Hay horarios específicos para soporte humano?

El chatbot funciona 24/7, pero el soporte humano está disponible en horario de oficina:
- Lunes a Viernes: 8:00 AM - 5:00 PM
- Sábados: 8:00 AM - 1:00 PM

= Contacto y Soporte

== Canales de Atención

=== Oficina de Admisiones
- *Teléfono:* (02) XXX-XXXX
- *Email:* admisiones\@ists.edu.ec
- *Horario:* Lunes a Viernes 8:00-17:00

=== Secretaría Académica
- *Teléfono:* (02) XXX-XXXX
- *Email:* secretaria\@ists.edu.ec
- *Horario:* Lunes a Viernes 8:00-17:00

=== Soporte Técnico del Chatbot
- *Email:* soporte\@ists.edu.ec
- *Asunto:* "Soporte Chatbot WhatsApp"

=== Redes Sociales Oficiales
- *Facebook:* /ISTSSudamericano
- *Instagram:* \@ists_oficial
- *Twitter:* \@ISTS_Ecuador
- *Website:* www.ists.edu.ec

== Ubicación del Instituto

// TODO: Add map showing institute location
// #figure(
//   image("images/ists_location_map.png", width: 90%),
//   caption: [Ubicación del Instituto Superior Tecnológico Sudamericano]
// )

*Dirección:*
Av. Principal N°123 y Calle Secundaria
Quito, Ecuador

*Referencias:*
Frente al parque central, junto al edificio municipal

*Transporte público:*
Líneas de bus: 10, 15, 23
Parada: "ISTS"

== Feedback y Sugerencias

Su opinión es importante. Para sugerencias sobre el chatbot:

+ Envíe un email a: feedback\@ists.edu.ec
+ Asunto: "Sugerencia Chatbot"
+ Incluya:
  - Su experiencia de uso
  - Qué funcionó bien
  - Qué podría mejorar
  - Nuevas funcionalidades que le gustaría

= Anexos

== Glosario de Términos

=== Chatbot
Programa de computadora que simula conversaciones humanas mediante inteligencia artificial.

=== RAG (Retrieval Augmented Generation)
Tecnología que combina búsqueda de información con generación de respuestas por IA.

=== WhatsApp Business API
Plataforma que permite a organizaciones comunicarse con usuarios a través de WhatsApp.

=== Base de Conocimiento
Colección de documentos e información institucional que el chatbot utiliza para responder.

=== Procesamiento de Lenguaje Natural (NLP)
Tecnología que permite a las computadoras entender y procesar lenguaje humano.

=== Vector Embeddings
Representación matemática de texto que permite búsquedas semánticas inteligentes.

== Lista de Comandos Útiles

Aunque el chatbot entiende lenguaje natural, estos comandos pueden ser útiles:

- *"Hola"* / *"Ayuda"*: Mensaje de bienvenida e instrucciones
- *"Menú"*: Lista de temas disponibles
- *"Carreras"*: Información sobre carreras disponibles
- *"Inscripciones"*: Proceso y fechas de inscripción
- *"Contacto"*: Información de contacto y ubicación
- *"Costos"*: Información sobre aranceles
- *"Requisitos"*: Documentos necesarios para matriculación

== Actualizaciones del Manual

*Versión 1.0* - Octubre 2025
- Versión inicial del manual de usuario
- Información completa sobre uso del chatbot
- Casos de uso y solución de problemas

*Próximas actualizaciones incluirán:*
- Nuevas funcionalidades del chatbot
- Panel de administración para usuarios autorizados
- Estadísticas y analytics de uso

---

#align(center)[
  #text(size: 10pt, style: "italic")[
    Este manual está sujeto a actualizaciones periódicas. \
    Consulte la versión más reciente en www.ists.edu.ec/manuales
  ]

  #v(1cm)

  #text(size: 9pt)[
    © 2025 Instituto Superior Tecnológico Sudamericano \
    Todos los derechos reservados
  ]
]
