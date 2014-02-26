# WinRM

This is a Go module for interacting with the Windows Remote Management system (WinRM)


Here is a quick usage example with Basic Authentication:

```go
package main


// We use a forked version of net/http and crypto/tls
// because the standard libs do now support renegotiation

import (
    "fmt"
    "launchpad.net/gwacl/fork/http"
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


And one with Client side certificates:

```Go
package main

// We use a forked version of net/http and crypto/tls
// because the standard libs do now support renegotiation

import (
    "fmt"
    "launchpad.net/gwacl/fork/http"
    "github.com/trobert2/winrm"
)


func main(){
    p := winrm.ShellParams{}
    Soap := winrm.SoapRequest{
        Endpoint:"https://192.168.100.155:5986/wsman",
        AuthType:"CertAuth",
        HttpInsecure:true,
        CertAuth: &winrm.CertificateCredentials{
            Cert: "/home/ubuntu/maas/SSL/certs/testing.pfx.pem",
            Key: "/home/ubuntu/maas/SSL/certs/dec.key",
        },
        HttpClient: &http.Client{},
    }
    cmdParam := winrm.CmdParams{
        Cmd: "dir",
        Args: "c:\\ /A",
    }

    v := &winrm.Envelope{}
    strdout, stderr, ret_code, err := v.RunCommand(p, cmdParam, Soap)
    fmt.Printf("Output:%s\nError: %s\nCode:%v\nERROR:%s\n", strdout, stderr, ret_code, err)
}
```