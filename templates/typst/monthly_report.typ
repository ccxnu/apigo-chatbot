// Monthly Analytics Report Template
// Data loaded from JSON string passed via --input data parameter

// Load data from JSON string input (convert string to bytes first)
#let data = json(bytes(sys.inputs.data))

#set document(
  title: "Reporte Mensual - " + data.month_year,
  author: "ISTS Chatbot",
  date: datetime.today()
)

#set page(
  paper: "a4",
  margin: (x: 2.5cm, y: 2cm),
  numbering: "1 / 1",
  number-align: center,
  header: align(right)[
    #text(size: 9pt, fill: gray)[
      Reporte de Analíticas del Chatbot
    ]
  ],
  footer: align(center)[
    #line(length: 100%, stroke: 0.5pt + gray)
    #text(size: 9pt, fill: gray)[
      Instituto Superior Tecnológico Sudamericano · #context counter(page).display("1 / 1")
    ]
  ]
)

#set text(
  font: "Arial",
  size: 11pt,
  lang: "es"
)

#set heading(numbering: "1.1")

// Title page
#align(center)[
  #v(3cm)
  // TODO: Add logo.png to templates/typst directory
  // #image("logo.png", width: 40%)
  #text(size: 20pt, weight: "bold", fill: blue.darken(30%))[
    INSTITUTO SUPERIOR TECNOLÓGICO SUDAMERICANO
  ]
  #v(1.5cm)
  #text(size: 24pt, weight: "bold")[
    Reporte Mensual de Analíticas
  ]
  #v(0.5cm)
  #text(size: 18pt)[
    Chatbot Institucional
  ]
  #v(0.5cm)
  #text(size: 16pt, fill: blue.darken(30%))[
    #data.month_year
  ]
  #v(2cm)
  #text(size: 12pt)[
    Generado: #data.generated_date
  ]
]

#pagebreak()

// Table of contents
#outline(
  title: "Índice",
  indent: auto
)

#pagebreak()

// ============================================
// EXECUTIVE SUMMARY
// ============================================

= Resumen Ejecutivo

Este reporte presenta las métricas clave del sistema de chatbot del Instituto Superior Tecnológico Sudamericano para el período de *#data.month_year*.

== Métricas Principales

#grid(
  columns: (1fr, 1fr),
  gutter: 1em,

  // Cost Card
  box(
    fill: blue.lighten(90%),
    inset: 1em,
    radius: 5pt,
    width: 100%,
  )[
    #text(size: 10pt, fill: gray)[*Costo Total*]
    #v(0.3em)
    #text(size: 20pt, weight: "bold", fill: blue.darken(30%))[
      \$#data.cost_this_month
    ]
    #v(0.2em)
    #text(size: 9pt, fill: gray)[
      #if data.cost_change != 0 [
        #if data.cost_change > 0 [
          ↑ #str(data.cost_change)% vs mes anterior
        ] else [
          ↓ #str(data.cost_change)% vs mes anterior
        ]
      ]
    ]
  ],

  // Active Users Card
  box(
    fill: green.lighten(90%),
    inset: 1em,
    radius: 5pt,
    width: 100%,
  )[
    #text(size: 10pt, fill: gray)[*Usuarios Activos*]
    #v(0.3em)
    #text(size: 20pt, weight: "bold", fill: green.darken(30%))[
      #str(data.active_users_this_month)
    ]
    #v(0.2em)
    #text(size: 9pt, fill: gray)[
      #str(data.new_users_this_month) nuevos usuarios
    ]
  ],
)

#v(1em)

#grid(
  columns: (1fr, 1fr),
  gutter: 1em,

  // Conversations Card
  box(
    fill: orange.lighten(90%),
    inset: 1em,
    radius: 5pt,
    width: 100%,
  )[
    #text(size: 10pt, fill: gray)[*Conversaciones*]
    #v(0.3em)
    #text(size: 20pt, weight: "bold", fill: orange.darken(30%))[
      #str(data.conversations_this_month)
    ]
    #v(0.2em)
    #text(size: 9pt, fill: gray)[
      #data.avg_messages_per_conversation mensajes promedio
    ]
  ],

  // Tokens Card
  box(
    fill: purple.lighten(90%),
    inset: 1em,
    radius: 5pt,
    width: 100%,
  )[
    #text(size: 10pt, fill: gray)[*Tokens Utilizados*]
    #v(0.3em)
    #text(size: 20pt, weight: "bold", fill: purple.darken(30%))[
      #data.tokens_this_month
    ]
    #v(0.2em)
    #text(size: 9pt, fill: gray)[
      \$#data.cost_per_conversation por conversación
    ]
  ],
)

#pagebreak()

// ============================================
// COST ANALYSIS
// ============================================

= Análisis de Costos

== Desglose de Costos

El costo total del mes fue de *\$#data.cost_this_month*, distribuido de la siguiente manera:

#table(
  columns: (2fr, 1fr, 1fr),
  align: (left, right, right),
  inset: 10pt,
  stroke: 0.5pt + gray,
  table.header(
    [*Concepto*], [*Costo*], [*Porcentaje*]
  ),

  [Costo LLM (Generación)], [\$#data.llm_cost], [#data.llm_cost_percent%],
  [Costo Embeddings], [\$#data.embedding_cost], [#data.embedding_cost_percent%],
  table.cell(fill: gray.lighten(80%), [*Total*]),
  table.cell(fill: gray.lighten(80%), [*\$#data.cost_this_month*]),
  table.cell(fill: gray.lighten(80%), [*100%*]),
)

== Detalles de Tokens

#table(
  columns: (2fr, 1fr, 1fr),
  align: (left, right, right),
  inset: 10pt,
  stroke: 0.5pt + gray,
  table.header(
    [*Tipo*], [*Cantidad*], [*Costo Unitario*]
  ),

  [Tokens de Entrada (Prompt)], [#data.prompt_tokens], [\$0.50/M],
  [Tokens de Salida (Completion)], [#data.completion_tokens], [\$1.50/M],
  [Tokens de Embeddings], [#data.embedding_tokens], [\$0.13/M],
  table.cell(fill: gray.lighten(80%), [*Total Tokens*]),
  table.cell(fill: gray.lighten(80%), [*#data.total_tokens*]),
  table.cell(fill: gray.lighten(80%), []),
)

== Métricas de Eficiencia

- *Costo por Conversación:* \$#data.cost_per_conversation
- *Tokens Promedio por Conversación:* #data.avg_tokens_per_conversation
- *Costo por Usuario Activo:* \$#data.cost_per_active_user

#if data.cost_projection != "" [
== Proyección de Fin de Mes

Basado en el uso actual, el costo estimado para fin de mes es de *\$#data.cost_projection*.
]

#pagebreak()

// ============================================
// USER ACTIVITY
// ============================================

= Actividad de Usuarios

== Estadísticas Generales

#table(
  columns: (2fr, 1fr),
  align: (left, right),
  inset: 10pt,
  stroke: 0.5pt + gray,

  [*Total de Usuarios (Histórico)*], [#str(data.total_users)],
  [*Usuarios Activos este Mes*], [#str(data.active_users_this_month)],
  [*Usuarios Nuevos*], [#str(data.new_users_this_month)],
  [*Usuarios Recurrentes*], [#str(data.returning_users_this_month)],
  [*Tasa de Retención*], [#data.retention_rate%],
)

== Usuarios por Rol

#table(
  columns: (2fr, 1fr, 1fr),
  align: (left, right, right),
  inset: 10pt,
  stroke: 0.5pt + gray,
  table.header(
    [*Rol*], [*Cantidad*], [*Porcentaje*]
  ),

  [Estudiantes], [#str(data.students_count)], [#data.students_percent%],
  [Docentes], [#str(data.professors_count)], [#data.professors_percent%],
  [Externos], [#str(data.external_count)], [#data.external_percent%],
)

== Nivel de Participación

- *Mensajes Promedio por Usuario:* #data.avg_messages_per_user
- *Sesiones Promedio por Usuario:* #data.avg_sessions_per_user
- *Duración Promedio de Sesión:* #data.avg_session_duration minutos


#if data.top_users.len() > 0 [
== Usuarios Más Activos (Top 10)

#table(
  columns: (1fr, 2fr, 1fr),
  align: (center, left, right),
  inset: 8pt,
  stroke: 0.5pt + gray,
  table.header(
    [*\#*], [*Usuario*], [*Mensajes*]
  ),
  ..for (index, user) in data.top_users.enumerate() {
    ([#(index + 1)], [#user.name], [#str(user.message_count)])
  }
)
]

#pagebreak()

// ============================================
// CONVERSATIONS
// ============================================

= Análisis de Conversaciones

== Volumen de Conversaciones

#table(
  columns: (2fr, 1fr),
  align: (left, right),
  inset: 10pt,
  stroke: 0.5pt + gray,

  [*Total de Conversaciones (Histórico)*], [#str(data.total_conversations)],
  [*Conversaciones este Mes*], [#str(data.conversations_this_month)],
  [*Conversaciones Activas*], [#str(data.active_conversations)],
  [*Promedio de Mensajes por Conversación*], [#data.avg_messages_per_conversation],
)

== Intervención del Administrador

#box(
  fill: yellow.lighten(80%),
  inset: 1em,
  radius: 5pt,
  width: 100%,
)[
  *Tasa de Intervención:* #data.admin_intervention_rate%

  De #str(data.conversations_this_month) conversaciones este mes, #str(data.conversations_with_admin_help) requirieron asistencia del administrador.

  #let intervention = float(data.admin_intervention_rate)
  #if intervention > 15.0 [
    ⚠ *Alerta:* La tasa de intervención es alta. Considere revisar y mejorar el contenido de la base de conocimientos.
  ]
]

== Estado de Conversaciones

#table(
  columns: (2fr, 1fr),
  align: (left, right),
  inset: 10pt,
  stroke: 0.5pt + gray,

  [*Conversaciones Bloqueadas*], [#str(data.blocked_conversations)],
  [*Conversaciones Temporales*], [#str(data.temporary_conversations)],
)

#pagebreak()

// ============================================
// MESSAGES
// ============================================

= Análisis de Mensajes

== Volumen de Mensajes

#table(
  columns: (2fr, 1fr),
  align: (left, right),
  inset: 10pt,
  stroke: 0.5pt + gray,

  [*Total de Mensajes este Mes*], [#str(data.messages_this_month)],
  [*Mensajes de Usuarios*], [#str(data.user_messages_this_month)],
  [*Respuestas del Bot*], [#str(data.bot_messages_this_month)],
  [*Mensajes de Administradores*], [#str(data.admin_messages_this_month)],
  [*Promedio Diario*], [#data.avg_messages_per_day],
)

== Distribución Horaria

#if data.peak_hour != 0 [
*Hora Pico:* #str(data.peak_hour):00 hrs (#str(data.peak_hour_count) mensajes)

#let peak = data.peak_hour
Los usuarios son más activos entre las #str(data.peak_hour):00 y las #(peak + 2):00 horas.
]

#pagebreak()

// ============================================
// TOP QUERIES
// ============================================

= Consultas Más Frecuentes

Las siguientes son las preguntas más realizadas por los usuarios durante este mes:


#if data.top_queries.len() > 0 [
#table(
  columns: (1fr, 3fr, 1fr, 1fr),
  align: (center, left, center, center),
  inset: 8pt,
  stroke: 0.5pt + gray,
  table.header(
    [*\#*], [*Consulta*], [*Veces*], [*Calidad*]
  ),
  ..for (index, query) in data.top_queries.enumerate() {
    let quality = if query.has_good_answer { "✓" } else { "⚠" }
    ([#(index + 1)], [#query.query_text], [#str(query.query_count)], [#quality])
  }
)

#v(1em)

*Leyenda:*
- ✓ = Respuesta de buena calidad (similaridad > 0.5)
- ⚠ = Necesita mejorar contenido (similaridad < 0.5)
]


#if data.queries_needing_attention.len() > 0 [
== Consultas que Requieren Atención

Las siguientes consultas tienen baja calidad de respuesta y deberían agregarse más contenido a la base de conocimientos:

#for query in data.queries_needing_attention [
- *"#query.query_text"* (#str(query.query_count) veces, similaridad: #query.avg_similarity)
]
]

#pagebreak()

// ============================================
// KNOWLEDGE BASE
// ============================================

= Base de Conocimientos

== Uso de la Base de Conocimientos


#if data.top_chunks.len() > 0 [
=== Fragmentos Más Utilizados

#table(
  columns: (1fr, 2fr, 1fr, 1fr),
  align: (center, left, center, center),
  inset: 8pt,
  stroke: 0.5pt + gray,
  table.header(
    [*\#*], [*Documento*], [*Usos*], [*Similaridad*]
  ),
  ..for (index, chunk) in data.top_chunks.enumerate() {
    ([#(index + 1)], [#chunk.document_title], [#str(chunk.usage_count)], [#chunk.avg_similarity])
  }
)
]

== Estadísticas Generales

#table(
  columns: (2fr, 1fr),
  align: (left, right),
  inset: 10pt,
  stroke: 0.5pt + gray,

  [*Total de Documentos*], [#str(data.total_documents)],
  [*Total de Fragmentos (Chunks)*], [#str(data.total_chunks)],
  [*Fragmentos Utilizados*], [#str(data.chunks_used)],
  [*Fragmentos No Utilizados*], [#str(data.chunks_never_used)],
  [*Tasa de Cobertura*], [#data.coverage_rate%],
)

#let coverage = float(data.coverage_rate)
#if coverage < 50.0 [
#box(
  fill: yellow.lighten(80%),
  inset: 1em,
  radius: 5pt,
  width: 100%,
)[
  ⚠ *Recomendación:* La tasa de cobertura es baja (#data.coverage_rate%). Esto significa que muchos fragmentos de la base de conocimientos no están siendo utilizados. Considere:

  - Revisar y eliminar contenido irrelevante
  - Mejorar la calidad de los embeddings
  - Agregar contenido que responda a las consultas más frecuentes
]
]

#pagebreak()

// ============================================
// SYSTEM PERFORMANCE
// ============================================

= Rendimiento del Sistema

== Tiempos de Respuesta

#table(
  columns: (2fr, 1fr),
  align: (left, right),
  inset: 10pt,
  stroke: 0.5pt + gray,

  [*Tiempo Promedio de Respuesta LLM*], [#data.avg_llm_response_time ms],
  [*Tiempo de Respuesta P95*], [#data.p95_response_time ms],
  [*Tiempo de Respuesta P99*], [#data.p99_response_time ms],
)

== Salud del Sistema

#table(
  columns: (2fr, 1fr),
  align: (left, right),
  inset: 10pt,
  stroke: 0.5pt + gray,

  [*Errores en Últimas 24h*], [#str(data.errors_last_24h)],
  [*Conversaciones Fallidas*], [#str(data.failed_conversations)],
  [*Disponibilidad*], [#data.uptime%],
)

#let errors = data.errors_last_24h
#if errors > 10 [
#box(
  fill: red.lighten(80%),
  inset: 1em,
  radius: 5pt,
  width: 100%,
)[
  ⚠ *Alerta:* Se detectaron #str(data.errors_last_24h) errores en las últimas 24 horas. Revise los logs del sistema.
]
]

#pagebreak()

// ============================================
// RECOMMENDATIONS
// ============================================

= Recomendaciones

Basado en los datos recopilados este mes, se sugieren las siguientes acciones:

== Optimización de Costos

#let cost_conv = float(data.cost_per_conversation)
#if cost_conv > 0.05 [
- El costo por conversación (\$#data.cost_per_conversation) es alto. Considere:
  - Reducir el tamaño del contexto del prompt
  - Limitar el número de chunks recuperados
  - Ajustar los parámetros del modelo (temperatura, max_tokens)
]

#let avg_tokens = float(data.avg_tokens_per_conversation)
#if avg_tokens > 5000 [
- El promedio de tokens por conversación (#data.avg_tokens_per_conversation) es elevado
- Revise si el sistema está enviando información innecesaria al LLM
]

== Mejora de Contenidos

#let intervention = float(data.admin_intervention_rate)
#if intervention > 15.0 [
- La tasa de intervención del administrador (#data.admin_intervention_rate%) indica que el bot no puede responder muchas consultas
- Priorice agregar contenido para las consultas más frecuentes sin buena respuesta
]


#if data.queries_needing_attention.len() > 0 [
- Agregue o mejore contenido para las siguientes consultas:
#for query in data.queries_needing_attention [
  - "#query.query_text"
]
]

== Base de Conocimientos

#let coverage_rec = float(data.coverage_rate)
#if coverage_rec < 40.0 [
- La cobertura de la base de conocimientos es baja (#data.coverage_rate%)
- Elimine fragmentos no utilizados y agregue contenido relevante
]

#let chunks_unused = data.chunks_never_used
#if chunks_unused > 0 [
- Hay #str(data.chunks_never_used) fragmentos que nunca han sido utilizados
- Revise y actualice o elimine este contenido
]

== Experiencia de Usuario

#let avg_msgs = float(data.avg_messages_per_conversation)
#if avg_msgs > 10.0 [
- El promedio de #data.avg_messages_per_conversation mensajes por conversación es alto
- Considere mejorar las respuestas para que sean más completas y directas
]

#pagebreak()

// ============================================
// CONCLUSIONS
// ============================================

= Conclusiones

== Resumen del Período

Durante #data.month_year, el chatbot del Instituto Superior Tecnológico Sudamericano atendió a *#str(data.active_users_this_month) usuarios activos* a través de *#str(data.conversations_this_month) conversaciones*, con un costo total de *\$#data.cost_this_month*.

== Indicadores Clave

#let satisfaction = 100.0 - float(data.admin_intervention_rate)
- ✓ *Usuarios satisfechos:* #calc.round(satisfaction, digits: 1)% de conversaciones resueltas sin intervención
- ✓ *Disponibilidad:* #data.uptime%
- ✓ *Eficiencia:* \$#data.cost_per_conversation por conversación

#let intervention_conc = float(data.admin_intervention_rate)
#if intervention_conc > 20.0 [
- ⚠ *Área de mejora:* Tasa de intervención del administrador alta
]

#let coverage_conc = float(data.coverage_rate)
#if coverage_conc < 50.0 [
- ⚠ *Área de mejora:* Cobertura de la base de conocimientos baja
]

== Próximos Pasos

+ Implementar las recomendaciones sugeridas en este reporte
+ Monitorear continuamente las métricas de costo y calidad
+ Actualizar la base de conocimientos mensualmente
+ Revisar consultas con baja calidad de respuesta

#v(2cm)

#align(center)[
  #line(length: 60%, stroke: 0.5pt + gray)
  #v(0.5cm)
  #text(size: 10pt, fill: gray)[
    Fin del Reporte

    Este reporte fue generado automáticamente por el sistema de analíticas del chatbot.

    Para más información, contacte al administrador del sistema.
  ]
]
