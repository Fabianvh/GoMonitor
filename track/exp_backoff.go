package track

import (
	"time"
)

// Functie is gebruikt als delayer, implementeer een backoff interval
type ExpBackoff struct {
	counter int
	base    int
}

func calculateExponential(base, counter int) int {
	if counter == 0 {
		return 1
	}
	return base * calculateExponential(base, counter-1)
}

// Functie geeft seconds terug
func (e *ExpBackoff) Delay() time.Duration {
	e.counter++
	return time.Duration(calculateExponential(e.base, e.counter)) * time.Second
}

// Functie geeft pointer terug naar ExpBackoff
func NewExpBackoff(base int) *ExpBackoff {
	return &ExpBackoff{
		base: base,
	}
}
