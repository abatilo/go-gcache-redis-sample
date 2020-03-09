package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup

	// This loop for simulating different channels being requested at the same time
	for z := 0; z < 120; z++ {
		wg.Add(1)
		go func(z int) {
			// This loop for simulating segments rolling over time
			var i int
			for i = 0; i < 20; i++ {
				ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
				defer cancel()
				path := fmt.Sprintf("http://localhost:8000/000%d_%d.mp4", z, i)
				// This loop for simulating the same segment getting requested repeatedly
				for j := 0; j < 10000; j++ {
					wg.Add(1)
					go func() {
						http.Get(path)
						wg.Done()
					}()
				}
				<-ctx.Done()
				fmt.Printf("Channel %d Segment %d done\n", z, i)
			}
			wg.Done()
		}(z)
	}

	wg.Wait()
}
