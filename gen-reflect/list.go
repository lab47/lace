package reflect

//go:generate go run ../pkg/pkgreflect/cmd/pkgreflect os os.go
//go:generate go run ../pkg/pkgreflect/cmd/pkgreflect time time.go
//go:generate go run ../pkg/pkgreflect/cmd/pkgreflect io io.go
//go:generate go run ../pkg/pkgreflect/cmd/pkgreflect bytes bytes.go
//go:generate go run ../pkg/pkgreflect/cmd/pkgreflect encoding/json encoding_json.go
//go:generate go run ../pkg/pkgreflect/cmd/pkgreflect -lace-name lace.lang github.com/lab47/lace/core/lang core_lang.go
//go:generate go run ../pkg/pkgreflect/cmd/pkgreflect -lace-name lace.rpc github.com/lab47/lace/pkg/rpc lace_rpc.go
