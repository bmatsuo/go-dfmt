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
		printf("%+v\n", nil)
		printf("%#+v\n", (*testStruct)(nil))
		printf("%+v\n", "abc")
		printf("%#+v\n", "abc")
		printf("%+v\n", 12+4i)
		printf("%#+v\n", 12+4i)
		printf("%+v\n", []string{"a", "b"})
		printf("%#+v\n", []string{"a", "b"})
		printf("%# +v\n", []string{"a", "b"})
		printf("%+v\n", map[string]int64{"a": 1, "b": 2})
		printf("%#+v\n", map[string]int64{"a": 1, "b": 2})
		printf("%# +v\n", map[string]int64{"a": 1, "b": 2})
		printf("%+v\n", &struct{ A, b int }{1, 2})
		printf("%#+v\n", &struct{ A, b int }{1, 2})
		printf("%# +v\n", &struct{ A, b int }{1, 2})
		printf("%+v\n", struct{ A, b int }{1, 2})
		printf("%#+v\n", struct{ A, b int }{1, 2})
		printf("%# +v\n", struct{ A, b int }{1, 2})
		printf("%#+v\n", testStruct{nil, map[string]int{"abc": 2}})
		printf("%#+v\n", &testStruct{nil, map[string]int{"abc": 2}})
		printf("%# 0v\n", &testStruct{&testStruct{}, map[string]int{"abc": 2}})

		var x ****testStruct
		x = new(***testStruct)
		*x = new(**testStruct)
		**x = new(*testStruct)
		printf("%+v\n", x)
		printf("%#+v\n", x)
	*/
}
