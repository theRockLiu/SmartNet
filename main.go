// SmartNet project main.go
package main

import (
	"fmt"
	"log"
	"net"
	"time"
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
	return iDone - 1, err
}

func handle_conn(conn net.Conn) {
	//
	defer conn.Close()
	//
	var bytesReadBuf [CONST_READ_BUF_LEN]byte
	i32Cnt, iDataIdx, iFreeIdx := int(0), int(0), int(0)
	var err error
	//
	for {

		if (CONST_READ_BUF_LEN - iFreeIdx) < CONST_MAX_PKG_LEN {
			copy(bytesReadBuf[:], bytesReadBuf[iDataIdx:iFreeIdx])
			iFreeIdx -= iDataIdx
			log.Println("free idx : ", iFreeIdx)
			iDataIdx = 0
		}

		i32Cnt, err = conn.Read(bytesReadBuf[iFreeIdx:])
		if err != nil {
			log.Fatal(err)
			return
		}
		iFreeIdx += i32Cnt
		log.Println("recv : ", i32Cnt, " bytes")

		if (iFreeIdx - iDataIdx) >= CONST_PKG_HDR_LEN {
			i32Cnt, err = handle_pkg(bytesReadBuf[iDataIdx:iFreeIdx])
			iDataIdx += i32Cnt
			log.Println("data idx : ", iDataIdx)
			if iDataIdx == iFreeIdx {
				iDataIdx = 0
				iFreeIdx = 0
			}
		}
	}
}

func client() {
	// connect to the server
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	// send the message
	//regPkg := RegPkg{32, 1, 1, "fd", "sd"}
	var bytesWriteBuf [1024 * 64]byte
	//bytesWriteBuf[0] = 32

	i32Cnt, err := conn.Write(bytesWriteBuf[:1024*64-2])
	//msg := "Hello World"
	fmt.Println("Sending", i32Cnt)
	//err = gob.NewEncoder(c).Encode(msg)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(1000000000000000000)
	//conn.Close()
}

func main() {
	go start_tcp_service(string(":9999"))
	go client()
	var input string
	fmt.Scanln(&input)
}
