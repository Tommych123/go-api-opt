package sum

import (
	"encoding/json"
	"testing"
	"bytes"
)

func benchPayload(tb testing.TB, count, repeat int) []byte {
	tb.Helper()

	nums := make([]int64, count)
	for i := range nums {
		nums[i] = int64(i % 1000)
	}

	// Генерируем JSON максимально похожий на реальный запрос
	b, err := json.Marshal(map[string]any{
		"numbers": nums,
		"repeat":  repeat,
	})
	if err != nil {
		tb.Fatalf("marshal payload: %v", err)
	}
	return b
}

// BenchmarkParseAndSum фиксирует baseline производительность суммирования
// Он будем сравнивать до и после оптимизации
func BenchmarkParseAndSum(b *testing.B) {
	payload := benchPayload(b, 1000, 200)

	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))

	for i := 0; i < b.N; i++ {
		if _, err := ParseAndSum(payload); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeAndSum(b *testing.B) {
	payload := benchPayload(b, 1000, 200)

	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))

	for i := 0; i < b.N; i++ {
		if _, err := DecodeAndSum(bytes.NewReader(payload)); err != nil {
			b.Fatal(err)
		}
	}
}