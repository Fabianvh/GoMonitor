package gossm

import (
	"fmt"
	"os"
	"time"

	"github.com/ssimunic/gossm/dial"
	"github.com/ssimunic/gossm/logger"
	"github.com/ssimunic/gossm/notify"
	"github.com/ssimunic/gossm/track"
)

type Monitor struct {
	// Laad instellingen en servers
	config *Config

	// Checker voor checks de servers
	checkerCh chan *Server

	// Ongebruikte functie voor het versturen van notifcaties, Notifcatie uit de code gebouwd. enkel krijg ik deze niet uit de code zonder andere functies onwerkbaar te krijgen.
	notifiers notify.Notifiers

	// Notifier voor het aangeven van server die niet beschikbaar zijn
	notifierCh chan *Server

	// Ook voor de notifier functie, funcitie word gebruikt om spam te voorkomen. Geeft delay op de notificer.
	notificationTracker map[*Server]*track.TimeTracker

	// Gebruikt voor het testen van de connecties richting de gecontroleerde servers
	dialer *dial.Dialer

	// Gebruikt dit channel voor het stoppen van de code
	stop chan struct{}
}

func NewMonitor(c *Config) *Monitor {
	m := &Monitor{
		config:              c,
		checkerCh:           make(chan *Server),
		notifiers:           c.Settings.Notifications.GetNotifiers(),
		notifierCh:          make(chan *Server),
		notificationTracker: make(map[*Server]*track.TimeTracker),
		dialer:              dial.NewDialer(c.Settings.Monitor.MaxConnections),
		stop:                make(chan struct{}),
	}
	m.initialize()
	return m
}

func (m *Monitor) initialize() {
	// Initaliseert notifcatie methodes, ongebruikt
	for _, notifier := range m.notifiers {
		if initializer, ok := notifier.(notify.Initializer); ok {
			logger.Logln("Initializing", initializer)
			initializer.Initialize()
		}
	}

	for _, server := range m.config.Servers {
		// Initaliseert notifactietracker, ongebruikt
		m.notificationTracker[server] = NewTrackerWithExpBackoff(m.config.Settings.Monitor.ExponentialBackoffSeconds)

		// Set standaard checkinterval and timeoutinterval voor de servers
		switch {
		case server.CheckInterval <= 0:
			server.CheckInterval = m.config.Settings.Monitor.CheckInterval
		case server.Timeout <= 0:
			server.Timeout = m.config.Settings.Monitor.Timeout
		}
	}
}

// Functie creeert een tracker met tijd als delayer
func NewTrackerWithExpBackoff(expBackoffSeconds int) *track.TimeTracker {
	return track.NewTracker(track.NewExpBackoff(expBackoffSeconds))
}

// Draait de monitoring voor altijd zonder onderbrekingen
func (m *Monitor) Run() {
	m.RunForSeconds(0)
}

// Draait de code met runforseconds of voor altijd als 0 is gegeven als argument
func (m *Monitor) RunForSeconds(runningSeconds int) {
	if runningSeconds != 0 {
		go func() {
			runningSecondsTime := time.Duration(runningSeconds) * time.Second
			<-time.After(runningSecondsTime)
			m.stop <- struct{}{}
		}()
	}

	for _, server := range m.config.Servers {
		go m.scheduleServer(server)
	}

	logger.Logln("Starting monitor.")
	m.monitor()
}

func (m *Monitor) scheduleServer(s *Server) {
	m.checkerCh <- s

	// Geeft periode aan
	tickerSeconds := time.NewTicker(time.Duration(s.CheckInterval) * time.Second)
	for range tickerSeconds.C {
		m.checkerCh <- s
	}
}

func (m *Monitor) monitor() {
	go m.listenForChecks()
	go m.listenForNotifications()

	// Wacht voor foutmelding en sluit daarna code af
	<-m.stop
	logger.Logln("Terminating.")
	os.Exit(0)
}

func (m *Monitor) listenForChecks() {
	for server := range m.checkerCh {
		m.checkServerStatus(server)
	}
}

func (m *Monitor) listenForNotifications() {
	for server := range m.notifierCh {
		timeTracker := m.notificationTracker[server]
		if timeTracker.IsReady() {
			nextDelay, nextTime := timeTracker.SetNext()
			logger.Logln("Sending notifications for", server)
			go m.notifiers.NotifyAll(fmt.Sprintf("%s (%s)", server.Name, server))
			logger.Logln("Next available notification for", server.String(), "in", nextDelay, "at", nextTime)
		}
	}
}

func (m *Monitor) checkServerStatus(server *Server) {
	// Functie checkt voor vrije plekken in de dail functie
	worker, output := m.dialer.NewWorker()
	go func() {
		logger.Logln("Checking", server)

		formattedAddress := fmt.Sprintf("%s:%d", server.IPAddress, server.Port)
		timeoutSeconds := time.Duration(server.Timeout) * time.Second
		worker <- dial.NetAddressTimeout{NetAddress: dial.NetAddress{Network: server.Protocol, Address: formattedAddress}, Timeout: timeoutSeconds}
		dialerStatus := <-output

		// Error handeling
		if !dialerStatus.Ok {
			logger.Logln(dialerStatus.Err)
			logger.Logln("ERROR", server)
			go func() {
				m.notifierCh <- server
			}()
			return
		}

		// Logger start
		logger.Logln("OK", server)
		// Reset tijd voor tracking van de server zelf
		if m.notificationTracker[server].HasBeenRan() {
			m.notificationTracker[server] = NewTrackerWithExpBackoff(m.config.Settings.Monitor.ExponentialBackoffSeconds)
		}
	}()
}
