package main
import (
	"os"
	"bufio"
)

func readConfig() []string{
	
    file, err := os.Open("config")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
	remotes = append(remotes,scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        panic(err)
    }
return remotes
}
