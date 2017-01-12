package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"github.com/jroimartin/gocui"
	"regexp"
	"sort"
)
func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}
func sendScriptAutomatically(g *gocui.Gui, v *gocui.View) error {
	scriptView, err := g.View("script");
	if(err != nil){
	return err
	}


		content := scriptView.ViewBuffer()
			
		var bestNode *Node
		var i = 0
		for _, node := range nodeList.nodes{
			if(i==0){
				bestNode = node
			}
			if(node.usage<bestNode.usage && node.working == false){
				bestNode = node
			}
			i++
		}
		output := sendScriptToRemote(bestNode,content)
		scriptView.Clear()
		fmt.Fprint(scriptView, "Job sent to ", bestNode.hostPort, "\n" )
		fmt.Fprint(scriptView,output)
	
	return nil
}
func sendScript(g *gocui.Gui, v *gocui.View) error {
	var l string
	scriptView, err := g.View("script");
	if(err != nil){
	return err
	}

	_, cy := v.Cursor()
	l, err = v.Line(cy)
	if err != nil{
		l = ""
	}
	var isValidHostPort = regexp.MustCompile(`(\w|\d)*[:]\d*`)

	if isValidHostPort.MatchString(l){

		hostPort := isValidHostPort.FindString(l)
		content := scriptView.ViewBuffer()
		go sendScriptToRemote((nodeList.nodes[hostPort]),content)
		//scriptView.Clear()
		//fmt.Fprintf(scriptView,output)
	}
	
	return nil
}
func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}
func getScriptName(g *gocui.Gui, v *gocui.View) error {

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		v.Editable = true
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("msg"); err != nil {
			
			return err
		}		
	}
	return nil
}
func loadScript(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	l, err = v.Line(cy)
	if err != nil{
		return nil
	}
	var scriptView *gocui.View
	scriptView, _ = g.View("script")
	scriptView.Clear()	
	b, err := ioutil.ReadFile(l[0:len(l)-1])
	if err != nil{
		fmt.Fprintf(scriptView, "%s", "Failed to load script! Check path and try again")
	}
	
	fmt.Fprintf(scriptView, "%s", b)
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("script"); err != nil {
		return err
	}
	
	

	return nil
}

func closeMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("script"); err != nil {
		return err
	}
	return nil
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "script" {
		_, err := g.SetCurrentView("main")
		return err
	}
	_, err := g.SetCurrentView("script")
	return err
}


func scanUsages(g *gocui.Gui){
	for{
		select{
			case <-time.After(1000 * time.Millisecond):
				    keys := make([]string, 0, len(nodeList.nodes))
				    for k,_ := range nodeList.nodes {
					keys = append(keys, k)
				    }
				sort.Strings(keys)
				for _, k := range keys{
					node := nodeList.nodes[k]
					node.usage = getRemoteMachineCpuUsage(node.hostPort)
				}
				v, err := g.View("main")
				v.Clear()
				g.Execute(func(g *gocui.Gui) error {
						for _, k := range keys{
							node := nodeList.nodes[k]
						if err != nil {
							// handle error
						}
						
						if node.usage == -1{
							fmt.Fprintf(v, "\033[31;4m%d %s\033[0m\n",node.class, node.hostPort +" Unable to reach server!")
						}else if(node.working){
							fmt.Fprint(v,node.class," ", node.hostPort +" WORKING\n")
						}else{
							fmt.Fprint(v,node.class," ", node.hostPort +" CPU Usage=", node.usage, "%\n")
						}					
						}
				return nil
				})
			}
	}
}
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", -1, -1, 45, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorWhite
		v.SelFgColor = gocui.ColorBlack
		for _, node := range nodeList.nodes{
			fmt.Fprintln(v,node.class, " ", node.hostPort, " ", node.usage)
		}
		if _, err := g.SetCurrentView("main"); err != nil {
				return err
		}
		
	}
	if v, err := g.SetView("script", 45, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//b, err := ioutil.ReadFile("Mark.Twain-Tom.Sawyer.txt")

		fmt.Fprintf(v, "%s", "[Enter] to load script")
		v.Editable = false
		v.Wrap = true
		
	}

	return nil
}

var nodeList NodeList
func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	nodeList = readConfig()
	go scanUsages(g)
	
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	
}
