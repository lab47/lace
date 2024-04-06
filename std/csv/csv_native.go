package csv

import (
	"encoding/csv"
	"io"
	"strings"

	. "github.com/lab47/lace/core"
)

func csvLazySeq(rdr *csv.Reader) (*LazySeq, error) {
	var c = func(env *Env, args []Object) (Object, error) {
		t, err := rdr.Read()
		if err == io.EOF {
			return EmptyList, nil
		}
		if err != nil {
			return nil, err
		}
		l, err := csvLazySeq(rdr)
		if err != nil {
			return nil, err
		}
		return NewConsSeq(MakeStringVector(t), l), nil
	}
	return NewLazySeq(Proc{Fn: c}), nil
}

func csvSeqOpts(env *Env, src Object, opts Map) (Object, error) {
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
		c, err := AssertChar(env, c, "comma must be a char")
		if err != nil {
			return nil, err
		}
		csvReader.Comma = c.Ch
	}
	if ok, c := opts.Get(MakeKeyword("comment")); ok {
		c, err := AssertChar(env, c, "comment must be a char")
		if err != nil {
			return nil, err
		}
		csvReader.Comment = c.Ch
	}
	if ok, c := opts.Get(MakeKeyword("fields-per-record")); ok {
		i, err := AssertInt(env, c, "fields-per-record must be an integer")
		if err != nil {
			return nil, err
		}
		csvReader.FieldsPerRecord = i.I
	}
	if ok, c := opts.Get(MakeKeyword("lazy-quotes")); ok {
		b, err := AssertBoolean(env, c, "lazy-quotes must be a boolean")
		if err != nil {
			return nil, err
		}
		csvReader.LazyQuotes = b.B
	}
	if ok, c := opts.Get(MakeKeyword("trim-leading-space")); ok {
		b, err := AssertBoolean(env, c, "trim-leading-space must be a boolean")
		if err != nil {
			return nil, err
		}
		csvReader.TrimLeadingSpace = b.B
	}
	return csvLazySeq(csvReader)
}

func sliceOfStrings(env *Env, obj Object) (res []string, err error) {
	sq, err := AssertSeqable(env, obj, "CSV record must be Seqable")
	if err != nil {
		return nil, err
	}
	s := sq.Seq()
	for !s.IsEmpty() {
		res = append(res, s.First().ToString(false))
		s = s.Rest()
	}
	return
}

func writeWriter(env *Env, wr io.Writer, data Seqable, opts Map) error {
	csvWriter := csv.NewWriter(wr)
	if ok, c := opts.Get(MakeKeyword("comma")); ok {
		c, err := AssertChar(env, c, "comma must be a char")
		if err != nil {
			return err
		}
		csvWriter.Comma = c.Ch
	}
	if ok, c := opts.Get(MakeKeyword("use-crlf")); ok {
		b, err := AssertBoolean(env, c, "use-crlf must be a boolean")
		if err != nil {
			return err
		}
		csvWriter.UseCRLF = b.B
	}
	s := data.Seq()
	for !s.IsEmpty() {
		sl, err := sliceOfStrings(env, s.First())
		if err != nil {
			return err
		}
		err = csvWriter.Write(sl)
		if err != nil {
			return err
		}
		s = s.Rest()
	}
	csvWriter.Flush()
	return nil
}

func write(env *Env, wr io.Writer, data Seqable, opts Map) (Object, error) {
	err := writeWriter(env, wr, data, opts)
	return NIL, err
}

func writeString(env *Env, data Seqable, opts Map) (string, error) {
	var b strings.Builder
	err := writeWriter(env, &b, data, opts)
	return b.String(), err
}
