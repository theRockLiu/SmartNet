// base.go
package base

type IBaseSession interface {
	WaitForWritingData() error
	Init() error
}

type STcpSession struct {
	mBytesWriteBuf [1024 * 64]byte
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

var GGlobalObj *SGlobalObj

func init() {
	GGlobalObj = new(SGlobalObj)
}
