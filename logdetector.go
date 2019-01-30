package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/fsnotify/fsnotify"
	_ "github.com/go-sql-driver/mysql"
)

const MYFILE = "C:/workspace_new/rec-file-checker/log/20190130.log"
const READ_BUFFER_LIMIT = 512
const FIND_KEYWORD = "ERROR"

var CONTROL = "" // make(chan string)
var sizeChk int64

func main() {
	match, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`, "2019-01-29 13:01:43")
	log.Println("match: ", match)

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
				// log.Println("event:", event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					// log.Println("-modified file:", event.Name)
					// log.Println("event.Op:", event.Op)
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
	checkFileErr(err)
	defer file.Close()

	beforeSize := sizeChk
	stat, err := os.Stat(fname)
	checkFileErr(err)

	// TODO Race Condition 에 대한 고려
	sizeChk = stat.Size()

	log.Println("sizeChk, beforeSize :", sizeChk, beforeSize)

	_, err = file.Seek(beforeSize, 0)
	checkFileErr(err)
	reader := bufio.NewReader(file)

	n := 1
	for n < 100 {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		// TODO 라인의 처음부분에 시간정보가 있는것을 regexp 로 걸른다.
		// TODO 키워드를 추출한다.
		// TODO Throttling Logic 추가 구현

		fmt.Printf(">> %s\n", string(line))
		n++
	}
}

func checkFileErr(e error) {
	if e != nil {
		panic(e)
	}
}

// func readFile(fname string) {

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

func dbCheckError(err error) {
	if err != nil {
		panic(err)
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
