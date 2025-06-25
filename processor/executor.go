package processor

import (
	"bufio"
	"io"
	"log"
	"sync"
)

// Process "Read + ProcessFn"
func ProcessFile(file io.Reader, numWorkers int, chunkSize int, processFn func(data []byte)) {
	scanner := bufio.NewScanner(file)
	chunkChan := make(chan *Chunk, numWorkers)
	var wg sync.WaitGroup

	// spin up a goroutine for each workeer
	for i := range numWorkers {
		wg.Add(1)
		// read a Chunk from "chunkChan"
		go func(workerId int) {
			defer wg.Done()
			// Process "N" payload from each "Chunk"
			for chunk := range chunkChan {
				for _, payload := range chunk.payloads {
					// call the Process Func
					processFn(payload.Data)
					releasePayload(payload)
				}
				releaseChunk(chunk)

			}
		}(i)
	}

	// Read data from the File
	// Put it in Channel as "Chunk"
	go func() {
		defer close(chunkChan)
		for {
			chunk := getChunk()
			for len(chunk.payloads) < chunkSize && scanner.Scan() {
				payload := getPayload()
				// Add to Payload Data
				payload.Data = append(payload.Data[:0], scanner.Bytes()...)
				// Add new Payload to Chunk
				chunk.payloads = append(chunk.payloads, payload)
			}

			// when atleast one Payload is added, put it in channel
			if len(chunk.payloads) > 0 {
				chunkChan <- chunk
			} else {
				// if we couldn't read the chunk, release the chunk we took from the pool
				releaseChunk(chunk)
				// Stop the loop
				break
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Scanner error: %v", err)
		}
	}()

	wg.Wait()
}
