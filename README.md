# gua

Somehow I saw this [repo](https://github.com/alexflint/go-arg) and want
to rewrite a same one using [ecp](https://github.com/wrfly/ecp).

If you want to convert a golang structure into some command line flags,
just `gua`.

## Example

```golang
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
```

```txt
➜ /tmp/example -h
Usage of /tmp/example:
  -cliFlags.Age string
  -cliFlags.Extra.Loc string
        location (default "home")
  -cliFlags.Extra.Valid string
  -name string
        just a name (default "wrfly")

➜ /tmp/example
&{Name:wrfly Age:0 Extra:{Loc:home Valid:false}}

➜ /tmp/example \
    -name frog \
    -cliFlags.Age 3 \
    -cliFlags.Extra.Loc pool \
    -cliFlags.Extra.Valid true
&{Name:frog Age:3 Extra:{Loc:pool Valid:true}}
```