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
	"os/exec"
	"fmt"
	"bytes"
	"strings"
	"time"
	"math/rand"
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

type imgReturn struct {
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

	var nameImage string
	t.Type = strings.ToLower(t.Type)
	t.Name = strings.ToLower(t.Name)
	nameImage=  strings.Replace(t.Name , "jpeg", "jpg", -1)
	nameImage =  strings.Replace(nameImage , ".jpeg", "", -1)
	nameImage =  strings.Replace(nameImage , ".jpg", "", -1)
	nameImage =  strings.Replace(nameImage , ".png", "", -1)
	nameImage = nameImage+strconv.FormatInt(time.Now().Unix(),10)+ "_" + strconv.Itoa(rand.Intn(6666))+ "_" + strconv.Itoa(rand.Intn(6666))
	nameImage = nameImage+"."+t.Type
	p, err := base64.StdEncoding.DecodeString(t.ImgBase64)
	if (err != nil){
		Info.Println("error")
		Info.Println(err)
	}
	f, err := os.Create("Media/Raw/"+nameImage)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(p); err != nil {
		panic(err)
	}
	fi, err := f.Stat()
	if err != nil {
		// Could not obtain stat, handle error
	}
        var size int64
	var base64Image string
	base64Image = ""
	size = fi.Size()
	fmt.Println("The file is %d bytes long", size)
	if (size < 3666000) {


		if err := f.Sync(); err != nil {
			panic(err)
		}
		fmt.Println("En cour")
		if (t.Type == "jpg") {
			out, err := exec.Command("/bin/bash", "-c", "App/bin/Release/./guetzli " + "Media/Raw/" + nameImage + " Media/G/" + nameImage).Output()
			if err != nil {
				fmt.Println("error occured")
				fmt.Printf("122")
				fmt.Printf("%s", err)
			}
			fmt.Println("Finish")
			buff, err := ioutil.ReadFile("Media/G/" + nameImage)
			if err != nil {
			}
			fmt.Printf("%s", out)

			base64Image = base64.StdEncoding.EncodeToString(buff)
		}
		if (t.Type == "png") {

			out, err := exec.Command("/bin/bash", "-c", "zopfli/./zopflipng " +"--iterations=1500" +" Media/Raw/" + nameImage + " Media/G/" + nameImage).Output()
			if err != nil {
				fmt.Println("error occured")
				fmt.Printf("138")
				fmt.Printf("%s", err)
			}
			fmt.Println("Finish")
			buff, err := ioutil.ReadFile("Media/G/" + nameImage)
			if err != nil {
				fmt.Println("error occured")
				fmt.Printf("146")
				fmt.Printf("%s", err)
			}
			fmt.Printf("%s", out)

			base64Image = base64.StdEncoding.EncodeToString(buff)
		}

		err = os.Remove("Media/G/" + nameImage)
		if err != nil {
			fmt.Println("error occured")
			fmt.Printf("154")
			fmt.Printf("%s", err)
		}

	}
	err = os.Remove("Media/Raw/" + nameImage)
	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)
	}


	data := &imgReturn{Type : "jpg", Name : t.Name, ImgBase64 : base64Image , Data: t.Data}
	data_json, err := json.Marshal(data)
	var jsonStr = []byte(data_json)
	req, err := http.NewRequest("POST", t.ReturnUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
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
	r.HandleFunc("/put_img/", put_id)
	http.Handle("/", r)
	Info.Println("port :9040")
	err3 := http.ListenAndServe(":9040", nil)
	if (err3 != nil){
		Info.Println("error")
		Info.Println(err3)
	}else {
		Info.Println("encour")
	}


}
