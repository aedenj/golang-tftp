package main

import (
  "fmt"
  "net"
  "log"
  "math/rand"
  "time"
  "strings"
  "strconv"
  "github.com/aedenj/golang-tftp/tftp"
)

func random(min, max int) int {
  return rand.Intn(max-min) + min
}

func main() {
  s, err := net.ResolveUDPAddr("udp", "0.0.0.0:3000")
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
  rand.Seed(time.Now().Unix())

  for {
    n, addr, err := conn.ReadFromUDP(buffer)
    fmt.Print("-> ", string(buffer[0:n-1]))

    if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
      fmt.Println("Exiting UDP server!")
      return
    }

    data := []byte(strconv.Itoa(random(1, 1001)))
    fmt.Printf("data: %s\n", string(data))
    _, err = conn.WriteToUDP(data, addr)
    if err != nil {
      fmt.Println(err)
      return
    }
  }
}
