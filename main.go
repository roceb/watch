package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/fsnotify/fsnotify"
)

var count = 1

func main() {
	// Watch for files
	//io.PipeReader()
	//WatchForFiles()
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var lFlag = flag.String("l", path, "location to run command on")
	var cFlag = flag.String("c", "go run .", "command to be run")
	flag.Parse()
	fmt.Printf("Location: %s\nCmd: %s", *lFlag, *cFlag)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	doneChan := make(chan bool)

	go func(doneChan chan bool) {
		defer func() {
			doneChan <- true
		}()

		err := WatchForFiles(*lFlag, *cFlag)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("File has been changed")
	}(doneChan)

	<-doneChan
	//fmt.Println(err)
}

func WatchForFiles(filePath string, cmd string) (err error) {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			fmt.Printf("%d", count)
			RunCommand(cmd)
			count++
			initialStat = stat
			continue
		}

		time.Sleep(1 * time.Second)
	}

}

func RunCommand(arg string) {
	if len(arg) < 2 {
		cmd := exec.Command("go", "build", ".")
		cmd.Run()

	}
}
