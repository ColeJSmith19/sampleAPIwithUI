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
	// var reader io.Reader

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// reader.Read(body)
	fmt.Println(string(body))

	// mail, e := mail.ReadMessage(reader)
	// if e != nil {
	// 	fmt.Fprintf(w, e.Error())
	// 	return
	// }

	// fmt.Println(mail.Header)

	// fmt.Fprintf(w, mail.Body)
}
