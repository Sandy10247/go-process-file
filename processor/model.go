package processor

import "sync"

// Core Structs
type Payload struct {
	Data []byte
}

type Chunk struct {
	payloads []*Payload
}

// #region Payload

// declare a package level Sync for "Payload"
var payloadPool = sync.Pool{
	New: func() any {
		return &Payload{}
	},
}

// take a "Payload" item from the pool
func getPayload() *Payload {
	if newPayload, ok := payloadPool.Get().(*Payload); ok {
		return newPayload
	}
	panic("unable to get Payload from pool")
}

// give back to the pool
func releasePayload(payload *Payload) {
	// clean up data
	payload.Data = payload.Data[:0]
	// give back the "Payload" took from "pool"
	payloadPool.Put(payload)
}

//	#endregion

// #region Chunk

// declare a package level Sync for "Chunk"
var chunkPool = sync.Pool{
	New: func() any {
		return &Chunk{
			payloads: make([]*Payload, 0),
		}
	},
}

// take a "Chunk" item from the pool
func getChunk() *Chunk {
	if newChunk, ok := chunkPool.Get().(*Chunk); ok {
		return newChunk
	}
	panic("unable to get Chunk from pool")
}

// give back to the pool
func releaseChunk(chunk *Chunk) {
	chunk.payloads = chunk.payloads[:0]
	chunkPool.Put(chunk)
}

// #endregion
