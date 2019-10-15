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

const _SubCMD = "_SubCMD"

type gua struct {
	keys map[string]*keyInfo

	mainSet *flag.FlagSet
	flags   []keyInfo

	subSets  map[string]*flag.FlagSet
	subFlags map[string][]keyInfo
}

type keyInfo struct {
	Name  string `json:"name,omitempty"`
	Usage string `json:"usage,omitempty"`
	Value string `json:"val,omitempty"`
	// bool
	IsBool  bool `json:"bool,omitempty"`
	BoolVal bool `json:"boolVal,omitempty"`
	// sub cmds
	ParentCmd string `json:"parent,omitempty"`
	SubCmd    bool   `json:"isSub,omitempty"`
	Cmd       bool   `json:"isCmd,omitempty"`
}

func (k *keyInfo) fullName() string {
	if k.SubCmd {
		return k.ParentCmd + "." + k.Name
	}
	return k.Name
}

func firstKeyLower(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

func (g *gua) getKey(pName, sName string, tag reflect.StructTag) string {
	info := &keyInfo{
		Usage: tag.Get("desc"),
		Name:  tag.Get("name"),
	}
	if _, ok := g.keys[info.fullName()]; ok {
		return info.fullName()
	}

	pName = firstKeyLower(pName)
	if info.Name == "" {
		info.Name = firstKeyLower(sName)
		if pName != "" {
			info.SubCmd = true
			info.ParentCmd = pName

			if _, ok := g.subSets[pName]; !ok {
				g.subSets[pName] = flag.NewFlagSet(pName, flag.ExitOnError)
			}
			if pInfo, ok := g.keys[pName]; ok {
				pInfo.Cmd = true
				g.keys[pName] = pInfo
			}
		}
	}

	if _, ok := g.keys[info.fullName()]; !ok {
		g.keys[info.fullName()] = info
	}

	return info.fullName()
}

func (g *gua) getFlagValue(field reflect.Value, name string) (string, bool) {
	info := g.keys[name]
	set := g.mainSet
	if info.SubCmd {
		set = g.subSets[info.ParentCmd]
	}

	val := set.Lookup(info.Name)
	if val != nil {
		str := val.Value.String()
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

var visitFunc = func(subCmds, flagUsages *[]string) func(f *flag.Flag) {
	return func(f *flag.Flag) {
		f.DefValue = strings.Trim(f.DefValue, "\"")
		var usage string
		switch {
		case f.Usage == _SubCMD:
			*subCmds = append(*subCmds, f.Name)
		case f.Usage == "" && f.DefValue == "":
			usage = fmt.Sprintf("  -%s", f.Name)
		case f.Usage != "" && f.DefValue == "":
			usage = fmt.Sprintf("  -%s\t%s", f.Name, f.Usage)
		case f.Usage == "" && f.DefValue != "":
			usage = fmt.Sprintf("  -%s\t[%s]", f.Name, f.DefValue)
		case f.Usage != "" && f.DefValue != "":
			usage = fmt.Sprintf("  -%s\t%s\t[%s]", f.Name, f.Usage, f.DefValue)
		}
		if usage != "" {
			*flagUsages = append(*flagUsages, usage)
		}
	}
}

// ParseWithFlagSet can use a flag set you give
func ParseWithFlagSet(c interface{}, f *flag.FlagSet) error {
	frog := gua{
		keys:     make(map[string]*keyInfo, 20),
		mainSet:  f,
		flags:    make([]keyInfo, 0, 10),
		subFlags: make(map[string][]keyInfo, 10),
		subSets:  make(map[string]*flag.FlagSet, 10),
	}

	frog.mainSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		w := tabwriter.NewWriter(
			os.Stderr, 10, 4, 3, ' ',
			tabwriter.StripEscape)
		subCmds := []string{}
		flagUsages := []string{}
		frog.mainSet.VisitAll(visitFunc(&subCmds, &flagUsages))

		if len(subCmds) != 0 {
			w.Write([]byte("cmds:\n"))
		}
		for _, subcmd := range subCmds {
			subCmds := []string{}
			flagUsages := []string{}
			frog.subSets[subcmd].VisitAll(visitFunc(&subCmds, &flagUsages))

			w.Write([]byte(fmt.Sprintf("  %s \n", subcmd)))
			for _, usage := range flagUsages {
				w.Write([]byte(fmt.Sprintf("    %s\n", usage)))
			}
		}

		w.Write([]byte("flags:\n"))
		for _, usage := range flagUsages {
			w.Write([]byte(usage + "\n"))
		}
		w.Flush()
	}

	ecp.GetKey = frog.getKey
	ecp.LookupValue = frog.getFlagValue

	if err := ecp.Parse(c, ""); err != nil {
		return err
	}

	var err error
	for _, fullKey := range ecp.List(c, "") {
		x := strings.Split(fullKey, "=")
		if len(x) == 1 {
			x = append(x, "")
		}
		name, value := x[0], x[1]
		info := frog.keys[name]
		info.Value = value

		// set bool flag
		info.BoolVal, err = ecp.GetBool(c, name)
		if err == nil {
			info.IsBool = true
		}

		if info.SubCmd {
			pName := info.ParentCmd
			// set sub cmd to main flag set
			frog.keys[pName] = &keyInfo{
				Name:      pName,
				Usage:     _SubCMD,
				ParentCmd: pName,
			}
			// append flags to this sub set
			frog.subFlags[pName] = append(frog.subFlags[pName], *info)
		}

		debugJSON("info: %s", info)
		frog.flags = append(frog.flags, *info)
	}

	debug("")

	// main flags
	for name, v := range frog.keys {
		if v.SubCmd {
			continue
		}
		if v.IsBool {
			v.BoolVal = *frog.mainSet.Bool(v.Name, v.BoolVal, v.Usage)
		} else {
			v.Value = *frog.mainSet.String(v.Name, v.Value, v.Usage)
		}
		debugJSON("main "+name+" %s", v)
	}

	debug("")
	// sub flags
	for subName, set := range frog.subSets {
		debugJSON("subSet %s", subName)
		for _, v := range frog.subFlags[subName] {
			if v.IsBool {
				v.BoolVal = *set.Bool(v.Name, v.BoolVal, v.Usage)
			} else {
				v.Value = *set.String(v.Name, v.Value, v.Usage)
			}
			debugJSON("%s", v)
		}
	}

	// parse flags
	frog.mainSet.Parse(os.Args[1:])

	if len(os.Args) >= 2 {
		cmd := os.Args[1]
		if set, ok := frog.subSets[cmd]; ok {
			set.Usage = frog.mainSet.Usage
			set.Parse(os.Args[2:])
		} else {
			fmt.Printf("sub command %s not found\n", cmd)
			frog.mainSet.Usage()
			os.Exit(-1)
		}
	}

	// parse from cli again
	return ecp.Parse(c, "")
}

// Parse the structure
func Parse(c interface{}) error {
	return ParseWithFlagSet(c, flag.CommandLine)
}
