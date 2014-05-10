// SmartNet project main.go
package main

import (
	"SmartNet/utils"
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
	ui32PkgLen int
	ui16Opcode uint16
	ui16Others uint16
}

type RegPkg struct {
	PkgHdr
	bytesName [10]byte
	bytesPwd  [10]byte
}

func (this RegPkg) get_size() int {
	var pkg RegPkg
	return int(unsafe.Sizeof(pkg))
}

var hdr_for_const PkgHdr

const (
	CONST_READ_BUF_LEN = 1024 * 64

	CONST_PKG_HDR_LEN = int(unsafe.Sizeof(hdr_for_const))
	CONST_MAX_PKG_LEN = 1024
)

type IBaseHandler interface {
	handle_pkg(*PkgHdr) error
}

type SRegHandler struct {
}

func (this SRegHandler) handle_pkg(pHdr *PkgHdr) error {
	log.Println("reg handler: ", *pHdr)

	return nil
}

var handlers [1024]IBaseHandler

const (
	OPCODE_REG_PKG = 1
)

func init() {
	handlers[OPCODE_REG_PKG] = SRegHandler{}
}

func handle_pkg(bs []byte) (iDone int, err error) {
	//get pkg header
	//slice struct:
	//array *uint8
	//len int
	//cap int
	pTmp := (**uint8)(unsafe.Pointer(&bs)) //get array addr
	ui32Len := len(bs)
	offset := 0

	for {

		if ui32Len < CONST_PKG_HDR_LEN {
			break
		}

		pHdr1 := (*PkgHdr)(unsafe.Pointer(*pTmp))
		log.Println(*pHdr1)
		pHdr := (*PkgHdr)(unsafe.Pointer(uintptr(unsafe.Pointer(*pTmp)) + uintptr(offset)))
		log.Println(*pHdr)

		if pHdr.ui32PkgLen > ui32Len {
			//not enough
			break
		}

		if handlers[pHdr.ui16Opcode].handle_pkg(pHdr) != nil {
			log.Fatal("god!\n")
			break
		}

		offset += pHdr.ui32PkgLen
		ui32Len -= pHdr.ui32PkgLen
	}

	return offset, err
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

		log.Println("recv buf : ", iDataIdx, ", ", iFreeIdx)
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

	var hdr PkgHdr
	log.Println("pkg hdr size is : ", unsafe.Sizeof(hdr))
	log.Println("pkg hdr align is : ", unsafe.Alignof(hdr.ui32PkgLen), ", ", unsafe.Alignof(hdr.ui16Opcode), ", ", unsafe.Alignof(hdr.ui16Others), ", ", unsafe.Alignof(hdr))
	log.Println("pkg hdr offset is : ", unsafe.Offsetof(hdr.ui32PkgLen), ", ", unsafe.Offsetof(hdr.ui16Opcode), ", ", unsafe.Offsetof(hdr.ui16Others))

	var pHdr *PkgHdr
	pHdr = (*PkgHdr)(unsafe.Pointer(&bytesWriteBuf))
	pHdr.ui16Opcode = OPCODE_REG_PKG
	pHdr.ui16Others = 1
	pHdr.ui32PkgLen = 16

	//bytesWriteBuf[0] = 123

	log.Println("client : ", *pHdr)

	//bytesWriteBuf[0] = 101
	//bytesWriteBuf[1] = 102

	i32Cnt, err := conn.Write(bytesWriteBuf[:16])
	pHdr.ui16Others = 2
	i32Cnt, err = conn.Write(bytesWriteBuf[:16])
	pHdr.ui16Others = 3
	i32Cnt, err = conn.Write(bytesWriteBuf[:16])

	pRegPkg := (*RegPkg)(unsafe.Pointer(&bytesWriteBuf))
	pRegPkg.ui32PkgLen = pRegPkg.get_size()
	pRegPkg.ui16Opcode = OPCODE_REG_PKG
	copy(pRegPkg.bytesName[:], "rock")
	copy(pRegPkg.bytesPwd[:], "pswd")
	utils.Write_all_the_data(conn, bytesWriteBuf[:pRegPkg.get_size()])

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

	var arr [1024]byte
	log.Printf("%p", &arr)
	log.Printf("%p", unsafe.Pointer(uintptr(unsafe.Pointer(&arr))+1)) //the address only increase 1 byte

	go start_tcp_service(string(":9999"))
	go client(string("127.0.0.1:9999"))
	var input string
	fmt.Scanln(&input)
}
