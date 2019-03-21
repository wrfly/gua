package gua

import (
	"flag"
	"strings"
	"testing"
)

func TestGua(t *testing.T) {
	type S struct {
		Str string `name:"sss" default:"str" desc:"just a string"`
		x   struct {
			INT int `default:"6"`
		}
		FFFFFF float64
	}
	s := new(S)

	Parse(s)

	t.Logf("%+v", s)
	flag.VisitAll(func(f *flag.Flag) {
		if strings.HasPrefix(f.Name, "test") {
			return
		}
		t.Logf("name=%s, value=%s, default=%s, desc=%s",
			f.Name, f.Value, f.DefValue, f.Usage)
	})
}
