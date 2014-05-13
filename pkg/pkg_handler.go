// pkg_handler
package pkg

import (
	"SmartNet/base"
	"log"
	"unsafe"
)

type IBaseHandler interface {
	HandlePkg(*PkgHdr) error
}

type SRegHandler struct {
}

func (this SRegHandler) HandlePkg(pHdr *PkgHdr) error {
	log.Println("reg handler: ", *pHdr)

	pRegPkg := (*RegPkg)(unsafe.Pointer(pHdr))

	log.Println(*pRegPkg)
	log.Println(string(pRegPkg.MbytesName[:]))
	log.Println(string(pRegPkg.MbytesPwd[:]))

	p := new(base.STcpSession)

	base.GGlobalObj.AddSession(1, (*base.IBaseSession)(unsafe.Pointer(p)))
	//base.GGlobalObj.AddSession(1, p)

	return nil
}

var Handlers [1024]IBaseHandler

const (
	OPCODE_REG_PKG = 1
)

func init() {
	Handlers[OPCODE_REG_PKG] = SRegHandler{}
}
