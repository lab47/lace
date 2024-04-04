package csv

import (
	"encoding/csv"
	"io"
	"strings"

	. "github.com/candid82/joker/core"
)

func csvLazySeq(rdr *csv.Reader) *LazySeq {
	var c = func(env *Env, args []Object) Object {
		t, err := rdr.Read()
		if err == io.EOF {
			return EmptyList
		}
		PanicOnErr(err)
		return NewConsSeq(MakeStringVector(t), csvLazySeq(rdr))
	}
	return NewLazySeq(Proc{Fn: c})
}

func csvSeqOpts(env *Env, src Object, opts Map) Object {
	var rdr io.Reader
	switch src := src.(type) {
	case String:
		rdr = strings.NewReader(src.S)
	case io.Reader:
		rdr = src
	default:
		panic(StubNewError("src must be a string or io.Reader"))
	}
	csvReader := csv.NewReader(rdr)
	csvReader.ReuseRecord = true
	if ok, c := opts.Get(MakeKeyword("comma")); ok {
		csvReader.Comma = AssertChar(env, c, "comma must be a char").Ch
	}
	if ok, c := opts.Get(MakeKeyword("comment")); ok {
		csvReader.Comment = AssertChar(env, c, "comment must be a char").Ch
	}
	if ok, c := opts.Get(MakeKeyword("fields-per-record")); ok {
		csvReader.FieldsPerRecord = AssertInt(env, c, "fields-per-record must be an integer").I
	}
	if ok, c := opts.Get(MakeKeyword("lazy-quotes")); ok {
		csvReader.LazyQuotes = AssertBoolean(env, c, "lazy-quotes must be a boolean").B
	}
	if ok, c := opts.Get(MakeKeyword("trim-leading-space")); ok {
		csvReader.TrimLeadingSpace = AssertBoolean(env, c, "trim-leading-space must be a boolean").B
	}
	return csvLazySeq(csvReader)
}

func sliceOfStrings(env *Env, obj Object) (res []string) {
	s := AssertSeqable(env, obj, "CSV record must be Seqable").Seq()
	for !s.IsEmpty() {
		res = append(res, s.First().ToString(false))
		s = s.Rest()
	}
	return
}

func writeWriter(env *Env, wr io.Writer, data Seqable, opts Map) {
	csvWriter := csv.NewWriter(wr)
	if ok, c := opts.Get(MakeKeyword("comma")); ok {
		csvWriter.Comma = AssertChar(env, c, "comma must be a char").Ch
	}
	if ok, c := opts.Get(MakeKeyword("use-crlf")); ok {
		csvWriter.UseCRLF = AssertBoolean(env, c, "use-crlf must be a boolean").B
	}
	s := data.Seq()
	for !s.IsEmpty() {
		err := csvWriter.Write(sliceOfStrings(env, s.First()))
		PanicOnErr(err)
		s = s.Rest()
	}
	csvWriter.Flush()
}

func write(env *Env, wr io.Writer, data Seqable, opts Map) Object {
	writeWriter(env, wr, data, opts)
	return NIL
}

func writeString(env *Env, data Seqable, opts Map) string {
	var b strings.Builder
	writeWriter(env, &b, data, opts)
	return b.String()
}
