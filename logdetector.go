package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	_ "github.com/go-sql-driver/mysql"
)

const (
	MAX_READ_LINE     = 100
	READ_BUFFER_LIMIT = 1024
)

var (
	sizeChk  int64 = 0
	sizeMap  sync.Map
	patterns []string
)

// var CONTROL = "" // make(chan string)

type Configuration struct {
	Patterns []string
	Logfiles []string
}

func main() {
	match, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`, "2019-01-29 13:01:43")
	log.Println("match: ", match)

	logs := getConfig("logfile")
	patterns = getConfig("pattern")

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

				if event.Op&fsnotify.Write == fsnotify.Write {
					readFile(event.Name)
				}

			case err := <-watcher.Errors:
				log.Println("error:", err)

			}
		}
	}()

	for idx := range logs {
		err := watcher.Add(logs[idx])
		if err != nil {
			log.Fatal(err)
		}
	}

	<-done

}

func readFile(fname string) {
	// file open
	file, err := os.Open(fname)
	checkFileErr(err)
	defer file.Close()

	// get previous file size (avoid race condition)
	// [GET]beforeSize := atomic.LoadInt64(&sizeChk)
	var beforeSize int64
	res, ok := sizeMap.Load(fname)
	if ok {
		beforeSize = res.(int64)
	} else {
		beforeSize = 0
	}

	// get current file size
	stat, err := os.Stat(fname)
	checkFileErr(err)
	// [SET]atomic.StoreInt64(&sizeChk, stat.Size())
	sizeMap.Store(fname, stat.Size())

	// log.Println("[DEBUG] sizeChk, beforeSize :", sizeChk, beforeSize)

	// file seek to
	_, err = file.Seek(beforeSize, 0)
	checkFileErr(err)

	// file read
	reader := bufio.NewReader(file)

	if beforeSize > 0 {
		n := 1
		for n < MAX_READ_LINE {
			line, _, err := reader.ReadLine()
			if err != nil {
				break
			}

			// TODO 라인의 처음부분에 시간정보가 있는것을 regexp 로 걸른다.
			// TODO Throttling Logic 추가 구현

			// find patterns from config file
			for idx := range patterns {

				// Detect keyword from log line
				if bytes.Contains(line, []byte(patterns[idx])) {
					fmt.Printf("DETECTED--> %s\n", string(line))

					// mysql db insert
					db, err := sql.Open("mysql", "root:!QAZ2wsx@tcp(127.0.0.1:3306)/test")
					if err != nil {
						log.Fatal(err)
					}
					defer db.Close()

					stmt, err := db.Prepare("INSERT INTO log_detect (category,content) values (?, ?)")
					checkFileErr(err)
					defer stmt.Close()

					_, err = stmt.Exec(patterns[idx], string(line))
					checkFileErr(err)
				}
			}

			n++
		}
	}
}

func getConfig(key string) []string {
	// from json file
	file, err := os.Open("conf.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	decErr := decoder.Decode(&configuration)
	if decErr != nil {
		panic(decErr)
	}

	if key == "pattern" {
		return configuration.Patterns
	} else {
		return configuration.Logfiles
	}
}

func checkFileErr(e error) {
	if e != nil {
		panic(e)
	}
}

func findKeywordUsingSplit(data string, keyword string) string {
	arr := strings.Split(data, "\n")

	for i, n := range arr {
		if strings.Index(n, keyword) != -1 {
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

// func readFileOld(fname string) {

// 	file, err := os.Open(fname)

// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()

// 	buf := make([]byte, READ_BUFFER_LIMIT)
// 	stat, err := os.Stat(fname)
// 	start := stat.Size() - READ_BUFFER_LIMIT

// 	_, err = file.ReadAt(buf, start)

// 	if err == nil {
// 		if bytes.Contains(buf, []byte(FIND_KEYWORD)) {

// 			strBuf := string(buf[:])
// 			findStr := findKeywordUsingSplit(strBuf, FIND_KEYWORD)
// 			// fmt.Printf("[%s]\n", findStr)
// 			tmp := findStr[0:23]

// 			if tmp != CONTROL {

// 				CONTROL = findStr[0:23]
// 				fmt.Println("start-----------------------------------------------")
// 				fmt.Printf("%s\n", buf)
// 				fmt.Println("end-------------------------------------------------")

// 				db, err := sql.Open("mysql", "root:!QAZ2wsx@tcp(127.0.0.1:3306)/test")
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 				defer db.Close()

// 				stmt, err := db.Prepare("INSERT INTO log_detect (category,content) values (?, ?)")
// 				dbCheckError(err)
// 				defer stmt.Close()

// 				_, err = stmt.Exec(tmp, strBuf)
// 				dbCheckError(err)
// 			}
// 		}
// 	}
// }
