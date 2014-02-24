package winrm

type EnvelopeAttrs struct {
    Xsd       string    `xml:"xmlns:xsd,attr,omitempty"`
    Xsi       string    `xml:"xmlns:xsi,attr,omitempty"`
    Rsp       string    `xml:"xmlns:rsp,attr,omitempty"`
    P         string    `xml:"xmlns:p,attr,omitempty"`
    W         string    `xml:"xmlns:w,attr,omitempty"`
    X         string    `xml:"xmlns:x,attr,omitempty"`
    A         string    `xml:"xmlns:a,attr,omitempty"`
    B         string    `xml:"xmlns:b,attr,omitempty"`
    Env       string    `xml:"xmlns:env,attr,omitempty"`
    Cfg       string    `xml:"xmlns:cfg,attr,omitempty"`
    N         string    `xml:"xmlns:n,attr,omitempty"`
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
    Address     ValueMustUnderstand  `xml:"a:Address"`
}

type Selector struct {
    Set     ValueName  `xml:"Selector"`
}

type Headers struct {
    To                  string                  `xml:"a:To"`
    OptionSet           *OptionSet              `xml:"w:OptionSet,omitempty"`
    ReplyTo             *ReplyAddress           `xml:"a:ReplyTo,omitempty"`
    MaxEnvelopeSize     *ValueMustUnderstand    `xml:"w:MaxEnvelopeSize,omitempty"`
    MessageID           string                  `xml:"a:MessageID,omitempty"`
    Locale              *LocaleAttr             `xml:"p:Locale,omitempty"`
    DataLocale          *LocaleAttr             `xml:"p:DataLocale,omitempty"`
    OperationTimeout    string                  `xml:"w:OperationTimeout"`
    ResourceURI         *ValueMustUnderstand    `xml:"w:ResourceURI,omitempty"`
    Action              *ValueMustUnderstand    `xml:"a:Action,omitempty"`
    SelectorSet         *Selector               `xml:"w:SelectorSet,omitempty"`
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

type EnvVariable struct {
    Value   string  `xml:",innerxml"`
    Name    string  `xml:"Name,attr"`
}

type Environment struct {
    Variable    []EnvVariable  `xml:"rsp:Variable"`
}

type Shell struct{
    InputStreams        string          `xml:"rsp:InputStreams,omitempty"`
    OutputStreams       string          `xml:"rsp:OutputStreams,omitempty"`
    WorkingDirectory    string          `xml:"rsp:WorkingDirectory,omitempty"`
    IdleTimeOut         string          `xml:"rsp:IdleTimeOut,omitempty"`
    Environment         *Environment    `xml:"rsp:Environment,omitempty"`
}

type BodyStruct struct {
    CommandLine     *Command    `xml:"rsp:CommandLine,omitempty"`
    Receive         *Receive    `xml:"rsp:Receive,omitempty"`
    Signal          *Signal     `xml:"rsp:Signal,omitempty"`
    Shell           *Shell      `xml:"rsp:Shell"`
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
