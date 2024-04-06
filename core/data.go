//go:build !gen_data
// +build !gen_data

package core

func ProcessReplData() {
	// Let MaybeLazy() handle initialization.
}

func ProcessLinterData(env *Env, dialect Dialect) error {
	if dialect == EDN {
		markLaceNamespacesAsUsed(env)
		return nil
	}

	ns := env.CoreNamespace

	if err := processInEnvInNS(env, ns, linter_allData); err != nil {
		return err
	}

	ns.Resolve("*loaded-libs*").Value = EmptySet()
	if dialect == JOKER {
		markLaceNamespacesAsUsed(env)
		return processInEnvInNS(env, ns, linter_laceData)
	}
	if err := processInEnvInNS(env, ns, linter_cljxData); err != nil {
		return err
	}
	switch dialect {
	case CLJ:
		if err := processInEnvInNS(env, ns, linter_cljData); err != nil {
			return err
		}
	case CLJS:
		if err := processInEnvInNS(env, ns, linter_cljsData); err != nil {
			return err
		}
	}
	removeLaceNamespaces(env)
	return nil
}
