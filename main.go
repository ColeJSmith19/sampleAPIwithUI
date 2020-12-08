package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/test", test)
	http.HandleFunc("/testpost", testPost)
	log.Fatal(http.ListenAndServe(":18080", nil))
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Belmont")
}

func testPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
