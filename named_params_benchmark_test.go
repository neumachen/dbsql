package sqlstmt

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkConvertNamedToPositionalParams(b *testing.B) {
	queryStatement := []byte(
		"SELECT * FROM users WHERE name = :name AND age > :age",
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ConvertNamedToPositionalParams(queryStatement)
	}
}

func ExampleConvertNamedToPositionalParams_benchmark() {
	benchmarkResult := testing.Benchmark(BenchmarkConvertNamedToPositionalParams)
	fmt.Println(
		strings.ReplaceAll(benchmarkResult.String(),
			"BenchmarkConvertNamedToPositionalParams",
			"ConvertNamedToPositionalParams"),
	)
}
