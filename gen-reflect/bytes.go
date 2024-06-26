package reflect

import (
	bytes "bytes"
	"reflect"

	"github.com/lab47/lace/pkg/pkgreflect"
)

func init() {
	Buffer_methods := map[string]pkgreflect.Func{}
	Reader_methods := map[string]pkgreflect.Func{}
	Buffer_methods["Bytes"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "[]byte", Doc: "Bytes returns a slice of length b.Len() holding the unread portion of the buffer.\nThe slice is valid for use only until the next buffer modification (that is,\nonly until the next call to a method like [Buffer.Read], [Buffer.Write], [Buffer.Reset], or [Buffer.Truncate]).\nThe slice aliases the buffer content at least until the next buffer modification,\nso immediate changes to the slice will affect the result of future reads."}
	Buffer_methods["AvailableBuffer"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "[]byte", Doc: "AvailableBuffer returns an empty buffer with b.Available() capacity.\nThis buffer is intended to be appended to and\npassed to an immediately succeeding [Buffer.Write] call.\nThe buffer is only valid until the next write operation on b."}
	Buffer_methods["String"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "string", Doc: "String returns the contents of the unread portion of the buffer\nas a string. If the [Buffer] is a nil pointer, it returns \"<nil>\".\n\nTo build strings more efficiently, see the strings.Builder type."}
	Buffer_methods["Len"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int", Doc: "Len returns the number of bytes of the unread portion of the buffer;\nb.Len() == len(b.Bytes())."}
	Buffer_methods["Cap"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int", Doc: "Cap returns the capacity of the buffer's underlying byte slice, that is, the\ntotal space allocated for the buffer's data."}
	Buffer_methods["Available"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int", Doc: "Available returns how many bytes are unused in the buffer."}
	Buffer_methods["Truncate"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "n", Tag: "int"}}, Tag: "any", Doc: "Truncate discards all but the first n unread bytes from the buffer\nbut continues to use the same allocated storage.\nIt panics if n is negative or greater than the length of the buffer."}
	Buffer_methods["Reset"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "Reset resets the buffer to be empty,\nbut it retains the underlying storage for use by future writes.\nReset is the same as [Buffer.Truncate](0)."}
	Buffer_methods["Grow"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "n", Tag: "int"}}, Tag: "any", Doc: "Grow grows the buffer's capacity, if necessary, to guarantee space for\nanother n bytes. After Grow(n), at least n bytes can be written to the\nbuffer without another allocation.\nIf n is negative, Grow will panic.\nIf the buffer can't grow it will panic with [ErrTooLarge]."}
	Buffer_methods["Write"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "p", Tag: "[]byte"}}, Tag: "any", Doc: "Write appends the contents of p to the buffer, growing the buffer as\nneeded. The return value n is the length of p; err is always nil. If the\nbuffer becomes too large, Write will panic with [ErrTooLarge]."}
	Buffer_methods["WriteString"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "s", Tag: "string"}}, Tag: "any", Doc: "WriteString appends the contents of s to the buffer, growing the buffer as\nneeded. The return value n is the length of s; err is always nil. If the\nbuffer becomes too large, WriteString will panic with [ErrTooLarge]."}
	Buffer_methods["ReadFrom"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "r", Tag: "io.Reader"}}, Tag: "any", Doc: "ReadFrom reads data from r until EOF and appends it to the buffer, growing\nthe buffer as needed. The return value n is the number of bytes read. Any\nerror except io.EOF encountered during the read is also returned. If the\nbuffer becomes too large, ReadFrom will panic with [ErrTooLarge]."}
	Buffer_methods["WriteTo"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "w", Tag: "io.Writer"}}, Tag: "any", Doc: "WriteTo writes data to w until the buffer is drained or an error occurs.\nThe return value n is the number of bytes written; it always fits into an\nint, but it is int64 to match the io.WriterTo interface. Any error\nencountered during the write is also returned."}
	Buffer_methods["WriteByte"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "c", Tag: "byte"}}, Tag: "error", Doc: "WriteByte appends the byte c to the buffer, growing the buffer as needed.\nThe returned error is always nil, but is included to match [bufio.Writer]'s\nWriteByte. If the buffer becomes too large, WriteByte will panic with\n[ErrTooLarge]."}
	Buffer_methods["WriteRune"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "r", Tag: "rune"}}, Tag: "any", Doc: "WriteRune appends the UTF-8 encoding of Unicode code point r to the\nbuffer, returning its length and an error, which is always nil but is\nincluded to match [bufio.Writer]'s WriteRune. The buffer is grown as needed;\nif it becomes too large, WriteRune will panic with [ErrTooLarge]."}
	Buffer_methods["Read"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "p", Tag: "[]byte"}}, Tag: "any", Doc: "Read reads the next len(p) bytes from the buffer or until the buffer\nis drained. The return value n is the number of bytes read. If the\nbuffer has no data to return, err is io.EOF (unless len(p) is zero);\notherwise it is nil."}
	Buffer_methods["Next"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "n", Tag: "int"}}, Tag: "[]byte", Doc: "Next returns a slice containing the next n bytes from the buffer,\nadvancing the buffer as if the bytes had been returned by [Buffer.Read].\nIf there are fewer than n bytes in the buffer, Next returns the entire buffer.\nThe slice is only valid until the next call to a read or write method."}
	Buffer_methods["ReadByte"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "ReadByte reads and returns the next byte from the buffer.\nIf no byte is available, it returns error io.EOF."}
	Buffer_methods["ReadRune"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "ReadRune reads and returns the next UTF-8-encoded\nUnicode code point from the buffer.\nIf no bytes are available, the error returned is io.EOF.\nIf the bytes are an erroneous UTF-8 encoding, it\nconsumes one byte and returns U+FFFD, 1."}
	Buffer_methods["UnreadRune"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "error", Doc: "UnreadRune unreads the last rune returned by [Buffer.ReadRune].\nIf the most recent read or write operation on the buffer was\nnot a successful [Buffer.ReadRune], UnreadRune returns an error.  (In this regard\nit is stricter than [Buffer.UnreadByte], which will unread the last byte\nfrom any read operation.)"}
	Buffer_methods["UnreadByte"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "error", Doc: "UnreadByte unreads the last byte returned by the most recent successful\nread operation that read at least one byte. If a write has happened since\nthe last read, if the last read returned an error, or if the read read zero\nbytes, UnreadByte returns an error."}
	Buffer_methods["ReadBytes"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "delim", Tag: "byte"}}, Tag: "any", Doc: "ReadBytes reads until the first occurrence of delim in the input,\nreturning a slice containing the data up to and including the delimiter.\nIf ReadBytes encounters an error before finding a delimiter,\nit returns the data read before the error and the error itself (often io.EOF).\nReadBytes returns err != nil if and only if the returned data does not end in\ndelim."}
	Buffer_methods["ReadString"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "delim", Tag: "byte"}}, Tag: "any", Doc: "ReadString reads until the first occurrence of delim in the input,\nreturning a string containing the data up to and including the delimiter.\nIf ReadString encounters an error before finding a delimiter,\nit returns the data read before the error and the error itself (often io.EOF).\nReadString returns err != nil if and only if the returned data does not end\nin delim."}
	Reader_methods["Len"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int", Doc: "Len returns the number of bytes of the unread portion of the\nslice."}
	Reader_methods["Size"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int64", Doc: "Size returns the original length of the underlying byte slice.\nSize is the number of bytes available for reading via [Reader.ReadAt].\nThe result is unaffected by any method calls except [Reader.Reset]."}
	Reader_methods["Read"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}}, Tag: "any", Doc: "Read implements the [io.Reader] interface."}
	Reader_methods["ReadAt"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}, {Name: "off", Tag: "int64"}}, Tag: "any", Doc: "ReadAt implements the [io.ReaderAt] interface."}
	Reader_methods["ReadByte"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "ReadByte implements the [io.ByteReader] interface."}
	Reader_methods["UnreadByte"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "error", Doc: "UnreadByte complements [Reader.ReadByte] in implementing the [io.ByteScanner] interface."}
	Reader_methods["ReadRune"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "ReadRune implements the [io.RuneReader] interface."}
	Reader_methods["UnreadRune"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "error", Doc: "UnreadRune complements [Reader.ReadRune] in implementing the [io.RuneScanner] interface."}
	Reader_methods["Seek"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "offset", Tag: "int64"}, {Name: "whence", Tag: "int"}}, Tag: "any", Doc: "Seek implements the [io.Seeker] interface."}
	Reader_methods["WriteTo"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "w", Tag: "io.Writer"}}, Tag: "any", Doc: "WriteTo implements the [io.WriterTo] interface."}
	Reader_methods["Reset"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}}, Tag: "any", Doc: "Reset resets the [Reader.Reader] to be reading from b."}
	pkgreflect.AddPackage("bytes", &pkgreflect.Package{
		Doc: "Package bytes implements functions for the manipulation of byte slices.",
		Types: map[string]pkgreflect.Type{
			"Buffer": {Doc: "", Value: reflect.TypeOf((*bytes.Buffer)(nil)).Elem(), Methods: Buffer_methods},
			"Reader": {Doc: "", Value: reflect.TypeOf((*bytes.Reader)(nil)).Elem(), Methods: Reader_methods},
		},

		Functions: map[string]pkgreflect.FuncValue{
			"Clone": {Doc: "Clone returns a copy of b[:len(b)].\nThe result may have additional unused capacity.\nClone(nil) returns nil.", Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.Clone)},

			"Compare": {Doc: "Compare returns an integer comparing two byte slices lexicographically.\nThe result will be 0 if a == b, -1 if a < b, and +1 if a > b.\nA nil argument is equivalent to an empty slice.", Args: []pkgreflect.Arg{{Name: "a", Tag: "[]byte"}, {Name: "b", Tag: "[]byte"}}, Tag: "int", Value: reflect.ValueOf(bytes.Compare)},

			"Contains": {Doc: "Contains reports whether subslice is within b.", Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}, {Name: "subslice", Tag: "[]byte"}}, Tag: "bool", Value: reflect.ValueOf(bytes.Contains)},

			"ContainsAny": {Doc: "ContainsAny reports whether any of the UTF-8-encoded code points in chars are within b.", Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}, {Name: "chars", Tag: "string"}}, Tag: "bool", Value: reflect.ValueOf(bytes.ContainsAny)},

			"ContainsFunc": {Doc: "ContainsFunc reports whether any of the UTF-8-encoded code points r within b satisfy f(r).", Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}, {Name: "f", Tag: "Unknown"}}, Tag: "bool", Value: reflect.ValueOf(bytes.ContainsFunc)},

			"ContainsRune": {Doc: "ContainsRune reports whether the rune is contained in the UTF-8-encoded byte slice b.", Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}, {Name: "r", Tag: "rune"}}, Tag: "bool", Value: reflect.ValueOf(bytes.ContainsRune)},

			"Count": {Doc: "Count counts the number of non-overlapping instances of sep in s.\nIf sep is an empty slice, Count returns 1 + the number of UTF-8-encoded code points in s.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "sep", Tag: "[]byte"}}, Tag: "int", Value: reflect.ValueOf(bytes.Count)},

			"Cut": {Doc: "Cut slices s around the first instance of sep,\nreturning the text before and after sep.\nThe found result reports whether sep appears in s.\nIf sep does not appear in s, cut returns s, nil, false.\n\nCut returns slices of the original slice s, not copies.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "sep", Tag: "[]byte"}}, Tag: "any", Value: reflect.ValueOf(bytes.Cut)},

			"CutPrefix": {Doc: "CutPrefix returns s without the provided leading prefix byte slice\nand reports whether it found the prefix.\nIf s doesn't start with prefix, CutPrefix returns s, false.\nIf prefix is the empty byte slice, CutPrefix returns s, true.\n\nCutPrefix returns slices of the original slice s, not copies.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "prefix", Tag: "[]byte"}}, Tag: "any", Value: reflect.ValueOf(bytes.CutPrefix)},

			"CutSuffix": {Doc: "CutSuffix returns s without the provided ending suffix byte slice\nand reports whether it found the suffix.\nIf s doesn't end with suffix, CutSuffix returns s, false.\nIf suffix is the empty byte slice, CutSuffix returns s, true.\n\nCutSuffix returns slices of the original slice s, not copies.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "suffix", Tag: "[]byte"}}, Tag: "any", Value: reflect.ValueOf(bytes.CutSuffix)},

			"Equal": {Doc: "Equal reports whether a and b\nare the same length and contain the same bytes.\nA nil argument is equivalent to an empty slice.", Args: []pkgreflect.Arg{{Name: "a", Tag: "[]byte"}, {Name: "b", Tag: "[]byte"}}, Tag: "bool", Value: reflect.ValueOf(bytes.Equal)},

			"EqualFold": {Doc: "EqualFold reports whether s and t, interpreted as UTF-8 strings,\nare equal under simple Unicode case-folding, which is a more general\nform of case-insensitivity.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "t", Tag: "[]byte"}}, Tag: "bool", Value: reflect.ValueOf(bytes.EqualFold)},

			"Fields": {Doc: "Fields interprets s as a sequence of UTF-8-encoded code points.\nIt splits the slice s around each instance of one or more consecutive white space\ncharacters, as defined by unicode.IsSpace, returning a slice of subslices of s or an\nempty slice if s contains only white space.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}}, Tag: "[][]byte", Value: reflect.ValueOf(bytes.Fields)},

			"FieldsFunc": {Doc: "FieldsFunc interprets s as a sequence of UTF-8-encoded code points.\nIt splits the slice s at each run of code points c satisfying f(c) and\nreturns a slice of subslices of s. If all code points in s satisfy f(c), or\nlen(s) == 0, an empty slice is returned.\n\nFieldsFunc makes no guarantees about the order in which it calls f(c)\nand assumes that f always returns the same value for a given c.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "f", Tag: "Unknown"}}, Tag: "[][]byte", Value: reflect.ValueOf(bytes.FieldsFunc)},

			"HasPrefix": {Doc: "HasPrefix reports whether the byte slice s begins with prefix.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "prefix", Tag: "[]byte"}}, Tag: "bool", Value: reflect.ValueOf(bytes.HasPrefix)},

			"HasSuffix": {Doc: "HasSuffix reports whether the byte slice s ends with suffix.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "suffix", Tag: "[]byte"}}, Tag: "bool", Value: reflect.ValueOf(bytes.HasSuffix)},

			"Index": {Doc: "Index returns the index of the first instance of sep in s, or -1 if sep is not present in s.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "sep", Tag: "[]byte"}}, Tag: "int", Value: reflect.ValueOf(bytes.Index)},

			"IndexAny": {Doc: "IndexAny interprets s as a sequence of UTF-8-encoded Unicode code points.\nIt returns the byte index of the first occurrence in s of any of the Unicode\ncode points in chars. It returns -1 if chars is empty or if there is no code\npoint in common.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "chars", Tag: "string"}}, Tag: "int", Value: reflect.ValueOf(bytes.IndexAny)},

			"IndexByte": {Doc: "IndexByte returns the index of the first instance of c in b, or -1 if c is not present in b.", Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}, {Name: "c", Tag: "byte"}}, Tag: "int", Value: reflect.ValueOf(bytes.IndexByte)},

			"IndexFunc": {Doc: "IndexFunc interprets s as a sequence of UTF-8-encoded code points.\nIt returns the byte index in s of the first Unicode\ncode point satisfying f(c), or -1 if none do.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "f", Tag: "Unknown"}}, Tag: "int", Value: reflect.ValueOf(bytes.IndexFunc)},

			"IndexRune": {Doc: "IndexRune interprets s as a sequence of UTF-8-encoded code points.\nIt returns the byte index of the first occurrence in s of the given rune.\nIt returns -1 if rune is not present in s.\nIf r is utf8.RuneError, it returns the first instance of any\ninvalid UTF-8 byte sequence.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "r", Tag: "rune"}}, Tag: "int", Value: reflect.ValueOf(bytes.IndexRune)},

			"Join": {Doc: "Join concatenates the elements of s to create a new byte slice. The separator\nsep is placed between elements in the resulting slice.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[][]byte"}, {Name: "sep", Tag: "[]byte"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.Join)},

			"LastIndex": {Doc: "LastIndex returns the index of the last instance of sep in s, or -1 if sep is not present in s.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "sep", Tag: "[]byte"}}, Tag: "int", Value: reflect.ValueOf(bytes.LastIndex)},

			"LastIndexAny": {Doc: "LastIndexAny interprets s as a sequence of UTF-8-encoded Unicode code\npoints. It returns the byte index of the last occurrence in s of any of\nthe Unicode code points in chars. It returns -1 if chars is empty or if\nthere is no code point in common.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "chars", Tag: "string"}}, Tag: "int", Value: reflect.ValueOf(bytes.LastIndexAny)},

			"LastIndexByte": {Doc: "LastIndexByte returns the index of the last instance of c in s, or -1 if c is not present in s.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "c", Tag: "byte"}}, Tag: "int", Value: reflect.ValueOf(bytes.LastIndexByte)},

			"LastIndexFunc": {Doc: "LastIndexFunc interprets s as a sequence of UTF-8-encoded code points.\nIt returns the byte index in s of the last Unicode\ncode point satisfying f(c), or -1 if none do.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "f", Tag: "Unknown"}}, Tag: "int", Value: reflect.ValueOf(bytes.LastIndexFunc)},

			"Map": {Doc: "Map returns a copy of the byte slice s with all its characters modified\naccording to the mapping function. If mapping returns a negative value, the character is\ndropped from the byte slice with no replacement. The characters in s and the\noutput are interpreted as UTF-8-encoded code points.", Args: []pkgreflect.Arg{{Name: "mapping", Tag: "Unknown"}, {Name: "s", Tag: "[]byte"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.Map)},

			"NewBuffer": {Doc: "NewBuffer creates and initializes a new [Buffer] using buf as its\ninitial contents. The new [Buffer] takes ownership of buf, and the\ncaller should not use buf after this call. NewBuffer is intended to\nprepare a [Buffer] to read existing data. It can also be used to set\nthe initial size of the internal buffer for writing. To do that,\nbuf should have the desired capacity but a length of zero.\n\nIn most cases, new([Buffer]) (or just declaring a [Buffer] variable) is\nsufficient to initialize a [Buffer].", Args: []pkgreflect.Arg{{Name: "buf", Tag: "[]byte"}}, Tag: "Buffer", Value: reflect.ValueOf(bytes.NewBuffer)},

			"NewBufferString": {Doc: "NewBufferString creates and initializes a new [Buffer] using string s as its\ninitial contents. It is intended to prepare a buffer to read an existing\nstring.\n\nIn most cases, new([Buffer]) (or just declaring a [Buffer] variable) is\nsufficient to initialize a [Buffer].", Args: []pkgreflect.Arg{{Name: "s", Tag: "string"}}, Tag: "Buffer", Value: reflect.ValueOf(bytes.NewBufferString)},

			"NewReader": {Doc: "NewReader returns a new [Reader.Reader] reading from b.", Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}}, Tag: "Reader", Value: reflect.ValueOf(bytes.NewReader)},

			"Repeat": {Doc: "Repeat returns a new byte slice consisting of count copies of b.\n\nIt panics if count is negative or if the result of (len(b) * count)\noverflows.", Args: []pkgreflect.Arg{{Name: "b", Tag: "[]byte"}, {Name: "count", Tag: "int"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.Repeat)},

			"Replace": {Doc: "Replace returns a copy of the slice s with the first n\nnon-overlapping instances of old replaced by new.\nIf old is empty, it matches at the beginning of the slice\nand after each UTF-8 sequence, yielding up to k+1 replacements\nfor a k-rune slice.\nIf n < 0, there is no limit on the number of replacements.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "old", Tag: "[]byte"}, {Name: "new", Tag: "[]byte"}, {Name: "n", Tag: "int"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.Replace)},

			"ReplaceAll": {Doc: "ReplaceAll returns a copy of the slice s with all\nnon-overlapping instances of old replaced by new.\nIf old is empty, it matches at the beginning of the slice\nand after each UTF-8 sequence, yielding up to k+1 replacements\nfor a k-rune slice.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "old", Tag: "[]byte"}, {Name: "new", Tag: "[]byte"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.ReplaceAll)},

			"Runes": {Doc: "Runes interprets s as a sequence of UTF-8-encoded code points.\nIt returns a slice of runes (Unicode code points) equivalent to s.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}}, Tag: "[]rune", Value: reflect.ValueOf(bytes.Runes)},

			"Split": {Doc: "Split slices s into all subslices separated by sep and returns a slice of\nthe subslices between those separators.\nIf sep is empty, Split splits after each UTF-8 sequence.\nIt is equivalent to SplitN with a count of -1.\n\nTo split around the first instance of a separator, see Cut.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "sep", Tag: "[]byte"}}, Tag: "[][]byte", Value: reflect.ValueOf(bytes.Split)},

			"SplitAfter": {Doc: "SplitAfter slices s into all subslices after each instance of sep and\nreturns a slice of those subslices.\nIf sep is empty, SplitAfter splits after each UTF-8 sequence.\nIt is equivalent to SplitAfterN with a count of -1.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "sep", Tag: "[]byte"}}, Tag: "[][]byte", Value: reflect.ValueOf(bytes.SplitAfter)},

			"SplitAfterN": {Doc: "SplitAfterN slices s into subslices after each instance of sep and\nreturns a slice of those subslices.\nIf sep is empty, SplitAfterN splits after each UTF-8 sequence.\nThe count determines the number of subslices to return:\n\n\tn > 0: at most n subslices; the last subslice will be the unsplit remainder.\n\tn == 0: the result is nil (zero subslices)\n\tn < 0: all subslices", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "sep", Tag: "[]byte"}, {Name: "n", Tag: "int"}}, Tag: "[][]byte", Value: reflect.ValueOf(bytes.SplitAfterN)},

			"SplitN": {Doc: "SplitN slices s into subslices separated by sep and returns a slice of\nthe subslices between those separators.\nIf sep is empty, SplitN splits after each UTF-8 sequence.\nThe count determines the number of subslices to return:\n\n\tn > 0: at most n subslices; the last subslice will be the unsplit remainder.\n\tn == 0: the result is nil (zero subslices)\n\tn < 0: all subslices\n\nTo split around the first instance of a separator, see Cut.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "sep", Tag: "[]byte"}, {Name: "n", Tag: "int"}}, Tag: "[][]byte", Value: reflect.ValueOf(bytes.SplitN)},

			"Title": {Doc: "Title treats s as UTF-8-encoded bytes and returns a copy with all Unicode letters that begin\nwords mapped to their title case.\n\nDeprecated: The rule Title uses for word boundaries does not handle Unicode\npunctuation properly. Use golang.org/x/text/cases instead.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.Title)},

			"ToLower": {Doc: "ToLower returns a copy of the byte slice s with all Unicode letters mapped to\ntheir lower case.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.ToLower)},

			"ToLowerSpecial": {Doc: "ToLowerSpecial treats s as UTF-8-encoded bytes and returns a copy with all the Unicode letters mapped to their\nlower case, giving priority to the special casing rules.", Args: []pkgreflect.Arg{{Name: "c", Tag: "unicode.SpecialCase"}, {Name: "s", Tag: "[]byte"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.ToLowerSpecial)},

			"ToTitle": {Doc: "ToTitle treats s as UTF-8-encoded bytes and returns a copy with all the Unicode letters mapped to their title case.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.ToTitle)},

			"ToTitleSpecial": {Doc: "ToTitleSpecial treats s as UTF-8-encoded bytes and returns a copy with all the Unicode letters mapped to their\ntitle case, giving priority to the special casing rules.", Args: []pkgreflect.Arg{{Name: "c", Tag: "unicode.SpecialCase"}, {Name: "s", Tag: "[]byte"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.ToTitleSpecial)},

			"ToUpper": {Doc: "ToUpper returns a copy of the byte slice s with all Unicode letters mapped to\ntheir upper case.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.ToUpper)},

			"ToUpperSpecial": {Doc: "ToUpperSpecial treats s as UTF-8-encoded bytes and returns a copy with all the Unicode letters mapped to their\nupper case, giving priority to the special casing rules.", Args: []pkgreflect.Arg{{Name: "c", Tag: "unicode.SpecialCase"}, {Name: "s", Tag: "[]byte"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.ToUpperSpecial)},

			"ToValidUTF8": {Doc: "ToValidUTF8 treats s as UTF-8-encoded bytes and returns a copy with each run of bytes\nrepresenting invalid UTF-8 replaced with the bytes in replacement, which may be empty.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "replacement", Tag: "[]byte"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.ToValidUTF8)},

			"Trim": {Doc: "Trim returns a subslice of s by slicing off all leading and\ntrailing UTF-8-encoded code points contained in cutset.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "cutset", Tag: "string"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.Trim)},

			"TrimFunc": {Doc: "TrimFunc returns a subslice of s by slicing off all leading and trailing\nUTF-8-encoded code points c that satisfy f(c).", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "f", Tag: "Unknown"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.TrimFunc)},

			"TrimLeft": {Doc: "TrimLeft returns a subslice of s by slicing off all leading\nUTF-8-encoded code points contained in cutset.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "cutset", Tag: "string"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.TrimLeft)},

			"TrimLeftFunc": {Doc: "TrimLeftFunc treats s as UTF-8-encoded bytes and returns a subslice of s by slicing off\nall leading UTF-8-encoded code points c that satisfy f(c).", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "f", Tag: "Unknown"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.TrimLeftFunc)},

			"TrimPrefix": {Doc: "TrimPrefix returns s without the provided leading prefix string.\nIf s doesn't start with prefix, s is returned unchanged.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "prefix", Tag: "[]byte"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.TrimPrefix)},

			"TrimRight": {Doc: "TrimRight returns a subslice of s by slicing off all trailing\nUTF-8-encoded code points that are contained in cutset.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "cutset", Tag: "string"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.TrimRight)},

			"TrimRightFunc": {Doc: "TrimRightFunc returns a subslice of s by slicing off all trailing\nUTF-8-encoded code points c that satisfy f(c).", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "f", Tag: "Unknown"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.TrimRightFunc)},

			"TrimSpace": {Doc: "TrimSpace returns a subslice of s by slicing off all leading and\ntrailing white space, as defined by Unicode.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.TrimSpace)},

			"TrimSuffix": {Doc: "TrimSuffix returns s without the provided trailing suffix string.\nIf s doesn't end with suffix, s is returned unchanged.", Args: []pkgreflect.Arg{{Name: "s", Tag: "[]byte"}, {Name: "suffix", Tag: "[]byte"}}, Tag: "[]byte", Value: reflect.ValueOf(bytes.TrimSuffix)},
		},

		Variables: map[string]pkgreflect.Value{
			"ErrTooLarge": {Doc: "", Value: reflect.ValueOf(&bytes.ErrTooLarge)},
		},

		Consts: map[string]pkgreflect.Value{
			"MinRead": {Doc: "", Value: reflect.ValueOf(bytes.MinRead)},
		},
	})
}
