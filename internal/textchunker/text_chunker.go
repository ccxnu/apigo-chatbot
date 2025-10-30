package textchunker

import (
	"strings"
	"unicode"
)

// ChunkText splits text into chunks with optional overlap
func ChunkText(text string, chunkSize, overlap int) []string {
	if chunkSize <= 0 {
		chunkSize = 1000 // Default chunk size
	}
	if overlap < 0 {
		overlap = 0
	}
	if overlap >= chunkSize {
		overlap = chunkSize / 2 // Ensure overlap is less than chunk size
	}

	// Clean up the text
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{}
	}

	// Split by sentences first for better chunk boundaries
	sentences := splitIntoSentences(text)

	var chunks []string
	var currentChunk strings.Builder
	var previousChunk string
	currentSize := 0

	for _, sentence := range sentences {
		sentenceLen := len(sentence)

		// If adding this sentence would exceed chunk size
		if currentSize > 0 && currentSize+sentenceLen > chunkSize {
			// Save current chunk
			chunk := currentChunk.String()
			chunks = append(chunks, strings.TrimSpace(chunk))

			// Start new chunk with overlap from previous chunk
			currentChunk.Reset()
			if overlap > 0 && len(previousChunk) > 0 {
				overlapText := getLastNChars(previousChunk, overlap)
				currentChunk.WriteString(overlapText)
				currentSize = len(overlapText)
			} else {
				currentSize = 0
			}

			previousChunk = chunk
		}

		// Add sentence to current chunk
		if currentChunk.Len() > 0 {
			currentChunk.WriteString(" ")
			currentSize++
		}
		currentChunk.WriteString(sentence)
		currentSize += sentenceLen
	}

	// Add the last chunk if not empty
	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return chunks
}

// splitIntoSentences splits text into sentences
func splitIntoSentences(text string) []string {
	var sentences []string
	var currentSentence strings.Builder

	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		currentSentence.WriteRune(r)

		// Check for sentence endings
		if r == '.' || r == '!' || r == '?' {
			// Look ahead to see if it's really the end of a sentence
			if i+1 < len(runes) {
				next := runes[i+1]
				// If followed by space and capital letter, it's likely a sentence end
				if unicode.IsSpace(next) {
					if i+2 < len(runes) && unicode.IsUpper(runes[i+2]) {
						sentence := strings.TrimSpace(currentSentence.String())
						if sentence != "" {
							sentences = append(sentences, sentence)
						}
						currentSentence.Reset()
					}
				}
			} else {
				// End of text
				sentence := strings.TrimSpace(currentSentence.String())
				if sentence != "" {
					sentences = append(sentences, sentence)
				}
				currentSentence.Reset()
			}
		} else if r == '\n' {
			// Treat newlines as potential sentence breaks
			if i+1 < len(runes) && runes[i+1] == '\n' {
				// Double newline is definitely a break
				sentence := strings.TrimSpace(currentSentence.String())
				if sentence != "" {
					sentences = append(sentences, sentence)
				}
				currentSentence.Reset()
				i++ // Skip the next newline
			}
		}
	}

	// Add any remaining text as a sentence
	if currentSentence.Len() > 0 {
		sentence := strings.TrimSpace(currentSentence.String())
		if sentence != "" {
			sentences = append(sentences, sentence)
		}
	}

	return sentences
}

// getLastNChars returns the last n characters from a string
func getLastNChars(s string, n int) string {
	if n <= 0 || len(s) == 0 {
		return ""
	}

	// Try to find a good breaking point (word boundary)
	if n >= len(s) {
		return s
	}

	start := len(s) - n
	substring := s[start:]

	// Try to start at a word boundary
	spaceIdx := strings.Index(substring, " ")
	if spaceIdx > 0 && spaceIdx < len(substring)/2 {
		return strings.TrimSpace(substring[spaceIdx:])
	}

	return strings.TrimSpace(substring)
}
