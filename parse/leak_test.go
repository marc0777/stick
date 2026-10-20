package parse

import (
	"runtime"
	"testing"
	"time"
)

// A parse error must not leave the lexer goroutine blocked on its next token.
func TestParseErrorDoesNotLeakLexer(t *testing.T) {
	before := runtime.NumGoroutine()
	for i := 0; i < 100; i++ {
		if _, err := Parse("{% endif %} and more {{ text }} {{ text }}"); err == nil {
			t.Fatal("expected a parse error")
		}
	}
	for deadline := time.Now().Add(time.Second); runtime.NumGoroutine() > before+5; {
		if time.Now().After(deadline) {
			t.Fatalf("goroutines: before %d, after %d", before, runtime.NumGoroutine())
		}
		time.Sleep(10 * time.Millisecond)
	}
}
