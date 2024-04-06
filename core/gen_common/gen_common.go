package gen_common

// Imports of std libraries required by core libraries go here, as
// both gen_data and gen_code pull in this package.
import (
	_ "github.com/lab47/lace/std/html"
	_ "github.com/lab47/lace/std/string"
)

type FileInfo struct {
	Name     string
	Filename string
}

/*
	The entries must be ordered such that a given namespace depends

/* only upon namespaces loaded above it. E.g. lace.template depends
/* on lace.walk, so is listed afterwards, not in alphabetical
/* order.
*/
var CoreSourceFiles []FileInfo = []FileInfo{
	{
		Name:     "<lace.core>",
		Filename: "core.clj",
	},
	{
		Name:     "<lace.repl>",
		Filename: "repl.clj",
	},
	{
		Name:     "<lace.walk>",
		Filename: "walk.clj",
	},
	{
		Name:     "<lace.template>",
		Filename: "template.clj",
	},
	{
		Name:     "<lace.test>",
		Filename: "test.clj",
	},
	{
		Name:     "<lace.set>",
		Filename: "set.clj",
	},
	{
		Name:     "<lace.tools.cli>",
		Filename: "tools_cli.clj",
	},
	{
		Name:     "<lace.hiccup>",
		Filename: "hiccup.clj",
	},
	{
		Name:     "<lace.pprint>",
		Filename: "pprint.clj",
	},
	{
		Name:     "<lace.better-cond>",
		Filename: "better_cond.clj",
	},
	{
		Name:     "<lace.core>",
		Filename: "linter_all.clj",
	},
	{
		Name:     "<lace.core>",
		Filename: "linter_lace.clj",
	},
	{
		Name:     "<lace.core>",
		Filename: "linter_cljx.clj",
	},
	{
		Name:     "<lace.core>",
		Filename: "linter_clj.clj",
	},
	{
		Name:     "<lace.core>",
		Filename: "linter_cljs.clj",
	},
}
