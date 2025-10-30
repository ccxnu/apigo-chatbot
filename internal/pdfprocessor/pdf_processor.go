package pdfprocessor

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

// ExtractTextFromBase64PDF extracts text from a base64-encoded PDF
func ExtractTextFromBase64PDF(base64Data string) (string, error) {
	// Decode base64 string
	pdfData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 PDF: %w", err)
	}

	// Try to extract text directly from PDF
	text, err := extractTextFromPDF(pdfData)
	if err != nil {
		return "", fmt.Errorf("failed to extract text from PDF: %w", err)
	}

	// If text extraction yields very little content, it might be a scanned PDF
	// In that case, attempt OCR (this is a simplified check)
	if len(strings.TrimSpace(text)) < 50 {
		// For OCR, we would need to convert PDF to images first
		// This is a placeholder - full OCR implementation would be more complex
		return text, nil // Return what we have for now
	}

	return text, nil
}

// extractTextFromPDF extracts text from PDF bytes using ledongthuc/pdf
func extractTextFromPDF(pdfData []byte) (string, error) {
	reader := bytes.NewReader(pdfData)

	// Create a ReaderAt from the byte slice
	pdfReader, err := pdf.NewReader(reader, int64(len(pdfData)))
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	var textBuilder strings.Builder
	numPages := pdfReader.NumPage()

	for i := 1; i <= numPages; i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			// Log error but continue with other pages
			continue
		}

		textBuilder.WriteString(text)
		textBuilder.WriteString("\n\n") // Add page separator
	}

	return textBuilder.String(), nil
}

// Note: OCR functionality for scanned PDFs can be added later by:
// 1. Installing Tesseract OCR (apt-get install tesseract-ocr libtesseract-dev libleptonica-dev)
// 2. Using github.com/otiai10/gosseract/v2 for OCR
// 3. Converting PDF pages to images first using a library like unipdf
// For now, this implementation focuses on text-based PDFs only
