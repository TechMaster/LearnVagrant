package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	result := "user site at http://techmaster.com/user\n"
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
	fmt.Println("Lắng nghe ở cổng 8004")
	log.Fatal(http.ListenAndServe(":8004", nil))
}
