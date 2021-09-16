package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	result := "SECRET at http://techmaster.com/secret\n"

	_, _ = fmt.Fprintf(w, result)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Lắng nghe ở cổng 8007")
	log.Fatal(http.ListenAndServe(":8007", nil))
}
