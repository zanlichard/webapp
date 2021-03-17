package toolkit

import "testing"

func TestBuiltinType2String(t *testing.T) {
	a := 1.4
	b := 3
	t.Logf("float 2 string:%s", ConvertToString(a))
	t.Logf("int 2 string:%s", ConvertToString(b))

}
