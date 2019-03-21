package main

import (
	"fmt"

	"github.com/wrfly/gua"
)

type cliFlags struct {
	Name  string `name:"name" default:"wrfly" desc:"just a name"`
	Age   int
	Extra struct {
		Loc   string `default:"home" desc:"location"`
		Valid bool
	}
}

func main() {
	cli := new(cliFlags)
	gua.Parse(cli)
	fmt.Printf("%+v\n", cli)
}
