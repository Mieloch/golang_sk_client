package main
import (
	"net"
	"encoding/binary"
	"bytes"
)
type Request struct{
	code int64
}
func createRequest(code int64) []byte{
		byteBuffer := new(bytes.Buffer)
		var request = Request{code}
		binary.Write(byteBuffer, binary.BigEndian, &request)
		return byteBuffer.Bytes()
}
func getRemoteMachineCpuUsage(hostPort string) uint64{
		conn, err := net.Dial("tcp", hostPort)
		if err != nil {
			return 0
		}
		var request = createRequest(0)
		conn.Write(request)
		buffer := make([]byte, 8)
		conn.Read(buffer)
		buf2 := bytes.NewBuffer(buffer)
    		data,_ := binary.ReadUvarint(buf2)
		return data
}
