package god

import (
	"testing"
	"time"
)

func BenchmarkTrace(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		func() {
			e := Trace("foo", "bar")
			defer e.End()
			time.Sleep(time.Microsecond)
		}()
	}
}
