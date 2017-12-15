package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type webhook struct{}

func (webhook) push(p pdf, j job) {
	req, err := http.NewRequest("POST", j.Webhook, bytes.NewBuffer(p.GetContent()))
	if err != nil {
		panic(err)
	}
	fmt.Println(req)
	req.Header.Set("X-Go-Pdf-Bot", "https://github.com/Kirouane/gopdfbot")
	req.Header.Set("Content-Type", "application/pdf")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
