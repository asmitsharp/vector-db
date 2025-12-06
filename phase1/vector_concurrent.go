package phase1

import (
	"sort"
	"sync"
)

// SearchConcurrent uses goroutine to parallelize the search result
func (db *VectorDB) SearchConcurrent(query []float64, k int) ([]SearchResult, error) {
	numWorkers := 8
	totalVectors := len(db.vectors)

	chunkSize := totalVectors / numWorkers

	resultChan := make(chan []SearchResult, numWorkers)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		workerId := i

		go func() {
			defer wg.Done()

			// Partioning
			start := workerId * chunkSize
			end := start + chunkSize

			if workerId == numWorkers-1 {
				end = totalVectors
			}

			var localResults []SearchResult

			//Process the chunk
			for j := start; j < end; j++ {
				vec := db.vectors[j]
				queryNorm := magnitude(query)
				dist := cosineSimilarityOptimized(query, vec.Data, vec.norm, queryNorm)

				localResults = append(localResults, SearchResult{
					ID:         vec.ID,
					Similarity: dist,
					Metadata:   vec.MetaData,
				})
			}

			// Instead of sending all items, we sort and send only Top K (Like a local Reduce)
			sort.Slice(localResults, func(i, j int) bool {
				return localResults[i].Similarity > localResults[j].Similarity
			})

			if len(localResults) > k {
				localResults = localResults[:k]
			}

			// send the small packet to the main channel
			resultChan <- localResults
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var globalResults []SearchResult

	for packet := range resultChan {
		globalResults = append(globalResults, packet...)
	}

	sort.Slice(globalResults, func(i, j int) bool {
		return globalResults[i].Similarity > globalResults[j].Similarity
	})

	if len(globalResults) > k {
		return globalResults[:k], nil
	}

	return globalResults, nil
}
