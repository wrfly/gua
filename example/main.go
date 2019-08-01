package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wrfly/gua"
)

type cliFlags struct {
	Name     string        `name:"name" default:"wrfly" desc:"just a name"`
	Age      int           `desc:"the age"`
	Slice    []string      `desc:"test string slice"`
	SliceInt []int         `desc:"test int slice" default:"1 2 3"`
	Time     time.Duration `desc:"test time duration"`
	Extra    struct {
		Loc   string `default:"home" desc:"location"`
		Valid bool   `default:"true"`
		Debug bool
	}
	Type string `desc:"A|B|C"`
}

func main() {
	cli := new(cliFlags)
	if err := gua.Parse(cli); err != nil {
		panic(err)
	}

	bs, _ := json.MarshalIndent(cli, "", "  ")
	fmt.Printf("%s\n", bs)
}
