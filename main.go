// SmartNet project main.go
package main

import (
	"SmartNet/base"
	"SmartNet/pkg"
	"SmartNet/utils"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
	"unsafe"
)

func StartTcpServerService(strLsnAddr string /*":9999"*/) {
	// listen on a port
	lsner, err := net.Listen("tcp", strLsnAddr)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("now is listening on : ", strLsnAddr)

	for {
		// accept a connection
		conn, err := lsner.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}

		//
		pSession := new(STcpSession)
		pSession.Start()
		// handle the connection
		//go RoutineReadFunc(conn)
	}
}

func handle_pkg(bs []byte) (int, error) {
	//get pkg header
	//slice struct:
	//array *uint8
	//len int
	//cap int
	pTmp := (**uint8)(unsafe.Pointer(&bs)) //get array addr
	ui32Len := len(bs)
	offset := 0
	var err error

	for {

		if ui32Len < pkg.CONST_PKG_HDR_LEN {
			break
		}

		pHdr := (*pkg.PkgHdr)(unsafe.Pointer(uintptr(unsafe.Pointer(*pTmp)) + uintptr(offset)))

		if pHdr.Mui32PkgLen > ui32Len {
			//not enough
			break
		}

		if pkg.Handlers[pHdr.Mui16Opcode].HandlePkg(pHdr) != nil {
			err = errors.New("handle pkg err!")
			log.Fatal("god!\n")
			break
		}

		offset += pHdr.Mui32PkgLen
		ui32Len -= pHdr.Mui32PkgLen
	}

	return offset, err
}

func start_tcp_client_service(strServerAddr string) {
	log.Println("start new tcp client service...")
	// connect to the server
	conn, err := net.Dial("tcp", strServerAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("established a new conn on ", strServerAddr)
	defer conn.Close()
	//sTimeOut := time.Time{5, 0, nil}
	//conn.SetWriteDeadline(sTimeOut)

	var bytesWriteBuf [pkg.CONST_WRITE_BUF_LEN]byte

	///for testing
	var hdr pkg.PkgHdr
	log.Println("pkg hdr size is : ", unsafe.Sizeof(hdr))
	log.Println("pkg hdr align is : ", unsafe.Alignof(hdr.Mui32PkgLen), ", ", unsafe.Alignof(hdr.Mui16Opcode), ", ", unsafe.Alignof(hdr.Mui16Others), ", ", unsafe.Alignof(hdr))
	log.Println("pkg hdr offset is : ", unsafe.Offsetof(hdr.Mui32PkgLen), ", ", unsafe.Offsetof(hdr.Mui16Opcode), ", ", unsafe.Offsetof(hdr.Mui16Others))

	//sign in server...
	pRegPkg := (*pkg.RegPkg)(unsafe.Pointer(&bytesWriteBuf))
	pRegPkg.Mui32PkgLen = pRegPkg.GetSize()
	pRegPkg.Mui16Opcode = pkg.OPCODE_REG_PKG
	copy(pRegPkg.MbytesName[:], "rock")
	copy(pRegPkg.MbytesPwd[:], "pswd")
	err = utils.WriteAllData(conn, bytesWriteBuf[:pRegPkg.GetSize()])

	//pRegPkg.Mui32PkgLen = pRegPkg.GetSize()
	//pRegPkg.Mui16Opcode = pkg.OPCODE_REG_PKG
	//copy(pRegPkg.MbytesName[:], "liuhy")
	//copy(pRegPkg.MbytesPwd[:], "haha")
	//err = utils.WriteAllData(conn, bytesWriteBuf[:pRegPkg.GetSize()])

	//pRegPkg.Mui32PkgLen = pRegPkg.GetSize()
	//pRegPkg.Mui16Opcode = pkg.OPCODE_REG_PKG
	//copy(pRegPkg.MbytesName[:], "hheheheheh")
	//copy(pRegPkg.MbytesPwd[:], "hehehehfdsfd")
	//err = utils.WriteAllData(conn, bytesWriteBuf[:pRegPkg.GetSize()])

	//msg := "Hello World"
	//fmt.Println("Sending", i32Cnt)
	//err = gob.NewEncoder(c).Encode(msg)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(1000000000000000000)
	//conn.Close()
}

func main() {

	var arr [1024]byte
	log.Printf("%p", &arr)
	log.Printf("%p", unsafe.Pointer(uintptr(unsafe.Pointer(&arr))+1)) //the address only increase 1 byte

	go StartTcpServerService(string(":9999"))
	go start_tcp_client_service(string("127.0.0.1:9999"))
	var input string
	fmt.Scanln(&input)
}
