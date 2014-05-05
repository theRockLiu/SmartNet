// SmartNet project main.go
package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

func start_tcp_service(strLsnAddr string /*":9999"*/) {
	// listen on a port
	lsner, err := net.Listen("tcp", strLsnAddr)
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		// accept a connection
		conn, err := lsner.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}

		// handle the connection
		go handle_conn(conn)
	}
}

type PkgHdr struct {
	ui32PkgLen uint32
	ui16Opcode uint16
	ui16Others uint16
}

type RegPkg struct {
	PkgHdr
	bytesName [10]byte
	bytesPwd  [10]byte
}

const (
	CONST_READ_BUF_LEN = 1024 * 64
	CONST_PKG_HDR_LEN  = 8
	CONST_MAX_PKG_LEN  = 1024
)

func handle_pkg(bs []byte) (iDone int, err error) {

	iDone = len(bs)
	return iDone, err
}

func handle_conn(conn net.Conn) {
	//
	defer conn.Close()
	//
	var arrReadBuf [CONST_READ_BUF_LEN]byte
	i32Cnt, i32DataHead, i32FreeHead := int(0), int(0), int(0)
	var err error
	//
	for {
		if (CONST_READ_BUF_LEN - i32FreeHead) < CONST_MAX_PKG_LEN {
			copy(arrReadBuf[:], arrReadBuf[i32DataHead:i32FreeHead])
			i32FreeHead -= i32DataHead
			i32DataHead = 0
		}

		i32Cnt, err = conn.Read(arrReadBuf[i32FreeHead:])
		if err != nil {
			log.Fatal(err)
			return
		}
		i32FreeHead += i32Cnt
		log.Println("recv : ", i32Cnt, " bytes")

		if (i32FreeHead - i32DataHead) >= CONST_PKG_HDR_LEN {
			i32Cnt, err = handle_pkg(arrReadBuf[i32DataHead:])
			i32DataHead += i32Cnt
		}
	}
}

func client() {
	// connect to the server
	c, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	// send the message
	msg := "Hello World"
	fmt.Println("Sending", msg)
	err = gob.NewEncoder(c).Encode(msg)
	if err != nil {
		fmt.Println(err)
	}
	c.Close()
}

func main() {
	go start_tcp_service(string(":9999"))
	go client()
	var input string
	fmt.Scanln(&input)
}
