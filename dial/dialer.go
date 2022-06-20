package dial

import (
	"net"
)

//Dailer controleerd connecties
type Dialer struct {
	semaphore chan struct{}
}

// Struct slaat informatie over de connectie
type Status struct {
	Ok  bool
	Err error
}

// NewDialer geeft pointer naar nieuwe dailer
func NewDialer(concurrentConnections int) *Dialer {
	return &Dialer{
		semaphore: make(chan struct{}, concurrentConnections),
	}
}

// Newworker is gebruikt voor het versturen van een address naar NetAddressTimeout voor het maken en ontvangen van de Dialerstatussen
func (d *Dialer) NewWorker() (chan<- NetAddressTimeout, <-chan Status) {
	netAddressTimeoutCh := make(chan NetAddressTimeout)
	dialerStatusCh := make(chan Status)

	d.semaphore <- struct{}{}
	go func() {
		netAddressTimeout := <-netAddressTimeoutCh
		conn, err := net.DialTimeout(netAddressTimeout.Network, netAddressTimeout.Address, netAddressTimeout.Timeout)

		dialerStatus := Status{}

		if err != nil {
			dialerStatus.Ok = false
			dialerStatus.Err = err
		} else {
			dialerStatus.Ok = true
			conn.Close()
		}
		dialerStatusCh <- dialerStatus
		<-d.semaphore
	}()

	return netAddressTimeoutCh, dialerStatusCh
}
