package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jroimartin/gocui"
	"regexp"

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
		bestNode := findBestNode()
		go sendScriptToRemote(g, bestNode,content)
	return nil
}
func getLine(v *gocui.View) string{
	_, cy := v.Cursor()
	l, err := v.Line(cy)
	if err != nil{
		l = ""
	}
	return l
}
func sendScript(g *gocui.Gui, v *gocui.View) error {
	scriptView, err := g.View("script");
	if(err != nil){
		return err
	}
	line := getLine(v)
	var isValidHostPort = regexp.MustCompile(`(\w|\d)*[:]\d*`)
	if isValidHostPort.MatchString(line){
		hostPort := isValidHostPort.FindString(line)
		content := scriptView.ViewBuffer()
		go sendScriptToRemote(g, (nodeList.nodes[hostPort]),content)
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
	fileContent, err := ioutil.ReadFile(l[0:len(l)-1])
	if err != nil{
		fmt.Fprintf(scriptView, "%s", "Failed to load script! Check path and try again")
	}
	fmt.Fprintf(scriptView, "%s", fileContent)
	closeMsg(g,v)
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
	if v.Name() == "main" {
		_, err := g.SetCurrentView("jobs")
		return err
	}
	if v.Name() == "jobs" {
		_, err := g.SetCurrentView("script")
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", -1, -1, 45, maxY/2); err != nil {
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
		fmt.Fprintf(v, "%s", "[Enter] to load script")
		v.Editable = false
		v.Wrap = true
		
	}
	if v, err := g.SetView("jobs", -1, maxY/2, 45, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintf(v, "%s", "jobs")
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
