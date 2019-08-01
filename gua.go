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

var splitter = "#|-_-|#" // tricky word

type gua struct {
	m map[string]key
	f *flag.FlagSet
}

func firstKeyLower(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

func (g *gua) getKey(pName, sName string, tag reflect.StructTag) string {
	desc := tag.Get("desc")
	key := tag.Get("name")
	pName = firstKeyLower(pName)
	sName = firstKeyLower(sName)
	if key == "" {
		if pName == "" {
			key = sName
		} else {
			key = pName + "." + sName
		}
	}
	if desc != "" {
		key += splitter + desc
	}
	return key
}

func (g *gua) getFlagValue(field reflect.Value, key string) (string, bool) {
	key = strings.Split(key, splitter)[0]
	f := g.m[key]
	v := g.f.Lookup(f.name)
	if v != nil {
		str := v.Value.String()
		str = strings.Trim(str, "\"")
		return str, true
	}
	return "", false
}

func init() {
	// reset lookup key function
	ecp.LookupKey = func(original, _, _ string) string { return original }
}

type key struct {
	name  string
	value *string
	desc  string

	isBool bool
	boolV  *bool
}

// ParseWithNew use a fresh new flag set
func ParseWithNew(c interface{}, name string) error {
	return ParseWithFlagSet(c, flag.NewFlagSet(name, flag.ExitOnError))
}

// ParseWithFlagSet can use a flag set you give
func ParseWithFlagSet(c interface{}, f *flag.FlagSet) error {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		w := tabwriter.NewWriter(
			os.Stderr, 10, 4, 3, ' ',
			tabwriter.StripEscape)
		f.VisitAll(func(f *flag.Flag) {
			f.DefValue = strings.Trim(f.DefValue, "\"")
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
		var (
			value, desc string
			isBool      bool
		)

		x := strings.Split(fullKey, "=")
		name := x[0]
		if len(x) == 2 {
			value = x[1]
		}

		// set bool flag
		boolV, err := ecp.GetBool(c, name)
		if err == nil {
			isBool = true
		}

		xx := strings.Split(x[0], splitter)
		if len(xx) >= 2 {
			name = xx[0]
			desc = xx[1]
		}

		frog.m[name] = key{
			name:   name,
			value:  &value,
			desc:   desc,
			isBool: isBool,
			boolV:  &boolV,
		}
	}

	for _, v := range frog.m {
		if v.isBool {
			v.boolV = f.Bool(v.name, *v.boolV, v.desc)
		} else {
			v.value = f.String(v.name, *v.value, v.desc)
		}
	}
	f.Parse(os.Args[1:])

	// parse from cli again
	return ecp.Parse(c, "")
}

// Parse the structure
func Parse(c interface{}) error {
	return ParseWithFlagSet(c, flag.CommandLine)
}
