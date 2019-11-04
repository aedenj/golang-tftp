package main

import (
  "fmt"
  "net"
  "log"
  "github.com/aedenj/golang-tftp/tftp"
)


func main() {
  uaddr, err := net.ResolveUDPAddr("udp", ":69")
  if err != nil {
    log.Fatal(err)
    return
  }

  conn, err := net.ListenUDP("udp", uaddr)
  if err != nil {
    log.Fatal(err)
    return
  }
  defer conn.Close()


  buffer := make([]byte, tftp.MaxPacketSize)
  for {
    n, _, err := conn.ReadFromUDP(buffer)
    if err != nil || n <= 0 {
      log.Fatal(err)
      return
    }

    req, err := tftp.ParsePacket(buffer)
    if err != nil || n <= 0 {
      // We got a bad request packet. We're not sending the ack
      // so the client will try again. This is consistent with the
      // first paragraph of section 2 in RFC 1350
      log.Printf("Bad request packet: %s", err)
      continue
    }

    go HandleRequest(uaddr,req)
  }

}

func HandleRequest(addr *net.UDPAddr, req tftp.Packet) {
  fmt.Println(req)
  reqpkt, ok := req.(*tftp.PacketRequest)
  if !ok {
    log.Printf("Invalid packet type for new connection!")
    return
  }

  //clientaddr, err := net.ResolveUDPAddr("udp", addr.String())
  //if err != nil {
    //log.Printf("Error: %s", err)
    //return
  //}

  switch reqpkt.Op {
    case tftp.OpWRQ:
      err := HandleWriteRequest(reqpkt, addr)
      if err != nil {
        log.Println("write request finished, with error:")
        log.Println(err)
      }
    default:
      log.Println("Invalid Packet Type!")
  }
  fmt.Println("Done Handling")
}

func HandleWriteRequest(wrq *tftp.PacketRequest, addr *net.UDPAddr) error {
  fmt.Println("Handle Write")
  listaddr, err := net.ResolveUDPAddr("udp", ":0")
  if err != nil {
    return err
  }

  conn, err := net.ListenUDP("udp", listaddr)
  if err != nil {
    log.Fatal(err)
    return err
  }
  defer conn.Close()
  // Connection directly to their open port
  //conn, err := net.DialUDP("udp", listaddr, addr)
  //if err != nil {
    //return err
  //}

  fmt.Println("WRITE ACK")
  ackPkt := &tftp.PacketData{BlockNum:uint16(0)}
  _, err = conn.WriteToUDP(ackPkt.Serialize(), addr)
  if err != nil {
    return err
  }

  buffer := make([]byte, tftp.MaxPacketSize)
  for {
    fmt.Println("READ DATA")
    n, _, err := conn.ReadFromUDP(buffer)
    if err != nil || n <= 0 {
      fmt.Println("ABOUT TO ERR")
      return err
    }

      fmt.Println("ABOUT TO PARSE")
    data, err := tftp.ParsePacket(buffer)
    if err != nil || n <= 0 {
      // We got a bad request packet. We're not sending the ack
      // so the client will try again. This is consistent with the
      // first paragraph of section 2 in RFC 1350
      log.Printf("Bad request packet: %s", err)
      continue
    }

    fmt.Println("ACTUAL DATA")
    fmt.Println(data)
  }


  if err := conn.Close(); err != nil {
    log.Print(err)
  }

  return nil
}
