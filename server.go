package gossm

import (
	"fmt"
)

type Servers []*Server

type Server struct {
	Name          string `json:"name"`
	IPAddress     string `json:"ipAddress"`
	Port          int    `json:"port"`
	Protocol      string `json:"protocol"`
	CheckInterval int    `json:"checkInterval"`
	Timeout       int    `json:"timeout"`
	// Struct met servernaam, ipadres, poort, welk protocol er gebruikt wordt, de check interval en de timeout interval
}

func (s *Server) String() string {
	return fmt.Sprintf("%s %s:%d", s.Protocol, s.IPAddress, s.Port)
}

func (s *Server) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}
