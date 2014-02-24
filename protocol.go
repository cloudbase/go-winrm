package winrm

import (
    "encoding/xml"
    // "io/ioutil"
    "fmt"
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

// Generate SOAP envelope headers
func (envelope *Envelope) GetSoapHeaders(params HeaderParams){
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
                Attr:"ShellID",
            },
        }
    }

    if params.MessageID == "" {
        uuid, err := uuid.NewV4()
        if err != nil{
            fmt.Printf("Error: %v\n", err)
        }
        params.MessageID = fmt.Sprintf("uuid:%s", uuid)
    }
    envelope.Headers.MessageID = params.MessageID
}

// TODO: Do a soap request and return ShellID
func (envelope *Envelope) GetShell(params ShellParams, soap SoapRequest) (*string, error){
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
    output, err := xml.MarshalIndent(envelope, "", "")
    if err != nil {
        fmt.Printf("error: %v\n", err)
    }
    // response from WinRM
    resp, err := soap.SendMessage(output)
    defer resp.Body.Close()
    if err != nil{
        fmt.Errorf("err")
        return nil, err
    }

    respObj, err := GetObjectFromXML(resp.Body)
    if err != nil {
        return nil, err
    }
    shellID := &respObj.Body.Shell.ShellId

    return shellID, err
}
