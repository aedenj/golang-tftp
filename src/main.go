package main

import (
  "fmt"
  "net"
  "log"
  "github.com/aedenj/golang-tftp/tftp"
)


func main() {
  s, err := net.ResolveUDPAddr("udp", ":3000")
  if err != nil {
    log.Fatal(err)
    return
  }

  conn, err := net.ListenUDP("udp", s)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer conn.Close()


  buffer := make([]byte, tftp.MaxPacketSize)
  //tftpChannel := make(chan tftp.PacketRequest)
  for {
    n, _, err := conn.ReadFromUDP(buffer)
    if err != nil || n <= 0 {
      log.Fatal(err)
      return
    }
    p, err := tftp.ParsePacket(buffer)
    if err != nil || n <= 0 {
      log.Fatal(err)
      return
    }

    go handle_write(p)
  }
}

func handle_write(req tftp.Packet) {
  fmt.Println(req)
}
