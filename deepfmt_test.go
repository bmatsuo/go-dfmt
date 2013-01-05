package deepfmt

import (
	"fmt"
	"testing"
)

type testStruct struct{ X, y interface{} }

// %v and %V as well as %#v and %#V output the same strings when desired
func TestCompatability(t *testing.T) {
	for i, test := range []interface{}{
		"abc",
		12 + 4i,
		[]string{"a", "b"},
		map[string]int64{"a": 1}, // map iteration order consistant
		&struct{ A, b int }{1, 2},
		struct{ A, b int }{1, 2},
		testStruct{1, map[string]int{"abc": 2}},
		&testStruct{1, map[string]int{"abc": 2}},
	} {
		std := fmt.Sprintf("%v", test)
		deep := fmt.Sprintf("%V", &Formatter{v: test})
		if std != deep {
			t.Errorf("[%d] mismatch %q != %q ", i, std, deep)
		}

		std = fmt.Sprintf("%#v", test)
		deep = fmt.Sprintf("%#V", &Formatter{v: test})
		if std != deep {
			t.Errorf("[%d] # mismatch %q != %q ", i, std, deep)
		}
	}

}

func TestDeepfmt(t *testing.T) {
	/*
		Printf("%V\n", "abc")
		Printf("%#V\n", "abc")
		Printf("%V\n", 12+4i)
		Printf("%#V\n", 12+4i)
		Printf("%V\n", []string{"a", "b"})
		Printf("%#V\n", []string{"a", "b"})
		Printf("%# V\n", []string{"a", "b"})
		Printf("%V\n", map[string]int64{"a": 1, "b": 2})
		Printf("%#V\n", map[string]int64{"a": 1, "b": 2})
		Printf("%# V\n", map[string]int64{"a": 1, "b": 2})
		Printf("%V\n", &struct{ A, b int }{1, 2})
		Printf("%#V\n", &struct{ A, b int }{1, 2})
		Printf("%# V\n", &struct{ A, b int }{1, 2})
		Printf("%V\n", struct{ A, b int }{1, 2})
		Printf("%#V\n", struct{ A, b int }{1, 2})
		Printf("%# V\n", struct{ A, b int }{1, 2})

		var x ****int
		x = new(***int)
		*x = new(**int)
		**x = new(*int)
		***x = new(int)
		Printf("%V\n", x)
		Printf("%#V\n", x)
	*/
}
