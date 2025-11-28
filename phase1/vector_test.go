package phase1

import (
	"math"
	"testing"
)

func TestDotProduct(t *testing.T) {
	a1 := []float64{1, 2, 3}
	b1 := []float64{4, 5, 6}

	result1 := dotProduct(a1, b1)
	expected1 := 32.0

	if result1 != float64(expected1) {
		t.Errorf("dotproduct(a1[1,2,3], b1[4,5,6]) = %f; but expected is %f", result1, expected1)
	}

	a2 := []float64{1, 0, 0}
	b2 := []float64{1, 0, 0}
	result2 := dotProduct(a2, b2)
	expected2 := 1.0

	if result2 != expected2 {
		t.Errorf("dotproduct(a2[1,0,0], b2[1,0,0]) = %f; but expected is %f", result2, expected2)
	}

}

func TestMagnitude(t *testing.T) {
	v1 := []float64{1, 0, 0}

	result1 := magnitude(v1)
	expected1 := 1.0

	if result1 != expected1 {
		t.Errorf("magnitude(v1[1,0,0]) = %f; but expected is %f", result1, expected1)
	}

	v2 := []float64{1, 2, 3}
	result2 := magnitude(v2)
	expected2 := math.Sqrt(14)

	if math.Abs(result2-expected2) > 0.00001 {
		t.Errorf("magnitude(v1[1,2,3]) = %f; but expected is %f", result2, expected2)
	}
}

func TestCosine(t *testing.T) {

	a1 := []float64{1, 2, 3}
	b1 := []float64{4, 5, 6}
	result1 := cosineSimilarity(a1, b1)
	expected1 := (32) / (math.Sqrt(14) * math.Sqrt(77))

	if math.Abs(result1-expected1) > 0.00001 {
		t.Errorf("cosineSimilarity(a1[1,2,3], b1[4,5,6]) = %f; but expected is %f", result1, expected1)
	}

	a2 := []float64{1, 0, 0}
	b2 := []float64{1, 0, 0}
	result2 := cosineSimilarity(a2, b2)

	if math.Abs(result2-1.0) > 0.0001 {
		t.Errorf("cosineSimilarity (identical) = %f; want 1.0", result2)
	}

	// Test 3: Perpendicular vectors should be 0.0
	a3 := []float64{1, 0}
	b3 := []float64{0, 1}
	result3 := cosineSimilarity(a3, b3)

	if math.Abs(result3-0.0) > 0.0001 {
		t.Errorf("cosineSimilarity (perpendicular) = %f; want 0.0", result3)
	}

	// Test 4: Opposite vectors should be -1.0
	a4 := []float64{1, 2}
	b4 := []float64{-1, -2}
	result4 := cosineSimilarity(a4, b4)

	if math.Abs(result4-(-1.0)) > 0.0001 {
		t.Errorf("cosineSimilarity (opposite) = %f; want -1.0", result4)
	}
}
