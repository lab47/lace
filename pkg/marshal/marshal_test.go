package marshal

import (
	"fmt"
	"testing"

	"github.com/lab47/lace/core"
	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	t.Run("handles types", func(t *testing.T) {
		var e core.Env

		must := func(obj core.Object, err error) core.Object {
			if err != nil {
				panic(err)
			}

			return obj
		}

		input := []core.Object{
			core.MakeSymbol("foo"),
			core.MakeSymbol("user/foo"),
			core.MakeKeyword("bar"),
			core.MakeInt(47),
			core.MakeBoolean(true),
			core.NIL,
			core.NewListFrom(
				core.MakeSymbol("+"),
				core.MakeInt(3),
				core.MakeInt(4),
			),
			core.NewVectorFrom(
				core.MakeSymbol("+"),
				core.MakeSymbol("-"),
				core.NewListFrom(
					core.MakeInt(7),
					core.MakeInt(8),
				),
			),
			must(core.NewHashMap(&e,
				core.MakeSymbol("name"),
				core.MakeString("lace"),
			)),
		}

		for _, d := range input {
			lbl := fmt.Sprintf("type %T", d)
			t.Run(lbl, func(t *testing.T) {
				r := require.New(t)
				b, err := Marshal(d)
				r.NoError(err)

				obj, err := Unmarshal(&e, b)
				r.NoError(err)

				ostr, err := d.ToString(&e, true)
				r.NoError(err)

				str, err := obj.ToString(&e, true)
				r.NoError(err)

				r.True(obj.Equals(&e, d), "didn't round trip: %s != %s", ostr, str)
			})
		}

	})

	t.Run("handles data that seems recursive but isn't", func(t *testing.T) {
		l := core.NewListFrom(core.MakeSymbol("+"))

		d := l.Cons(l)

		var e core.Env

		r := require.New(t)
		b, err := Marshal(d)
		r.NoError(err)

		obj, err := Unmarshal(&e, b)
		r.NoError(err)

		ostr, err := d.ToString(&e, true)
		r.NoError(err)

		str, err := obj.ToString(&e, true)
		r.NoError(err)

		r.True(obj.Equals(&e, d), "didn't round trip: %s != %s", ostr, str)
	})

	t.Run("maintains identity", func(t *testing.T) {
		l := core.NewListFrom(core.MakeSymbol("+"))

		d := core.NewListFrom(l, l)

		var e core.Env

		r := require.New(t)
		b, err := Marshal(d)
		r.NoError(err)

		obj, err := Unmarshal(&e, b)
		r.NoError(err)

		seq := obj.(core.Seq)

		l1, err := seq.First(&e)
		r.NoError(err)

		l2, err := core.Second(&e, seq)
		r.NoError(err)

		var col Collection
		err = decoder.Unmarshal(b, &col)
		r.NoError(err)

		var out any
		err = decoder.Unmarshal(col.Values[0], &out)
		r.NoError(err)

		r.IsType(Collection{}, out)

		err = decoder.Unmarshal(col.Values[1], &out)
		r.NoError(err)

		r.IsType(Ref{}, out)

		// Check identity across the unmarshalling
		r.True(l1 == l2)

		ostr, err := d.ToString(&e, true)
		r.NoError(err)

		str, err := obj.ToString(&e, true)
		r.NoError(err)

		r.True(obj.Equals(&e, d), "didn't round trip: %s != %s", ostr, str)
	})
}
