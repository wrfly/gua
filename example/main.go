package main

import (
	"fmt"
	"time"

	"github.com/wrfly/gua"
)

type cliFlags struct {
	Name     string `name:"name" default:"wrfly" desc:"just a name"`
	Age      int
	Slice    []string      `desc:"test string slice"`
	SliceInt []int         `desc:"test int slice"`
	Time     time.Duration `desc:"test time duration"`
	Extra    struct {
		Loc   string `default:"home" desc:"location"`
		Valid bool   `default:"true"`
	}
}

func main() {
	cli := new(cliFlags)
	if err := gua.Parse(cli); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", *cli)
}
