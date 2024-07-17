package similarity

import "math"

func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		panic("Vectors must be the same length")
	}
	var dotProduct, magnitudeA, magnitudeB float64
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		magnitudeA += a[i] * a[i]
		magnitudeB += b[i] * b[i]
	}
	magnitudeA = math.Sqrt(magnitudeA)
	magnitudeB = math.Sqrt(magnitudeB)
	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}
	return dotProduct / (magnitudeA * magnitudeB)
}

func FindMostSimilarVector(questionVector []float64, embeddingArrs [][]float64) (int, float64) {
	var maxSimilarity float64
	var mostSimilarIndex int
	for i, vector := range embeddingArrs {
		similarity := CosineSimilarity(questionVector, vector)
		if similarity > maxSimilarity || i == 0 {
			maxSimilarity = similarity
			mostSimilarIndex = i
		}
	}
	return mostSimilarIndex, maxSimilarity
}
