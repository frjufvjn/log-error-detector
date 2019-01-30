package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const FILENAME = "C:/workspace_new/log-error-detector/web/index.html" // file name

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
Ref.
	http://pyrasis.com/book/GoForTheReallyImpatient/Unit50/05
	https://gobyexample.com/reading-files
*/

func main() {
	dat, err := ioutil.ReadFile(FILENAME)
	check(err)
	fmt.Print(string(dat))

	f, err := os.Open(FILENAME)
	check(err)
	defer f.Close()

	b1 := make([]byte, 5)
	n1, err := f.Read(b1)
	check(err)
	fmt.Printf("%d bytes: %s\n", n1, string(b1))

	o2, err := f.Seek(6, 0)
	check(err)
	b2 := make([]byte, 2)
	n2, err := f.Read(b2)
	check(err)
	fmt.Printf("%d bytes @ %d: %s\n", n2, o2, string(b2))

	o3, err := f.Seek(6, 0)
	check(err)
	b3 := make([]byte, 2)
	n3, err := io.ReadAtLeast(f, b3, 2)
	check(err)
	fmt.Printf("%d bytes @ %d: %s\n", n3, o3, string(b3))

	// seek 하고나서 ReadLine 하는 예제
	_, err = f.Seek(400, 0)
	check(err)

	r4 := bufio.NewReader(f)
	b4, err := r4.Peek(15)

	check(err)
	fmt.Printf("5 bytes: %s\n", string(b4))

	n := 1
	for n < 100 {
		line, _, err := r4.ReadLine()
		if err != nil {
			break
		}
		fmt.Printf(">> %s\n", string(line))
		// fmt.Println(">>", isPrefix)
		n++
	}

}
