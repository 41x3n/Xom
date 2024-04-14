package shared

import (
	"fmt"
	"io"
	"log"
)

// CloseResponseBody closes the response body and handles any error
func CloseResponseBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
