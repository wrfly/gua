# guá

Somehow I saw this [repo](https://github.com/alexflint/go-arg) and want
to rewrite a same one using [ecp](https://github.com/wrfly/ecp).

If you want to convert a golang structure into some command line flags,
just `guá`.

It's small and useful for some simple command line tools.

## Example

```golang
package main

import (
    "encoding/json"
    "fmt"
    "time"

    "github.com/wrfly/gua"
)

type cliFlags struct {
    Name     string        `name:"nnnnname" default:"wrfly" desc:"just a name"`
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
```

```bash
# run the example
./gua
{
  "Name": "wrfly",
  "Age": 0,
  "Slice": null,
  "SliceInt": [
    1,
    2,
    3
  ],
  "Time": 0,
  "Extra": {
    "Loc": "home",
    "Valid": false,
    "Debug": false
  },
  "Type": ""
}


# show some help message
./gua -h
Usage of ./gua:
 -age           the age
 -extra.debug   [false]
 -extra.loc     location   [home]
 -extra.valid   [false]
 -nnnnname      just a name   [wrfly]
 -slice         test string slice
 -sliceInt      test int slice   [1 2 3]
 -time          test time duration
 -type          A|B|C


# add some flags
./gua -age 18 -extra.debug -nnnnname gua \
    -slice "hello world" -sliceInt "1 3 5 7" \
    -time 1m -type C
{
  "Name": "gua",
  "Age": 18,
  "Slice": [
    "hello",
    "world"
  ],
  "SliceInt": [
    1,
    3,
    5,
    7
  ],
  "Time": 60000000000,
  "Extra": {
    "Loc": "home",
    "Valid": false,
    "Debug": true
  },
  "Type": "C"
}
```
