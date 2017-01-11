package main
import (
	"net"
	"encoding/binary"
	"bytes"



)
type Request struct{
	code int32
	buffer [1024] byte
}
func createRequest(code int32, payload string) []byte{
		byteBuffer := new(bytes.Buffer)
		var a [1024]byte
		copy(a[:], payload)
		var request = Request{code,a}
		binary.Write(byteBuffer, binary.BigEndian, &request)
		return byteBuffer.Bytes()
}
func getRemoteMachineCpuUsage(hostPort string) uint64{
		conn, err := net.Dial("tcp", hostPort)
		if err != nil {
			return 0
		}
		var request = createRequest(0, "")
		conn.Write(request)
		buffer := make([]byte, 8)
		conn.Read(buffer)
		buf2 := bytes.NewBuffer(buffer)
    		data,_ := binary.ReadUvarint(buf2)
		conn.Close()
		return data
}
func sendScriptToRemote(hostPort string, script string) string{
		conn, err := net.Dial("tcp", hostPort)
		if err != nil {
			return "connection fail"
		}
		var request = createRequest(1, script)
		conn.Write(request)
		buffer := make([]byte, 1024)
		n,err := conn.Read(buffer)
		if err != nil {
			return "fail"
		}
		conn.Close()
		s := string(buffer[:n])
		return s
}
