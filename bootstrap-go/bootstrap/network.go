package bootstrap

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func sendMessage(conn net.Conn, msg Message, encrypt EncryptFunc, sign SignFunc) error {
	serialized, err := serializeMessage(msg, encrypt, sign)
	if err != nil {
		return err
	}
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(serialized)))
	_, err = conn.Write(lenBuf)
	if err != nil {
		return err
	}
	_, err = conn.Write(serialized)
	if err != nil {
		return err
	}
	return nil
}

func receiveMessage(conn net.Conn, decrypt DecryptFunc) (Message, error) {
	lenBuf := make([]byte, 4)
	lenRead, err := io.ReadFull(conn, lenBuf)
	if err != nil {
		return nil, fmt.Errorf("read message error: %v", err)
	}
	if lenRead != len(lenBuf) {
		return nil, fmt.Errorf("read %d bytes instead of %d", lenRead, len(lenBuf))
	}
	msgLen := binary.BigEndian.Uint32(lenBuf)
	msg := make([]byte, msgLen)
	lenRead, err = io.ReadFull(conn, msg)
	if err != nil {
		return nil, fmt.Errorf("read message error: %v", err)
	}
	if uint32(lenRead) != msgLen {
		return nil, fmt.Errorf("read %d bytes instead of %d", lenRead, msgLen)
	}
	return deserializeMessage(msg, decrypt)
}
