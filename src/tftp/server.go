package tftp

import (
  "net"
  "os"
  "os/signal"

  "log"
)

type Server struct {
  Conn *net.UDPConn
}

func NewServer() (*Server) {
  return &Server{}
}

func (s *Server) Listen() (error) {
  udpAddress, err := net.ResolveUDPAddr("udp", ":3000")
  if err != nil {
    return err
  }

  conn, err := net.ListenUDP("udp", udpAddress)
  if err != nil {
    return err
  }

  s.Conn = conn
  return nil
}

func (s *Server) SetupInterruptHandler() {
  // handle ctrl-c
  go func() {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    for sig := range c {
      log.Println("Received", sig, ", exiting")
      os.Exit(0)
    }
  }()
}

