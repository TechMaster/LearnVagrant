package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	result := "main site at http://techmaster.com\n"
	out, err := exec.Command("hostname").Output()
	if err != nil {
		log.Fatal(err)
	}
	result += string(out)

	//----
	out, err = exec.Command("hostname", "-i").Output()
	if err != nil {
		log.Fatal(err)
	}
	result += string(out)

	//----
	out, err = exec.Command("cat", "/etc/os-release").Output()
	if err != nil {
		log.Fatal(err)
	}
	result += string(out)
	fmt.Println(time.Now())
	_, _ = fmt.Fprintf(w, result)
}

func bloghandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now())
	_, _ = fmt.Fprintf(w, "This is blog in main site")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/blog", bloghandler)
	fmt.Println("Lắng nghe ở cổng 8001")
	log.Fatal(http.ListenAndServe(":8001", nil))
}
