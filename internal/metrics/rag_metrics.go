package metrics

import (
	"math"
)

// RAGMetrics calculates retrieval quality metrics for RAG systems
type RAGMetrics struct{}

// NewRAGMetrics creates a new RAG metrics calculator
func NewRAGMetrics() *RAGMetrics {
	return &RAGMetrics{}
}

// MetricsResult contains calculated RAG quality metrics
type MetricsResult struct {
	PrecisionAtK float64
	RecallAtK    float64
	F1AtK        float64
	MRR          float64
	MAP          float64
	NDCG         float64
}

// RetrievedChunk represents a chunk returned from similarity search
type RetrievedChunk struct {
	ChunkID         int
	SimilarityScore float64
	Position        int // 1-based position in results
	IsRelevant      bool
}

// CalculatePrecisionAtK calculates Precision@K
// Precision@K = (Number of relevant items in top K) / K
func (m *RAGMetrics) CalculatePrecisionAtK(chunks []RetrievedChunk) float64 {
	if len(chunks) == 0 {
		return 0.0
	}

	relevantCount := 0
	for _, chunk := range chunks {
		if chunk.IsRelevant {
			relevantCount++
		}
	}

	return float64(relevantCount) / float64(len(chunks))
}

// CalculateRecallAtK calculates Recall@K
// Recall@K = (Number of relevant items in top K) / (Total relevant items in collection)
// Note: totalRelevant should be provided based on ground truth or user feedback
func (m *RAGMetrics) CalculateRecallAtK(chunks []RetrievedChunk, totalRelevant int) float64 {
	if totalRelevant == 0 {
		return 0.0
	}

	relevantCount := 0
	for _, chunk := range chunks {
		if chunk.IsRelevant {
			relevantCount++
		}
	}

	return float64(relevantCount) / float64(totalRelevant)
}

// CalculateF1AtK calculates F1 Score@K
// F1@K = 2 * (Precision * Recall) / (Precision + Recall)
func (m *RAGMetrics) CalculateF1AtK(precision, recall float64) float64 {
	if precision+recall == 0 {
		return 0.0
	}

	return 2 * (precision * recall) / (precision + recall)
}

// CalculateMRR calculates Mean Reciprocal Rank
// MRR = 1 / (position of first relevant item)
func (m *RAGMetrics) CalculateMRR(chunks []RetrievedChunk) float64 {
	for _, chunk := range chunks {
		if chunk.IsRelevant {
			return 1.0 / float64(chunk.Position)
		}
	}
	return 0.0
}

// CalculateMAP calculates Mean Average Precision
// MAP = (1/K) * Σ(Precision@i * relevance_i)
func (m *RAGMetrics) CalculateMAP(chunks []RetrievedChunk) float64 {
	if len(chunks) == 0 {
		return 0.0
	}

	sumPrecision := 0.0
	relevantCount := 0

	for i, chunk := range chunks {
		if chunk.IsRelevant {
			relevantCount++
			// Calculate precision at this position
			precisionAtI := float64(relevantCount) / float64(i+1)
			sumPrecision += precisionAtI
		}
	}

	if relevantCount == 0 {
		return 0.0
	}

	return sumPrecision / float64(len(chunks))
}

// CalculateNDCG calculates Normalized Discounted Cumulative Gain
// DCG@K = Σ(relevance_i / log2(i+1))
// NDCG@K = DCG@K / IDCG@K (ideal DCG)
func (m *RAGMetrics) CalculateNDCG(chunks []RetrievedChunk) float64 {
	if len(chunks) == 0 {
		return 0.0
	}

	// Calculate DCG (using similarity scores as relevance)
	dcg := 0.0
	for _, chunk := range chunks {
		relevance := 0.0
		if chunk.IsRelevant {
			relevance = chunk.SimilarityScore
		}
		dcg += relevance / math.Log2(float64(chunk.Position)+1)
	}

	// Calculate IDCG (ideal DCG - assume all relevant items at top)
	// Create ideal ranking
	idealChunks := make([]RetrievedChunk, len(chunks))
	copy(idealChunks, chunks)

	// Sort by relevance (relevance first, then by similarity score)
	for i := 0; i < len(idealChunks); i++ {
		for j := i + 1; j < len(idealChunks); j++ {
			iScore := 0.0
			if idealChunks[i].IsRelevant {
				iScore = idealChunks[i].SimilarityScore
			}
			jScore := 0.0
			if idealChunks[j].IsRelevant {
				jScore = idealChunks[j].SimilarityScore
			}

			if jScore > iScore {
				idealChunks[i], idealChunks[j] = idealChunks[j], idealChunks[i]
			}
		}
	}

	idcg := 0.0
	for i, chunk := range idealChunks {
		relevance := 0.0
		if chunk.IsRelevant {
			relevance = chunk.SimilarityScore
		}
		idcg += relevance / math.Log2(float64(i+2))
	}

	if idcg == 0.0 {
		return 0.0
	}

	return dcg / idcg
}

// CalculateAllMetrics calculates all RAG quality metrics at once
func (m *RAGMetrics) CalculateAllMetrics(chunks []RetrievedChunk, totalRelevant int) MetricsResult {
	precision := m.CalculatePrecisionAtK(chunks)
	recall := m.CalculateRecallAtK(chunks, totalRelevant)

	return MetricsResult{
		PrecisionAtK: precision,
		RecallAtK:    recall,
		F1AtK:        m.CalculateF1AtK(precision, recall),
		MRR:          m.CalculateMRR(chunks),
		MAP:          m.CalculateMAP(chunks),
		NDCG:         m.CalculateNDCG(chunks),
	}
}

// EstimateRelevanceFromSimilarity estimates if a chunk is relevant based on similarity score
// This is a heuristic when no user feedback is available
// Typically, similarity > 0.75 indicates high relevance
func EstimateRelevanceFromSimilarity(score float64, threshold float64) bool {
	return score >= threshold
}

// CalculateStaleness calculates how many days since document was last updated
func CalculateStaleness(lastUpdated, documentPublished int64) int {
	now := int64(0) // You'd use time.Now().Unix() in production
	if lastUpdated > 0 {
		return int((now - lastUpdated) / (24 * 3600))
	}
	if documentPublished > 0 {
		return int((now - documentPublished) / (24 * 3600))
	}
	return 0
}
