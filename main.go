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

    v := &winrm.Envelope{}
    shell, _ := v.GetShell(p, Soap)
    cmdParam := winrm.CmdParams{
        ShellID: shell,
        Cmd: "dir",
        Args: "c:\\ /A",
    }
    ret, _ := v.SendCommand(cmdParam, Soap)
    strdout, stderr, ret_code, _ := v.GetCommandOutput(shell, ret, Soap)
    _ = v.CleanupShell(shell, ret, Soap)
    _ = v.CloseShell(shell, Soap)
    fmt.Printf("Output:%s\nError: %s\nCode:%i\n", strdout, stderr, ret_code)
}