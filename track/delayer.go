package track

import (
	"time"
)

// Delay functie voor het vertragen van een functie/actie
type Delayer interface {
	Delay() time.Duration
}
