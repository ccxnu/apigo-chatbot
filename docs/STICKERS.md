# 🎨 Soporte de Stickers en el Chatbot

## ✅ Lo que ya funciona

### Emojis
El chatbot **ya puede usar emojis** en todas sus respuestas. Los emojis están incluidos en el system prompt y el LLM los usará automáticamente para hacer las respuestas más amigables.

Ejemplos de emojis que usa:
- 👋 Saludos
- 📅 Horarios
- 🎓 Carreras
- 📋 Trámites
- 😊 Felicidad
- 🏫 Instituto
- 💬 Conversación

### Stickers
El código **ya soporta enviar stickers**, pero necesitas configurarlos.

## 🚀 Cómo usar Stickers

### 1. Método Básico - Código Manual

Puedes enviar stickers directamente desde el código usando el cliente de WhatsApp:

```go
// En cualquier handler o servicio
err := whatsappClient.SendSticker(chatID, "https://tu-servidor.com/stickers/feliz.webp")
```

### 2. Método Avanzado - Parámetros del Sistema

La migración `000026_add_sticker_support.up.sql` crea un parámetro `BOT_STICKERS` con URLs de stickers predefinidos.

**Estructura del parámetro:**
```json
{
  "enabled": true,
  "stickers": {
    "welcome": "https://example.com/stickers/welcome.webp",
    "thanks": "https://example.com/stickers/thanks.webp",
    "thinking": "https://example.com/stickers/thinking.webp",
    "happy": "https://example.com/stickers/happy.webp",
    "confused": "https://example.com/stickers/confused.webp"
  }
}
```

**Actualiza las URLs** con tus propios stickers alojados.

### 3. Crear Stickers Personalizados

**Requisitos:**
- Formato: **WebP**
- Tamaño recomendado: 512x512 píxeles
- Tamaño máximo: 100 KB
- Fondo: Transparente (recomendado)

**Herramientas para crear stickers:**
1. [sticker.ly](https://sticker.ly/) - Crear desde imágenes
2. [WA Sticker Maker](https://play.google.com/store/apps/details?id=com.marsvard.stickermakerforwhatsapp) - App móvil
3. Photoshop/GIMP + plugin WebP

### 4. Alojar los Stickers

Opciones para alojar tus stickers:

**Opción A: Servidor Propio**
```bash
# En tu servidor
mkdir -p /var/www/chatbot-stickers
# Sube tus .webp files
# Configura nginx para servir estático
```

**Opción B: CDN/Storage**
- AWS S3 + CloudFront
- Google Cloud Storage
- Cloudinary
- imgbb (gratis)

**Opción C: GitHub (para pruebas)**
```bash
# Crea un repo público
# Sube tus .webp
# Usa URLs raw: https://raw.githubusercontent.com/user/repo/main/sticker.webp
```

### 5. Configurar las URLs en la Base de Datos

```sql
UPDATE cht_parameters
SET prm_data = '{
  "enabled": true,
  "stickers": {
    "welcome": "https://tu-cdn.com/welcome.webp",
    "thanks": "https://tu-cdn.com/thanks.webp",
    "thinking": "https://tu-cdn.com/thinking.webp",
    "happy": "https://tu-cdn.com/happy.webp",
    "confused": "https://tu-cdn.com/confused.webp"
  }
}'::jsonb
WHERE prm_code = 'BOT_STICKERS';
```

## 💡 Ideas de Uso

### Respuestas Contextuales
Puedes hacer que el bot envíe stickers basándose en el contexto:

- **Bienvenida** → Sticker de saludo
- **Respuesta exitosa** → Sticker feliz
- **No encontró info** → Sticker confundido
- **Procesando** → Sticker pensando
- **Despedida** → Sticker de adiós

### Ejemplo en Handler

```go
// En el RAG handler, después de generar respuesta
if strings.Contains(strings.ToLower(userMessage), "hola") {
    // Envía sticker de bienvenida
    stickerURL := h.getParam("BOT_STICKERS", "")
    // Parse JSON y obtener URL del sticker "welcome"
    // Luego enviar: client.SendSticker(chatID, url)
}
```

## 🔧 Próximos Pasos

Para hacer que el bot use stickers automáticamente, necesitarías:

1. **Modificar el RAG Handler** para detectar contextos
2. **Leer el parámetro BOT_STICKERS** desde la cache
3. **Enviar el sticker apropiado** según el contexto
4. **Registrar en la BD** el mensaje tipo 'sticker'

## ⚠️ Limitaciones Actuales

- Los stickers deben estar alojados en URLs públicas
- WhatsApp requiere formato WebP
- No hay auto-conversión de formatos (PNG/JPG → WebP)
- No hay validación de URLs en el código actual

## 📚 Referencias

- [WhatsApp Sticker Format](https://faq.whatsapp.com/general/how-to-create-stickers-for-whatsapp)
- [whatsmeow Documentation](https://pkg.go.dev/go.mau.fi/whatsmeow)
- [WebP Converter](https://cloudconvert.com/webp-converter)
