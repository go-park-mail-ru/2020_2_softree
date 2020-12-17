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

func easyjson22b96abDecodeServerCanalPkgDomainEntity(in *jlexer.Lexer, out *Wallets) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Wallets, 0, 1)
			} else {
				*out = Wallets{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Wallet
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson22b96abEncodeServerCanalPkgDomainEntity(out *jwriter.Writer, in Wallets) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Wallets) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson22b96abEncodeServerCanalPkgDomainEntity(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Wallets) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson22b96abEncodeServerCanalPkgDomainEntity(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Wallets) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson22b96abDecodeServerCanalPkgDomainEntity(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Wallets) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson22b96abDecodeServerCanalPkgDomainEntity(l, v)
}
func easyjson22b96abDecodeServerCanalPkgDomainEntity1(in *jlexer.Lexer, out *Wallet) {
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
		case "title":
			out.Title = string(in.String())
		case "value":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Value).UnmarshalJSON(data))
			}
		case "UserId":
			out.UserId = int64(in.Int64())
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
func easyjson22b96abEncodeServerCanalPkgDomainEntity1(out *jwriter.Writer, in Wallet) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix[1:])
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"value\":"
		out.RawString(prefix)
		out.Raw((in.Value).MarshalJSON())
	}
	{
		const prefix string = ",\"UserId\":"
		out.RawString(prefix)
		out.Int64(int64(in.UserId))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Wallet) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson22b96abEncodeServerCanalPkgDomainEntity1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Wallet) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson22b96abEncodeServerCanalPkgDomainEntity1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Wallet) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson22b96abDecodeServerCanalPkgDomainEntity1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Wallet) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson22b96abDecodeServerCanalPkgDomainEntity1(l, v)
}
