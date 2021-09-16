package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	result := "video service at http://techmaster.com/video\n"
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

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Lắng nghe ở cổng 8005")
	log.Fatal(http.ListenAndServe(":8005", nil))
}
