package test

import (
	"testing"
)

func BenchmarkPayloadSizeZero(b *testing.B) {

	client, payload, results, _ := prepareRequest(1000000, 0)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sendRequest("http://localhost:8888/", payload, uint(i), client, results)
	}
}

func BenchmarkPayloadSizeTen(b *testing.B) {

	client, payload, results, _ := prepareRequest(1000000, 10)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sendRequest("http://localhost:8888/", payload, uint(i), client, results)
	}
}

func BenchmarkPayloadSizeHundred(b *testing.B) {

	client, payload, results, _ := prepareRequest(1000000, 100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sendRequest("http://localhost:8888/", payload, uint(i), client, results)
	}
}
