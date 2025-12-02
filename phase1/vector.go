package phase1

import (
	"fmt"
	"math"
	"sort"
)

// vector represents a single vector with its metadata
type Vector struct {
	ID       string
	Data     []float64
	MetaData map[string]string
	norm     float64 // cache magnitude of vector for performance
}

type SearchResult struct {
	ID         string
	Similarity float64
	Metadata   map[string]string
	Vector     []float64
}

func NewVector(id string, data []float64, metadata map[string]string) *Vector {
	return &Vector{
		ID:       id,
		Data:     data,
		MetaData: metadata,
		norm:     magnitude(data),
	}
}

// in-memmory store for slice of vector pointers
type VectorDB struct {
	vectors []*Vector
}

// creates empty vector database
func NewVectorDB() *VectorDB {
	return &VectorDB{
		vectors: make([]*Vector, 0),
	}
}

// Size returns the number of vector stored
func (db *VectorDB) Size() int {
	return len(db.vectors)
}

// Insert add a new vector to the database
func (db *VectorDB) Insert(id string, data []float64, metadata map[string]string) error {
	if len(data) == 0 {
		return fmt.Errorf("cannot create vector : data slice is empty")
	}

	for i := range db.vectors {
		if db.vectors[i].ID == id {
			return fmt.Errorf("vector with id %s already exists", id)
		}
	}

	if len(db.vectors) > 0 {
		expectedDim := len(db.vectors[0].Data)
		if expectedDim != len(data) {
			return fmt.Errorf("dimension mismatch: got %d, expected %d", len(data), expectedDim)
		}
	}

	if metadata == nil {
		metadata = make(map[string]string)
	}

	v := NewVector(id, data, metadata)
	db.vectors = append(db.vectors, v)

	return nil
}

// Search finds the k most similar vectors to the query
func (db *VectorDB) Search(query []float64, k int) ([]SearchResult, error) {
	if len(query) == 0 {
		return nil, fmt.Errorf("Your query is empty.")
	}

	if len(db.vectors) == 0 {
		return []SearchResult{}, nil
	}

	expectedDim := len(db.vectors[0].Data)
	if expectedDim != len(query) {
		return nil, fmt.Errorf("query dimension mismatch: got %d, expected %d", len(query), expectedDim)
	}

	queryNorm := magnitude(query)

	result := make([]SearchResult, len(db.vectors))

	for i, v := range db.vectors {
		similarity := cosineSimilarityOptimized(query, v.Data, queryNorm, v.norm)

		result[i] = SearchResult{
			ID:         v.ID,
			Similarity: similarity,
			Metadata:   v.MetaData,
			Vector:     v.Data,
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Similarity > result[j].Similarity
	})

	if k > len(result) {
		k = len(result)
	}

	return result[:k], nil
}

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

func cosineSimilarityOptimized(a, b []float64, norm, queryNorm float64) float64 {
	dotP := dotProduct(a, b)
	if norm == 0 || queryNorm == 0 {
		return 0.0
	}

	return dotP / (norm * queryNorm)
}
