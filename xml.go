package main

import (
    "fmt"
    "encoding/xml"
    "os"
    )

type EnvelopeAttrs struct {
    Xsd       string    `xml:"xmlns:xsd,attr"`
    Xsi       string    `xml:"xmlns:xsi,attr"`
    Rsp       string    `xml:"xmlns:rsp,attr"`
    P         string    `xml:"xmlns:p,attr"`
    W         string    `xml:"xmlns:w,attr"`
    X         string    `xml:"xmlns:x,attr"`
    A         string    `xml:"xmlns:a,attr"`
    B         string    `xml:"xmlns:b,attr"`
    Env       string    `xml:"xmlns:env,attr"`
    Cfg       string    `xml:"xmlns:cfg,attr"`
    N         string    `xml:"xmlns:n,attr"`
}

var Namespaces EnvelopeAttrs = EnvelopeAttrs{
    Xsd:"http://www.w3.org/2001/XMLSchema",
    Xsi:"http://www.w3.org/2001/XMLSchema-instance",
    Rsp:"http://schemas.microsoft.com/wbem/wsman/1/windows/shell",
    P:"http://schemas.microsoft.com/wbem/wsman/1/wsman.xsd",
    W:"http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
    X:"http://schemas.xmlsoap.org/ws/2004/09/transfer",
    A:"http://schemas.xmlsoap.org/ws/2004/08/addressing",
    B:"http://schemas.dmtf.org/wbem/wsman/1/cimbinding.xsd",
    Env:"http://www.w3.org/2003/05/soap-envelope",
    Cfg:"http://schemas.microsoft.com/wbem/wsman/1/config",
    N:"http://schemas.xmlsoap.org/ws/2004/09/enumeration",
}

type ValueName struct {
    Value       string  `xml:",innerxml"`
    Attr        string  `xml:"Name,attr"`
}

type OptionSet struct {
    Option      []ValueName  `xml:"w:Option"`
}

type ValueMustUnderstand struct {
    Value     string  `xml:",innerxml"`
    Attr      string  `xml:"mustUnderstand,attr"`
}

type LocaleAttr struct {
    MustUnderstand  string  `xml:"mustUnderstand,attr"`
    Lang            string  `xml:"xml:lang,attr"`   
}

type ReplyAddress struct {
    Address     ValueMustUnderstand  `xml:"Address"`
}

type Headers struct {
    To                  string                  `xml:"env:Header>a:To"`
    OptionSet           *OptionSet              `xml:"env:Header>w:OptionSet,omitempty"`
    ReplyTo             *ReplyAddress           `xml:"env:Header>a:ReplyTo,omitempty"`
    MaxEnvelopeSize     *ValueMustUnderstand    `xml:"env:Header>w:MaxEnvelopeSize,omitempty"`
    MessageID           string                  `xml:"env:Header>a:MessageID"`
    Locale              *LocaleAttr             `xml:"env:Header>p:Locale,omitempty"`
    DataLocale          *LocaleAttr             `xml:"env:Header>p:DataLocale,omitempty"`
    OperationTimeout    string                  `xml:"env:Header>w:OperationTimeout"`
    ResourceURI         *ValueMustUnderstand    `xml:"env:Header>w:ResourceURI,omitempty"`
    Action              *ValueMustUnderstand    `xml:"env:Header>w:Action,omitempty"`
    SelectorSet         *ValueName              `xml:"env:Header>w:SelectorSet,omitempty"`
}

type Command struct {
    Command     string  `xml:"rsp:Command"`
    Arguments   string  `xml:"rsp:Arguments,omitempty"`

}

type DesiredStreamProps struct{
    Value       string  `xml:",innerxml"`
    Attr        string  `xml:"CommandId,attr"`
}

type Receive struct {
    DesiredStream   DesiredStreamProps  `xml:"rsp:DesiredStream"`
}

type Signal struct {
    Attr    string  `xml:"CommandId,attr"`
    Code    string  `xml:"rsp:Code"`
}

type Shell struct{
    InputStreams    string  `xml:"rsp:InputStreams,omitempty"`
    OutputStreams   string  `xml:"rsp:OutputStreams,omitempty"`
}

type BodyStruct struct {
    CommandLine     *Command    `xml:"rsp:CommandLine,omitempty"`
    Receive         *Receive    `xml:"rsp:Receive,omitempty"`
    Signal          *Signal     `xml:"rsp:Signal,omitempty"`
    Shell           *Shell      `xml:"rsp:Shell"`
}

var Body BodyStruct = BodyStruct{
    // CommandLine:&Command{
    //     Command:"dir",
    //     // Arguments:"C:\\",
    // },
    // Receive:&Receive{
    //     DesiredStream:DesiredStreamProps{
    //         Value:"stdout stderr",
    //         Attr:"sfsdfds",
    //     },
    // },
    Shell:&Shell{
        InputStreams:"stdin",
        OutputStreams:"stdout stderr",
    },
}

var Head Headers = Headers {
    OperationTimeout:"PT60S",
    To:"http://windows-host:5985/wsman",
    MessageID:"uuid:safgfdh",
    OptionSet:&OptionSet{
            []ValueName{
                ValueName{Attr:"WINRS_NOPROFILE", Value:"FALSE"},
                ValueName{Attr:"WINRS_CODEPAGE", Value:"437"},
            },
    },
    SelectorSet:&ValueName{Value:"dskjfgf", Attr:"ShellID"},
    ReplyTo:&ReplyAddress{
        ValueMustUnderstand{
            Value:"http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
            Attr:"true",
        },
    },
    MaxEnvelopeSize:&ValueMustUnderstand{
        Value:"153600",
        Attr:"true",
    },
    ResourceURI:&ValueMustUnderstand{
        Value:"http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd",
        Attr:"true",
    },
    DataLocale:&LocaleAttr{
        MustUnderstand:"true",
        Lang:"en-US",
    },
    Locale:&LocaleAttr{
        MustUnderstand:"true",
        Lang:"en-US",
    },
    Action:&ValueMustUnderstand{
        Value:"http://schemas.xmlsoap.org/ws/2004/09/transfer/Delete",
        Attr:"true",
    },
}


type TestEnv struct {
    XMLName   xml.Name  `xml:"env:Envelope"`
    EnvelopeAttrs
    Headers     *Headers    `xml:"env:Header,omitempty"`
    Body        *BodyStruct `xml:"Body,omitempty"`
}

func main(){
    v := &TestEnv{EnvelopeAttrs:Namespaces, Headers:&Head, Body:&Body}
    output, err := xml.MarshalIndent(v, "  ", "    ")
    if err != nil {
        fmt.Printf("error: %v\n", err)
    }
    os.Stdout.Write(output)
}