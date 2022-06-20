package track

import (
	"time"
)

// Struct word gebruikt om tijd te volgen
type TimeTracker struct {
	delayer  Delayer
	nextTime time.Time
	// Counter word gebruikt om te onthouden hoevaak een delayer heeft gedraaid
	counter int
}

// Isready kijk of de huidige tijd gelijk is na nextime
// Bij eerste opstart zet hij nettime op 0
func (t *TimeTracker) IsReady() bool {
	if time.Now().After(t.nextTime) {
		return true
	}
	return false
}

// SetNext update nextTime gebaseerd op delay implementatie
// geeft delay en time terug welke later veranderd in huidige tijd
func (t *TimeTracker) SetNext() (time.Duration, time.Time) {
	t.counter++
	nextDelay := t.delayer.Delay()
	t.nextTime = time.Now().Add(nextDelay)
	return nextDelay, t.nextTime
}

// NewTrucker geeft pointer terug met nieuwe TimeTracker en zet delay
func NewTracker(delayer Delayer) *TimeTracker {
	return &TimeTracker{delayer: delayer}
}

// HasBeenRan kijk hoevaak de tijd delay heeft gedraaid
func (t *TimeTracker) HasBeenRan() bool {
	if t.counter > 0 {
		return true
	}
	return false
}
