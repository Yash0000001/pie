package concurrency

import (
	"fmt"
	"sync"

	"github.com/yash0000001/pie/internal/counter"
)

func ProcessPaths(rootArr []string) int {
	var wg sync.WaitGroup
	results := make(chan int)

	for _, val := range rootArr {
		wg.Add(1) //number of go routines spawn
		go func(path string) {
			defer wg.Done()
			total, err := counter.ProcessPath(path)
			if err != nil {
				fmt.Printf("❌ Error on %s: %v\n", path, err)
				results <- 0
				return
			}
			results <- total
		}(val)
	}

	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()
	grandTotal := 0
	for total := range results {
		grandTotal += total
	}

	return grandTotal
}
