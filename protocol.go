package winrm

import (
    "encoding/xml"
    // "io/ioutil"
    "fmt"
    "errors"
    "github.com/nu7hatch/gouuid"
)

type Envelope struct {
    XMLName         xml.Name    `xml:"env:Envelope"`
    EnvelopeAttrs
    Headers         *Headers    `xml:"env:Header,omitempty"`
    Body            *BodyStruct `xml:"env:Body,omitempty"`
}

type ShellParams struct {
    IStream     string
    OStream     string
    WorkingDir  string
    EnvVars     *Environment
    NoProfile   bool
    Codepage    string
}

type HeaderParams struct {
    ResourceURI         string
    Action              string
    ShellID             string
    MessageID           string
}

type CmdParams struct {
    ShellID     string
    Cmd         string
    Args        string
    Timeout     string
}

func (envelope *Envelope) RunCommand(shellParams ShellParams, params CmdParams, soap SoapRequest) (string, string, int, error){
    shell, err_shell := envelope.GetShell(shellParams, soap)
    if err_shell != nil {
        return "", "", 0, err_shell
    }
    params.ShellID = shell
    commID, commErr := envelope.SendCommand(params, soap)
    if commErr != nil {
        return "", "", 0, commErr
    }
    strdout, stderr, ret_code, err := envelope.GetCommandOutput(params.ShellID, commID, soap)
    if err != nil{
        return "", "", 0, err
    }
    err_clean := envelope.CleanupShell(params.ShellID, commID, soap)
    if err_clean != nil {
        return "", "", 0, err
    }
    err_close := envelope.CloseShell(params.ShellID, soap)
    if err_close != nil {
        return "", "", 0, err
    }
    return strdout, stderr, ret_code, err
}

// Generate SOAP envelope headers
func (envelope *Envelope) GetSoapHeaders(params HeaderParams) (error){
    envelope.Headers = &Headers{
        OperationTimeout:"PT60S",
        To:"http://windows-host:5985/wsman",
        ReplyTo:&ReplyAddress{
            ValueMustUnderstand{
                Value:"http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
                Attr:"true",
            },
        },
        DataLocale:&LocaleAttr{
            MustUnderstand:"true",
            Lang:"en-US",
        },
        Locale:&LocaleAttr{
            MustUnderstand:"true",
            Lang:"en-US",
        },
        MaxEnvelopeSize:&ValueMustUnderstand{
            Value:"153600",
            Attr:"true",
        },
    }

    if params.ResourceURI != "" {
        envelope.Headers.ResourceURI = &ValueMustUnderstand{
            Value:params.ResourceURI,
            Attr:"true",
        }
    }

    if params.Action != "" {
        envelope.Headers.Action = &ValueMustUnderstand{
            Value:params.Action,
            Attr:"true",
        }
    }

    if params.ShellID != "" {
        envelope.Headers.SelectorSet = &Selector{
            ValueName{
                Value:params.ShellID,
                Attr:"ShellId",
            },
        }
    }

    if params.MessageID == "" {
        uuid, err := uuid.NewV4()
        if err != nil{
            return err
        }
        params.MessageID = fmt.Sprintf("uuid:%s", uuid)
    }
    envelope.Headers.MessageID = params.MessageID
    return nil
}

// TODO: Do a soap request and return ShellID
func (envelope *Envelope) GetShell(params ShellParams, soap SoapRequest) (string, error){
    HeadParams := HeaderParams {
        ResourceURI: "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd",
        Action: "http://schemas.xmlsoap.org/ws/2004/09/transfer/Create",
    }
    // envelope := Envelope{}
    envelope.GetSoapHeaders(HeadParams)

    if params.Codepage == "" {
        params.Codepage = "437"
    }
    envelope.Headers.OptionSet = &OptionSet{
        []ValueName{
            ValueName{Attr:"WINRS_NOPROFILE", Value:"FALSE"},
            ValueName{Attr:"WINRS_CODEPAGE", Value:params.Codepage},
        },
    }
    var Body BodyStruct = BodyStruct{}
    var ShellVars Shell = Shell {}

    if params.IStream == "" {
        ShellVars.InputStreams = "stdin"
    }else{
        ShellVars.InputStreams = params.IStream
    }

    if params.OStream == "" {
        ShellVars.OutputStreams = "stdout stderr"
    }else{
        ShellVars.OutputStreams = params.OStream
    }

    if params.EnvVars != nil {
        ShellVars.Environment = params.EnvVars
    }

    // send request to WinRm
    Body.Shell = &ShellVars
    envelope.Body = &Body
    envelope.EnvelopeAttrs = Namespaces

    // response from WinRM
    resp, err := soap.SendMessage(envelope)
    if err != nil{
        return "", err
    }
    defer resp.Body.Close()

    respObj, err := GetObjectFromXML(resp.Body)
    if err != nil {
        return "", err
    }
    shellID := respObj.Body.Shell.ShellId

    return shellID, err
}

func (envelope *Envelope) SendCommand(params CmdParams, soap SoapRequest) (string, error){
    HeadParams := HeaderParams {
        ResourceURI: "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd",
        Action: "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Command",
    }
    if params.ShellID == "" {
        return "", errors.New("Invalid ShellId")
    }
    HeadParams.ShellID = params.ShellID
    envelope.GetSoapHeaders(HeadParams)

    envelope.Headers.OptionSet = &OptionSet{
        []ValueName{
            ValueName{Attr:"WINRS_CONSOLEMODE_STDIN", Value:"TRUE"},
            ValueName{Attr:"WINRS_SKIP_CMD_SHELL", Value:"FALSE"},
        },
    }

    if params.Timeout == "" {
        envelope.Headers.OperationTimeout = "PT3600S"
    }else{
        envelope.Headers.OperationTimeout = params.Timeout
    }
    // var Body BodyStruct = BodyStruct{}

    envelope.EnvelopeAttrs = Namespaces
    if params.Cmd == "" {
        return "", errors.New("Invalid command")
    }
    envelope.Body = &BodyStruct{
        CommandLine: &Command{
            Command: params.Cmd,
        },
    }

    if params.Args != "" {
        envelope.Body.CommandLine.Arguments = params.Args
    }

    // fmt.Printf("%s\n", output)
    resp, err := soap.SendMessage(envelope)
    if err != nil{
        return "", err
    }
    defer resp.Body.Close()

    respObj, err := GetObjectFromXML(resp.Body)
    if err != nil{
        return "", err
    }
    // contents, _ := ioutil.ReadAll(resp.Body)
    // fmt.Printf("REQ:%s\n\nRESP:%s\n\nSHELL:%s\n\n", output, contents, shellID)
    return respObj.Body.CommandResponse.CommandId, nil
}

func (envelope *Envelope) GetCommandOutput(shellID, commandID string, soap SoapRequest) (string, string, int, error){
    HeadParams := HeaderParams {
        ResourceURI: "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd",
        Action: "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Receive",
        ShellID: shellID,
    }
    envelope.GetSoapHeaders(HeadParams)
    envelope.EnvelopeAttrs = Namespaces
    envelope.Body = &BodyStruct{
        Receive: &Receive{
            DesiredStream: DesiredStreamProps{
                Value: "stdout stderr",
                Attr: commandID,
            },
        },
    }

    resp, err := soap.SendMessage(envelope)
    if err != nil{
        return "", "", 0, err
    }
    defer resp.Body.Close()

    stdout, stderr, retCode := ParseCommandOutput(resp.Body)
    // fmt.Printf("%s\n", output)
    return stdout, stderr, retCode, nil
}

func (envelope *Envelope) CleanupShell(shellID, commandID string, soap SoapRequest) (error){
    HeadParams := HeaderParams {
        ResourceURI: "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd",
        Action: "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Signal",
        ShellID: shellID,
    }
    envelope.GetSoapHeaders(HeadParams)
    envelope.EnvelopeAttrs = Namespaces
    sig := Signal{
        Attr: commandID,
        Code: "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/signal/terminate",
    }
    envelope.Body = &BodyStruct{
        Signal: &sig,
    }

    resp, err := soap.SendMessage(envelope)
    if err != nil{
        return err
    }
    defer resp.Body.Close()
    // contents, err2 := ioutil.ReadAll(resp.Body)
    // fmt.Printf("%s --> %s", contents, err2)
    if resp.StatusCode != 200 {
        return errors.New(fmt.Sprintf("Remote host returned error status code: %d", resp.StatusCode))
    }
    return nil
}

func (envelope *Envelope) CloseShell(shellID string, soap SoapRequest) (error){
    HeadParams := HeaderParams {
        ResourceURI: "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd",
        Action: "http://schemas.xmlsoap.org/ws/2004/09/transfer/Delete",
        ShellID: shellID,
    }
    var Body BodyStruct = BodyStruct{}

    envelope.EnvelopeAttrs = Namespaces
    envelope.GetSoapHeaders(HeadParams)
    envelope.Body = &Body

    resp, err := soap.SendMessage(envelope)
    // contents, err := ioutil.ReadAll(resp.Body)
    // fmt.Printf("REQ:%s\n\nRESP:%s\n\nSHELL:%s\n\n", output, contents, shellID)
    if err != nil{
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        return errors.New("Remote host returned error status code")
    }
    return nil
}