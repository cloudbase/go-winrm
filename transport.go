package winrm

import (
    "net/http"
    "strings"
    "crypto/tls"
    "fmt"
    "bytes"
    // "errors"
)

func GetHttpHeader () (map[string]string) {
    header := make(map[string]string)
    header["Content-Type"] = "application/soap+xml;charset=UTF-8"
    header["User-Agent"] = "Python WinRM client"
    return header
}

func HttpBasicAuth (url, username, pass string, data []byte, insecure bool) (*http.Response, error){
    protocol := strings.Split(url, ":")
    if protocol[0] != "http" && protocol[0] != "https"{
        fmt.Errorf("Invalid protocol. Expected http or https")
    }

    header := GetHttpHeader()
    
    client := &http.Client{}
    // Ignore SSL certificate errors
    if protocol[0] == "https" {
        tr := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
        }
        client.Transport = tr
    }
    body := bytes.NewBuffer(data)
    req, err := http.NewRequest("POST", url, body)
    req.ContentLength = int64(len(data))
    req.SetBasicAuth(username, pass)

    for k, v := range header {
        req.Header.Add(k, v)
    }

    resp, err := client.Do(req)
    return resp, err
}
