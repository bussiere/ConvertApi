package main

import (
	"bytes"
	"net/http"
	"encoding/json"
	"encoding/base64"
	"io/ioutil"
)

type img struct {
ReturnUrl string`json:"returnUrl"`
Type string `json:"type"`
Name string `json:"name"`
ImgBase64 string `json:"imgBase64"`
Data string `json:"Data"`

}


func main(){
	var url string
	url = "http://localhost:9040/put_img/"
	var base64Image string

	buff, err := ioutil.ReadFile("test.png")
	if err != nil {
	}

	base64Image  = base64.StdEncoding.EncodeToString(buff)

	data := &img{Type : "png", ReturnUrl : "http://localhost:9060/get_img/", Name : "toto.png", ImgBase64 : base64Image, Data : "id" }
	data_json, err := json.Marshal(data)
	var jsonStr = []byte(data_json)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}