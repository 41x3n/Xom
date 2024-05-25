package shared

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

// CloseResponseBody closes the response body and handles any error
func CloseResponseBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Printf("Error while closing the file: %v", err)
	}
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// GetGoroutineID util function to printout goroutine id
func GetGoroutineID() {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	fmt.Println("GOROUTINE ID: ", idField)
}
