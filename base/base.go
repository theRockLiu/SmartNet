// base.go
package base

import (
	"net"
)

const (
	CONST_SESSION_WRITE_BUF_LEN = 1024 * 64
	CONST_SESSION_WRITE_BUF_LEN = 1024 * 64
)

type IBaseSession interface {
	HandleWriting() error
	HandleReading() error
	Start() error
	Stop() error
}

type STcpSession struct {
	MBytesWriteBuf [CONST_SESSION_WRITE_BUF_LEN]byte
	MBytesReadBuf  [CONST_SESSION_READ_BUF_LEN]byte
}

type SGlobalObj struct {
	MMapSessions map[int32](*IBaseSession)
}

func (this *SGlobalObj) AddSession(i32Id int32, sess *IBaseSession) error {
	this.MMapSessions[i32Id] = sess

	return nil
}

func (this *STcpSession) Init() error {

	return nil
}

func (this *STcpSession) WaitForWritngData() error {

	return nil
}

func (this *STcpSession) Start() error {

	go func(*STcpSession) {
		/*read data and handling*/

	}(this)

}

var GGlobalObj *SGlobalObj

func init() {
	GGlobalObj = new(SGlobalObj)
}

func RoutineReadFunc(conn net.Conn) {
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
