package reflect

import (
	json "encoding/json"
	"reflect"

	"github.com/lab47/lace/pkg/pkgreflect"
)

type UnmarshalerImpl struct {
	UnmarshalJSONFn func([]byte) error
}

func (s *UnmarshalerImpl) UnmarshalJSON(a0 []byte) error {
	return s.UnmarshalJSONFn(a0)
}

type MarshalerImpl struct {
	MarshalJSONFn func() ([]byte, error)
}

func (s *MarshalerImpl) MarshalJSON() ([]byte, error) {
	return s.MarshalJSONFn()
}

func init() {
	InvalidUnmarshalError_methods := map[string]pkgreflect.Func{}
	Number_methods := map[string]pkgreflect.Func{}
	UnmarshalFieldError_methods := map[string]pkgreflect.Func{}
	UnmarshalTypeError_methods := map[string]pkgreflect.Func{}
	Unmarshaler_methods := map[string]pkgreflect.Func{}
	InvalidUTF8Error_methods := map[string]pkgreflect.Func{}
	Marshaler_methods := map[string]pkgreflect.Func{}
	MarshalerError_methods := map[string]pkgreflect.Func{}
	UnsupportedTypeError_methods := map[string]pkgreflect.Func{}
	UnsupportedValueError_methods := map[string]pkgreflect.Func{}
	SyntaxError_methods := map[string]pkgreflect.Func{}
	Decoder_methods := map[string]pkgreflect.Func{}
	Delim_methods := map[string]pkgreflect.Func{}
	Encoder_methods := map[string]pkgreflect.Func{}
	RawMessage_methods := map[string]pkgreflect.Func{}
	Token_methods := map[string]pkgreflect.Func{}
	UnmarshalTypeError_methods["Error"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: ""}
	UnmarshalFieldError_methods["Error"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: ""}
	InvalidUnmarshalError_methods["Error"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: ""}
	Number_methods["String"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: "String returns the literal text of the number."}
	Number_methods["Float64"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "Float64 returns the number as a float64."}
	Number_methods["Int64"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "Int64 returns the number as an int64."}
	UnsupportedTypeError_methods["Error"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: ""}
	UnsupportedValueError_methods["Error"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: ""}
	InvalidUTF8Error_methods["Error"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: ""}
	MarshalerError_methods["Error"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: ""}
	MarshalerError_methods["Unwrap"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "error", Doc: "Unwrap returns the underlying error."}
	SyntaxError_methods["Error"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: ""}
	Decoder_methods["UseNumber"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "UseNumber causes the Decoder to unmarshal a number into an interface{} as a\n[Number] instead of as a float64."}
	Decoder_methods["DisallowUnknownFields"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "DisallowUnknownFields causes the Decoder to return an error when the destination\nis a struct and the input contains object keys which do not match any\nnon-ignored, exported fields in the destination."}
	Decoder_methods["Decode"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "v", Tag: "any"}}, Tag: "error", Doc: "Decode reads the next JSON-encoded value from its\ninput and stores it in the value pointed to by v.\n\nSee the documentation for [Unmarshal] for details about\nthe conversion of JSON into a Go value."}
	Decoder_methods["Buffered"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "io.Reader", Doc: "Buffered returns a reader of the data remaining in the Decoder's\nbuffer. The reader is valid until the next call to [Decoder.Decode]."}
	Encoder_methods["Encode"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "v", Tag: "any"}}, Tag: "error", Doc: "Encode writes the JSON encoding of v to the stream,\nfollowed by a newline character.\n\nSee the documentation for [Marshal] for details about the\nconversion of Go values to JSON."}
	Encoder_methods["SetIndent"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "prefix", Tag: "string"}, {Name: "indent", Tag: "string"}}, Tag: "any", Doc: "SetIndent instructs the encoder to format each subsequent encoded\nvalue as if indented by the package-level function Indent(dst, src, prefix, indent).\nCalling SetIndent(\"\", \"\") disables indentation."}
	Encoder_methods["SetEscapeHTML"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "on", Tag: "bool"}}, Tag: "any", Doc: "SetEscapeHTML specifies whether problematic HTML characters\nshould be escaped inside JSON quoted strings.\nThe default behavior is to escape &, <, and > to \\u0026, \\u003c, and \\u003e\nto avoid certain safety problems that can arise when embedding JSON in HTML.\n\nIn non-HTML settings where the escaping interferes with the readability\nof the output, SetEscapeHTML(false) disables this behavior."}
	RawMessage_methods["MarshalJSON"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "MarshalJSON returns m as the JSON encoding of m."}
	RawMessage_methods["UnmarshalJSON"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "data", Tag: "[]byte"}}, Tag: "error", Doc: "UnmarshalJSON sets *m to a copy of data."}
	Delim_methods["String"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: ""}
	Decoder_methods["Token"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "Token returns the next JSON token in the input stream.\nAt the end of the input stream, Token returns nil, [io.EOF].\n\nToken guarantees that the delimiters [ ] { } it returns are\nproperly nested and matched: if Token encounters an unexpected\ndelimiter in the input, it will return an error.\n\nThe input stream consists of basic JSON values—bool, string,\nnumber, and null—along with delimiters [ ] { } of type [Delim]\nto mark the start and end of arrays and objects.\nCommas and colons are elided."}
	Decoder_methods["More"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "bool", Doc: "More reports whether there is another element in the\ncurrent array or object being parsed."}
	Decoder_methods["InputOffset"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int64", Doc: "InputOffset returns the input stream byte offset of the current decoder position.\nThe offset gives the location of the end of the most recently returned token\nand the beginning of the next token."}
	pkgreflect.AddPackage("encoding/json", &pkgreflect.Package{
		Doc: "Package json implements encoding and decoding of JSON as defined in RFC 7159.",
		Types: map[string]pkgreflect.Type{
			"Decoder":               {Doc: "", Value: reflect.TypeOf((*json.Decoder)(nil)).Elem(), Methods: Decoder_methods},
			"Delim":                 {Doc: "", Value: reflect.TypeOf((*json.Delim)(nil)).Elem(), Methods: Delim_methods},
			"Encoder":               {Doc: "", Value: reflect.TypeOf((*json.Encoder)(nil)).Elem(), Methods: Encoder_methods},
			"InvalidUTF8Error":      {Doc: "", Value: reflect.TypeOf((*json.InvalidUTF8Error)(nil)).Elem(), Methods: InvalidUTF8Error_methods},
			"InvalidUnmarshalError": {Doc: "", Value: reflect.TypeOf((*json.InvalidUnmarshalError)(nil)).Elem(), Methods: InvalidUnmarshalError_methods},
			"Marshaler":             {Doc: "", Value: reflect.TypeOf((*json.Marshaler)(nil)).Elem(), Methods: Marshaler_methods},
			"MarshalerError":        {Doc: "", Value: reflect.TypeOf((*json.MarshalerError)(nil)).Elem(), Methods: MarshalerError_methods},
			"Number":                {Doc: "", Value: reflect.TypeOf((*json.Number)(nil)).Elem(), Methods: Number_methods},
			"RawMessage":            {Doc: "", Value: reflect.TypeOf((*json.RawMessage)(nil)).Elem(), Methods: RawMessage_methods},
			"SyntaxError":           {Doc: "", Value: reflect.TypeOf((*json.SyntaxError)(nil)).Elem(), Methods: SyntaxError_methods},
			"Token":                 {Doc: "", Value: reflect.TypeOf((*json.Token)(nil)).Elem(), Methods: Token_methods},
			"UnmarshalFieldError":   {Doc: "", Value: reflect.TypeOf((*json.UnmarshalFieldError)(nil)).Elem(), Methods: UnmarshalFieldError_methods},
			"UnmarshalTypeError":    {Doc: "", Value: reflect.TypeOf((*json.UnmarshalTypeError)(nil)).Elem(), Methods: UnmarshalTypeError_methods},
			"Unmarshaler":           {Doc: "", Value: reflect.TypeOf((*json.Unmarshaler)(nil)).Elem(), Methods: Unmarshaler_methods},
			"UnsupportedTypeError":  {Doc: "", Value: reflect.TypeOf((*json.UnsupportedTypeError)(nil)).Elem(), Methods: UnsupportedTypeError_methods},
			"UnsupportedValueError": {Doc: "", Value: reflect.TypeOf((*json.UnsupportedValueError)(nil)).Elem(), Methods: UnsupportedValueError_methods},
			"UnmarshalerImpl":       {Doc: `Struct version of interface Unmarshaler for implementation`, Value: reflect.TypeFor[UnmarshalerImpl]()},
			"MarshalerImpl":         {Doc: `Struct version of interface Marshaler for implementation`, Value: reflect.TypeFor[MarshalerImpl]()},
		},

		Functions: map[string]pkgreflect.FuncValue{
			"Compact": {Doc: "Compact appends to dst the JSON-encoded src with\ninsignificant space characters elided.", Args: []pkgreflect.Arg{{Name: "dst", Tag: "bytes.Buffer"}, {Name: "src", Tag: "[]byte"}}, Tag: "error", Value: reflect.ValueOf(json.Compact)},

			"HTMLEscape": {Doc: "HTMLEscape appends to dst the JSON-encoded src with <, >, &, U+2028 and U+2029\ncharacters inside string literals changed to \\u003c, \\u003e, \\u0026, \\u2028, \\u2029\nso that the JSON will be safe to embed inside HTML <script> tags.\nFor historical reasons, web browsers don't honor standard HTML\nescaping within <script> tags, so an alternative JSON encoding must be used.", Args: []pkgreflect.Arg{{Name: "dst", Tag: "bytes.Buffer"}, {Name: "src", Tag: "[]byte"}}, Tag: "any", Value: reflect.ValueOf(json.HTMLEscape)},

			"Indent": {Doc: "Indent appends to dst an indented form of the JSON-encoded src.\nEach element in a JSON object or array begins on a new,\nindented line beginning with prefix followed by one or more\ncopies of indent according to the indentation nesting.\nThe data appended to dst does not begin with the prefix nor\nany indentation, to make it easier to embed inside other formatted JSON data.\nAlthough leading space characters (space, tab, carriage return, newline)\nat the beginning of src are dropped, trailing space characters\nat the end of src are preserved and copied to dst.\nFor example, if src has no trailing spaces, neither will dst;\nif src ends in a trailing newline, so will dst.", Args: []pkgreflect.Arg{{Name: "dst", Tag: "bytes.Buffer"}, {Name: "src", Tag: "[]byte"}, {Name: "prefix", Tag: "string"}, {Name: "indent", Tag: "string"}}, Tag: "error", Value: reflect.ValueOf(json.Indent)},

			"Marshal": {Doc: "Marshal returns the JSON encoding of v.\n\nMarshal traverses the value v recursively.\nIf an encountered value implements [Marshaler]\nand is not a nil pointer, Marshal calls [Marshaler.MarshalJSON]\nto produce JSON. If no [Marshaler.MarshalJSON] method is present but the\nvalue implements [encoding.TextMarshaler] instead, Marshal calls\n[encoding.TextMarshaler.MarshalText] and encodes the result as a JSON string.\nThe nil pointer exception is not strictly necessary\nbut mimics a similar, necessary exception in the behavior of\n[Unmarshaler.UnmarshalJSON].\n\nOtherwise, Marshal uses the following type-dependent default encodings:\n\nBoolean values encode as JSON booleans.\n\nFloating point, integer, and [Number] values encode as JSON numbers.\nNaN and +/-Inf values will return an [UnsupportedValueError].\n\nString values encode as JSON strings coerced to valid UTF-8,\nreplacing invalid bytes with the Unicode replacement rune.\nSo that the JSON will be safe to embed inside HTML <script> tags,\nthe string is encoded using [HTMLEscape],\nwhich replaces \"<\", \">\", \"&\", U+2028, and U+2029 are escaped\nto \"\\u003c\",\"\\u003e\", \"\\u0026\", \"\\u2028\", and \"\\u2029\".\nThis replacement can be disabled when using an [Encoder],\nby calling [Encoder.SetEscapeHTML](false).\n\nArray and slice values encode as JSON arrays, except that\n[]byte encodes as a base64-encoded string, and a nil slice\nencodes as the null JSON value.\n\nStruct values encode as JSON objects.\nEach exported struct field becomes a member of the object, using the\nfield name as the object key, unless the field is omitted for one of the\nreasons given below.\n\nThe encoding of each struct field can be customized by the format string\nstored under the \"json\" key in the struct field's tag.\nThe format string gives the name of the field, possibly followed by a\ncomma-separated list of options. The name may be empty in order to\nspecify options without overriding the default field name.\n\nThe \"omitempty\" option specifies that the field should be omitted\nfrom the encoding if the field has an empty value, defined as\nfalse, 0, a nil pointer, a nil interface value, and any empty array,\nslice, map, or string.\n\nAs a special case, if the field tag is \"-\", the field is always omitted.\nNote that a field with name \"-\" can still be generated using the tag \"-,\".\n\nExamples of struct field tags and their meanings:\n\n\t// Field appears in JSON as key \"myName\".\n\tField int `json:\"myName\"`\n\n\t// Field appears in JSON as key \"myName\" and\n\t// the field is omitted from the object if its value is empty,\n\t// as defined above.\n\tField int `json:\"myName,omitempty\"`\n\n\t// Field appears in JSON as key \"Field\" (the default), but\n\t// the field is skipped if empty.\n\t// Note the leading comma.\n\tField int `json:\",omitempty\"`\n\n\t// Field is ignored by this package.\n\tField int `json:\"-\"`\n\n\t// Field appears in JSON as key \"-\".\n\tField int `json:\"-,\"`\n\nThe \"string\" option signals that a field is stored as JSON inside a\nJSON-encoded string. It applies only to fields of string, floating point,\ninteger, or boolean types. This extra level of encoding is sometimes used\nwhen communicating with JavaScript programs:\n\n\tInt64String int64 `json:\",string\"`\n\nThe key name will be used if it's a non-empty string consisting of\nonly Unicode letters, digits, and ASCII punctuation except quotation\nmarks, backslash, and comma.\n\nEmbedded struct fields are usually marshaled as if their inner exported fields\nwere fields in the outer struct, subject to the usual Go visibility rules amended\nas described in the next paragraph.\nAn anonymous struct field with a name given in its JSON tag is treated as\nhaving that name, rather than being anonymous.\nAn anonymous struct field of interface type is treated the same as having\nthat type as its name, rather than being anonymous.\n\nThe Go visibility rules for struct fields are amended for JSON when\ndeciding which field to marshal or unmarshal. If there are\nmultiple fields at the same level, and that level is the least\nnested (and would therefore be the nesting level selected by the\nusual Go rules), the following extra rules apply:\n\n1) Of those fields, if any are JSON-tagged, only tagged fields are considered,\neven if there are multiple untagged fields that would otherwise conflict.\n\n2) If there is exactly one field (tagged or not according to the first rule), that is selected.\n\n3) Otherwise there are multiple fields, and all are ignored; no error occurs.\n\nHandling of anonymous struct fields is new in Go 1.1.\nPrior to Go 1.1, anonymous struct fields were ignored. To force ignoring of\nan anonymous struct field in both current and earlier versions, give the field\na JSON tag of \"-\".\n\nMap values encode as JSON objects. The map's key type must either be a\nstring, an integer type, or implement [encoding.TextMarshaler]. The map keys\nare sorted and used as JSON object keys by applying the following rules,\nsubject to the UTF-8 coercion described for string values above:\n  - keys of any string type are used directly\n  - [encoding.TextMarshalers] are marshaled\n  - integer keys are converted to strings\n\nPointer values encode as the value pointed to.\nA nil pointer encodes as the null JSON value.\n\nInterface values encode as the value contained in the interface.\nA nil interface value encodes as the null JSON value.\n\nChannel, complex, and function values cannot be encoded in JSON.\nAttempting to encode such a value causes Marshal to return\nan [UnsupportedTypeError].\n\nJSON cannot represent cyclic data structures and Marshal does not\nhandle them. Passing cyclic structures to Marshal will result in\nan error.", Args: []pkgreflect.Arg{{Name: "v", Tag: "any"}}, Tag: "any", Value: reflect.ValueOf(json.Marshal)},

			"MarshalIndent": {Doc: "MarshalIndent is like [Marshal] but applies [Indent] to format the output.\nEach JSON element in the output will begin on a new line beginning with prefix\nfollowed by one or more copies of indent according to the indentation nesting.", Args: []pkgreflect.Arg{{Name: "v", Tag: "any"}, {Name: "prefix", Tag: "string"}, {Name: "indent", Tag: "string"}}, Tag: "any", Value: reflect.ValueOf(json.MarshalIndent)},

			"NewDecoder": {Doc: "NewDecoder returns a new decoder that reads from r.\n\nThe decoder introduces its own buffering and may\nread data from r beyond the JSON values requested.", Args: []pkgreflect.Arg{{Name: "r", Tag: "io.Reader"}}, Tag: "Decoder", Value: reflect.ValueOf(json.NewDecoder)},

			"NewEncoder": {Doc: "NewEncoder returns a new encoder that writes to w.", Args: []pkgreflect.Arg{{Name: "w", Tag: "io.Writer"}}, Tag: "Encoder", Value: reflect.ValueOf(json.NewEncoder)},

			"Unmarshal": {Doc: "Unmarshal parses the JSON-encoded data and stores the result\nin the value pointed to by v. If v is nil or not a pointer,\nUnmarshal returns an [InvalidUnmarshalError].\n\nUnmarshal uses the inverse of the encodings that\n[Marshal] uses, allocating maps, slices, and pointers as necessary,\nwith the following additional rules:\n\nTo unmarshal JSON into a pointer, Unmarshal first handles the case of\nthe JSON being the JSON literal null. In that case, Unmarshal sets\nthe pointer to nil. Otherwise, Unmarshal unmarshals the JSON into\nthe value pointed at by the pointer. If the pointer is nil, Unmarshal\nallocates a new value for it to point to.\n\nTo unmarshal JSON into a value implementing [Unmarshaler],\nUnmarshal calls that value's [Unmarshaler.UnmarshalJSON] method, including\nwhen the input is a JSON null.\nOtherwise, if the value implements [encoding.TextUnmarshaler]\nand the input is a JSON quoted string, Unmarshal calls\n[encoding.TextUnmarshaler.UnmarshalText] with the unquoted form of the string.\n\nTo unmarshal JSON into a struct, Unmarshal matches incoming object\nkeys to the keys used by [Marshal] (either the struct field name or its tag),\npreferring an exact match but also accepting a case-insensitive match. By\ndefault, object keys which don't have a corresponding struct field are\nignored (see [Decoder.DisallowUnknownFields] for an alternative).\n\nTo unmarshal JSON into an interface value,\nUnmarshal stores one of these in the interface value:\n\n  - bool, for JSON booleans\n  - float64, for JSON numbers\n  - string, for JSON strings\n  - []interface{}, for JSON arrays\n  - map[string]interface{}, for JSON objects\n  - nil for JSON null\n\nTo unmarshal a JSON array into a slice, Unmarshal resets the slice length\nto zero and then appends each element to the slice.\nAs a special case, to unmarshal an empty JSON array into a slice,\nUnmarshal replaces the slice with a new empty slice.\n\nTo unmarshal a JSON array into a Go array, Unmarshal decodes\nJSON array elements into corresponding Go array elements.\nIf the Go array is smaller than the JSON array,\nthe additional JSON array elements are discarded.\nIf the JSON array is smaller than the Go array,\nthe additional Go array elements are set to zero values.\n\nTo unmarshal a JSON object into a map, Unmarshal first establishes a map to\nuse. If the map is nil, Unmarshal allocates a new map. Otherwise Unmarshal\nreuses the existing map, keeping existing entries. Unmarshal then stores\nkey-value pairs from the JSON object into the map. The map's key type must\neither be any string type, an integer, implement [json.Unmarshaler], or\nimplement [encoding.TextUnmarshaler].\n\nIf the JSON-encoded data contain a syntax error, Unmarshal returns a [SyntaxError].\n\nIf a JSON value is not appropriate for a given target type,\nor if a JSON number overflows the target type, Unmarshal\nskips that field and completes the unmarshaling as best it can.\nIf no more serious errors are encountered, Unmarshal returns\nan [UnmarshalTypeError] describing the earliest such error. In any\ncase, it's not guaranteed that all the remaining fields following\nthe problematic one will be unmarshaled into the target object.\n\nThe JSON null value unmarshals into an interface, map, pointer, or slice\nby setting that Go value to nil. Because null is often used in JSON to mean\n“not present,” unmarshaling a JSON null into any other Go type has no effect\non the value and produces no error.\n\nWhen unmarshaling quoted strings, invalid UTF-8 or\ninvalid UTF-16 surrogate pairs are not treated as an error.\nInstead, they are replaced by the Unicode replacement\ncharacter U+FFFD.", Args: []pkgreflect.Arg{{Name: "data", Tag: "[]byte"}, {Name: "v", Tag: "any"}}, Tag: "error", Value: reflect.ValueOf(json.Unmarshal)},

			"Valid": {Doc: "Valid reports whether data is a valid JSON encoding.", Args: []pkgreflect.Arg{{Name: "data", Tag: "[]byte"}}, Tag: "bool", Value: reflect.ValueOf(json.Valid)},
		},

		Variables: map[string]pkgreflect.Value{},

		Consts: map[string]pkgreflect.Value{},
	})
}
