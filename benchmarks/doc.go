// Package benchmarks contains an optional gateway load-benchmark suite.
//
// The suite is gated behind the `benchmark` build tag so normal
// `go test ./...` / CI runs stay fast. Execute it with:
//
//	make benchmark
package benchmarks
