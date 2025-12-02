package main

import (
	"fmt"
	"log"

	"github.com/ashmitsharp/vector-db/phase1"
)

func main() {
	fmt.Println("=== Simple Vector Database Demo ===")

	// Create database
	db := phase1.NewVectorDB()

	// Insert some sample vectors (imagine these are text embeddings)
	vectors := map[string]struct {
		data     []float64
		metadata map[string]string
	}{
		"cat": {
			data:     []float64{0.8, 0.9, 0.1, 0.2},
			metadata: map[string]string{"type": "animal", "category": "pet"},
		},
		"dog": {
			data:     []float64{0.75, 0.85, 0.15, 0.25},
			metadata: map[string]string{"type": "animal", "category": "pet"},
		},
		"lion": {
			data:     []float64{0.7, 0.8, 0.2, 0.3},
			metadata: map[string]string{"type": "animal", "category": "wild"},
		},
		"car": {
			data:     []float64{0.1, 0.2, 0.9, 0.8},
			metadata: map[string]string{"type": "vehicle"},
		},
		"truck": {
			data:     []float64{0.15, 0.25, 0.85, 0.75},
			metadata: map[string]string{"type": "vehicle"},
		},
	}

	// Insert all vectors
	for id, v := range vectors {
		if err := db.Insert(id, v.data, v.metadata); err != nil {
			log.Fatalf("Failed to insert %s: %v", id, err)
		}
		fmt.Printf("✓ Inserted: %s\n", id)
	}

	fmt.Printf("\nDatabase size: %d vectors\n\n", db.Size())

	// Search 1: Query similar to "cat"
	fmt.Println("Query 1: Similar to 'cat' [0.78, 0.88, 0.12, 0.22]")
	query1 := []float64{0.78, 0.88, 0.12, 0.22}
	results1, _ := db.Search(query1, 3)
	printResults(results1)

	// Search 2: Query similar to "car"
	fmt.Println("\nQuery 2: Similar to 'car' [0.12, 0.22, 0.88, 0.78]")
	query2 := []float64{0.12, 0.22, 0.88, 0.78}
	results2, _ := db.Search(query2, 3)
	printResults(results2)
}

func printResults(results []phase1.SearchResult) {
	for i, r := range results {
		fmt.Printf("  %d. %-10s (similarity: %.4f) %v\n",
			i+1, r.ID, r.Similarity, r.Metadata)
	}
}
