package main

import (
  "io"
)


func NewServer(readHandler func(filename string, rf io.ReaderFrom) error,
  writeHandler func(filename string, wt io.WriterTo) error) *Server {
  s := &Server{
    readHandler:       readHandler,
    writeHandler:      writeHandler,
  }

  return s
}

type Server struct {
  readHandler  func(filename string, rf io.ReaderFrom) error
  writeHandler func(filename string, wt io.WriterTo) error
}
