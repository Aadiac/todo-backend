package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func ErrHandle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func PrintRawData(r *http.Request) {
	bytes, _ := io.ReadAll(r.Body)
	fmt.Println("Raw data : ", string(bytes))
}

func PreflightCheck(w http.ResponseWriter ,r *http.Request){
	if r.Method == "OPTIONS"{
		w.WriteHeader(http.StatusOK)
		return
	}
}