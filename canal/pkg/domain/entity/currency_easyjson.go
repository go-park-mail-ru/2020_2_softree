// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package entity

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonE5a98965DecodeServerCanalPkgDomainEntity(in *jlexer.Lexer, out *Currency) {
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
		case "base":
			out.Base = string(in.String())
		case "title":
			out.Title = string(in.String())
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
func easyjsonE5a98965EncodeServerCanalPkgDomainEntity(out *jwriter.Writer, in Currency) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"base\":"
		out.RawString(prefix[1:])
		out.String(string(in.Base))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Currency) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE5a98965EncodeServerCanalPkgDomainEntity(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Currency) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE5a98965EncodeServerCanalPkgDomainEntity(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Currency) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE5a98965DecodeServerCanalPkgDomainEntity(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Currency) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE5a98965DecodeServerCanalPkgDomainEntity(l, v)
}
func easyjsonE5a98965DecodeServerCanalPkgDomainEntity1(in *jlexer.Lexer, out *Currencies) {
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
		case "Currencies":
			if in.IsNull() {
				in.Skip()
				out.Currencies = nil
			} else {
				in.Delim('[')
				if out.Currencies == nil {
					if !in.IsDelim(']') {
						out.Currencies = make([]Currency, 0, 2)
					} else {
						out.Currencies = []Currency{}
					}
				} else {
					out.Currencies = (out.Currencies)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Currency
					(v1).UnmarshalEasyJSON(in)
					out.Currencies = append(out.Currencies, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjsonE5a98965EncodeServerCanalPkgDomainEntity1(out *jwriter.Writer, in Currencies) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Currencies\":"
		out.RawString(prefix[1:])
		if in.Currencies == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Currencies {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Currencies) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE5a98965EncodeServerCanalPkgDomainEntity1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Currencies) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE5a98965EncodeServerCanalPkgDomainEntity1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Currencies) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE5a98965DecodeServerCanalPkgDomainEntity1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Currencies) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE5a98965DecodeServerCanalPkgDomainEntity1(l, v)
}
