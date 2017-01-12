package main
import (
	"os"
	"bufio"
	"strings"
	"strconv"
)
type Node struct{
	class int
	hostPort string
	usage int32
	working bool
}
type NodeList struct{
	nodes map[string]*Node
}
func readConfig() NodeList{
    nodes := make(map[string]*Node)
    file, err := os.Open("config")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
	configLine := scanner.Text()
	split := strings.Split(configLine,";")	
	class, _ := strconv.ParseInt(split[0],10,64)	
	node := Node{int(class),split[1], 0, false}	
	nodes[node.hostPort] = &node
    }

    if err := scanner.Err(); err != nil {
        panic(err)
    }
return NodeList{nodes}
}
