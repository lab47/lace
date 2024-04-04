//go:build !gen_data
// +build !gen_data

package core

var haveSetCoreNamespaces bool

func SetupGlobalEnvCoreData() {
	// Let MaybeLazy() handle initialization.
	if !haveSetCoreNamespaces {
		setCoreNamespaces(GLOBAL_ENV)
		haveSetCoreNamespaces = true
	}
}

func ProcessReplData() {
	// Let MaybeLazy() handle initialization.
}

func ProcessLinterData(dialect Dialect) {
	if dialect == EDN {
		markJokerNamespacesAsUsed()
		return
	}
	processData(linter_allData)
	GLOBAL_ENV.CoreNamespace.Resolve("*loaded-libs*").Value = EmptySet()
	if dialect == JOKER {
		markJokerNamespacesAsUsed()
		processData(linter_jokerData)
		return
	}
	processData(linter_cljxData)
	switch dialect {
	case CLJ:
		processData(linter_cljData)
	case CLJS:
		processData(linter_cljsData)
	}
	removeJokerNamespaces()
}
