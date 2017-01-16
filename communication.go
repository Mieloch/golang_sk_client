package main
import (
	"net"
	"encoding/binary"
	"bytes"
	"time"
	"github.com/jroimartin/gocui"
	"fmt"
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
func getRemoteMachineCpuUsage(hostPort string) (ret int32){
		conn, err := net.Dial("tcp", hostPort)
		if err != nil {
			return -1
		}
		var request = createRequest(0, "")
		conn.Write(request)
		buffer := make([]byte, 8)
		conn.Read(buffer)
		buf2 := bytes.NewBuffer(buffer)
		binary.Read(buf2, binary.LittleEndian, &ret)
    		//data,_ := binary.ReadVarint(buf2)
		conn.Close()
		return
}
func scanUsages(g *gocui.Gui){
	for{
		select{
			case <-time.After(1000 * time.Millisecond):
				nodesHostPorts := getNodeHostPorts()
				updateNodeUsage(nodesHostPorts)
				v, _ := g.View("main")
				v.Clear()
				g.Execute(func(g *gocui.Gui) error {
					for _, k := range nodesHostPorts{
						node := nodeList.nodes[k]			
						if node.usage == -1{
				fmt.Fprintf(v, "\033[31;4m%d %s\033[0m\n",node.class, node.hostPort +" Unable to reach server!\n")
						}else if(node.working){
							fmt.Fprintf(v,"\033[33;4m%d %s\033[0m\n", node.class, " " + node.hostPort +" WORKING")
						}else{
							fmt.Fprintf(v,"\033[32;4m%d %s %d%%\033[0m\n", node.class," " + node.hostPort +" CPU Usage=", node.usage)
						}					
					}
					return nil
				})
			}
	}
}

func sendScriptToRemote(g *gocui.Gui, node *Node, script string) string{
		conn, err := net.Dial("tcp", node.hostPort)
		if err != nil {
			return "connection fail"
		}
		var request = createRequest(1, script)
		conn.Write(request)
		node.working = true
		buffer := make([]byte, 1024)
		n,err := conn.Read(buffer)
		node.working = false
		if err != nil {
			return "fail"
		}
		conn.Close()
		s := string(buffer[:n])
		g.Execute(func(g *gocui.Gui) error {
			v, _ := g.View("jobs")
			fmt.Fprint(v,node.hostPort + "job execution time = " + s +"s\n")
			return nil
		})
		return s
}
