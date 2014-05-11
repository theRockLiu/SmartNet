// smart_pkg.go
package pkg

import (
	"unsafe"
)

type PkgHdr struct {
	Mui32PkgLen int
	Mui16Opcode uint16
	Mui16Others uint16
}

type RegPkg struct {
	PkgHdr
	MbytesName [10]byte
	MbytesPwd  [10]byte
}

func (this RegPkg) GetSize() int {
	var pkg RegPkg
	return int(unsafe.Sizeof(pkg))
}

var hdr_for_const PkgHdr

const (
	CONST_READ_BUF_LEN = 1024 * 64

	CONST_PKG_HDR_LEN = int(unsafe.Sizeof(hdr_for_const))
	CONST_MAX_PKG_LEN = 1024
)
