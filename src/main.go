package main

import (
  "net"
  "log"
  "reflect"
  "github.com/aedenj/golang-tftp/tftp"
)

func main() {
  srv := tftp.NewServer()

  srv.Listen()
  log.Println("Listening on ", srv.Conn.LocalAddr())
  defer srv.Conn.Close()

  srv.SetupInterruptHandler()

  buffer := make([]byte, tftp.MaxPacketSize)
  for {
    n, caddr, err := srv.Conn.ReadFromUDP(buffer)
    if err != nil || n <= 0 {
      log.Fatal(err)
      return
    }

    req, err := tftp.ParsePacket(buffer)
    if err != nil {
      // We got a bad request packet. We're not sending the ack
      // so the client will try again. This is consistent with the
      // first paragraph of section 2 in RFC 1350
      log.Println("Bad request packet: ", err)
      continue
    }

    go HandleRequest(caddr, req)
  }
}

func HandleRequest(caddr *net.UDPAddr, req tftp.Packet) {
  reqPkt, ok := req.(*tftp.PacketRequest)
  if !ok {
    log.Printf("Invalid packet type for new connection!")
    return
  }

  switch reqPkt.Op {
    case tftp.OpWRQ:
      err := HandleWriteRequest(caddr)
      if err != nil {
        log.Println("write request finished, with error:")
        log.Println(err)
      }
    default:
      log.Println("Invalid Packet Type.")
  }

}

func HandleWriteRequest(caddr *net.UDPAddr) (error){
  log.Println("HANDLE WRITE")

  waddr, err := net.ResolveUDPAddr("udp", ":0")
  if err != nil {
    log.Println("Failed getting UDP address to write to: ", err)
    return err
  }

  conn, err := net.DialUDP("udp", waddr, caddr)
  if err != nil {
    log.Println("Error dialing UDP to client: ", err)
    return err
  }
  defer conn.Close()

  ack := &tftp.PacketData{}
  conn.Write(ack.Serialize())
  buffer := make([]byte, tftp.MaxPacketSize)
  for {
    _, _, err = conn.ReadFromUDP(buffer)
    if err != nil {
      log.Println("Failed to read data from udp connection")
      return err
    }
    pkt, _ := tftp.ParsePacket(buffer)

    log.Println(reflect.TypeOf(pkt).String())
    conn.Write(ack.Serialize())
  }

  return nil
}
