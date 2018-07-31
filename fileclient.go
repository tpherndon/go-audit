package main

import (
	"bufio"
	"bytes"
	"syscall"
)

type FileClient struct {
	s *bufio.Scanner
}

func NewFileClient(s *bufio.Scanner) (*FileClient, error) {
	f := &FileClient{
		s: s,
	}

	return f, nil
}

func (f *FileClient) Receive() (*syscall.NetlinkMessage, error) {
	f.s.Scan()
	if err := f.s.Err(); err != nil {
		return nil, err
	}

	// construct syscall.NetlinkMessage
	line := f.s.Bytes()
	words := bytes.Fields(line)
	var typ, seq, pid []byte
	for _, word := range words {
		kv := bytes.Split(word, []byte("="))
		key := string(kv[0])
		value := kv[1]
		switch key {
		case "type":
			typ = value
		case "pid":
			pid = value
		case "msg":
			parts := bytes.Split(value, []byte(":"))
			seq = bytes.TrimRight(parts[1], ")")
		}
	}

	msg := &syscall.NetlinkMessage{
		Header: syscall.NlMsghdr{
			Len:   Endianness.Uint32(line),
			Type:  Endianness.Uint16(typ),
			Flags: Endianness.Uint16([]byte("")),
			Seq:   Endianness.Uint32(seq),
			Pid:   Endianness.Uint32(pid),
		},
		Data: line,
	}
	return msg, nil
}
