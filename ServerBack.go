package main

import (

	"log"
	"io"
	"io/ioutil"
	"os"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"encoding/json"
	"encoding/base64"
	"fmt"
)


var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)


type img struct {
	ReturnUrl string `json:"returnUrl"`
	Type string `json:"type"`
	Name string `json:"name"`
	ImgBase64 string `json:"imgBase64"`
	Data string `json:"data"`

}



type return_data struct {
	Code int
	Data string


}

func Init(
traceHandle io.Writer,
infoHandle io.Writer,
warningHandle io.Writer,
errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func processImage(t img) {

	p, err := base64.StdEncoding.DecodeString(t.ImgBase64)
	if (err != nil){
		Info.Println("error")
		Info.Println(err)
	}
	f, err := os.Create("Media/Final/"+t.Name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(p); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}
	fmt.Println(t.Data)
}

func put_id(w http.ResponseWriter, r *http.Request) {
	var data_send return_data
	data_send.Code = 200
	data_json, err := json.Marshal(data_send)
	if (err != nil){
		Info.Println("error")
		Info.Println(err)
	}
	decoder := json.NewDecoder(r.Body)
	var t img
	err2 := decoder.Decode(&t)
	go processImage(t)
	if (err2 != nil){
		Info.Println("error")
		Info.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data_json))) //len(dec)
	w.Write(data_json)
}


func main() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)


	Info.Println("1")
	r := mux.NewRouter()
	Info.Println("2")
	r.HandleFunc("/get_img/", put_id)
	http.Handle("/", r)
	Info.Println("port :9060")
	err3 := http.ListenAndServe(":9060", nil)
	if (err3 != nil){
		Info.Println("error")
		Info.Println(err3)
	}else {
		Info.Println("encour")
	}


}
