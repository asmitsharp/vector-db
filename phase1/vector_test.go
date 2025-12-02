package phase1

import (
	"fmt"
	"math"
	"math/rand"
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

func TestVectorDB(t *testing.T) {
	// Create new database
	db := NewVectorDB()

	// Test 1: Insert vectors
	err := db.Insert("cat", []float64{0.8, 0.9, 0.1}, map[string]string{"type": "animal"})
	if err != nil {
		t.Fatalf("Failed to insert cat: %v", err)
	}

	err = db.Insert("dog", []float64{0.7, 0.85, 0.15}, map[string]string{"type": "animal"})
	if err != nil {
		t.Fatalf("Failed to insert dog: %v", err)
	}

	err = db.Insert("car", []float64{0.1, 0.2, 0.9}, map[string]string{"type": "vehicle"})
	if err != nil {
		t.Fatalf("Failed to insert car: %v", err)
	}

	// Test 2: Check size
	if db.Size() != 3 {
		t.Errorf("Expected size 3, got %d", db.Size())
	}

	// Test 3: Search for something similar to cat
	query := []float64{0.78, 0.87, 0.12}
	results, err := db.Search(query, 2)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	// Should find cat first (most similar)
	if results[0].ID != "cat" {
		t.Errorf("Expected 'cat' as top result, got '%s'", results[0].ID)
	}

	// Should find dog second
	if results[1].ID != "dog" {
		t.Errorf("Expected 'dog' as second result, got '%s'", results[1].ID)
	}

	t.Logf("Search results:")
	for i, r := range results {
		t.Logf("  %d. %s (similarity: %.4f)", i+1, r.ID, r.Similarity)
	}
}

func TestVectorDBErrors(t *testing.T) {
	db := NewVectorDB()

	// Test: Empty vector should fail
	err := db.Insert("test", []float64{}, nil)
	if err == nil {
		t.Error("Expected error for empty vector")
	}

	// Test: Duplicate ID should fail
	db.Insert("dup", []float64{1, 2, 3}, nil)
	err = db.Insert("dup", []float64{4, 5, 6}, nil)
	if err == nil {
		t.Error("Expected error for duplicate ID")
	}

	// Test: Dimension mismatch should fail
	db.Insert("first", []float64{1, 2, 3}, nil)
	err = db.Insert("second", []float64{1, 2}, nil) // Wrong dimension!
	if err == nil {
		t.Error("Expected error for dimension mismatch")
	}
}

func BenchmarkSearch(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		db := createTestDB(size, 128)
		query := randomVector(128)

		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				db.Search(query, 10)
			}
		})
	}
}

// Helper: create test database with n random vectors
func createTestDB(n int, dim int) *VectorDB {
	db := NewVectorDB()
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("vec_%d", i)
		data := randomVector(dim)
		db.Insert(id, data, nil)
	}
	return db
}

// Helper: generate random vector
func randomVector(dim int) []float64 {
	v := make([]float64, dim)
	for i := range v {
		v[i] = rand.Float64()
	}
	return v
}
