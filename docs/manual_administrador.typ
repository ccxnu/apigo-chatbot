#set document(
  title: "Manual de Administrador - Chatbot ISTS",
  author: "Pablo Moisés Cuenca Cuenca",
  date: datetime.today(),
)

#set page(
  paper: "a5",
  margin: (top: 1.5cm, right: 1.5cm, bottom: 1.5cm, left: 1.5cm),
  numbering: "1",
)

#set text(
  font: "Times New Roman",
  size: 9pt,
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

  #upper[*Manual de Administrador*]

  #v(3em)

  *Autor:*

  Pablo Moisés Cuenca Cuenca

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

Bienvenido al *Manual del Administrador del Sistema de Chatbot ISTS*. Este documento está diseñado para guiarle en el uso del panel de administración web del chatbot institucional.

El panel administrativo le permite gestionar completamente el chatbot: configurar parámetros, cargar documentos para el conocimiento del bot, conectar WhatsApp y monitorear estadísticas del sistema.

== Objetivos del Manual

Este manual tiene los siguientes objetivos:

- Explicar cómo acceder y navegar el panel de administración
- Guiar en la configuración inicial del sistema
- Enseñar a gestionar el conocimiento del chatbot (documentos RAG)
- Mostrar cómo conectar y administrar WhatsApp
- Explicar la gestión de conversaciones con usuarios
- Describir la administración de usuarios y parámetros del sistema

== A Quién Está Dirigido

Este manual está dirigido a:

- *Administradores del sistema* responsables de la configuración
- *Personal de soporte* que responde consultas de usuarios
- *Coordinadores académicos* que gestionan información institucional
- *Personal de TI* encargado del mantenimiento

== Requisitos Previos

Para utilizar el panel de administración, necesita:

- Una computadora con navegador web moderno (Chrome, Firefox, Edge)
- Conexión a Internet estable
- Credenciales de acceso (usuario y contraseña)
- Permisos de administrador asignados por el área de TI

== Acceso al Sistema

El panel de administración está disponible en:

*URL de acceso:* `http://localhost:6600` (servidor local)

*Navegadores soportados:*
- Google Chrome 90+
- Mozilla Firefox 88+

#pagebreak()

= Inicio de Sesión

== Acceder al Panel

=== Paso 1: Abrir la Aplicación

+ Abra su navegador web
+ Ingrese la URL del panel de administración
+ Espere a que cargue la página de inicio de sesión

#figure(
  image("assets/adm_iniciar_sesion.png", width: 80%),
  caption: [Página de inicio de sesión]
)

=== Paso 2: Ingresar Credenciales

+ En el campo *"Usuario"*, ingrese su nombre de usuario
+ En el campo *"Contraseña"*, ingrese su contraseña
+ Haga clic en el botón *"Iniciar Sesión"*

*Nota:* Las credenciales son proporcionadas por el administrador del sistema. Si no tiene credenciales, contacte al área de TI.

// TODO: Add screenshot of filled login form
#figure(
  image("assets/adm_ingreso_credenciales.png", width: 50%),
  caption: [Formulario de inicio de sesión completado]
)

=== Paso 3: Acceso al Dashboard

Una vez autenticado exitosamente, será redirigido automáticamente al *Panel de Control* (Dashboard).

// TODO: Add screenshot of dashboard after login
#figure(
  image("assets/admin_panel.png", width: 100%),
  caption: [Panel de control principal]
)

== Recuperar Contraseña

Si olvidó su contraseña:

+ En la página de inicio de sesión, haga clic en *"¿Olvidaste tu contraseña?"*
+ Ingrese su nombre de usuario o correo electrónico
+ Recibirá instrucciones por correo para restablecer su contraseña
+ Siga el enlace en el correo y establezca una nueva contraseña

// TODO: Add screenshot of password recovery page
#figure(
  image("assets/adm_forgot_password.png", width: 50%),
  caption: [Página de recuperación de contraseña]
)

== Cerrar Sesión

Para cerrar su sesión de forma segura:

+ Haga clic en su avatar o nombre en la esquina superior derecha
+ En el menú desplegable, seleccione *"Cerrar Sesión"*
+ Será redirigido a la página de inicio de sesión

*Importante:* Siempre cierre sesión al terminar, especialmente en computadoras compartidas.

// TODO: Add screenshot of user menu with logout option
#figure(
  image("assets/adm_cerrar_sesion.png", width: 50%),
  caption: [Menú de usuario con opción de cerrar sesión]
)

#pagebreak()

= Navegación del Panel

== Estructura del Panel

El panel de administración está compuesto por:

=== Barra Superior (Header)

- Logo institucional (esquina izquierda)
- Navegación de migas de pan (breadcrumbs)
- Botón de tema (cambiar entre modo claro/oscuro)
- Menú de usuario con avatar (esquina derecha)

=== Menú Lateral (Sidebar)

El menú lateral contiene las siguientes secciones principales:

*General:*
- Panel de Control (Dashboard)
- Chats
- Estadísticas

*Gestión de Conocimiento:*
- RAG (Documentos y Chunks)

*WhatsApp:*
- Conexión WhatsApp

*Sistema:*
- Usuarios
- Parámetros

// TODO: Add screenshot of sidebar menu
#figure(
  image("assets/adm_sidebar.png", width: 30%),
  caption: [Menú de navegación lateral]
)

=== Área de Contenido

El área central muestra el contenido de la sección seleccionada. Cambia dinámicamente según la opción elegida en el menú lateral.

== Cambiar Tema Visual

El panel soporta modo claro y oscuro:

+ Haga clic en el icono de sol/luna en la barra superior
+ El tema cambiará inmediatamente
+ Su preferencia se guarda automáticamente

// TODO: Add screenshot showing both light and dark themes
#figure(
  image("assets/adm_theme_toogle.png", width: 35%),
  caption: [Comparación de tema claro y oscuro]
)

#figure(
  image("assets/adm_dark_theme.png", width: 95%),
  caption: [Tema oscuro]
)

== Colapsar el Menú Lateral

Para tener más espacio de trabajo:

+ Haga clic en el icono de hamburguesa (tres líneas) en la parte superior del sidebar
+ El menú se colapsará mostrando solo iconos
+ Haga clic nuevamente para expandirlo

#pagebreak()

= Panel de Control (Dashboard)

== Vista General

El Dashboard proporciona una visión general del estado del sistema con métricas clave y accesos rápidos.

// TODO: Add screenshot of complete dashboard
#figure(
  image("assets/admin_panel.png", width: 100%),
  caption: [Vista completa del panel de control]
)

== Secciones del Dashboard

=== Métricas de Servicio

Tarjetas superiores que muestran:

*Costo del Servicio:*
- Costo total acumulado de uso de APIs (LLM, embeddings)
- Indicador de incremento o decremento porcentual

*Velocidad Promedio:*
- Tiempo promedio de respuesta del chatbot
- Medido en milisegundos

*Tokens Consumidos:*
- Total de tokens procesados por el LLM
- Indicador de uso mensual

*Usuarios Activos:*
- Cantidad de usuarios únicos que han interactuado
- Período: últimos 30 días

// TODO: Add screenshot of metrics cards
// #figure(
//   image("images/admin_metrics_cards.png", width: 100%),
//   caption: [Tarjetas de métricas del servicio]
// )

=== Gráficos de Actividad

*Gráfico de Resumen:*
- Visualización de mensajes por día
- Comparación de mensajes entrantes vs. salientes
- Rango: últimos 7 días

*Gráfico de Ventas/Conversiones:* (si aplica)
- Seguimiento de conversiones o matrículas generadas
- Gráfico de barras por mes

// TODO: Add screenshot of activity charts
// #figure(
//   image("images/admin_dashboard_charts.png", width: 100%),
//   caption: [Gráficos de actividad del sistema]
// )

#pagebreak()

= Gestión de Documentos RAG

== ¿Qué es RAG?

RAG (Retrieval Augmented Generation) es el sistema de conocimiento del chatbot. Permite que el bot responda preguntas basándose en información que usted ingresa manualmente.

*Flujo:*
+ Usted crea documentos con título y categoría
+ Agrega fragmentos de texto (chunks) al documento
+ Los chunks se indexan con inteligencia artificial
+ Cuando un usuario pregunta, el bot busca chunks relevantes
+ El bot genera una respuesta usando esos chunks

== Acceder a Gestión RAG

+ En el menú lateral, haga clic en *"RAG"*
+ Se abrirá la vista de gestión de documentos

// TODO: Add screenshot of RAG main page
#figure(
  image("assets/adm_document.png", width: 100%),
  caption: [Página principal de gestión RAG]
)

== Pestañas de RAG

El módulo RAG tiene tres pestañas:

+ *Documentos*: Crear y gestionar documentos base
+ *Chunks*: Agregar y editar fragmentos de texto
+ *Estadísticas*: Métricas de uso de chunks

=== Pestaña Documentos

==== Crear un Nuevo Documento

+ Haga clic en el botón *"+ Nuevo Documento"*
+ Se abrirá un formulario de creación

// TODO: Add screenshot of create document dialog
#figure(
  image("assets/adm_documentos.png", width: 50%),
  caption: [Formulario de creación de documento]
)

+ Complete los siguientes campos:

  *Título:* Nombre descriptivo del documento
  - Ejemplo: "Proceso de Matrícula 2025"

  *Categoría:* Seleccione la categoría apropiada
  - Información Académica
  - Procesos Administrativos
  - Requisitos de Admisión
  - Reglamentos
  - Otros

  *Resumen:* Breve descripción del contenido (opcional)
  - Ejemplo: "Pasos y requisitos para la matrícula ordinaria"

  *Fuente:* Origen de la información (opcional)
  - Ejemplo: "Secretaría Académica"

+ Haga clic en *"Crear"*
+ El documento aparecerá en la lista

*Nota:* Este documento es solo un contenedor. Debe agregar chunks (fragmentos de texto) en la pestaña de Chunks.

==== Gestionar Documentos

*Ver detalles de un documento:*
+ Haga clic en el documento en la lista
+ Se mostrará información detallada:
  - Título y categoría
  - Resumen y fuente
  - Fecha de creación
  - Número de chunks asociados

// TODO: Add screenshot of document details panel
// #figure(
//   image("images/admin_rag_document_details.png", width: 80%),
//   caption: [Panel de detalles de documento]
// )

*Editar metadatos:*
+ Haga clic en el icono de lápiz (editar)
+ Modifique el título o categoría
+ Haga clic en *"Guardar"*

*Eliminar documento:*
+ Haga clic en el icono de papelera (eliminar)
+ Confirme la eliminación en el diálogo
+ *Advertencia:* Esto también eliminará todos los chunks asociados

*Nota importante:* Los documentos son solo contenedores organizacionales. El contenido real está en los chunks que debe agregar manualmente.

// TODO: Add screenshot of delete confirmation
#figure(
  image("assets/adm_eliminar_doc.png", width: 45%),
  caption: [Confirmación de eliminación de documento]
)

*Buscar documentos:*
+ Use la barra de búsqueda en la parte superior
+ Puede buscar por título, categoría o contenido
+ Los resultados se filtran en tiempo real

=== Pestaña Chunks

Esta pestaña muestra todos los fragmentos de texto y permite agregar nuevos.

==== Crear un Nuevo Chunk

+ Haga clic en el botón *"+ Nuevo Chunk"*
+ Complete el formulario:

  *Documento:* Seleccione el documento al que pertenece

  *Contenido:* Ingrese el texto del fragmento
  - Tamaño recomendado: 500-1500 caracteres
  - Sea específico y completo en cada chunk
  - Un chunk debe contener información coherente

+ Haga clic en *"Crear"*
+ El sistema automáticamente:
  - Genera el embedding vectorial
  - Indexa el chunk para búsqueda semántica
  - Crea registro de estadísticas

// TODO: Add screenshot of create chunk dialog
#figure(
  image("assets/adm_create_chunk.png", width: 80%),
  caption: [Formulario de creación de chunk/fragmento]
)

*Ejemplo de buen chunk:*
```
Para matricularse debe seguir estos pasos:
1. Ingresar al sistema con su número de cédula
2. Seleccionar la carrera y horario
3. Subir documentos requeridos (cédula, certificado de votación)
4. Realizar el pago de matrícula ($150)
5. Confirmar su matrícula en secretaría

Fechas: Del 15 al 30 de enero de 2025
Contacto: matriculas@ists.edu.ec
```

==== Listar Chunks

La tabla muestra todos los chunks existentes.

// TODO: Add screenshot of chunks tab
#figure(
  image("assets/adm_chunks.png", width: 100%),
  caption: [Vista de chunks indexados]
)

==== Columnas de la Tabla de Chunks

- *ID*: Identificador único del chunk
- *Documento*: Documento fuente del que proviene
- *Contenido*: Preview del texto del chunk
// - *Tokens*: Cantidad de tokens en el chunk
- *Índice*: Posición del chunk en el documento original
- *Acciones*: Ver, editar, eliminar

==== Ver Contenido Completo

+ Haga clic en el icono de ojo (ver) en un chunk
+ Se abrirá un panel lateral mostrando:
  - Contenido completo del chunk
  - Metadatos asociados
  - Documento fuente

// TODO: Add screenshot of chunk content viewer
// #figure(
//   image("images/admin_rag_chunk_viewer.png", width: 90%),
//   caption: [Visor de contenido de chunk]
// )

==== Editar un Chunk

Puede editar el contenido de un chunk:

+ Haga clic en el icono de lápiz (editar)
+ Modifique el texto en el editor
+ Haga clic en *"Guardar"*
+ El embedding se regenerará automáticamente

#figure(
  image("assets/adm_editar_chunk.png", width: 90%),
  caption: [Editar un fragmento]
)

*Importante:* Al editar, el sistema recalcula el embedding vectorial, por lo que puede tomar unos segundos.

==== Filtrar Chunks

*Por documento:*
+ Use el selector *"Filtrar por documento"*
+ Seleccione el documento deseado
+ Solo se mostrarán chunks de ese documento

*Por búsqueda:*
+ Use la barra de búsqueda
+ Ingrese palabras clave
+ Se mostrarán chunks que contengan esas palabras

=== Pestaña Estadísticas

Muestra métricas de uso de cada chunk:
// #figure(
//   image("images/admin_rag_create_chunk.png", width: 80%),
//   caption: [Formulario de creación de chunk]
// )

// TODO: Add screenshot of statistics tab
#figure(
  image("assets/adm_estadisticas.png", width: 100%),
  caption: [Estadísticas de uso de chunks]
)

==== Métricas Disponibles

*Por cada chunk:*
- *Consultas*: Veces que el chunk fue usado en respuestas
- *Última consulta*: Fecha de último uso
- *Similitud promedio*: Qué tan relevante fue en promedio
- *Tasa de uso*: Porcentaje de uso respecto al total

==== Chunks Más Utilizados

+ Haga clic en *"Ordenar por uso"*
+ Los chunks más consultados aparecerán primero
+ Esto indica qué información es más solicitada por usuarios

*Utilidad:*
- Identificar temas populares
- Priorizar actualización de información frecuente
- Detectar gaps de información (chunks poco usados)


== Mejores Prácticas RAG

=== Calidad del Contenido

*Al crear chunks:*
+ Use información oficial y actualizada
+ Revise ortografía y gramática
+ Sea específico y completo
+ Incluya contexto suficiente en cada chunk
+ Un chunk debe ser autocontenido (comprensible por sí solo)

=== Organización

*Categorización:*
+ Asigne categorías consistentes a documentos
+ Use nombres descriptivos para documentos
+ Agrupe chunks relacionados en el mismo documento

*Tamaño de chunks:*
+ Muy cortos (\<300 chars): Poco contexto para el LLM
+ Recomendado (500-1500 chars): Balance entre precisión y contexto
+ Muy largos (\>2000 chars): Información irrelevante puede incluirse

*Actualización:*
+ Revise y actualice chunks trimestralmente
+ Elimine información desactualizada
+ Edite chunks obsoletos con información nueva

=== Monitoreo

*Revisión regular:*
+ Revise estadísticas mensualmente
+ Identifique chunks con baja tasa de uso
+ Mejore o elimine contenido irrelevante

#pagebreak()

= Conexión de WhatsApp

== Vista General WhatsApp

El módulo WhatsApp permite conectar una cuenta de WhatsApp Business al chatbot para que pueda enviar y recibir mensajes.

// TODO: Add screenshot of WhatsApp main page
#figure(
  image("assets/adm_whatsapp.png", width: 100%),
  caption: [Página principal de WhatsApp]
)

== Conectar WhatsApp por Primera Vez

=== Requisitos Previos

- Una cuenta de WhatsApp Business en un teléfono móvil
- El teléfono con conexión a Internet
- Permisos de administrador en el panel

=== Pasos para Conectar

==== Paso 1: Iniciar Conexión

+ En el menú lateral, haga clic en *"WhatsApp"*
+ Haga clic en el botón *"Iniciar Conexión"* o *"Conectar WhatsApp"*

// TODO: Add screenshot of initial WhatsApp connection page
#figure(
  image("assets/adm_session.png", width: 80%),
  caption: [Página inicial de conexión WhatsApp]
)

==== Paso 2: Escanear Código QR

+ El sistema generará un código QR
+ El código se actualiza automáticamente cada 30 segundos si no es escaneado

// TODO: Add screenshot of QR code displayed
#figure(
  image("assets/adm_qrcode.png", width: 70%),
  caption: [Código QR para vincular WhatsApp]
)

+ En su teléfono móvil:
  - Abra WhatsApp Business
  - Toque el menú (tres puntos verticales)
  - Seleccione *"Dispositivos vinculados"*
  - Toque *"Vincular un dispositivo"*
  - Escanee el código QR mostrado en el panel

// TODO: Add screenshot of phone scanning QR (optional, can be drawn)
// #figure(
//   image("assets/adm_link_device.jpg", width: 30%),
//   caption: [Escanear QR desde WhatsApp móvil]
// )

==== Paso 3: Confirmación de Conexión

+ Una vez escaneado, el estado cambiará a *"Conectando..."*
+ En pocos segundos, el estado mostrará *"Conectado"*
+ Verá los detalles de la cuenta vinculada:
  - Número de teléfono
  - Nombre de la cuenta
  - Fecha de conexión
  - Estado: Activo

// TODO: Add screenshot of successful connection
#figure(
  image("assets/adm_link_device.jpg", width: 30%),
  caption: [Escanear QR desde WhatsApp móvil]
)

*Importante:* Mantenga el teléfono con WhatsApp conectado a Internet para que el chatbot funcione.

== Estado de Conexión

=== Indicadores de Estado

El panel muestra el estado actual de la conexión:

*Estados posibles:*

- *Desconectado* (rojo): Sin conexión activa
- *Conectando* (amarillo): En proceso de vinculación
- *Conectado* (verde): Funcionando correctamente
- *Error* (rojo): Problema de conexión

// TODO: Add screenshot showing different connection states
#figure(
  image("assets/adm_estado_conexion.png", width: 45%),
  caption: [Diferentes estados de conexión]
)

=== Monitorear Conexión

*Panel de estado muestra:*
- Tiempo de conexión activa
- Último mensaje recibido/enviado
- Cantidad de mensajes del día
- Latencia de conexión

== Desconectar WhatsApp

Si necesita desconectar la sesión:

+ Haga clic en el botón *"Desconectar"*
+ Aparecerá un diálogo de confirmación
+ Confirme la desconexión

// TODO: Add screenshot of disconnect confirmation
#figure(
  image("assets/adm_estado_desconectado.png", width: 45%),
  caption: [Confirmación de desconexión]
)

*Nota:* Esto cerrará la sesión en el servidor, pero el dispositivo seguirá vinculado en WhatsApp. Para desvincular completamente, debe hacerlo desde el teléfono.

== Reconectar WhatsApp

Si la conexión se pierde (por ejemplo, teléfono sin Internet):

+ El estado mostrará *"Desconectado"*
+ Haga clic en *"Reconectar"*
+ Si el teléfono sigue vinculado, se reconectará automáticamente
+ Si no, deberá escanear un nuevo código QR

#pagebreak()

= Gestión de Conversaciones (Chats)

== Vista General de Chats

El módulo de Chats permite ver y gestionar todas las conversaciones entre usuarios y el chatbot, con una interfaz similar a WhatsApp.

// TODO: Add screenshot of chats main interface
#figure(
  image("assets/adm_chatbot.png", width: 50%),
  caption: [Interfaz principal de gestión de chats]
)

== Estructura de la Interfaz

La interfaz de chats está dividida en dos paneles:

*Panel Izquierdo:* Lista de conversaciones
*Panel Derecho:* Mensajes de la conversación seleccionada

=== Panel de Lista de Conversaciones

==== Elementos de Cada Conversación

Cada conversación muestra:
- *Avatar*: Foto o inicial del usuario
- *Nombre*: Nombre del contacto o número de teléfono
- *Último mensaje*: Preview del mensaje más reciente
- *Hora*: Timestamp del último mensaje
- *Indicadores*:
  - Badge numérico: Cantidad de mensajes no leídos
  - Check azul: Última respuesta del administrador
  - Icono de pin: Conversación fijada

// TODO: Add screenshot of conversation list item
// #figure(
//   image("images/admin_chats_conversation_item.png", width: 70%),
//   caption: [Elemento de conversación en la lista]
// )

==== Buscar Conversaciones

+ Use la barra de búsqueda en la parte superior del panel izquierdo
+ Puede buscar por:
  - Nombre del usuario
  - Número de teléfono
  - Contenido de mensajes
+ Los resultados se filtran en tiempo real

// TODO: Add screenshot of search in conversations
#figure(
  image("assets/adm_search.png", width: 60%),
  caption: [Búsqueda de conversaciones]
)

== Ver Conversación

=== Abrir una Conversación

+ Haga clic en cualquier conversación de la lista
+ El panel derecho mostrará el historial completo de mensajes

// TODO: Add screenshot of opened conversation
#figure(
  image("assets/adm_chat_mesage.png", width: 40%),
  caption: [Conversación abierta con historial de mensajes]
)

=== Elementos del Panel de Conversación

==== Cabecera de Conversación

En la parte superior se muestra:
- *Avatar y nombre* del usuario
- *Número de teléfono*
- *Estado*: En línea / Última vez activo
- *Botones de acción*:
  - Ver información del usuario (icono "i")
  - Bloquear usuario
  - Eliminar conversación
  - Marcar como no leída

// TODO: Add screenshot of conversation header
// #figure(
//   image("images/admin_chats_header.png", width: 100%),
//   caption: [Cabecera de conversación con información del usuario]
// )

=== Enviar un Mensaje

+ En el área de composición (parte inferior), escriba su mensaje
+ Haga clic en el botón de enviar (icono de avión de papel) o presione Enter

// TODO: Add screenshot of message composer
#figure(
  image("assets/adm_send_mesage.png", width: 50%),
  caption: [Área de composición de mensajes]
)

*Características del compositor:*
- Soporte para texto con formato (negritas, cursivas)
- Emojis (haga clic en el icono de emoji)
- Adjuntar archivos (imágenes, documentos)
- Mensajes de múltiples líneas (Shift + Enter para nueva línea)

=== Bloquear Usuario

Si un usuario envía spam o contenido inapropiado:

+ Abra la conversación del usuario
+ Haga clic en el menú de opciones (tres puntos)
+ Seleccione *"Bloquear usuario"*
+ Confirme la acción en el diálogo

// TODO: Add screenshot of block confirmation
// #figure(
//   image("images/admin_chats_block.png", width: 60%),
//   caption: [Confirmación de bloqueo de usuario]
// )

*Efecto:*
- El usuario no recibirá más respuestas del bot
- Sus mensajes seguirán llegando pero marcados como "bloqueado"
- La conversación aparece en el filtro "Bloqueadas"

*Desbloquear:*
+ Vaya a *"Filtros"* → *"Bloqueadas"*
+ Abra la conversación del usuario bloqueado
+ Seleccione *"Desbloquear usuario"*

=== Eliminar Conversación

+ Haga clic derecho en la conversación
+ Seleccione *"Eliminar conversación"*
+ Confirme la eliminación

*Advertencia:* Esta acción es irreversible. Todo el historial de mensajes se eliminará permanentemente.

// TODO: Add screenshot of delete confirmation
// #figure(
//   image("images/admin_chats_delete.png", width: 60%),
//   caption: [Confirmación de eliminación de conversación]
// )

== Panel de Información del Usuario

Al hacer clic en el icono "i" en la cabecera, se abre un panel lateral con información del usuario:

// TODO: Add screenshot of user info panel
// #figure(
//   image("images/admin_chats_user_info.png", width: 80%),
//   caption: [Panel de información del usuario]
// )

*Información mostrada:*
- Nombre completo
- Número de teléfono
- Correo electrónico (si está registrado)
- Rol (estudiante, aspirante, padre, etc.)
- Estado de verificación
- Fecha de primer contacto
- Total de mensajes intercambiados
- Última actividad

*Acciones disponibles:*
- Ver historial completo de interacciones

#pagebreak()

= Administración de Parámetros

== ¿Qué son los Parámetros?

Los parámetros son configuraciones del sistema que controlan el comportamiento del chatbot. Permiten modificar ajustes sin necesidad de cambiar código.

*Ejemplos de parámetros:*
- Horarios de atención
- Mensajes de bienvenida
- Configuración de LLM (temperatura, max tokens)
- URLs de APIs externas
- Límites de uso

== Acceder a Parámetros

+ En el menú lateral, vaya a *"Sistema"* → *"Parámetros"*
+ Se mostrará la lista de todos los parámetros del sistema

// TODO: Add screenshot of parameters list
#figure(
  image("assets/adm_parameters.png", width: 100%),
  caption: [Lista de parámetros del sistema]
)

== Estructura de Parámetros

Cada parámetro tiene:
- *Nombre*: Descripción legible del grupo (ej: "Temperatura del LLM")
- *Código*: Identificador único (ej: `LLM_TEMPERATURE`)
- *Datos*: Valor del parámetro en formato JSON
- *Descripción*: Explicación de para qué sirve
- *Estado*: Activo o Inactivo

== Buscar Parámetros

+ Use la barra de búsqueda en la parte superior
+ Puede buscar por:
  - Código del parámetro
  - Nombre
  - Descripción
+ Use filtros para mostrar solo activos/inactivos
+ Use el filtro por nombre

== Ver Detalles de un Parámetro

+ Haga clic en cualquier parámetro de la lista
+ Se abrirá un panel con información completa

// TODO: Add screenshot of parameter details

*Información mostrada:*
- Todos los campos del parámetro
- Valor actual (JSON formateado)
- Historial de cambios (últimas 10 modificaciones)
- Validación del formato JSON

== Editar un Parámetro

+ Seleccione el parámetro que desea editar
+ Haga clic en el botón *"Editar"* (icono de lápiz)
+ Se abrirá un diálogo de edición

#figure(
  image("assets/adm_param_edit.png", width: 50%),
  caption: [Vista detallada de un parámetro]
)

*Campos editables:*

- *Nombre*: Puede modificarse
- *Descripción*: Actualice si es necesario
- *Datos (JSON)*: Editor con validación de sintaxis
- *Estado*: Activo/Inactivo

*Editor JSON:*
- Resaltado de sintaxis
- Validación en tiempo real
- Formateo automático (botón "Formatear")
- Detecta errores de sintaxis

*Ejemplo de datos:*
```json
{
  "temperature": 0.7,
  "maxTokens": 2000,
  "model": "llama-3.3-70b"
}
```

*Pasos:*
+ Modifique los valores necesarios
+ Asegúrese de que el JSON sea válido (sin errores en rojo)
+ Haga clic en *"Guardar"*
+ Confirme los cambios

*Advertencia:* Algunos parámetros afectan el comportamiento crítico del sistema. Edite con precaución.

== Crear un Nuevo Parámetro

+ Haga clic en el botón *"+ Nuevo Parámetro"*
+ Complete el formulario:

// TODO: Add screenshot of create parameter dialog
#figure(
  image("assets/adm_create_parametro.png", width: 50%),
  caption: [Formulario de creación de parámetro]
)

*Campos requeridos:*

- *Nombre*:
  - Descripción breve y clara
  - Ejemplo: "Horario de Atención"

- *Código*:
  - Identificador único en mayúsculas
  - Solo letras, números y guiones bajos
  - Ejemplo: `HORARIO_ATENCION`

- *Datos (JSON)*:
  - Valor inicial del parámetro
  - Debe ser JSON válido

- *Descripción*:
  - Explique para qué sirve el parámetro
  - Documente el formato esperado de los datos

*Ejemplo completo:*

Código: `HORARIO_ATENCION`

Nombre: `Horario de Atención`

Datos:
```json
{
  "lunes_viernes": {
    "inicio": "08:00",
    "fin": "18:00"
  },
  "sabado": {
    "inicio": "08:00",
    "fin": "13:00"
  },
  "domingo": {
    "activo": false
  }
}
```

Descripción: `Define los horarios en que el chatbot responde automáticamente`

+ Haga clic en *"Crear"*
+ El parámetro aparecerá en la lista

== Desactivar un Parámetro

Si necesita desactivar temporalmente un parámetro sin eliminarlo:

+ Abra el parámetro
+ Haga clic en *"Editar"*
+ Cambie el estado a *"Inactivo"*
+ Guarde los cambios

*Efecto:*
- El parámetro sigue en la base de datos
- El sistema usa valores por defecto
- Aparece marcado como inactivo en la lista

== Eliminar un Parámetro

+ Seleccione el parámetro
+ Haga clic en el botón *"Eliminar"* (icono de papelera)
+ Confirme la eliminación

// TODO: Add screenshot of delete parameter confirmation
#figure(
  image("assets/adm_param_eliminate.png", width: 60%),
  caption: [Confirmación de eliminación de parámetro]
)

*Advertencia:*
- Solo elimine parámetros personalizados
- La eliminación es permanente

== Parámetros Importantes del Sistema

=== LLM_CONFIG

Configuración del modelo de lenguaje:

```json
{
  "provider": "groq",
  "model": "llama-3.3-70b-versatile",
  "temperature": 0.7,
  "maxTokens": 2000
}
```

*Campos:*
- `provider`: "groq" u "openai"
- `model`: Nombre del modelo
- `temperature`: 0.0 (preciso) a 1.0 (creativo)
- `maxTokens`: Longitud máxima de respuesta

=== EMBEDDING_CONFIG

Configuración de embeddings para búsqueda:

```json
{
  "provider": "openai",
  "model": "text-embedding-3-small",
  "dimensions": 1536
}
```

=== WELCOME_MESSAGE

Mensaje de bienvenida del chatbot:

```json
{
  "text": "¡Hola! Soy el asistente virtual del ISTS. ¿En qué puedo ayudarte?",
  "showMenu": true,
  "menuOptions": [
    "Información de carreras",
    "Proceso de inscripción",
    "Costos y becas",
    "Ubicación y contacto"
  ]
}
```

=== RAG_SEARCH_CONFIG

Configuración de búsqueda semántica:

```json
{
  "similarityThreshold": 0.7,
  "maxResults": 5,
  "contextWindow": 3
}
```

*Campos:*
- `similarityThreshold`: Mínimo de similitud (0.0 a 1.0)
- `maxResults`: Cantidad de chunks a recuperar
- `contextWindow`: Chunks vecinos a incluir


== Recargar Caché de Parámetros

El sistema cachea parámetros en memoria para mayor rendimiento. Después de editar parámetros críticos:

+ Haga clic en el botón *"Recargar Caché"*
+ El sistema actualizará la cache
+ Los cambios se aplicarán inmediatamente

*Cuándo recargar:*
- Después de modificar parámetros de configuración LLM
- Al cambiar configuración de embeddings
- Cuando los cambios no se reflejan inmediatamente

#figure(
  image("assets/adm_refresh.png", width: 60%),
  caption: [Recargar parámetros]
)

#pagebreak()

= Gestión de Usuarios

== Acceder a Gestión de Usuarios

+ En el menú lateral, vaya a *"Sistema"* → *"Usuarios"*
+ Se mostrará la lista de usuarios administradores

// TODO: Add screenshot of users list
// #figure(
//   image("images/admin_users_list.png", width: 100%),
//   caption: [Lista de usuarios administradores]
// )

== Lista de Usuarios

La tabla muestra:
- *Nombre de usuario*: Login del usuario
- *Nombre completo*: Nombre real del administrador
- *Email*: Correo electrónico
- *Rol*: Admin, Super Admin, Moderador
- *Estado*: Activo/Inactivo
- *Último acceso*: Fecha del último login
- *Acciones*: Editar, desactivar, eliminar

== Crear Nuevo Usuario

+ Haga clic en el botón *"+ Nuevo Usuario"*
+ Complete el formulario:

// TODO: Add screenshot of create user form
// #figure(
//   image("images/admin_users_create.png", width: 80%),
//   caption: [Formulario de creación de usuario]
// )

*Campos requeridos:*

- *Nombre de usuario*:
  - Identificador único para login
  - Sin espacios, solo letras y números
  - Mínimo 4 caracteres
  - Ejemplo: `jperez`

- *Nombre completo*:
  - Nombre real del usuario
  - Ejemplo: `Juan Pérez`

- *Email*:
  - Correo electrónico válido
  - Se usará para recuperación de contraseña
  - Debe ser único en el sistema

- *Contraseña*:
  - Mínimo 8 caracteres
  - Debe incluir: mayúsculas, minúsculas, números
  - Ejemplo: `Clave123`

- *Confirmar contraseña*:
  - Debe coincidir con la contraseña

- *Rol*:
  - Seleccione el nivel de acceso (Admin por defecto)

*Ejemplo:*
```
Usuario: mrodriguez
Nombre: María Rodríguez
Email: mrodriguez@ists.edu.ec
Contraseña: AdminISTS2025
Rol: Admin
```

+ Haga clic en *"Crear Usuario"*
+ El usuario recibirá un correo de bienvenida con sus credenciales

== Editar Usuario

+ Haga clic en el icono de lápiz (editar) en la fila del usuario
+ Modifique los campos necesarios:

*Campos editables:*
- Nombre completo
- Email
- Rol
- Estado (Activo/Inactivo)

*No editable:*
- Nombre de usuario (es el identificador)

+ Haga clic en *"Guardar Cambios"*

// TODO: Add screenshot of edit user dialog
// #figure(
//   image("images/admin_users_edit.png", width: 80%),
//   caption: [Diálogo de edición de usuario]
// )

== Restablecer Contraseña

Si un usuario olvidó su contraseña:

*Opción 1: Usuario lo hace desde login*
+ El usuario hace clic en "¿Olvidaste tu contraseña?" en el login
+ Ingresa su email
+ Recibe enlace de restablecimiento

*Opción 2: Admin restablece manualmente*
+ En la lista de usuarios, haga clic en el menú de opciones (tres puntos)
+ Seleccione *"Restablecer contraseña"*
+ Elija:
  - Generar contraseña temporal y enviarla por email
  - Establecer contraseña manualmente

// TODO: Add screenshot of password reset dialog
// #figure(
//   image("images/admin_users_reset_password.png", width: 70%),
//   caption: [Diálogo de restablecimiento de contraseña]
// )

== Desactivar Usuario

Para deshabilitar temporalmente un usuario sin eliminarlo:

+ Haga clic en el toggle de estado del usuario
+ El usuario cambiará a *"Inactivo"*
+ El usuario no podrá iniciar sesión hasta reactivarlo

*Efecto:*
- Si el usuario está conectado, se cerrará su sesión
- No podrá iniciar sesión nuevamente
- Los datos y permisos se conservan

*Reactivar:*
+ Haga clic nuevamente en el toggle
+ El estado cambiará a *"Activo"*

== Eliminar Usuario

+ Haga clic en el icono de papelera (eliminar)
+ Confirme la eliminación en el diálogo

// TODO: Add screenshot of delete user confirmation
// #figure(
//   image("images/admin_users_delete.png", width: 60%),
//   caption: [Confirmación de eliminación de usuario]
// )

*Advertencia:*
- Esta acción es irreversible
- Todo el historial de acciones del usuario se mantendrá pero sin vínculo
- No puede eliminar su propio usuario
- No puede eliminar al último Super Admin

#pagebreak()

= Estadísticas y Reportes

== Acceder a Estadísticas

+ En el menú lateral, haga clic en *"Estadísticas"*
+ Se mostrará el panel de analytics

// TODO: Add screenshot of statistics main page
#figure(
  image("assets/adm_estadisticas.png", width: 100%),
  caption: [Panel principal de estadísticas]
)

== Panel de Estadísticas

=== KPIs Principales

Tarjetas superiores con métricas clave:

*Total de Mensajes:*
- Mensajes procesados en el período seleccionado
- Comparación con período anterior

*Usuarios Únicos:*
- Cantidad de usuarios diferentes que interactuaron
- Tendencia de crecimiento

*Tiempo de Respuesta Promedio:*
- Latencia promedio del chatbot
- Medida en milisegundos

*Tasa de Satisfacción:*
- Porcentaje de consultas resueltas exitosamente
- Basado en feedback de usuarios

// TODO: Add screenshot of KPI cards
// #figure(
//   image("images/admin_statistics_kpis.png", width: 100%),
//   caption: [Indicadores clave de rendimiento]
// )

=== Gráficos de Tendencias

*Mensajes por Día:*
- Gráfico de líneas mostrando volumen diario
- Diferenciación entre mensajes recibidos y enviados
- Rango configurable (7, 30, 90 días)

*Mensajes por Hora:*
- Gráfico de barras con distribución horaria
- Identifica horas pico de actividad
- Útil para planificar horarios de atención

*Temas Más Consultados:*
- Gráfico de pastel o barras
- Top 10 temas según chunks más usados
- Ayuda a identificar necesidades de información

// TODO: Add screenshot of trend charts
// #figure(
//   image("images/admin_statistics_trends.png", width: 100%),
//   caption: [Gráficos de tendencias de uso]
// )

=== Estadísticas de Conversaciones

*Por estado:*
- Abiertas: Conversaciones activas
- Cerradas: Conversaciones finalizadas
- Abandonadas: Sin respuesta del usuario >24h

*Por origen:*
- WhatsApp
- Web (si aplica)

*Duración promedio:*
- Tiempo desde primer mensaje hasta resolución
- Cantidad promedio de mensajes por conversación

=== Rendimiento del Sistema

*Uso de Recursos:*
- Tokens LLM consumidos
- Llamadas a API de embeddings
- Costos asociados

*Precisión RAG:*
- Tasa de éxito de búsqueda semántica
- Chunks promedio por respuesta
- Similitud promedio de chunks usados

*Errores:*
- Cantidad de errores por tipo
- Timeouts
- Fallos de API

// TODO: Add screenshot of system performance metrics
// #figure(
//   image("images/admin_statistics_performance.png", width: 100%),
//   caption: [Métricas de rendimiento del sistema]
// )

== Filtros y Rangos de Tiempo

=== Selector de Período

+ Use el selector de rango en la parte superior
+ Opciones predefinidas:
  - Hoy
  - Últimos 7 días
  - Últimos 30 días
  - Este mes
  - Mes anterior
  - Personalizado

*Rango personalizado:*
+ Seleccione "Personalizado"
+ Elija fecha de inicio y fin
+ Haga clic en "Aplicar"

// TODO: Add screenshot of date range selector
// #figure(
//   image("images/admin_statistics_date_range.png", width: 70%),
//   caption: [Selector de rango de fechas]
// )

=== Filtros Adicionales

*Por usuario:*
- Filtre estadísticas de un usuario específico
- Útil para analizar casos particulares

*Por tipo de consulta:*
- Información académica
- Procesos administrativos
- Ubicación y contacto
- Otros

*Por resultado:*
- Exitosas (bot respondió correctamente)
- Derivadas a humano
- Sin respuesta satisfactoria

== Exportar Reportes

Para generar un reporte en PDF:

+ Haga clic en el botón *"Generar Reporte PDF"*
+ Seleccione el período de tiempo deseado
+ El sistema generará un documento PDF con:
  - KPIs generales
  - Gráficos de tendencias
  - Estadísticas principales
+ El PDF se descargará automáticamente

#pagebreak()

= Configuración del Sistema

== Acceder a Configuración

+ Haga clic en su avatar en la esquina superior derecha
+ Seleccione *"Configuración"* del menú
+ O vaya a *"Sistema"* → *"Configuración"* en el menú lateral

// TODO: Add screenshot of settings page
// #figure(
//   image("images/admin_settings_main.png", width: 100%),
//   caption: [Página de configuración del sistema]
// )

== Secciones de Configuración

=== General

*Información Institucional:*
- Nombre de la institución
- Logo (cargar nueva imagen)
- Colores del tema personalizado
- Zona horaria
- Idioma del sistema

// TODO: Add screenshot of general settings
// #figure(
//   image("images/admin_settings_general.png", width: 90%),
//   caption: [Configuración general del sistema]
// )

*Cambiar logo:*
+ Haga clic en *"Cambiar logo"*
+ Seleccione una imagen (PNG, JPG, SVG)
+ Recomendado: 200x200px, fondo transparente
+ Haga clic en *"Subir"*
+ El logo se actualizará en todo el sistema

=== Seguridad

*Políticas de Contraseña:*
- Longitud mínima: 8 caracteres
- Requiere mayúsculas, minúsculas y números
- Cambio recomendado cada 90 días

*Sesiones:*
- Tiempo de inactividad: 30 minutos
- Duración máxima: 8 horas

// TODO: Add screenshot of security settings
// #figure(
//   image("images/admin_settings_security.png", width: 90%),
//   caption: [Configuración de seguridad]
// )

=== Copia de Seguridad

Los backups automáticos se realizan diariamente e incluyen:
- Base de datos completa
- Documentos RAG cargados
- Configuración del sistema

// TODO: Add screenshot of backup settings
// #figure(
//   image("images/admin_settings_backup.png", width: 80%),
//   caption: [Configuración de copias de seguridad]
// )

Para realizar un backup manual, contacte al área de TI.

#pagebreak()

= Solución de Problemas

== Problemas Comunes

=== No Puedo Iniciar Sesión

*Síntomas:*
- Mensaje "Credenciales incorrectas"
- Mensaje "Usuario no encontrado"

*Soluciones:*

+ Verifique que está escribiendo correctamente usuario y contraseña
  - Revise que Bloq Mayús esté desactivado
  - Copie y pegue credenciales si las tiene guardadas

+ Use la opción "¿Olvidaste tu contraseña?"
  - Ingrese su email
  - Revise su bandeja de entrada y spam
  - Siga las instrucciones del correo

+ Contacte al administrador del sistema
  - Es posible que su cuenta esté desactivada
  - El administrador puede restablecer su contraseña

=== El Código QR de WhatsApp No Funciona

*Síntomas:*
- El código QR no se genera
- Al escanear, WhatsApp muestra error
- Conexión no se establece

*Soluciones:*

+ Actualice la página del navegador
  - Presione F5 o Ctrl+R
  - El código QR se regenerará

+ Verifique que su teléfono esté conectado a Internet
  - Use WiFi estable
  - Verifique que WhatsApp funcione normalmente

+ Use WhatsApp Business
  - La conexión solo funciona con WhatsApp Business
  - Verifique que tiene la app correcta instalada

+ Solicite un nuevo código QR
  - Haga clic en "Generar nuevo QR"
  - Los códigos expiran después de 30 segundos

+ Reinicie la conexión
  - Haga clic en "Desconectar" y luego "Conectar"
  - Intente el proceso nuevamente

=== Los Chunks No se Crean o Indexan

*Síntomas:*
- Error al crear chunk
- Chunk creado pero no aparece en búsquedas
- Mensaje de error de embedding

*Soluciones:*

+ Verifique el contenido del chunk
  - No debe estar vacío
  - Mínimo recomendado: 100 caracteres
  - Evite solo números o símbolos

+ Espere unos segundos
  - La generación de embedding puede tomar tiempo
  - Refresque la página después de 10 segundos

+ Verifique que seleccionó un documento
  - Todo chunk debe pertenecer a un documento
  - Cree el documento primero si no existe

+ Contacte soporte técnico
  - Puede haber problema con el servicio de embeddings
  - Proporcione el contenido del chunk que intentó crear

=== El Chatbot No Responde Correctamente

*Síntomas:*
- Respuestas irrelevantes o incorrectas
- El bot dice "no tengo información sobre eso"
- Respuestas muy genéricas

*Soluciones:*

+ Revise el conocimiento RAG
  - Verifique que hay documentos procesados
  - Asegúrese de que los documentos contengan la información necesaria
  - Agregue más documentos si faltan temas

+ Contacte al área de TI
  - Puede requerir ajuste de parámetros del sistema
  - Es posible que necesite verificar la configuración de APIs

=== El Panel Se Ve Lento o No Carga

*Síntomas:*
- Páginas tardan mucho en cargar
- Botones no responden
- Navegación está congelada

*Soluciones:*

+ Actualice el navegador
  - Presione F5 o Ctrl+R
  - O Ctrl+Shift+R para forzar recarga completa

+ Limpie la caché del navegador
  - Presione Ctrl+Shift+Delete
  - Seleccione "Caché e imágenes"
  - Limpie

+ Verifique su conexión a Internet
  - Haga un test de velocidad
  - Reinicie su router si es necesario

+ Use un navegador moderno actualizado
  - Chrome, Firefox, Edge en sus últimas versiones
  - Evite navegadores antiguos o no soportados

+ Cierre pestañas innecesarias
  - Demasiadas pestañas pueden ralentizar el navegador
  - Cierre otras aplicaciones pesadas

+ Contacte al administrador de TI
  - Puede haber problemas en el servidor
  - Solicite verificación del estado del sistema

== Mensajes de Error

=== "Sesión Expirada"

*Causa:* Su sesión de login ha caducado por inactividad.

*Solución:*
+ Haga clic en "Aceptar"
+ Será redirigido al login
+ Vuelva a iniciar sesión
+ Para sesiones más largas, ajuste el timeout en "Configuración" → "Seguridad"

=== "Sin Permisos"

*Causa:* Está intentando acceder a una función para la que no tiene permisos.

*Solución:*
+ Verifique su rol de usuario
+ Contacte a un Super Admin para solicitar permisos adicionales
+ Use solo las funciones para las que está autorizado

=== "Error de Conexión"

*Causa:* El navegador no puede conectarse al servidor.

*Solución:*
+ Verifique su conexión a Internet
+ Refresque la página
+ Verifique que el servidor esté en línea (contacte TI)
+ Si persiste, limpie caché del navegador

=== "Archivo Demasiado Grande"

*Causa:* El archivo que intenta subir excede el límite.

*Solución:*
+ Verifique el tamaño del archivo
+ Comprima el archivo si es posible
+ Divida documentos muy grandes en partes
+ Contacte al administrador para aumentar el límite si es necesario

== Soporte Técnico

Si los problemas persisten después de probar las soluciones:

=== Recopilar Información

Antes de contactar soporte, reúna:
+ Descripción detallada del problema
+ Pasos para reproducir el error
+ Mensaje de error exacto (captura de pantalla)
+ Navegador y versión que está usando
+ Fecha y hora cuando ocurrió el problema
+ Su usuario y rol

=== Contactar Soporte

*Email:* #link("soporte-tecnico@ists.edu.ec")

*Asunto:* "Soporte Chatbot - [Breve descripción]"

*Contenido del correo:*
```
Usuario: [su usuario]
Rol: [su rol]
Navegador: [ej: Chrome 120]
Fecha/Hora del problema: [ej: 27/10/2025 10:30 AM]

Descripción del problema:
[Describa lo que estaba haciendo cuando ocurrió]

Pasos para reproducir:
1. [Primer paso]
2. [Segundo paso]
3. [Tercer paso]

Mensaje de error:
[Copie el mensaje exacto o adjunte captura]

Archivos adjuntos:
[Capturas de pantalla si aplica]
```

=== Reportar Bugs

Si descubre un error del sistema:

+ Anote todos los detalles del error
+ Tome capturas de pantalla
+ Verifique si puede reproducirlo consistentemente
+ Reporte a soporte técnico
+ No intente "arreglar" el error usted mismo si no es parte de TI

#pagebreak()

= Mejores Prácticas

== Seguridad

=== Contraseñas Seguras

*Recomendaciones:*
+ Use contraseñas únicas para cada sistema
+ Mínimo 12 caracteres con combinación de:
  - Mayúsculas y minúsculas
  - Números
  - Caracteres especiales (@, \#, \$, %, etc.)

+ No use información personal (nombre, fecha de nacimiento)
+ No comparta su contraseña con nadie
+ Cambie su contraseña periódicamente (cada 90 días)
+ Use un gestor de contraseñas si es posible

*Ejemplo de contraseña segura:*
`Ch4tb0t!STS@2025`

=== Protección de Sesión

+ Siempre cierre sesión al terminar
+ No deje la sesión abierta en computadoras compartidas
+ Bloquee su computadora al alejarse (Win+L en Windows)
+ No permita que el navegador guarde la contraseña en computadoras públicas

=== Verificación de Actividad

+ Revise periódicamente su actividad en el sistema
+ Reporte inmediatamente actividad sospechosa
+ Si detecta un login no autorizado, cambie su contraseña inmediatamente

== Gestión de Contenido

=== Calidad de Documentos RAG

*Antes de cargar:*
+ Revise que el documento esté actualizado
+ Verifique ortografía y gramática
+ Use formato claro con títulos y secciones
+ Evite imágenes escaneadas de baja calidad

*Organización:*
+ Use nombres descriptivos para documentos
+ Categorice apropiadamente
+ Evite duplicar información entre documentos
+ Mantenga documentos relacionados en la misma categoría

*Mantenimiento:*
+ Revise documentos trimestralmente
+ Actualice información desactualizada
+ Elimine documentos obsoletos
+ Registre fecha de última revisión en la descripción

=== Gestión de Conversaciones

*Respuestas Manuales:*
+ Sea profesional y cortés
+ Use gramática y ortografía correcta
+ Responda de manera clara y concisa
+ Personalice la respuesta (use el nombre del usuario)
+ Proporcione información verificada

*Tiempos de respuesta:*
+ Responda mensajes prioritarios en menos de 1 hora
+ Mensajes normales en menos de 24 horas
+ Use estados para indicar que está trabajando en una respuesta

*Escalamiento:*
+ Si no puede resolver, derive a la persona apropiada
+ No invente respuestas si no sabe
+ Documente consultas complejas para crear FAQs

== Monitoreo y Mantenimiento

=== Revisión Diaria

Al iniciar su jornada:
+ Revise el Dashboard para métricas del día anterior
+ Verifique estado de conexión WhatsApp
+ Revise mensajes no leídos en Chats
+ Identifique y responda consultas urgentes

=== Revisión Semanal

Una vez por semana:
+ Revise estadísticas de la semana
+ Identifique temas más consultados
+ Actualice documentos RAG si es necesario
+ Verifique logs para detectar errores recurrentes
+ Revise y responda feedback de usuarios

=== Revisión Mensual

Una vez al mes:
+ Genere reporte mensual completo
+ Analice tendencias y patrones
+ Planifique mejoras basadas en datos
+ Actualice conocimiento obsoleto
+ Revise y optimice parámetros del sistema
+ Capacite al equipo sobre nuevas funcionalidades

== Privacidad y Datos

=== Protección de Datos de Usuarios

+ No comparta información personal de usuarios fuera del sistema
+ No use datos de conversaciones para propósitos no autorizados
+ Respete la privacidad de las comunicaciones
+ Elimine conversaciones según política de retención

=== Cumplimiento

+ Cumpla con normativas de protección de datos
+ Mantenga confidencialidad de información sensible
+ No exporte datos sin autorización
+ Use cifrado al compartir información sensible

#pagebreak()

= Apéndices

== Glosario

*API (Application Programming Interface):*
Interfaz que permite la comunicación entre diferentes sistemas de software.

*Chunk:*
Fragmento de texto indexado del sistema RAG. Un documento se divide en múltiples chunks.

*Dashboard:*
Panel de control con visualización de métricas y estadísticas clave.

*Embedding:*
Representación vectorial de texto que permite búsquedas semánticas.

*JSON (JavaScript Object Notation):*
Formato de datos estructurados usado para almacenar configuración.

*LLM (Large Language Model):*
Modelo de inteligencia artificial que genera respuestas en lenguaje natural.

*QR Code:*
Código bidimensional que contiene información, usado para vincular WhatsApp.

*RAG (Retrieval Augmented Generation):*
Sistema que combina búsqueda de información con generación de respuestas por IA.

*Similaridad Semántica:*
Medida de qué tan relacionados están dos textos en significado.

*Token:*
Unidad de texto procesada por el LLM (aproximadamente 4 caracteres).

*Webhook:*
URL que recibe notificaciones automáticas cuando ocurre un evento.

== Atajos de Teclado

*Navegación General:*
- `Ctrl + K`: Abrir búsqueda rápida
- `Esc`: Cerrar diálogos y modales
- `Alt + 1-9`: Ir a sección del menú (1=Dashboard, 2=Chats, etc.)

*Gestión de Chats:*
- `Ctrl + Enter`: Enviar mensaje
- `Shift + Enter`: Nueva línea en mensaje
- `Ctrl + F`: Buscar en conversaciones
- `↑` / `↓`: Navegar entre conversaciones

*Edición:*
- `Ctrl + S`: Guardar cambios
- `Ctrl + Z`: Deshacer
- `Ctrl + Y`: Rehacer

*Sistema:*
- `F5`: Actualizar página
- `Ctrl + Shift + R`: Forzar recarga completa

== Recursos Adicionales

=== Documentación

- Manual de Usuario (para usuarios finales del chatbot)
- Manual del Programador (para desarrolladores del sistema)
- Documentación de API: `http://localhost:8080/docs`

=== Videos Tutoriales

(Enlaces a videos cuando estén disponibles)
- Configuración Inicial del Sistema
- Gestión de Documentos RAG
- Conexión de WhatsApp
- Uso del Panel de Chats

=== FAQ del Sistema

*¿Puedo tener múltiples sesiones abiertas?*
Sí, puede iniciar sesión desde diferentes dispositivos simultáneamente.

*¿Cuánto tiempo se guardan las conversaciones?*
Por defecto, las conversaciones se mantienen indefinidamente. Puede configurar retención automática.

*¿Puedo recuperar documentos eliminados?*
No, la eliminación es permanente. Use con precaución.

*¿El sistema tiene límites de uso?*
Sí, hay límites de API según el plan contratado (Groq, OpenAI).

*¿Puedo personalizar los mensajes automáticos?*
Sí, en "WhatsApp" → "Configuración" → "Mensajes Automáticos".

== Actualizaciones del Manual

*Versión 1.0 - Octubre 2025*
- Versión inicial del manual de administrador
- Cobertura completa de funcionalidades MVP
- Guías paso a paso con capturas de pantalla

*Próximas actualizaciones incluirán:*
- Nuevas funcionalidades añadidas al sistema
- Videos tutoriales integrados
- Casos de uso avanzados
- Automatizaciones y workflows

== Contacto

=== Soporte Técnico
- *Email:* #link("tecnico@ists.edu.ec")
- *Horario:* Lunes a Viernes, 8:00 AM - 6:00 PM

=== Capacitación
- *Email:* #link("capacitacion@ists.edu.ec")
- Solicite sesiones de capacitación para su equipo

---

#align(center)[
  #text(size: 10pt, style: "italic")[
    Este manual está sujeto a actualizaciones periódicas. \
    Consulte la versión más reciente en el panel de administración.
  ]

  #v(1cm)

  #text(size: 9pt)[
    © 2025 Instituto Superior Tecnológico Sudamericano \
    Todos los derechos reservados
  ]
]
