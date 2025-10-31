package mailer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	d "api-chatbot/domain"
)

// OTPMailer implements domain.OTPMailer interface
type OTPMailer struct {
	httpClient     d.HTTPClient
	tikeeURL       string
	senderEmail    string
	paramCache     d.ParameterCache
	contextTimeout time.Duration
}

func NewOTPMailer(
	httpClient d.HTTPClient,
	tikeeURL string,
	senderEmail string,
	paramCache d.ParameterCache,
	timeout time.Duration,
) *OTPMailer {
	return &OTPMailer{
		httpClient:     httpClient,
		tikeeURL:       tikeeURL,
		senderEmail:    senderEmail,
		paramCache:     paramCache,
		contextTimeout: timeout,
	}
}

// EmailRequest represents the request structure for Tikee email API
type EmailRequest struct {
	IDUser          string   `json:"idUser"`
	IDInstitution   int      `json:"idInstitution"`
	IDSolution      int      `json:"idSolution"`
	IDRequest       string   `json:"idRequest"`
	ProcessDate     string   `json:"processDate"`
	Process         string   `json:"process"`
	Tipo            string   `json:"tipo"`
	Prioridad       string   `json:"prioridad"`
	Destinos        []string `json:"destinos"`
	Plantilla       string   `json:"plantilla"`
	ListaValsEmail  []string `json:"lista_vals_email"`
	ListaVarsEmail  []string `json:"lista_vars_email"`
	Asunto          string   `json:"asunto"`
	NomEnvia        string   `json:"nom_envia"`
	CorreoEnvia     string   `json:"correo_envia"`
	Intentos        int      `json:"intentos"`
}

// SendOTPEmail sends an OTP code to the user's email
func (m *OTPMailer) SendOTPEmail(ctx context.Context, email, name, otpCode, userType string) error {
	// Get email template from parameters
	templateParam, exists := m.paramCache.Get("EMAIL_OTP_TEMPLATE")
	var template string
	var subject string

	if exists {
		if data, err := templateParam.GetDataAsMap(); err == nil {
			if tmpl, ok := data["html"].(string); ok {
				template = tmpl
			}
			if subj, ok := data["subject"].(string); ok {
				subject = subj
			}
		}
	}

	// Fallback template if not found in parameters
	if template == "" {
		template = m.getDefaultOTPTemplate()
	}

	if subject == "" {
		subject = "Código de verificación - Chatbot ISTS"
	}

	// Get current time
	currentTime := time.Now()
	fecha := currentTime.Format("2006-01-02")
	hora := currentTime.Format("15:04:05")

	// Determine user type in Spanish
	userTypeText := "usuario"
	switch userType {
	case "institute":
		userTypeText = "miembro de la institución"
	case "external":
		userTypeText = "usuario externo"
	}

	// Build email request (matching TypeScript defaults)
	emailReq := EmailRequest{
		IDUser:         email,
		IDInstitution:  0,
		IDSolution:     0,
		IDRequest:      fmt.Sprintf("otp-%d", time.Now().Unix()),
		ProcessDate:    currentTime.Format("2006-01-02 15:04:05"),
		Process:        "CHATBOT_OTP_VERIFICATION",
		Tipo:           "EMAIL",
		Prioridad:      "MEDIA", // Changed from ALTA to match TypeScript default
		Destinos:       []string{email},
		Plantilla:      template,
		ListaValsEmail: []string{name, email, otpCode, fecha, hora, userTypeText},
		ListaVarsEmail: []string{"[nom_usuario]", "[nom_email]", "[codigo_otp]", "[fecha]", "[hora]", "[tipo_usuario]"},
		Asunto:         subject,
		NomEnvia:       "automatizaciones@tikee.tech", // AWS SES verified email
		CorreoEnvia:    "automatizaciones@tikee.tech", // AWS SES verified email - must match NomEnvia
		Intentos:       0, // Changed from 1 to match TypeScript default
	}

	// Log the request for debugging
	slog.Debug("Sending OTP email request",
		"email", email,
		"destinos", emailReq.Destinos,
		"process", emailReq.Process,
	)

	// Create HTTP request - HTTPClient will marshal the body automatically
	req := d.HTTPRequest{
		URL:    m.tikeeURL,
		Method: "POST",
		AdditionalHeaders: []d.HTTPHeader{
			{Key: "Content-Type", Value: "application/json"},
		},
		Body: emailReq, // Pass the struct directly, not marshaled bytes
	}

	// Send request
	var response map[string]interface{}
	if err := m.httpClient.Do(ctx, req, &response); err != nil {
		slog.Error("Failed to send OTP email via Tikee API",
			"error", err,
			"email", email,
		)
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	// Log full response for debugging
	slog.Debug("Tikee API response", "response", response)

	// Check response
	if code, ok := response["code"].(string); ok && code == "COD_ERR" {
		errMsg := "Unknown error"
		if msg, ok := response["message"].(string); ok {
			errMsg = msg
		}
		// Also check for "info" field which may contain more details
		if info, ok := response["info"].(string); ok && info != "" {
			errMsg = info
		}
		slog.Error("Tikee API returned error",
			"code", code,
			"message", errMsg,
			"fullResponse", response,
			"email", email,
		)
		return fmt.Errorf("tikee API error: %s", errMsg)
	}

	slog.Info("OTP email sent successfully",
		"email", email,
		"userType", userType,
	)

	return nil
}

// getDefaultOTPTemplate returns the default HTML template for OTP emails
func (m *OTPMailer) getDefaultOTPTemplate() string {
	return `<table bgcolor="#fefefe" border="0" cellpadding="5" cellspacing="0" width="100%">
  <tbody>
    <tr align="center">
      <td>
        <table width="581" cellspacing="10" cellpadding="0">
          <tbody>
            <tr>
              <td style="font-family:Arial;font-size:12px;color:#333333">
                <p style="text-align:justify;margin-top:30px">
                  <b style="font-family:Georgia;font-size:12pt">Saludos [nom_usuario]</b>:
                </p>

                <p style="text-align:justify;margin-top:20px">
                  Has solicitado registrarte en el asistente virtual del Instituto Superior Tecnológico Sucre como <b>[tipo_usuario]</b>.
                </p>

                <table border="0" style="margin:30px 0px" width="100%" align="center" cellpadding="0" cellspacing="0">
                  <tbody>
                    <tr>
                      <td style="font-family:nunito sans,sans-serif;font-size:13pt;color:#5f5b5b;text-align:left;padding:5px"><b>Información:</b></td>
                      <td style="font-family:nunito sans,sans-serif;font-size:13pt;font-weight:300;color:#5f5b5b;text-align:right;padding:5px">Verificación de registro</td>
                    </tr>
                    <tr>
                      <td style="font-family:nunito sans,sans-serif;font-size:13pt;color:#5f5b5b;text-align:left;padding:5px"><b>Usuario:</b></td>
                      <td style="font-family:nunito sans,sans-serif;font-size:13pt;font-weight:300;color:#5f5b5b;text-align:right;padding:5px">[nom_email]</td>
                    </tr>
                    <tr>
                      <td style="font-family:nunito sans,sans-serif;font-size:13pt;color:#5f5b5b;text-align:left;padding:5px"><b>Código de verificación:</b></td>
                      <td style="font-family:nunito sans,sans-serif;font-size:24pt;font-weight:bold;color:#1a73e8;text-align:right;padding:15px;background-color:#f0f0f0;border-radius:5px">[codigo_otp]</td>
                    </tr>
                    <tr>
                      <td style="font-family:nunito sans,sans-serif;font-size:13pt;color:#5f5b5b;text-align:left;padding:5px"><b>Fecha y hora:</b></td>
                      <td style="font-family:nunito sans,sans-serif;font-size:13pt;font-weight:300;color:#5f5b5b;text-align:right;padding:5px">[fecha] a las [hora]</td>
                    </tr>
                  </tbody>
                </table>

                <p style="text-align:justify;background-color:#fff3cd;padding:15px;border-left:4px solid:#ffcc00;margin:20px 0">
                  <b>⚠ Importante:</b> Este código expirará en 10 minutos. No compartas este código con nadie.
                </p>

                <p style="text-align:justify">
                  Para completar tu registro, envía este código de 6 dígitos por WhatsApp al chatbot del Instituto.
                </p>

                <p style="text-align:justify;margin-top:20px">
                  Si no solicitaste este registro, por favor ignora este mensaje y contacta al departamento de sistemas.
                </p>

                <p style="text-align:center;margin-top:40px;font-size:10pt;color:#888">
                  <i>Este es un mensaje automático, por favor no responder.</i>
                </p>
              </td>
            </tr>
          </tbody>
        </table>
      </td>
    </tr>
  </tbody>
</table>`
}

func getDefaultTemplate() string {
	return `<table bgcolor="#fefefe" border="0" cellpadding="5" cellspacing="0" width="100%">
  <tbody>
    <tr align="center">
      <td>
        <table width="581" cellspacing="10" cellpadding="0">
          <tbody>
            <tr>
              <td style="font-family:Arial;font-size:12px;color:#333333">
                <p style="text-align:justify;margin-top:30px">
                  <b style="font-family:Georgia;font-size:12pt">Saludos [nom_usuario]</b>:
                </p>

                <p style="text-align:justify;margin-top:20px">
                  Has solicitado registrarte en el asistente virtual del Instituto Superior Tecnológico Sucre como <b>[tipo_usuario]</b>.
                </p>

                <table border="0" style="margin:30px 0px" width="100%" align="center" cellpadding="0" cellspacing="0">
                  <tbody>
                    <tr>
                      <td style="font-family:nunito sans,sans-serif;font-size:13pt;color:#5f5b5b;text-align:left;padding:5px"><b>Información:</b></td>
                      <td style="font-family:nunito sans,sans-serif;font-size:13pt;font-weight:300;color:#5f5b5b;text-align:right;padding:5px">Verificación de registro</td>
                    </tr>
                    <tr>
                      <td style="font-family:nunito sans,sans-serif;font-size:13pt;color:#5f5b5b;text-align:left;padding:5px"><b>Usuario:</b></td>
                      <td style="font-family:nunito sans,sans-serif;font-size:13pt;font-weight:300;color:#5f5b5b;text-align:right;padding:5px">[nom_email]</td>
                    </tr>
                    <tr>
                      <td style="font-family:nunito sans,sans-serif;font-size:13pt;color:#5f5b5b;text-align:left;padding:5px"><b>Código de verificación:</b></td>
                      <td style="font-family:nunito sans,sans-serif;font-size:24pt;font-weight:bold;color:#1a73e8;text-align:right;padding:15px;background-color:#f0f0f0;border-radius:5px">[codigo_otp]</td>
                    </tr>
                    <tr>
                      <td style="font-family:nunito sans,sans-serif;font-size:13pt;color:#5f5b5b;text-align:left;padding:5px"><b>Fecha y hora:</b></td>
                      <td style="font-family:nunito sans,sans-serif;font-size:13pt;font-weight:300;color:#5f5b5b;text-align:right;padding:5px">[fecha] a las [hora]</td>
                    </tr>
                  </tbody>
                </table>

                <p style="text-align:justify;background-color:#fff3cd;padding:15px;border-left:4px solid:#ffcc00;margin:5px 0">
                  <b>⚠ Importante:</b> Este código expirará en 10 minutos. No compartas este código con nadie.
                </p>

                <p style="text-align:justify">
                  Para completar tu registro, envía este código de 6 dígitos por WhatsApp al chatbot del Instituto.
                </p>

                <p style="text-align:justify;margin-top:20px">
                  Si no solicitaste este registro, por favor ignora este mensaje y contacta al departamento de sistemas.
                </p>

                <p style="text-align:center;margin-top:40px;font-size:10pt;color:#888">
                  <i>Este es un mensaje automático, por favor no responder.</i>
                </p>
              </td>
            </tr>
          </tbody>
        </table>
      </td>
    </tr>
  </tbody>
</table>`
}
