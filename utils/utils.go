// utils.go
package utils

import (
	"net"
)

func Write_all_the_data(conn net.Conn, bytesWriteBuf []byte) error {

	iStart, iRet, iLen := 0, 0, len(bytesWriteBuf)
	var err error

	for {

		if iLen == 0 {
			break
		}

		iRet, err = conn.Write(bytesWriteBuf[iStart:])
		if err != nil {
			return err
		}

		iLen -= iRet
		iStart += iRet

	}

	return err

}
