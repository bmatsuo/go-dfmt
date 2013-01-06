package deepfmt

import (
	"fmt"
	"testing"
)

type testStruct struct {
	X *testStruct
	y interface{}
}

// %v and %+v as well as %#v and %#V output the same strings when desired
func TestCompatability(t *testing.T) {
	for i, test := range []interface{}{
		nil,
		(*testStruct)(nil),
		"abc",
		12 + 4i,
		[]string{"a", "b"},
		map[string]int64{"a": 1}, // map iteration order consistant
		map[string]interface{}{"a": &testStruct{}},
		[]interface{}{&testStruct{}},
		&struct{ A, b int }{1, 2},
		struct{ A, b int }{1, 2},
		testStruct{nil, map[string]int{"abc": 2}},
		&testStruct{nil, map[string]int{"abc": 2}},
		testStruct{&testStruct{}, map[string]int{"abc": 2}},
		&testStruct{&testStruct{}, map[string]int{"abc": 2}},
	} {
		std := fmt.Sprintf("%v", test)
		deep := fmt.Sprintf("%0v", NewFormatter(test))
		if std != deep {
			t.Errorf("[%d] # mismatch %q != %q ", i, std, deep)
		}
		std = fmt.Sprintf("%#v", test)
		deep = fmt.Sprintf("%#0v", NewFormatter(test))
		if std != deep {
			t.Errorf("[%d] # mismatch %q != %q ", i, std, deep)
		}
	}
}

func TestDeepfmt(t *testing.T) {
	/*
		Printf("%+v\n", nil)
		Printf("%#+v\n", (*testStruct)(nil))
		Printf("%+v\n", "abc")
		Printf("%#+v\n", "abc")
		Printf("%+v\n", 12+4i)
		Printf("%#+v\n", 12+4i)
		Printf("%+v\n", []string{"a", "b"})
		Printf("%#+v\n", []string{"a", "b"})
		Printf("%# +v\n", []string{"a", "b"})
		Printf("%+v\n", map[string]int64{"a": 1, "b": 2})
		Printf("%#+v\n", map[string]int64{"a": 1, "b": 2})
		Printf("%# +v\n", map[string]int64{"a": 1, "b": 2})
		Printf("%+v\n", &struct{ A, b int }{1, 2})
		Printf("%#+v\n", &struct{ A, b int }{1, 2})
		Printf("%# +v\n", &struct{ A, b int }{1, 2})
		Printf("%+v\n", struct{ A, b int }{1, 2})
		Printf("%#+v\n", struct{ A, b int }{1, 2})
		Printf("%# +v\n", struct{ A, b int }{1, 2})
		Printf("%#+v\n", testStruct{nil, map[string]int{"abc": 2}})
		Printf("%#+v\n", &testStruct{nil, map[string]int{"abc": 2}})
		Printf("%# 0v\n", &testStruct{&testStruct{}, map[string]int{"abc": 2}})

		var x ****testStruct
		x = new(***testStruct)
		*x = new(**testStruct)
		**x = new(*testStruct)
		Printf("%+v\n", x)
		Printf("%#+v\n", x)
	*/
}
