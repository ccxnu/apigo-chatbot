package pdfprocessor

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"regexp"
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

	// Clean up the extracted text (remove unwanted line breaks)
	cleanedText := cleanPDFText(text)

	// If text extraction yields very little content, it might be a scanned PDF
	// In that case, attempt OCR (this is a simplified check)
	if len(strings.TrimSpace(cleanedText)) < 50 {
		// For OCR, we would need to convert PDF to images first
		// This is a placeholder - full OCR implementation would be more complex
		return cleanedText, nil // Return what we have for now
	}

	return cleanedText, nil
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

// cleanPDFText cleans up text extracted from PDF by normalizing line breaks
func cleanPDFText(text string) string {
	// Replace Windows line endings with Unix
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`[ \t]+`)
	text = spaceRegex.ReplaceAllString(text, " ")

	// Handle hyphenated words at line breaks (e.g., "exam-\nple" -> "example")
	hyphenRegex := regexp.MustCompile(`-\s*\n\s*`)
	text = hyphenRegex.ReplaceAllString(text, "")

	// Replace single line breaks with spaces (joining lines in same paragraph)
	// But preserve double line breaks (paragraph separators)
	text = regexp.MustCompile(`([^\n])\n([^\n])`).ReplaceAllString(text, "$1 $2")

	// Normalize multiple consecutive line breaks to double line break (paragraph separator)
	multiLineRegex := regexp.MustCompile(`\n{3,}`)
	text = multiLineRegex.ReplaceAllString(text, "\n\n")

	// Trim leading/trailing whitespace from each paragraph
	paragraphs := strings.Split(text, "\n\n")
	for i, para := range paragraphs {
		paragraphs[i] = strings.TrimSpace(para)
	}
	text = strings.Join(paragraphs, "\n\n")

	// Final cleanup: remove leading/trailing whitespace
	text = strings.TrimSpace(text)

	return text
}

// Note: OCR functionality for scanned PDFs can be added later by:
// 1. Installing Tesseract OCR (apt-get install tesseract-ocr libtesseract-dev libleptonica-dev)
// 2. Using github.com/otiai10/gosseract/v2 for OCR
// 3. Converting PDF pages to images first using a library like unipdf
// For now, this implementation focuses on text-based PDFs only
