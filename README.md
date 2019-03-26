# guá

Somehow I saw this [repo](https://github.com/alexflint/go-arg) and want
to rewrite a same one using [ecp](https://github.com/wrfly/ecp).

If you want to convert a golang structure into some command line flags,
just `guá`.

## Example

```golang
package main

import (
    "flag"
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
```

```txt
➜ /tmp/example -h
Usage of /tmp/example:
 -cliFlags.Age
 -cliFlags.Extra.Loc     location   [home]
 -cliFlags.Extra.Valid   [true]
 -cliFlags.Slice         test string slice
 -cliFlags.SliceInt      test int slice
 -cliFlags.Time          test time duration
 -name                   just a name   [wrfly]

➜ /tmp/example
{Name:wrfly Age:0 Slice:[] SliceInt:[] Time:0s Extra:{Loc:home Valid:true}}

➜ /tmp/example \  
    -name frog \
    -cliFlags.Age 3 \
    -cliFlags.Extra.Loc pool \
    -cliFlags.Extra.Valid true \
    -cliFlags.Slice "a b c" \
    -cliFlags.SliceInt "1 2 3" \
    -cliFlags.Time "365d"
{Name:frog Age:3 Slice:[a b c] SliceInt:[1 2 3] Time:8760h0m0s Extra:{Loc:pool Valid:true}}
```
