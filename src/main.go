package main

import (
  "net"
  "log"
  "os"
  "os/signal"

  "github.com/aedenj/golang-tftp/tftp"
)

func main() {
  udpAddress, err := net.ResolveUDPAddr("udp", ":3000")
  if err != nil {
    log.Fatal(err)
  }

  conn, err := net.ListenUDP("udp", udpAddress)
  if err != nil {
    log.Fatal(err)
  }

  log.Println("Listening on %v\n", conn.LocalAddr())

  // handle ctrl-c
  go func() {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    for sig := range c {
      log.Println("Received %v, exiting", sig)
      os.Exit(0)
    }
  }()

  buffer := make([]byte, tftp.MaxPacketSize)
  for {
    n, caddr, err := conn.ReadFromUDP(buffer)
    if err != nil || n <= 0 {
      log.Fatal(err)
      return
    }

    req, err := tftp.ParsePacket(buffer)
    if err != nil {
      // We got a bad request packet. We're not sending the ack
      // so the client will try again. This is consistent with the
      // first paragraph of section 2 in RFC 1350
      log.Println("Bad request packet: %s", err)
      continue
    }

    go HandleRequest(conn, caddr, req)
  }
}

func HandleRequest(conn *net.UDPConn, caddr *net.UDPAddr, req tftp.Packet) {
  waddr, err := net.ResolveUDPAddr("udp", ":0")
  if err != nil {
    log.Println("Failed getting UDP address to write to: ", err)
    return
  }

  conn, err = net.DialUDP("udp", waddr, caddr)
  if err != nil {
    log.Println("Error dialing UDP to client: ", err)
    return
  }
  defer conn.Close()

  ack := &tftp.PacketData{BlockNum: 0}
  conn.Write(ack.Serialize())

  buf := make([]byte, tftp.MaxPacketSize)
  _, _, err = conn.ReadFromUDP(buf)
  if err != nil {
    log.Println("Failed to read data from udp connection")
    return
  }
}

