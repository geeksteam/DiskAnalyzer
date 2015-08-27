package diskanalyzer

import (
	_ "fmt"
	"io"
	"log"
	"net/http"
	"testing"
	_ "time"
)

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	panic("heyhey")
	io.WriteString(w, "hello, world!\n")
}

func TestTest(t *testing.T) {
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
