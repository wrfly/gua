package gua

import (
	"encoding/json"
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
	m map[string]keyInfo
	f *flag.FlagSet
}

type keyInfo struct {
	Name    string `json:"0,omitempty"`
	Desc    string `json:"1,omitempty"`
	Value   string `json:"2,omitempty"`
	IsBool  bool   `json:"3,omitempty"`
	BoolVal bool   `json:"4,omitempty"`
	SubCmd  bool   `json:"5,omitempty"`
}

func (k *keyInfo) Encode() string {
	bs, _ := json.Marshal(k)
	return string(bs)
}

func (k *keyInfo) Decode(str string) {
	json.Unmarshal([]byte(str), k)
}

func firstKeyLower(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

func (g *gua) getKey(pName, sName string, tag reflect.StructTag) string {
	desc := tag.Get("desc")
	name := tag.Get("name")
	pName = firstKeyLower(pName)
	sName = firstKeyLower(sName)
	if name == "" {
		name = sName
		if pName != "" {
			x := new(keyInfo)
			x.Decode(pName)
			name = x.Name + "." + sName
		}
	}

	key := &keyInfo{
		Name: name,
		Desc: desc,
	}
	return key.Encode()
}

func (g *gua) getFlagValue(field reflect.Value, key string) (string, bool) {
	key = strings.Split(key, splitter)[0]
	f := g.m[key]
	v := g.f.Lookup(f.Name)
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

	frog := gua{make(map[string]keyInfo, 20), f}

	ecp.GetKey = frog.getKey
	ecp.LookupValue = frog.getFlagValue

	if err := ecp.Parse(c, ""); err != nil {
		return err
	}

	var err error
	for _, fullKey := range ecp.List(c, "") {
		x := strings.Split(fullKey, "=")
		info := new(keyInfo)
		info.Decode(x[0])
		if len(x) == 2 {
			info.Value = x[1]
		}

		// set bool flag
		info.BoolVal, err = ecp.GetBool(c, info.Name)
		if err == nil {
			info.IsBool = true
		}

		frog.m[info.Name] = *info
	}

	for _, v := range frog.m {
		if v.IsBool {
			v.BoolVal = *f.Bool(v.Name, v.BoolVal, v.Desc)
		} else {
			v.Value = *f.String(v.Name, v.Value, v.Desc)
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
