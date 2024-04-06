package core

type NativeSetup func(env *Env) error

var NativeRegistry = map[string]NativeSetup{}

func AddNativeNamespace(name string, setup NativeSetup) {
	NativeRegistry[name] = setup
}

func PopulateNativeNamespacesToEnv(env *Env) error {
	for _, setup := range NativeRegistry {
		err := setup(env)
		if err != nil {
			return err
		}
	}

	return nil
}
