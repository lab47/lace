package reflect

import (
	io "io"
	"reflect"

	"github.com/lab47/lace/pkg/pkgreflect"
)

type ByteReaderImpl struct {
	ReadByteFn func() (byte, error)
}

func (s *ByteReaderImpl) ReadByte() (byte, error) {
	return s.ReadByteFn()
}

type ByteScannerImpl struct {
	ReadByteFn   func() (byte, error)
	UnreadByteFn func() error
}

func (s *ByteScannerImpl) ReadByte() (byte, error) {
	return s.ReadByteFn()
}
func (s *ByteScannerImpl) UnreadByte() error {
	return s.UnreadByteFn()
}

type ByteWriterImpl struct {
	WriteByteFn func(byte) error
}

func (s *ByteWriterImpl) WriteByte(a0 byte) error {
	return s.WriteByteFn(a0)
}

type CloserImpl struct {
	CloseFn func() error
}

func (s *CloserImpl) Close() error {
	return s.CloseFn()
}

type ReadCloserImpl struct {
	CloseFn func() error
	ReadFn  func([]byte) (int, error)
}

func (s *ReadCloserImpl) Close() error {
	return s.CloseFn()
}
func (s *ReadCloserImpl) Read(a0 []byte) (int, error) {
	return s.ReadFn(a0)
}

type ReadSeekCloserImpl struct {
	CloseFn func() error
	ReadFn  func([]byte) (int, error)
	SeekFn  func(int64, int) (int64, error)
}

func (s *ReadSeekCloserImpl) Close() error {
	return s.CloseFn()
}
func (s *ReadSeekCloserImpl) Read(a0 []byte) (int, error) {
	return s.ReadFn(a0)
}
func (s *ReadSeekCloserImpl) Seek(a0 int64, a1 int) (int64, error) {
	return s.SeekFn(a0, a1)
}

type ReadSeekerImpl struct {
	ReadFn func([]byte) (int, error)
	SeekFn func(int64, int) (int64, error)
}

func (s *ReadSeekerImpl) Read(a0 []byte) (int, error) {
	return s.ReadFn(a0)
}
func (s *ReadSeekerImpl) Seek(a0 int64, a1 int) (int64, error) {
	return s.SeekFn(a0, a1)
}

type ReadWriteCloserImpl struct {
	CloseFn func() error
	ReadFn  func([]byte) (int, error)
	WriteFn func([]byte) (int, error)
}

func (s *ReadWriteCloserImpl) Close() error {
	return s.CloseFn()
}
func (s *ReadWriteCloserImpl) Read(a0 []byte) (int, error) {
	return s.ReadFn(a0)
}
func (s *ReadWriteCloserImpl) Write(a0 []byte) (int, error) {
	return s.WriteFn(a0)
}

type ReadWriteSeekerImpl struct {
	ReadFn  func([]byte) (int, error)
	SeekFn  func(int64, int) (int64, error)
	WriteFn func([]byte) (int, error)
}

func (s *ReadWriteSeekerImpl) Read(a0 []byte) (int, error) {
	return s.ReadFn(a0)
}
func (s *ReadWriteSeekerImpl) Seek(a0 int64, a1 int) (int64, error) {
	return s.SeekFn(a0, a1)
}
func (s *ReadWriteSeekerImpl) Write(a0 []byte) (int, error) {
	return s.WriteFn(a0)
}

type ReadWriterImpl struct {
	ReadFn  func([]byte) (int, error)
	WriteFn func([]byte) (int, error)
}

func (s *ReadWriterImpl) Read(a0 []byte) (int, error) {
	return s.ReadFn(a0)
}
func (s *ReadWriterImpl) Write(a0 []byte) (int, error) {
	return s.WriteFn(a0)
}

type ReaderImpl struct {
	ReadFn func([]byte) (int, error)
}

func (s *ReaderImpl) Read(a0 []byte) (int, error) {
	return s.ReadFn(a0)
}

type ReaderAtImpl struct {
	ReadAtFn func([]byte, int64) (int, error)
}

func (s *ReaderAtImpl) ReadAt(a0 []byte, a1 int64) (int, error) {
	return s.ReadAtFn(a0, a1)
}

type ReaderFromImpl struct {
	ReadFromFn func(io.Reader) (int64, error)
}

func (s *ReaderFromImpl) ReadFrom(a0 io.Reader) (int64, error) {
	return s.ReadFromFn(a0)
}

type RuneReaderImpl struct {
	ReadRuneFn func() (rune, int, error)
}

func (s *RuneReaderImpl) ReadRune() (rune, int, error) {
	return s.ReadRuneFn()
}

type RuneScannerImpl struct {
	ReadRuneFn   func() (rune, int, error)
	UnreadRuneFn func() error
}

func (s *RuneScannerImpl) ReadRune() (rune, int, error) {
	return s.ReadRuneFn()
}
func (s *RuneScannerImpl) UnreadRune() error {
	return s.UnreadRuneFn()
}

type SeekerImpl struct {
	SeekFn func(int64, int) (int64, error)
}

func (s *SeekerImpl) Seek(a0 int64, a1 int) (int64, error) {
	return s.SeekFn(a0, a1)
}

type StringWriterImpl struct {
	WriteStringFn func(string) (int, error)
}

func (s *StringWriterImpl) WriteString(a0 string) (int, error) {
	return s.WriteStringFn(a0)
}

type WriteCloserImpl struct {
	CloseFn func() error
	WriteFn func([]byte) (int, error)
}

func (s *WriteCloserImpl) Close() error {
	return s.CloseFn()
}
func (s *WriteCloserImpl) Write(a0 []byte) (int, error) {
	return s.WriteFn(a0)
}

type WriteSeekerImpl struct {
	SeekFn  func(int64, int) (int64, error)
	WriteFn func([]byte) (int, error)
}

func (s *WriteSeekerImpl) Seek(a0 int64, a1 int) (int64, error) {
	return s.SeekFn(a0, a1)
}
func (s *WriteSeekerImpl) Write(a0 []byte) (int, error) {
	return s.WriteFn(a0)
}

type WriterImpl struct {
	WriteFn func([]byte) (int, error)
}

func (s *WriterImpl) Write(a0 []byte) (int, error) {
	return s.WriteFn(a0)
}

type WriterAtImpl struct {
	WriteAtFn func([]byte, int64) (int, error)
}

func (s *WriterAtImpl) WriteAt(a0 []byte, a1 int64) (int, error) {
	return s.WriteAtFn(a0, a1)
}

type WriterToImpl struct {
	WriteToFn func(io.Writer) (int64, error)
}

func (s *WriterToImpl) WriteTo(a0 io.Writer) (int64, error) {
	return s.WriteToFn(a0)
}

func init() {
	ByteReader_methods := map[string]pkgreflect.Func{}
	ByteScanner_methods := map[string]pkgreflect.Func{}
	ByteWriter_methods := map[string]pkgreflect.Func{}
	Closer_methods := map[string]pkgreflect.Func{}
	LimitedReader_methods := map[string]pkgreflect.Func{}
	OffsetWriter_methods := map[string]pkgreflect.Func{}
	ReadCloser_methods := map[string]pkgreflect.Func{}
	ReadSeekCloser_methods := map[string]pkgreflect.Func{}
	ReadSeeker_methods := map[string]pkgreflect.Func{}
	ReadWriteCloser_methods := map[string]pkgreflect.Func{}
	ReadWriteSeeker_methods := map[string]pkgreflect.Func{}
	ReadWriter_methods := map[string]pkgreflect.Func{}
	Reader_methods := map[string]pkgreflect.Func{}
	ReaderAt_methods := map[string]pkgreflect.Func{}
	ReaderFrom_methods := map[string]pkgreflect.Func{}
	RuneReader_methods := map[string]pkgreflect.Func{}
	RuneScanner_methods := map[string]pkgreflect.Func{}
	SectionReader_methods := map[string]pkgreflect.Func{}
	Seeker_methods := map[string]pkgreflect.Func{}
	StringWriter_methods := map[string]pkgreflect.Func{}
	WriteCloser_methods := map[string]pkgreflect.Func{}
	WriteSeeker_methods := map[string]pkgreflect.Func{}
	Writer_methods := map[string]pkgreflect.Func{}
	WriterAt_methods := map[string]pkgreflect.Func{}
	WriterTo_methods := map[string]pkgreflect.Func{}
	PipeReader_methods := map[string]pkgreflect.Func{}
	PipeWriter_methods := map[string]pkgreflect.Func{}
	LimitedReader_methods["Read"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "p", Tag: "[]byte"}}, Tag: "any", Doc: ""}
	SectionReader_methods["Read"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "p", Tag: "[]byte"}}, Tag: "any", Doc: ""}
	SectionReader_methods["Seek"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "offset", Tag: "int64"}, {Name: "whence", Tag: "int"}}, Tag: "any", Doc: ""}
	SectionReader_methods["ReadAt"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "p", Tag: "[]byte"}, {Name: "off", Tag: "int64"}}, Tag: "any", Doc: ""}
	SectionReader_methods["Size"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "int64", Doc: "Size returns the size of the section in bytes."}
	SectionReader_methods["Outer"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "any", Doc: "Outer returns the underlying [ReaderAt] and offsets for the section.\n\nThe returned values are the same that were passed to [NewSectionReader]\nwhen the [SectionReader] was created."}
	OffsetWriter_methods["Write"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "p", Tag: "[]byte"}}, Tag: "any", Doc: ""}
	OffsetWriter_methods["WriteAt"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "p", Tag: "[]byte"}, {Name: "off", Tag: "int64"}}, Tag: "any", Doc: ""}
	OffsetWriter_methods["Seek"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "offset", Tag: "int64"}, {Name: "whence", Tag: "int"}}, Tag: "any", Doc: ""}
	PipeReader_methods["Read"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "data", Tag: "[]byte"}}, Tag: "any", Doc: "Read implements the standard Read interface:\nit reads data from the pipe, blocking until a writer\narrives or the write end is closed.\nIf the write end is closed with an error, that error is\nreturned as err; otherwise err is EOF."}
	PipeReader_methods["Close"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "error", Doc: "Close closes the reader; subsequent writes to the\nwrite half of the pipe will return the error [ErrClosedPipe]."}
	PipeReader_methods["CloseWithError"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "err", Tag: "error"}}, Tag: "error", Doc: "CloseWithError closes the reader; subsequent writes\nto the write half of the pipe will return the error err.\n\nCloseWithError never overwrites the previous error if it exists\nand always returns nil."}
	PipeWriter_methods["Write"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "data", Tag: "[]byte"}}, Tag: "any", Doc: "Write implements the standard Write interface:\nit writes data to the pipe, blocking until one or more readers\nhave consumed all the data or the read end is closed.\nIf the read end is closed with an error, that err is\nreturned as err; otherwise err is [ErrClosedPipe]."}
	PipeWriter_methods["Close"] = pkgreflect.Func{Args: []pkgreflect.Arg{}, Tag: "error", Doc: "Close closes the writer; subsequent reads from the\nread half of the pipe will return no bytes and EOF."}
	PipeWriter_methods["CloseWithError"] = pkgreflect.Func{Args: []pkgreflect.Arg{{Name: "err", Tag: "error"}}, Tag: "error", Doc: "CloseWithError closes the writer; subsequent reads from the\nread half of the pipe will return no bytes and the error err,\nor EOF if err is nil.\n\nCloseWithError never overwrites the previous error if it exists\nand always returns nil."}
	pkgreflect.AddPackage("io", &pkgreflect.Package{
		Doc: "Package io provides basic interfaces to I/O primitives.",
		Types: map[string]pkgreflect.Type{
			"ByteReader":          {Doc: "", Value: reflect.TypeOf((*io.ByteReader)(nil)).Elem(), Methods: ByteReader_methods},
			"ByteScanner":         {Doc: "", Value: reflect.TypeOf((*io.ByteScanner)(nil)).Elem(), Methods: ByteScanner_methods},
			"ByteWriter":          {Doc: "", Value: reflect.TypeOf((*io.ByteWriter)(nil)).Elem(), Methods: ByteWriter_methods},
			"Closer":              {Doc: "", Value: reflect.TypeOf((*io.Closer)(nil)).Elem(), Methods: Closer_methods},
			"LimitedReader":       {Doc: "", Value: reflect.TypeOf((*io.LimitedReader)(nil)).Elem(), Methods: LimitedReader_methods},
			"OffsetWriter":        {Doc: "", Value: reflect.TypeOf((*io.OffsetWriter)(nil)).Elem(), Methods: OffsetWriter_methods},
			"PipeReader":          {Doc: "", Value: reflect.TypeOf((*io.PipeReader)(nil)).Elem(), Methods: PipeReader_methods},
			"PipeWriter":          {Doc: "", Value: reflect.TypeOf((*io.PipeWriter)(nil)).Elem(), Methods: PipeWriter_methods},
			"ReadCloser":          {Doc: "", Value: reflect.TypeOf((*io.ReadCloser)(nil)).Elem(), Methods: ReadCloser_methods},
			"ReadSeekCloser":      {Doc: "", Value: reflect.TypeOf((*io.ReadSeekCloser)(nil)).Elem(), Methods: ReadSeekCloser_methods},
			"ReadSeeker":          {Doc: "", Value: reflect.TypeOf((*io.ReadSeeker)(nil)).Elem(), Methods: ReadSeeker_methods},
			"ReadWriteCloser":     {Doc: "", Value: reflect.TypeOf((*io.ReadWriteCloser)(nil)).Elem(), Methods: ReadWriteCloser_methods},
			"ReadWriteSeeker":     {Doc: "", Value: reflect.TypeOf((*io.ReadWriteSeeker)(nil)).Elem(), Methods: ReadWriteSeeker_methods},
			"ReadWriter":          {Doc: "", Value: reflect.TypeOf((*io.ReadWriter)(nil)).Elem(), Methods: ReadWriter_methods},
			"Reader":              {Doc: "", Value: reflect.TypeOf((*io.Reader)(nil)).Elem(), Methods: Reader_methods},
			"ReaderAt":            {Doc: "", Value: reflect.TypeOf((*io.ReaderAt)(nil)).Elem(), Methods: ReaderAt_methods},
			"ReaderFrom":          {Doc: "", Value: reflect.TypeOf((*io.ReaderFrom)(nil)).Elem(), Methods: ReaderFrom_methods},
			"RuneReader":          {Doc: "", Value: reflect.TypeOf((*io.RuneReader)(nil)).Elem(), Methods: RuneReader_methods},
			"RuneScanner":         {Doc: "", Value: reflect.TypeOf((*io.RuneScanner)(nil)).Elem(), Methods: RuneScanner_methods},
			"SectionReader":       {Doc: "", Value: reflect.TypeOf((*io.SectionReader)(nil)).Elem(), Methods: SectionReader_methods},
			"Seeker":              {Doc: "", Value: reflect.TypeOf((*io.Seeker)(nil)).Elem(), Methods: Seeker_methods},
			"StringWriter":        {Doc: "", Value: reflect.TypeOf((*io.StringWriter)(nil)).Elem(), Methods: StringWriter_methods},
			"WriteCloser":         {Doc: "", Value: reflect.TypeOf((*io.WriteCloser)(nil)).Elem(), Methods: WriteCloser_methods},
			"WriteSeeker":         {Doc: "", Value: reflect.TypeOf((*io.WriteSeeker)(nil)).Elem(), Methods: WriteSeeker_methods},
			"Writer":              {Doc: "", Value: reflect.TypeOf((*io.Writer)(nil)).Elem(), Methods: Writer_methods},
			"WriterAt":            {Doc: "", Value: reflect.TypeOf((*io.WriterAt)(nil)).Elem(), Methods: WriterAt_methods},
			"WriterTo":            {Doc: "", Value: reflect.TypeOf((*io.WriterTo)(nil)).Elem(), Methods: WriterTo_methods},
			"ByteReaderImpl":      {Doc: `Struct version of interface ByteReader for implementation`, Value: reflect.TypeFor[ByteReaderImpl]()},
			"ByteScannerImpl":     {Doc: `Struct version of interface ByteScanner for implementation`, Value: reflect.TypeFor[ByteScannerImpl]()},
			"ByteWriterImpl":      {Doc: `Struct version of interface ByteWriter for implementation`, Value: reflect.TypeFor[ByteWriterImpl]()},
			"CloserImpl":          {Doc: `Struct version of interface Closer for implementation`, Value: reflect.TypeFor[CloserImpl]()},
			"ReadCloserImpl":      {Doc: `Struct version of interface ReadCloser for implementation`, Value: reflect.TypeFor[ReadCloserImpl]()},
			"ReadSeekCloserImpl":  {Doc: `Struct version of interface ReadSeekCloser for implementation`, Value: reflect.TypeFor[ReadSeekCloserImpl]()},
			"ReadSeekerImpl":      {Doc: `Struct version of interface ReadSeeker for implementation`, Value: reflect.TypeFor[ReadSeekerImpl]()},
			"ReadWriteCloserImpl": {Doc: `Struct version of interface ReadWriteCloser for implementation`, Value: reflect.TypeFor[ReadWriteCloserImpl]()},
			"ReadWriteSeekerImpl": {Doc: `Struct version of interface ReadWriteSeeker for implementation`, Value: reflect.TypeFor[ReadWriteSeekerImpl]()},
			"ReadWriterImpl":      {Doc: `Struct version of interface ReadWriter for implementation`, Value: reflect.TypeFor[ReadWriterImpl]()},
			"ReaderImpl":          {Doc: `Struct version of interface Reader for implementation`, Value: reflect.TypeFor[ReaderImpl]()},
			"ReaderAtImpl":        {Doc: `Struct version of interface ReaderAt for implementation`, Value: reflect.TypeFor[ReaderAtImpl]()},
			"ReaderFromImpl":      {Doc: `Struct version of interface ReaderFrom for implementation`, Value: reflect.TypeFor[ReaderFromImpl]()},
			"RuneReaderImpl":      {Doc: `Struct version of interface RuneReader for implementation`, Value: reflect.TypeFor[RuneReaderImpl]()},
			"RuneScannerImpl":     {Doc: `Struct version of interface RuneScanner for implementation`, Value: reflect.TypeFor[RuneScannerImpl]()},
			"SeekerImpl":          {Doc: `Struct version of interface Seeker for implementation`, Value: reflect.TypeFor[SeekerImpl]()},
			"StringWriterImpl":    {Doc: `Struct version of interface StringWriter for implementation`, Value: reflect.TypeFor[StringWriterImpl]()},
			"WriteCloserImpl":     {Doc: `Struct version of interface WriteCloser for implementation`, Value: reflect.TypeFor[WriteCloserImpl]()},
			"WriteSeekerImpl":     {Doc: `Struct version of interface WriteSeeker for implementation`, Value: reflect.TypeFor[WriteSeekerImpl]()},
			"WriterImpl":          {Doc: `Struct version of interface Writer for implementation`, Value: reflect.TypeFor[WriterImpl]()},
			"WriterAtImpl":        {Doc: `Struct version of interface WriterAt for implementation`, Value: reflect.TypeFor[WriterAtImpl]()},
			"WriterToImpl":        {Doc: `Struct version of interface WriterTo for implementation`, Value: reflect.TypeFor[WriterToImpl]()},
		},

		Functions: map[string]pkgreflect.FuncValue{
			"Copy": {Doc: "Copy copies from src to dst until either EOF is reached\non src or an error occurs. It returns the number of bytes\ncopied and the first error encountered while copying, if any.\n\nA successful Copy returns err == nil, not err == EOF.\nBecause Copy is defined to read from src until EOF, it does\nnot treat an EOF from Read as an error to be reported.\n\nIf src implements [WriterTo],\nthe copy is implemented by calling src.WriteTo(dst).\nOtherwise, if dst implements [ReaderFrom],\nthe copy is implemented by calling dst.ReadFrom(src).", Args: []pkgreflect.Arg{{Name: "dst", Tag: "Writer"}, {Name: "src", Tag: "Reader"}}, Tag: "any", Value: reflect.ValueOf(io.Copy)},

			"CopyBuffer": {Doc: "CopyBuffer is identical to Copy except that it stages through the\nprovided buffer (if one is required) rather than allocating a\ntemporary one. If buf is nil, one is allocated; otherwise if it has\nzero length, CopyBuffer panics.\n\nIf either src implements [WriterTo] or dst implements [ReaderFrom],\nbuf will not be used to perform the copy.", Args: []pkgreflect.Arg{{Name: "dst", Tag: "Writer"}, {Name: "src", Tag: "Reader"}, {Name: "buf", Tag: "[]byte"}}, Tag: "any", Value: reflect.ValueOf(io.CopyBuffer)},

			"CopyN": {Doc: "CopyN copies n bytes (or until an error) from src to dst.\nIt returns the number of bytes copied and the earliest\nerror encountered while copying.\nOn return, written == n if and only if err == nil.\n\nIf dst implements [ReaderFrom], the copy is implemented using it.", Args: []pkgreflect.Arg{{Name: "dst", Tag: "Writer"}, {Name: "src", Tag: "Reader"}, {Name: "n", Tag: "int64"}}, Tag: "any", Value: reflect.ValueOf(io.CopyN)},

			"LimitReader": {Doc: "LimitReader returns a Reader that reads from r\nbut stops with EOF after n bytes.\nThe underlying implementation is a *LimitedReader.", Args: []pkgreflect.Arg{{Name: "r", Tag: "Reader"}, {Name: "n", Tag: "int64"}}, Tag: "Reader", Value: reflect.ValueOf(io.LimitReader)},

			"MultiReader": {Doc: "MultiReader returns a Reader that's the logical concatenation of\nthe provided input readers. They're read sequentially. Once all\ninputs have returned EOF, Read will return EOF.  If any of the readers\nreturn a non-nil, non-EOF error, Read will return that error.", Args: []pkgreflect.Arg{{Name: "readers", Tag: "Unknown"}}, Tag: "Reader", Value: reflect.ValueOf(io.MultiReader)},

			"MultiWriter": {Doc: "MultiWriter creates a writer that duplicates its writes to all the\nprovided writers, similar to the Unix tee(1) command.\n\nEach write is written to each listed writer, one at a time.\nIf a listed writer returns an error, that overall write operation\nstops and returns the error; it does not continue down the list.", Args: []pkgreflect.Arg{{Name: "writers", Tag: "Unknown"}}, Tag: "Writer", Value: reflect.ValueOf(io.MultiWriter)},

			"NewOffsetWriter": {Doc: "NewOffsetWriter returns an [OffsetWriter] that writes to w\nstarting at offset off.", Args: []pkgreflect.Arg{{Name: "w", Tag: "WriterAt"}, {Name: "off", Tag: "int64"}}, Tag: "OffsetWriter", Value: reflect.ValueOf(io.NewOffsetWriter)},

			"NewSectionReader": {Doc: "NewSectionReader returns a [SectionReader] that reads from r\nstarting at offset off and stops with EOF after n bytes.", Args: []pkgreflect.Arg{{Name: "r", Tag: "ReaderAt"}, {Name: "off", Tag: "int64"}, {Name: "n", Tag: "int64"}}, Tag: "SectionReader", Value: reflect.ValueOf(io.NewSectionReader)},

			"NopCloser": {Doc: "NopCloser returns a [ReadCloser] with a no-op Close method wrapping\nthe provided [Reader] r.\nIf r implements [WriterTo], the returned [ReadCloser] will implement [WriterTo]\nby forwarding calls to r.", Args: []pkgreflect.Arg{{Name: "r", Tag: "Reader"}}, Tag: "ReadCloser", Value: reflect.ValueOf(io.NopCloser)},

			"Pipe": {Doc: "Pipe creates a synchronous in-memory pipe.\nIt can be used to connect code expecting an [io.Reader]\nwith code expecting an [io.Writer].\n\nReads and Writes on the pipe are matched one to one\nexcept when multiple Reads are needed to consume a single Write.\nThat is, each Write to the [PipeWriter] blocks until it has satisfied\none or more Reads from the [PipeReader] that fully consume\nthe written data.\nThe data is copied directly from the Write to the corresponding\nRead (or Reads); there is no internal buffering.\n\nIt is safe to call Read and Write in parallel with each other or with Close.\nParallel calls to Read and parallel calls to Write are also safe:\nthe individual calls will be gated sequentially.", Args: []pkgreflect.Arg{}, Tag: "any", Value: reflect.ValueOf(io.Pipe)},

			"ReadAll": {Doc: "ReadAll reads from r until an error or EOF and returns the data it read.\nA successful call returns err == nil, not err == EOF. Because ReadAll is\ndefined to read from src until EOF, it does not treat an EOF from Read\nas an error to be reported.", Args: []pkgreflect.Arg{{Name: "r", Tag: "Reader"}}, Tag: "any", Value: reflect.ValueOf(io.ReadAll)},

			"ReadAtLeast": {Doc: "ReadAtLeast reads from r into buf until it has read at least min bytes.\nIt returns the number of bytes copied and an error if fewer bytes were read.\nThe error is EOF only if no bytes were read.\nIf an EOF happens after reading fewer than min bytes,\nReadAtLeast returns [ErrUnexpectedEOF].\nIf min is greater than the length of buf, ReadAtLeast returns [ErrShortBuffer].\nOn return, n >= min if and only if err == nil.\nIf r returns an error having read at least min bytes, the error is dropped.", Args: []pkgreflect.Arg{{Name: "r", Tag: "Reader"}, {Name: "buf", Tag: "[]byte"}, {Name: "min", Tag: "int"}}, Tag: "any", Value: reflect.ValueOf(io.ReadAtLeast)},

			"ReadFull": {Doc: "ReadFull reads exactly len(buf) bytes from r into buf.\nIt returns the number of bytes copied and an error if fewer bytes were read.\nThe error is EOF only if no bytes were read.\nIf an EOF happens after reading some but not all the bytes,\nReadFull returns [ErrUnexpectedEOF].\nOn return, n == len(buf) if and only if err == nil.\nIf r returns an error having read at least len(buf) bytes, the error is dropped.", Args: []pkgreflect.Arg{{Name: "r", Tag: "Reader"}, {Name: "buf", Tag: "[]byte"}}, Tag: "any", Value: reflect.ValueOf(io.ReadFull)},

			"TeeReader": {Doc: "TeeReader returns a [Reader] that writes to w what it reads from r.\nAll reads from r performed through it are matched with\ncorresponding writes to w. There is no internal buffering -\nthe write must complete before the read completes.\nAny error encountered while writing is reported as a read error.", Args: []pkgreflect.Arg{{Name: "r", Tag: "Reader"}, {Name: "w", Tag: "Writer"}}, Tag: "Reader", Value: reflect.ValueOf(io.TeeReader)},

			"WriteString": {Doc: "WriteString writes the contents of the string s to w, which accepts a slice of bytes.\nIf w implements [StringWriter], [StringWriter.WriteString] is invoked directly.\nOtherwise, [Writer.Write] is called exactly once.", Args: []pkgreflect.Arg{{Name: "w", Tag: "Writer"}, {Name: "s", Tag: "string"}}, Tag: "any", Value: reflect.ValueOf(io.WriteString)},
		},

		Variables: map[string]pkgreflect.Value{
			"Discard":          {Doc: "", Value: reflect.ValueOf(&io.Discard)},
			"EOF":              {Doc: "", Value: reflect.ValueOf(&io.EOF)},
			"ErrClosedPipe":    {Doc: "", Value: reflect.ValueOf(&io.ErrClosedPipe)},
			"ErrNoProgress":    {Doc: "", Value: reflect.ValueOf(&io.ErrNoProgress)},
			"ErrShortBuffer":   {Doc: "", Value: reflect.ValueOf(&io.ErrShortBuffer)},
			"ErrShortWrite":    {Doc: "", Value: reflect.ValueOf(&io.ErrShortWrite)},
			"ErrUnexpectedEOF": {Doc: "", Value: reflect.ValueOf(&io.ErrUnexpectedEOF)},
		},

		Consts: map[string]pkgreflect.Value{
			"SeekCurrent": {Doc: "", Value: reflect.ValueOf(io.SeekCurrent)},
			"SeekEnd":     {Doc: "", Value: reflect.ValueOf(io.SeekEnd)},
			"SeekStart":   {Doc: "", Value: reflect.ValueOf(io.SeekStart)},
		},
	})
}
