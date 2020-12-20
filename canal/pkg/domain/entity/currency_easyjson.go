// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package entity

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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
		case "value":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Value).UnmarshalJSON(data))
			}
		case "updated_at":
			if in.IsNull() {
				in.Skip()
				out.UpdatedAt = nil
			} else {
				if out.UpdatedAt == nil {
					out.UpdatedAt = new(timestamppb.Timestamp)
				}
				easyjsonE5a98965DecodeGoogleGolangOrgProtobufTypesKnownTimestamppb(in, out.UpdatedAt)
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
func easyjsonE5a98965EncodeServerCanalPkgDomainEntity(out *jwriter.Writer, in Currency) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Base != "" {
		const prefix string = ",\"base\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Base))
	}
	{
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"value\":"
		out.RawString(prefix)
		out.Raw((in.Value).MarshalJSON())
	}
	if in.UpdatedAt != nil {
		const prefix string = ",\"updated_at\":"
		out.RawString(prefix)
		easyjsonE5a98965EncodeGoogleGolangOrgProtobufTypesKnownTimestamppb(out, *in.UpdatedAt)
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
func easyjsonE5a98965DecodeGoogleGolangOrgProtobufTypesKnownTimestamppb(in *jlexer.Lexer, out *timestamppb.Timestamp) {
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
		case "seconds":
			out.Seconds = int64(in.Int64())
		case "nanos":
			out.Nanos = int32(in.Int32())
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
func easyjsonE5a98965EncodeGoogleGolangOrgProtobufTypesKnownTimestamppb(out *jwriter.Writer, in timestamppb.Timestamp) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Seconds != 0 {
		const prefix string = ",\"seconds\":"
		first = false
		out.RawString(prefix[1:])
		out.Int64(int64(in.Seconds))
	}
	if in.Nanos != 0 {
		const prefix string = ",\"nanos\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.Nanos))
	}
	out.RawByte('}')
}
func easyjsonE5a98965DecodeServerCanalPkgDomainEntity1(in *jlexer.Lexer, out *Currencies) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Currencies, 0, 1)
			} else {
				*out = Currencies{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Currency
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
func easyjsonE5a98965EncodeServerCanalPkgDomainEntity1(out *jwriter.Writer, in Currencies) {
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
