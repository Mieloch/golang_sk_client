package main
import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"sort"
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
func findBestNode() *Node{
	var bestNode *Node
		var i = 0
		for _, node := range nodeList.nodes{
			if(i==0){
				bestNode = node
			}
			if(node.usage<bestNode.usage && node.working == false && node.usage > -1){
				bestNode = node
			}
			i++
		}
	return bestNode
}
func getNodeHostPorts() []string{
	keys := make([]string, 0, len(nodeList.nodes))
	for k,_ := range nodeList.nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
func updateNodeUsage(nodesHostPorts []string){
	for _, k := range nodesHostPorts{
		node := nodeList.nodes[k]
		node.usage = getRemoteMachineCpuUsage(node.hostPort)
	}
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
