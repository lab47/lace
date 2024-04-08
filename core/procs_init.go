//go:build !fast_init
// +build !fast_init

package core

var privateMeta Map

func init() {
	v, err := EmptyArrayMap().Assoc(criticalKeywords.private, Boolean{B: true})
	if err != nil {
		panic(err)
	}

	privateMeta = v.(Map)
}

func initEnv(env *Env) error {
	var err error
	intern := func(env *Env, name string, proc ProcFn, procName string) {
		if err != nil {
			return
		}

		var vr *Var
		vr, err = env.CoreNamespace.Intern(env, MakeSymbol(name))
		if err != nil {
			return
		}

		vr.Value = Proc{Fn: proc, Name: procName}
		vr.isPrivate = true
		vr.meta = privateMeta
	}

	env.CoreNamespace.InternVar(env, "*assert*", Boolean{B: true},
		MakeMeta(nil, "When set to logical false, assert is a noop. Defaults to true.", "1.0"))

	intern(env, "list__", procList, "procList")
	intern(env, "cons__", procCons, "procCons")
	intern(env, "first__", procFirst, "procFirst")
	intern(env, "next__", procNext, "procNext")
	intern(env, "rest__", procRest, "procRest")
	intern(env, "conj__", procConj, "procConj")
	intern(env, "seq__", procSeq, "procSeq")
	intern(env, "instance?__", procIsInstance, "procIsInstance")
	intern(env, "assoc__", procAssoc, "procAssoc")
	intern(env, "meta__", procMeta, "procMeta")
	intern(env, "with-meta__", procWithMeta, "procWithMeta")
	intern(env, "=__", procEquals, "procEquals")
	intern(env, "count__", procCount, "procCount")
	intern(env, "subvec__", procSubvec, "procSubvec")
	intern(env, "cast__", procCast, "procCast")
	intern(env, "vec__", procVec, "procVec")
	intern(env, "hash-map__", procHashMap, "procHashMap")
	intern(env, "hash-set__", procHashSet, "procHashSet")
	intern(env, "str__", procStr, "procStr")
	intern(env, "symbol__", procSymbol, "procSymbol")
	intern(env, "gensym__", procGensym, "procGensym")
	intern(env, "keyword__", procKeyword, "procKeyword")
	intern(env, "apply__", procApply, "procApply")
	intern(env, "lazy-seq__", procLazySeq, "procLazySeq")
	intern(env, "delay__", procDelay, "procDelay")
	intern(env, "force__", procForce, "procForce")
	intern(env, "identical__", procIdentical, "procIdentical")
	intern(env, "compare__", procCompare, "procCompare")
	intern(env, "zero?__", procIsZero, "procIsZero")
	intern(env, "int__", procInt, "procInt")
	intern(env, "nth__", procNth, "procNth")
	intern(env, "<__", procLt, "procLt")
	intern(env, "<=__", procLte, "procLte")
	intern(env, ">__", procGt, "procGt")
	intern(env, ">=__", procGte, "procGte")
	intern(env, "==__", procEq, "procEq")
	intern(env, "inc'__", procIncEx, "procIncEx")
	intern(env, "inc__", procInc, "procInc")
	intern(env, "dec'__", procDecEx, "procDecEx")
	intern(env, "dec__", procDec, "procDec")
	intern(env, "add'__", procAddEx, "procAddEx")
	intern(env, "add__", procAdd, "procAdd")
	intern(env, "multiply'__", procMultiplyEx, "procMultiplyEx")
	intern(env, "multiply__", procMultiply, "procMultiply")
	intern(env, "divide__", procDivide, "procDivide")
	intern(env, "subtract'__", procSubtractEx, "procSubtractEx")
	intern(env, "subtract__", procSubtract, "procSubtract")
	intern(env, "max__", procMax, "procMax")
	intern(env, "min__", procMin, "procMin")
	intern(env, "pos__", procIsPos, "procIsPos")
	intern(env, "neg__", procIsNeg, "procIsNeg")
	intern(env, "quot__", procQuot, "procQuot")
	intern(env, "rem__", procRem, "procRem")
	intern(env, "bit-not__", procBitNot, "procBitNot")
	intern(env, "bit-and__", procBitAnd, "procBitAnd")
	intern(env, "bit-or__", procBitOr, "procBitOr")
	intern(env, "bit-xor_", procBitXor, "procBitXor")
	intern(env, "bit-and-not__", procBitAndNot, "procBitAndNot")
	intern(env, "bit-clear__", procBitClear, "procBitClear")
	intern(env, "bit-set__", procBitSet, "procBitSet")
	intern(env, "bit-flip__", procBitFlip, "procBitFlip")
	intern(env, "bit-test__", procBitTest, "procBitTest")
	intern(env, "bit-shift-left__", procBitShiftLeft, "procBitShiftLeft")
	intern(env, "bit-shift-right__", procBitShiftRight, "procBitShiftRight")
	intern(env, "unsigned-bit-shift-right__", procUnsignedBitShiftRight, "procUnsignedBitShiftRight")
	intern(env, "peek__", procPeek, "procPeek")
	intern(env, "pop__", procPop, "procPop")
	intern(env, "contains?__", procContains, "procContains")
	intern(env, "get__", procGet, "procGet")
	intern(env, "dissoc__", procDissoc, "procDissoc")
	intern(env, "disj__", procDisj, "procDisj")
	intern(env, "find__", procFind, "procFind")
	intern(env, "keys__", procKeys, "procKeys")
	intern(env, "vals__", procVals, "procVals")
	intern(env, "rseq__", procRseq, "procRseq")
	intern(env, "name__", procName, "procName")
	intern(env, "namespace__", procNamespace, "procNamespace")
	intern(env, "find-var__", procFindVar, "procFindVar")
	intern(env, "sort__", procSort, "procSort")
	intern(env, "eval__", procEval, "procEval")
	intern(env, "type__", procType, "procType")
	intern(env, "num__", procNumber, "procNumber")
	intern(env, "double__", procDouble, "procDouble")
	intern(env, "char__", procChar, "procChar")
	intern(env, "boolean__", procBoolean, "procBoolean")
	intern(env, "numerator__", procNumerator, "procNumerator")
	intern(env, "denominator__", procDenominator, "procDenominator")
	intern(env, "bigint__", procBigInt, "procBigInt")
	intern(env, "bigfloat__", procBigFloat, "procBigFloat")
	intern(env, "pr__", procPr, "procPr")
	intern(env, "pprint__", procPprint, "procPprint")
	intern(env, "newline__", procNewline, "procNewline")
	intern(env, "flush__", procFlush, "procFlush")
	intern(env, "read__", procRead, "procRead")
	intern(env, "read-line__", procReadLine, "procReadLine")
	intern(env, "reader-read-line__", procReaderReadLine, "procReaderReadLine")
	intern(env, "read-string__", procReadString, "procReadString")
	intern(env, "nano-time__", procNanoTime, "procNanoTime")
	intern(env, "macroexpand-1__", procMacroexpand1, "procMacroexpand1")
	intern(env, "load-string__", procLoadString, "procLoadString")
	intern(env, "find-ns__", procFindNamespace, "procFindNamespace")
	intern(env, "create-ns__", procCreateNamespace, "procCreateNamespace")
	intern(env, "inject-ns__", procInjectNamespace, "procInjectNamespace")
	intern(env, "remove-ns__", procRemoveNamespace, "procRemoveNamespace")
	intern(env, "all-ns__", procAllNamespaces, "procAllNamespaces")
	intern(env, "ns-name__", procNamespaceName, "procNamespaceName")
	intern(env, "ns-map__", procNamespaceMap, "procNamespaceMap")
	intern(env, "ns-unmap__", procNamespaceUnmap, "procNamespaceUnmap")
	intern(env, "var-ns__", procVarNamespace, "procVarNamespace")
	intern(env, "ns-initialized?__", procIsNamespaceInitialized, "procIsNamespaceInitialized")
	intern(env, "refer__", procRefer, "procRefer")
	intern(env, "alias__", procAlias, "procAlias")
	intern(env, "ns-aliases__", procNamespaceAliases, "procNamespaceAliases")
	intern(env, "ns-unalias__", procNamespaceUnalias, "procNamespaceUnalias")
	intern(env, "var-get__", procVarGet, "procVarGet")
	intern(env, "var-set__", procVarSet, "procVarSet")
	intern(env, "ns-resolve__", procNsResolve, "procNsResolve")
	intern(env, "array-map__", procArrayMap, "procArrayMap")
	intern(env, "buffer__", procBuffer, "procBuffer")
	intern(env, "buffered-reader__", procBufferedReader, "procBufferedReader")
	intern(env, "ex-info__", procExInfo, "procExInfo")
	intern(env, "ex-data__", procExData, "procExData")
	intern(env, "ex-cause__", procExCause, "procExCause")
	intern(env, "ex-message__", procExMessage, "procExMessage")
	intern(env, "regex__", procRegex, "procRegex")
	intern(env, "re-seq__", procReSeq, "procReSeq")
	intern(env, "re-find__", procReFind, "procReFind")
	intern(env, "rand__", procRand, "procRand")
	intern(env, "special-symbol?__", procIsSpecialSymbol, "procIsSpecialSymbol")
	intern(env, "subs__", procSubs, "procSubs")
	intern(env, "intern__", procIntern, "procIntern")
	intern(env, "set-meta__", procSetMeta, "procSetMeta")
	intern(env, "atom__", procAtom, "procAtom")
	intern(env, "deref__", procDeref, "procDeref")
	intern(env, "swap__", procSwap, "procSwap")
	intern(env, "swap-vals__", procSwapVals, "procSwapVals")
	intern(env, "reset__", procReset, "procReset")
	intern(env, "reset-vals__", procResetVals, "procResetVals")
	intern(env, "alter-meta__", procAlterMeta, "procAlterMeta")
	intern(env, "reset-meta__", procResetMeta, "procResetMeta")
	intern(env, "empty__", procEmpty, "procEmpty")
	intern(env, "bound?__", procIsBound, "procIsBound")
	intern(env, "format__", procFormat, "procFormat")
	intern(env, "load-file__", procLoadFile, "procLoadFile")
	intern(env, "load-lib-from-path__", procLoadLibFromPath, "procLoadLibFromPath")
	intern(env, "reduce-kv__", procReduceKv, "procReduceKv")
	intern(env, "slurp__", procSlurp, "procSlurp")
	intern(env, "spit__", procSpit, "procSpit")
	intern(env, "shuffle__", procShuffle, "procShuffle")
	intern(env, "realized?__", procIsRealized, "procIsRealized")
	intern(env, "derive-info__", procDeriveInfo, "procDeriveInfo")
	intern(env, "lace-version__", procLaceVersion, "procLaceVersion")

	intern(env, "hash__", procHash, "procHash")

	intern(env, "index-of__", procIndexOf, "procIndexOf")
	intern(env, "lib-path__", procLibPath, "procLibPath")
	intern(env, "intern-fake-var__", procInternFakeVar, "procInternFakeVar")
	intern(env, "parse__", procParse, "procParse")
	intern(env, "inc-problem-count__", procIncProblemCount, "procIncProblemCount")
	intern(env, "types__", procTypes, "procTypes")
	intern(env, "go__", procGo, "procGo")
	intern(env, "<!__", procReceive, "procReceive")
	intern(env, ">!__", procSend, "procSend")
	intern(env, "chan__", procCreateChan, "procCreateChan")
	intern(env, "close!__", procCloseChan, "procCloseChan")

	intern(env, "go-spew__", procGoSpew, "procGoSpew")
	intern(env, "verbosity-level__", procVerbosityLevel, "procVerbosityLevel")
	intern(env, "exit__", procExit, "procExit")

	return err
}

func lateInitializations() {
	// none needed for !fast_init
}
