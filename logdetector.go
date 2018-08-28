package main

import (
   	"bytes"
   	"fmt"
   	"log"
   	"os"
   	"strings"

   	"github.com/fsnotify/fsnotify"
   	"database/sql"
   	_ "github.com/go-sql-driver/mysql"
)

const MYFILE = "/somepath/logfile.log"
const READ_BUFFER_LIMIT = 2048
const FIND_KEYWORD = "ERROR"

var CONTROL = "" // make(chan string)

func main() {

   	watcher, err := fsnotify.NewWatcher()

   	if err != nil {
       	log.Fatal(err)
   	}

   	defer watcher.Close()

   	done := make(chan bool)
   
   	go func() {
       	for 	{
           	select {
           	case event := <-watcher.Events:
               // log.Println("event:", event)

               if event.Op&fsnotify.Write == fsnotify.Write {
                   // log.Println("-----------------------modified file:", event.Name)
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

   	buf := make([]byte, READ_BUFFER_LIMIT)
   	stat, err := os.Stat(fname)
   	start := stat.Size() - READ_BUFFER_LIMIT
   	_, err = file.ReadAt(buf, start)

   	if err == nil {
       	if bytes.Contains(buf, []byte(FIND_KEYWORD)) {

           	strBuf := string(buf[:])
           	findStr := findKeywordUsingSplit(strBuf, FIND_KEYWORD)
           	tmp := findStr[0:19]

           	if tmp != CONTROL {
               	CONTROL = findStr[0:19]
fmt.Println("start-----------------------------------------------")
                           	fmt.Printf("%s\n", buf)
fmt.Println("end-------------------------------------------------")
               
				db, err := sql.Open("mysql", "root:passwd@tcp(127.0.0.1:3306)/test")
				if err != nil {
                    log.Fatal(err)
				}
				defer db.Close()

				stmt, err := db.Prepare("INSERT INTO log_detect (category,content) values (?, ?)")
				dbCheckError(err)
				defer stmt.Close()

				_, err = stmt.Exec(tmp, strBuf)
				dbCheckError(err)
			}
		}
	}
}

func dbCheckError(err error) {
   if err != nil {
       panic(err)
   }
}

func findKeywordUsingSplit(data string, keyword string) string {
   arr := strings.Split(data, "\n")

   for i, n := range arr {
       if strings.Index(n,keyword) != -1 {
           return arr[i]
       }
   }
   return ""
}

func Find(data []string, keyword string) int {
   for i, n := range data {
       if keyword == n {
           return i
       }
   }
   return -1
}

func bytesToString(data []byte) string {
   return string(data[:])
}

func indexOf(slice []string, item string) int {
   for i, _ := range slice {
       if slice[i] == item {
           return i
       }
   }
   return -1
}
