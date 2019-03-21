package gua

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/wrfly/ecp"
)

type gua struct {
	m map[string]key
}

func (g *gua) getKey(parentName, structName string, tag reflect.StructTag) (key string) {
	desc := tag.Get("desc")
	key = tag.Get("name")
	if key == "" {
		key = fmt.Sprintf("%s.%s", parentName, structName)
	}
	if desc != "" {
		key += "|" + desc
	}
	return
}

func (g *gua) getFlagValue(field reflect.Value, key string) (value string, exist bool) {
	key = strings.Split(key, "|")[0]
	f := g.m[key]
	v := flag.Lookup(f.name)
	if v != nil {
		return v.Value.String(), true
	}
	return "", false
}

type key struct {
	name  string
	value *string
	desc  string
}

func Parse(c interface{}) {
	frog := gua{
		m: make(map[string]key, 20),
	}

	ecp.GetKey = frog.getKey
	ecp.LookupValue = frog.getFlagValue

	ecp.Default(c)

	fullName := reflect.TypeOf(c).String()
	i := strings.LastIndex(fullName, ".")
	structName := fullName[i+1:]

	ecp.Parse(c, structName)

	for _, fullKey := range ecp.List(c, structName) {
		var value, desc string

		x := strings.Split(fullKey, "=")
		name := x[0]
		if len(x) == 2 {
			value = x[1]
		}
		xx := strings.Split(x[0], "|")
		if len(xx) == 2 {
			name = xx[0]
			desc = xx[1]
		}

		frog.m[name] = key{
			name:  name,
			value: &value,
			desc:  desc,
		}
	}

	for _, v := range frog.m {
		v.value = flag.String(v.name, *v.value, v.desc)
	}
	flag.Parse()

	// parse from cli again
	ecp.Parse(c, structName)
}
