package log

import "testing"

func init() {
	Init("test", LogLevelDEBUG)
	InitAsync()
}

func BenchmarkLog(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Log("asdf")
		}
	})
}

func BenchmarkLogDebug(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Debug("asdf")
	}
}
