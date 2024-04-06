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
		Filename: "core.joke",
	},
	{
		Name:     "<lace.repl>",
		Filename: "repl.joke",
	},
	{
		Name:     "<lace.walk>",
		Filename: "walk.joke",
	},
	{
		Name:     "<lace.template>",
		Filename: "template.joke",
	},
	{
		Name:     "<lace.test>",
		Filename: "test.joke",
	},
	{
		Name:     "<lace.set>",
		Filename: "set.joke",
	},
	{
		Name:     "<lace.tools.cli>",
		Filename: "tools_cli.joke",
	},
	{
		Name:     "<lace.core>",
		Filename: "linter_all.joke",
	},
	{
		Name:     "<lace.core>",
		Filename: "linter_lace.joke",
	},
	{
		Name:     "<lace.core>",
		Filename: "linter_cljx.joke",
	},
	{
		Name:     "<lace.core>",
		Filename: "linter_clj.joke",
	},
	{
		Name:     "<lace.core>",
		Filename: "linter_cljs.joke",
	},
	{
		Name:     "<lace.hiccup>",
		Filename: "hiccup.joke",
	},
	{
		Name:     "<lace.pprint>",
		Filename: "pprint.joke",
	},
	{
		Name:     "<lace.better-cond>",
		Filename: "better_cond.joke",
	},
}
