package main
import (
    "fmt"
    "log"
    "os"
    "github.com/fsnotify/fsnotify"
)

const MYFILE = "logfile.log"

func main() {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()
    done := make(chan bool)
    go func() {
        for {
            select {
            case event := <-watcher.Events:
                log.Println("event:", event)
                if event.Op&fsnotify.Write == fsnotify.Write {
                    log.Println("modified file:", event.Name)
                    readFile(MYFILE)
                }
            case err := <-watcher.Errors:
                log.Println("error:", err)
            }
        }
    }()
    err = watcher.Add(MYFILE)
    if err != nil {
        log.Fatal(err)
    }
    <-done
}

func readFile(fname string) {
    file, err := os.Open(fname)
    if err != nil {
        panic(err)
    }
    defer file.Close()
    buf := make([]byte, 1024)
    stat, err := os.Stat(fname)
    start := stat.Size() - 1024
    _, err = file.ReadAt(buf, start)
    if err == nil {
        fmt.Printf("%s\n", buf)
    }
}
