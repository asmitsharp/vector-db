package phase1

import "math"

// dot product a.b = a1*b1 + a2*b2 .... + an*bn

func dotProduct(a, b []float64) float64 {
	if len(a) != len(b) {
		panic("vectors must be of the same length")
	}
	var result float64

	for i := range a {
		result += a[i] * b[i]
	}

	return result

}

// magnitude of the vector

func magnitude(v []float64) float64 {
	var magnitude float64
	for i := range v {
		magnitude += v[i] * v[i]
	}

	return math.Sqrt(magnitude)
}

//cosiesimilarity
// 1.0 = identical dir
// -1.0 = opposite dir
// 0.0 = perpendicular (unrelated)

func cosineSimilarity(a, b []float64) float64 {
	var cosine float64

	dotP := dotProduct(a, b)
	magnitudeA := magnitude(a)
	magnitudeB := magnitude(b)

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0.0
	}

	cosine = dotP / (magnitudeA * magnitudeB)

	return cosine
}
