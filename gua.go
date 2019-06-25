package gua

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/wrfly/ecp"
)

type gua struct {
	m map[string]key
	f *flag.FlagSet
}

func (g *gua) getKey(pName, sName string, tag reflect.StructTag) string {
	desc := tag.Get("desc")
	key := tag.Get("name")
	if key == "" {
		if pName == "" {
			key = sName
		} else {
			key = pName + "." + sName
		}
	}
	if desc != "" {
		key += "|" + desc
	}
	return strings.ToLower(key)
}

func (g *gua) getFlagValue(field reflect.Value, key string) (string, bool) {
	key = strings.Split(key, "|")[0]
	f := g.m[key]
	v := g.f.Lookup(f.name)
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

func ParseWithFlagSet(c interface{}, f *flag.FlagSet) error {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		w := tabwriter.NewWriter(
			os.Stderr, 10, 4, 3, ' ',
			tabwriter.StripEscape)
		f.VisitAll(func(f *flag.Flag) {
			var format string
			switch {
			case f.Usage == "" && f.DefValue == "":
				format = fmt.Sprintf(" -%s", f.Name)
			case f.Usage != "" && f.DefValue == "":
				format = fmt.Sprintf(" -%s\t%s", f.Name, f.Usage)
			case f.Usage == "" && f.DefValue != "":
				format = fmt.Sprintf(" -%s\t[%s]", f.Name, f.DefValue)
			case f.Usage != "" && f.DefValue != "":
				format = fmt.Sprintf(" -%s\t%s\t[%s]", f.Name, f.Usage, f.DefValue)
			}
			w.Write([]byte(format + "\n"))
		})
		w.Flush()
	}
	f.Usage = flag.Usage

	frog := gua{make(map[string]key, 20), f}

	ecp.GetKey = frog.getKey
	ecp.LookupValue = frog.getFlagValue

	if err := ecp.Parse(c, ""); err != nil {
		return err
	}

	for _, fullKey := range ecp.List(c, "") {
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
		v.value = f.String(v.name, *v.value, v.desc)
	}
	f.Parse(os.Args[1:])

	// parse from cli again
	return ecp.Parse(c, "")
}

// Parse the structure
func Parse(c interface{}) error {
	return ParseWithFlagSet(c, flag.CommandLine)
}
