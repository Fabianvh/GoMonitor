package dial

import (
	"time"
)

// Struct word gebruikt voor het transporteren van informatie richting de dialer
type NetAddress struct {
	Network string
	Address string
}

// NetAddressTimeout is een NetAdress voor Timeouts
type NetAddressTimeout struct {
	NetAddress
	Timeout time.Duration
}
