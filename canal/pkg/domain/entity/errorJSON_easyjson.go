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

func easyjsonCdb719ceDecodeServerCanalPkgDomainEntity(in *jlexer.Lexer, out *ErrorJSON) {
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
		case "name":
			if in.IsNull() {
				in.Skip()
				out.Name = nil
			} else {
				in.Delim('[')
				if out.Name == nil {
					if !in.IsDelim(']') {
						out.Name = make([]string, 0, 4)
					} else {
						out.Name = []string{}
					}
				} else {
					out.Name = (out.Name)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Name = append(out.Name, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "email":
			if in.IsNull() {
				in.Skip()
				out.Email = nil
			} else {
				in.Delim('[')
				if out.Email == nil {
					if !in.IsDelim(']') {
						out.Email = make([]string, 0, 4)
					} else {
						out.Email = []string{}
					}
				} else {
					out.Email = (out.Email)[:0]
				}
				for !in.IsDelim(']') {
					var v2 string
					v2 = string(in.String())
					out.Email = append(out.Email, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "password":
			if in.IsNull() {
				in.Skip()
				out.Password = nil
			} else {
				in.Delim('[')
				if out.Password == nil {
					if !in.IsDelim(']') {
						out.Password = make([]string, 0, 4)
					} else {
						out.Password = []string{}
					}
				} else {
					out.Password = (out.Password)[:0]
				}
				for !in.IsDelim(']') {
					var v3 string
					v3 = string(in.String())
					out.Password = append(out.Password, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "oldPassword":
			if in.IsNull() {
				in.Skip()
				out.OldPassword = nil
			} else {
				in.Delim('[')
				if out.OldPassword == nil {
					if !in.IsDelim(']') {
						out.OldPassword = make([]string, 0, 4)
					} else {
						out.OldPassword = []string{}
					}
				} else {
					out.OldPassword = (out.OldPassword)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.OldPassword = append(out.OldPassword, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "non_field_errors":
			if in.IsNull() {
				in.Skip()
				out.NonFieldError = nil
			} else {
				in.Delim('[')
				if out.NonFieldError == nil {
					if !in.IsDelim(']') {
						out.NonFieldError = make([]string, 0, 4)
					} else {
						out.NonFieldError = []string{}
					}
				} else {
					out.NonFieldError = (out.NonFieldError)[:0]
				}
				for !in.IsDelim(']') {
					var v5 string
					v5 = string(in.String())
					out.NonFieldError = append(out.NonFieldError, v5)
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
func easyjsonCdb719ceEncodeServerCanalPkgDomainEntity(out *jwriter.Writer, in ErrorJSON) {
	out.RawByte('{')
	first := true
	_ = first
	if len(in.Name) != 0 {
		const prefix string = ",\"name\":"
		first = false
		out.RawString(prefix[1:])
		{
			out.RawByte('[')
			for v6, v7 := range in.Name {
				if v6 > 0 {
					out.RawByte(',')
				}
				out.String(string(v7))
			}
			out.RawByte(']')
		}
	}
	if len(in.Email) != 0 {
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v8, v9 := range in.Email {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.String(string(v9))
			}
			out.RawByte(']')
		}
	}
	if len(in.Password) != 0 {
		const prefix string = ",\"password\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v10, v11 := range in.Password {
				if v10 > 0 {
					out.RawByte(',')
				}
				out.String(string(v11))
			}
			out.RawByte(']')
		}
	}
	if len(in.OldPassword) != 0 {
		const prefix string = ",\"oldPassword\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v12, v13 := range in.OldPassword {
				if v12 > 0 {
					out.RawByte(',')
				}
				out.String(string(v13))
			}
			out.RawByte(']')
		}
	}
	if len(in.NonFieldError) != 0 {
		const prefix string = ",\"non_field_errors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v14, v15 := range in.NonFieldError {
				if v14 > 0 {
					out.RawByte(',')
				}
				out.String(string(v15))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ErrorJSON) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonCdb719ceEncodeServerCanalPkgDomainEntity(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ErrorJSON) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonCdb719ceEncodeServerCanalPkgDomainEntity(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ErrorJSON) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonCdb719ceDecodeServerCanalPkgDomainEntity(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ErrorJSON) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonCdb719ceDecodeServerCanalPkgDomainEntity(l, v)
}
