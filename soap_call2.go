package main

import (
    "flag"
    "net/http"
    "fmt"
    "bytes"
    )

func GetPlainTextHeader () (map[string]string) {
    header := make(map[string]string)
    header["Content-Type"] = "application/soap+xml;charset=UTF-8"
    header["User-Agent"] = "Python WinRM client"
    return header
    }

func CreateSoapRequest (x map[string]string, message string) (*http.Request) {
    buffer := &bytes.Buffer{}
    req, err := http.NewRequest("POST", "http://192.168.122.196:5985/wsman", buffer)
    if err != nil{
        println("ERROR:", err)
        }
    for k, v := range x {
        req.Header.Add(k, v)
        }
    
    req.ContentLength = int64(len(message))
    req.SetBasicAuth("Administrator", "Passw0rd")
    return req
}

func main(){
   flag.Parse()
    head := GetPlainTextHeader()
    request := CreateSoapRequest(head, "dasdadasd")
    client := http.Client{}
    resp, err := client.Do(request)
    if err != nil{}
    fmt.Println(resp, "\n")
    fmt.Println(resp.StatusCode)
    return
}

