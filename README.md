# WinRM

This is a Go module for interacting with the Windows Remote Management system (WinRM)


Here is a quick usage example:

```go
package main


import (
    "fmt"
    "net/http"
    "github.com/trobert2/winrm"
)


func main(){
    p := winrm.ShellParams{}
    Soap := winrm.SoapRequest{
        Endpoint:"https://192.168.100.154:5986/wsman",
        Username:"Administrator",
        Passwd:"Passw0rd",
        AuthType:"BasicAuth",
        HttpInsecure:true,
        HttpClient: &http.Client{},
    }
    cmdParam := winrm.CmdParams{
        Cmd: "dir",
        Args: "c:\\ /A",
    }

    v := &winrm.Envelope{}
    strdout, stderr, ret_code, _ := v.RunCommand(p, cmdParam, Soap)
    fmt.Printf("Output:%s\nError: %s\nCode:%v\n", strdout, stderr, ret_code)
}
```
