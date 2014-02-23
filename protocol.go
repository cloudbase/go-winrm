package winrm

import (
    "encoding/xml"
    "fmt"
    "github.com/nu7hatch/gouuid"
)

type Envelope struct {
    XMLName         xml.Name  `xml:"env:Envelope"`
    EnvelopeAttrs
    Headers         *Headers    `xml:"env:Header,omitempty"`
    Body            *BodyStruct `xml:"Body,omitempty"`
}

type HeaderParams struct {
    ResourceURI         string
    Action              string
    ShellID             string
    MessageID           string
}

func GetSoapHeaders(params HeaderParams) Headers{
    var Head Headers = Headers {
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
        Head.ResourceURI = &ValueMustUnderstand{
            Value:params.ResourceURI,
            Attr:"true",
        }
    }

    if params.Action != "" {
        Head.Action = &ValueMustUnderstand{
            Value:params.Action,
            Attr:"true",
        }
    }

    if params.ShellID != "" {
        Head.SelectorSet = &Selector{
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
    Head.MessageID = params.MessageID
    
    return Head
}

type ShellParams struct {
    IStream     string
    OStream     string
    WorkingDir  string
    EnvVars     *Environment
    NoProfile   bool
    Codepage    string
}

// TODO: Do a soap request and return ShellID
func GetShell(params ShellParams) []byte {
    HeadParams := HeaderParams {
        ResourceURI: "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd",
        Action: "http://schemas.xmlsoap.org/ws/2004/09/transfer/Create",
    }
    var Head Headers = GetSoapHeaders(HeadParams)

    if params.Codepage == "" {
        params.Codepage = "437"
    }
    Head.OptionSet = &OptionSet{
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

    Body.Shell = &ShellVars
    v := &Envelope{EnvelopeAttrs:Namespaces, Headers:&Head, Body:&Body}
    output, err := xml.MarshalIndent(v, "  ", "    ")
    if err != nil {
        fmt.Printf("error: %v\n", err)
    }
    return output
}