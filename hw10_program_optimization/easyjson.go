package hw10programoptimization

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyJsonDecode(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Email":
			out.Email = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func easyJsonEncode(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Email\":"
		out.RawString(prefix[1:])
		out.String(string(in.Email))
	}
	out.RawByte('}')
}

func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyJsonEncode(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyJsonEncode(w, v)
}

func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyJsonDecode(&r, v)
	return r.Error()
}

func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyJsonDecode(l, v)
}
