package sqldb

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkPrepareStatement(b *testing.B) {
	queryStatement := `SELECT * FROM users WHERE name = :name AND age > :age`

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = PrepareStatement(queryStatement)
	}
}

func ExamplePrepareStatement_benchmark() {
	benchmarkResult := testing.Benchmark(BenchmarkPrepareStatement)
	fmt.Println(
		strings.ReplaceAll(benchmarkResult.String(),
			"BenchmarkPrepareStatement",
			"PrepareStatement"),
	)
}
