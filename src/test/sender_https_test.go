package test

import (
	"testing"
)

func BenchmarkHTTPSPayloadSizeZero(b *testing.B) {

	client, payload, results, _ := prepareHTTPSRequest(1000000, 0)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sendHTTPSRequest("https://localhost:9999/", payload, uint(i), client, results)
	}
}

func BenchmarkHTTPSPayloadSizeTen(b *testing.B) {

	client, payload, results, _ := prepareHTTPSRequest(1000000, 10)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sendHTTPSRequest("https://localhost:9999/", payload, uint(i), client, results)
	}
}

func BenchmarkHTTPSPayloadSizeHundred(b *testing.B) {

	client, payload, results, _ := prepareHTTPSRequest(1000000, 100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sendHTTPSRequest("https://localhost:9999/", payload, uint(i), client, results)
	}
}
