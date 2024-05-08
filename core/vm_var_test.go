package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompilerVars(t *testing.T) {
	t.Run("assigns args the correct positions", func(t *testing.T) {
		r := require.New(t)

		foo := MakeSymbol("foo")
		bar := MakeSymbol("bar")
		fr := newFrame([]Symbol{foo, bar})

		r.Equal(0, fr.lookup(foo).index)
		r.Equal(1, fr.lookup(bar).index)
	})

	t.Run("can promote a parent fn var to an upval", func(t *testing.T) {
		r := require.New(t)

		foo := MakeSymbol("foo")
		bar := MakeSymbol("bar")
		fr := newFrame([]Symbol{foo, bar})

		parent := fr.lookup(foo)

		chfr := fr.childFrame(nil)

		vb := chfr.lookup(foo)

		r.True(vb.upval)

		r.True(parent != vb)

		r.Exactly(vb, chfr.top.bindings[foo.name])

		t.Run("looking up again returns same varbind", func(t *testing.T) {
			vb2 := chfr.lookup(foo)
			r.Exactly(vb, vb2)
		})
	})

	t.Run("let bindings can be promoted to upvals", func(t *testing.T) {
		r := require.New(t)

		fr := newFrame(nil)

		foo := MakeSymbol("foo")

		lf := fr.bindingFrame(1)

		lvf := lf.set(foo)

		chfr := fr.childFrame(nil)

		chvb := chfr.lookup(foo)

		r.True(lvf.upval)

		chfr.assignUpvals()

		r.Equal(0, chvb.upidx)

		fr.popBindingFrame()

		fr.assignUpvals()

		r.Equal(0, lvf.upidx)
	})

	t.Run("imported upvals come before let upvals", func(t *testing.T) {
		r := require.New(t)

		a := MakeSymbol("a")

		fr := newFrame([]Symbol{a})

		foo := MakeSymbol("foo")

		topa := fr.lookup(a)

		lf := fr.bindingFrame(1)

		lvf := lf.set(foo)

		chfr := fr.childFrame(nil)

		cha := chfr.lookup(a)
		r.True(topa.upval)

		chvb := chfr.lookup(foo)
		r.True(lvf.upval)

		chfr.assignUpvals()

		r.Equal(0, cha.upidx)
		r.Equal(1, chvb.upidx)

		fr.popBindingFrame()

		fr.assignUpvals()

		r.Equal(0, topa.upidx)
		r.Equal(1, lvf.upidx)
	})

	t.Run("accessing a closure var from deeper first, then shallow doesn't make 2 vars", func(t *testing.T) {
		r := require.New(t)

		foo := MakeSymbol("foo")

		fr := newFrame([]Symbol{foo})

		chfr := fr.childFrame(nil)

		f1 := chfr.bindingFrame(0)

		f2 := chfr.bindingFrame(0)

		l1 := chfr.lookup(foo)

		r.Empty(f1.bindings)
		r.Nil(f2.bindings["foo"])

		chfr.popBindingFrame()
		r.Exactly(f1, chfr.top)

		l2 := chfr.lookup(foo)

		r.Nil(f1.bindings["foo"])

		// Compare as pointers
		r.True(l1 == l2)
	})

	t.Run("using instructions are updated when upvals are assigned", func(t *testing.T) {
		r := require.New(t)

		foo := MakeSymbol("foo")

		fr := newFrame([]Symbol{foo})

		l := fr.lookup(foo)

		si := l.read()
		r.Equal(GetLocal, si.Op)

		chfr := fr.childFrame(nil)

		b := chfr.lookup(foo)

		i := b.read()
		r.Equal(GetUpval, i.Op)
		r.Equal(int32(-1), i.A0)

		chfr.assignUpvals()

		r.Equal(GetUpval, i.Op)
		r.Equal(int32(0), i.A0)

		ri := l.refUpval()
		r.Equal(RefUpval, ri.Op)
		r.Equal(int32(-1), ri.A0)

		fr.assignUpvals()

		r.Equal(GetUpval, si.Op)
		r.Equal(int32(0), si.A0)

		r.Equal(RefUpval, ri.Op)
		r.Equal(int32(0), si.A0)
	})

	t.Run("upval indexes are local to each function", func(t *testing.T) {
		r := require.New(t)

		foo := MakeSymbol("foo")
		bar := MakeSymbol("bar")

		fr := newFrame([]Symbol{foo, bar})

		ch1 := fr.childFrame(nil)
		ch1.lookup(foo)
		b1 := ch1.lookup(bar)

		ch2 := ch1.childFrame(nil)
		b2 := ch2.lookup(bar)

		ch2.assignUpvals()

		r.Equal(0, b2.upidx)

		ri := b1.refUpval()

		ch1.assignUpvals()

		r.Equal(1, b1.upidx)
		r.Equal(int32(1), ri.A0)
	})

	t.Run("frame reports the sames of the upvals that need to be imported", func(t *testing.T) {
		r := require.New(t)

		foo := MakeSymbol("foo")

		fr := newFrame([]Symbol{foo})

		ch := fr.childFrame(nil)
		ch.lookup(foo)

		imported := ch.importedNames()

		r.Len(imported, 1)
		r.Equal("foo", imported[0].Name())
	})

	t.Run("closing a child frame adds upvals to the parent", func(t *testing.T) {
		r := require.New(t)

		foo := MakeSymbol("foo")

		fr := newFrame([]Symbol{foo})

		ch1 := fr.childFrame(nil)
		ch2 := ch1.childFrame(nil)

		ch2.lookup(foo)

		ch1vbs := ch2.closeFrame()

		r.Len(ch1vbs, 1)

		r.Equal("foo", ch1vbs[0].Name())
	})

	t.Run("peer scopes reuse same local slots (critical for args", func(t *testing.T) {
		r := require.New(t)

		foo := MakeSymbol("foo")

		fr := newFrame(nil)

		s1 := fr.bindingFrame(1)
		r.Equal(0, s1.set(foo).index)

		fr.popBindingFrame()

		s2 := fr.bindingFrame(1)
		r.Equal(0, s2.set(foo).index)

		fr.popBindingFrame()

	})

}
