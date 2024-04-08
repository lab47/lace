(ns
  ^{:go-imports ["io"]
    :doc "Provides basic interfaces to I/O primitives."}
  io)

(defn ^Int copy
  "Copies from src to dst until either EOF is reached on src or an error occurs.
  Returns the number of bytes copied or throws an error.
  src must be IOReader, e.g. as returned by lace.os/open.
  dst must be IOWriter, e.g. as returned by lace.os/create."
  {:added "1.0"
  :go "! n, err := io.Copy(dst, src); _res := int(n)"} ;; TODO: 32-bit issue
  [^IOWriter dst ^IOReader src])

(defn pipe
  "Pipe creates a synchronous in-memory pipe. It can be used to connect code expecting an IOReader
  with code expecting an IOWriter.
  Returns a vector [reader, writer]."
  {:added "1.0"
   :go "pipe()"}
  [])

(defn close
  "Closes f (IOWriter, IOReader, or File) if possible. Otherwise throws an error."
  {:added "1.0"
   :go "close(f)"}
  [^Object f])