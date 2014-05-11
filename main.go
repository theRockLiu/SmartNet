// SmartNet project main.go
package main

import (
	"SmartNet/pkg"
	"SmartNet/utils"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
	"unsafe"
)

func start_tcp_service(strLsnAddr string /*":9999"*/) {
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

		// handle the connection
		go handle_conn(conn)
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

func handle_conn(conn net.Conn) {
	log.Println("established a new conn, remote addr is : ", conn.RemoteAddr().String())
	//
	defer conn.Close()
	//
	var bytesReadBuf [pkg.CONST_READ_BUF_LEN]byte
	i32Cnt, iDataIdx, iFreeIdx := int(0), int(0), int(0)
	var err error
	//
	for {

		if (pkg.CONST_READ_BUF_LEN - iFreeIdx) < pkg.CONST_MAX_PKG_LEN {
			copy(bytesReadBuf[:], bytesReadBuf[iDataIdx:iFreeIdx])
			iFreeIdx -= iDataIdx
			iDataIdx = 0
			log.Println("move data to read buf head, data idx is 0 and free idx is ", iFreeIdx)
		}

		i32Cnt, err = conn.Read(bytesReadBuf[iFreeIdx:])
		if err != nil {
			log.Fatal(err)
			return
		}
		iFreeIdx += i32Cnt
		log.Println("conn read new bytes, now read buf's data idx is ", iDataIdx, ", and free idx is ", iFreeIdx)

		if (iFreeIdx - iDataIdx) >= pkg.CONST_PKG_HDR_LEN {
			i32Cnt, err = handle_pkg(bytesReadBuf[iDataIdx:iFreeIdx])
			iDataIdx += i32Cnt
			log.Println("hande pkg is done, now data idx is ", iDataIdx, " and free idx is ", iFreeIdx)
			if iDataIdx == iFreeIdx {
				iDataIdx = 0
				iFreeIdx = 0
			}
		}
	}
}

func client(strServerAddr string) {
	// connect to the server
	conn, err := net.Dial("tcp", strServerAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// send the message
	//regPkg := RegPkg{32, 1, 1, "fd", "sd"}
	var bytesWriteBuf [1024 * 64]byte
	//bytesWriteBuf[0] = 32

	var hdr pkg.PkgHdr
	log.Println("pkg hdr size is : ", unsafe.Sizeof(hdr))
	log.Println("pkg hdr align is : ", unsafe.Alignof(hdr.Mui32PkgLen), ", ", unsafe.Alignof(hdr.Mui16Opcode), ", ", unsafe.Alignof(hdr.Mui16Others), ", ", unsafe.Alignof(hdr))
	log.Println("pkg hdr offset is : ", unsafe.Offsetof(hdr.Mui32PkgLen), ", ", unsafe.Offsetof(hdr.Mui16Opcode), ", ", unsafe.Offsetof(hdr.Mui16Others))

	//var pHdr *PkgHdr
	//pHdr = (*PkgHdr)(unsafe.Pointer(&bytesWriteBuf))
	//pHdr.ui16Opcode = OPCODE_REG_PKG
	//pHdr.ui16Others = 1
	//pHdr.ui32PkgLen = 16

	////bytesWriteBuf[0] = 123

	//log.Println("client : ", *pHdr)

	////bytesWriteBuf[0] = 101
	////bytesWriteBuf[1] = 102

	//i32Cnt, err := conn.Write(bytesWriteBuf[:16])
	//pHdr.ui16Others = 2
	//i32Cnt, err = conn.Write(bytesWriteBuf[:16])
	//pHdr.ui16Others = 3
	//i32Cnt, err = conn.Write(bytesWriteBuf[:16])

	pRegPkg := (*pkg.RegPkg)(unsafe.Pointer(&bytesWriteBuf))
	pRegPkg.Mui32PkgLen = pRegPkg.GetSize()
	pRegPkg.Mui16Opcode = pkg.OPCODE_REG_PKG
	copy(pRegPkg.MbytesName[:], "rock")
	copy(pRegPkg.MbytesPwd[:], "pswd")
	err = utils.WriteAllData(conn, bytesWriteBuf[:pRegPkg.GetSize()])

	pRegPkg.Mui32PkgLen = pRegPkg.GetSize()
	pRegPkg.Mui16Opcode = pkg.OPCODE_REG_PKG
	copy(pRegPkg.MbytesName[:], "liuhy")
	copy(pRegPkg.MbytesPwd[:], "haha")
	err = utils.WriteAllData(conn, bytesWriteBuf[:pRegPkg.GetSize()])

	pRegPkg.Mui32PkgLen = pRegPkg.GetSize()
	pRegPkg.Mui16Opcode = pkg.OPCODE_REG_PKG
	copy(pRegPkg.MbytesName[:], "hheheheheh")
	copy(pRegPkg.MbytesPwd[:], "hehehehfdsfd")
	err = utils.WriteAllData(conn, bytesWriteBuf[:pRegPkg.GetSize()])

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

	go start_tcp_service(string(":9999"))
	go client(string("127.0.0.1:9999"))
	var input string
	fmt.Scanln(&input)
}
