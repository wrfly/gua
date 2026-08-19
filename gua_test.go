package gua

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestGua(t *testing.T) {
	type Frog struct {
		Name     string `name:"frog.name" default:"gua" desc:"just a frog"`
		Location string
		Age      uint8
		Alive    bool `default:"true"`
		Extra    struct {
			HaveGF bool `desc:"a single frog"`
		}

		// basic types
		String  string
		Bool    bool
		Int     int
		Int8    int8
		Int32   int32
		Int64   int64
		Uint    uint
		Uint8   uint8
		Uint16  uint16
		Uint32  uint32
		Uint64  uint64
		Float32 float32
		Float64 float64

		// pointers
		StringPtr  *string
		BoolPtr    *bool
		IntPtr     *int
		Int8Ptr    *int8
		Int32Ptr   *int32
		Int64Ptr   *int64
		UintPtr    *uint
		Uint8Ptr   *uint8
		Uint16Ptr  *uint16
		Uint32Ptr  *uint32
		Uint64Ptr  *uint64
		Float32Ptr *float32
		Float64Ptr *float64

		// slices
		StringSlice  []string
		BoolSlice    []bool
		IntSlice     []int
		Int8Slice    []int8
		Int32Slice   []int32
		Int64Slice   []int64
		UintSlice    []uint
		Uint8Slice   []uint8
		Uint16Slice  []uint16
		Uint32Slice  []uint32
		Uint64Slice  []uint64
		Float32Slice []float32
		Float64Slice []float64
	}
	frog := new(Frog)

	newSet := flag.NewFlagSet(
		fmt.Sprint(time.Now().Nanosecond()),
		flag.ContinueOnError)

	if err := ParseWithFlagSet(frog, newSet); err != nil {
		t.Fatal("???")
	}
	newSet.Usage()

	// test error
	type ParseError struct {
		Int int `default:"int"`
	}
	e := new(ParseError)
	if err := Parse(e); err == nil {
		t.Fatal("???")
	}

}

// the flags actually have to end up in the structure: the test above only
// checks that parsing does not fail, which kept passing through an ecp
// upgrade that changed how values are set
func TestParseFlags(t *testing.T) {
	type conf struct {
		Name  string `default:"wrfly" desc:"a name"`
		Age   int
		Alive bool `default:"true"`
		Debug bool
		Slice []string
		Ints  []int `default:"1 2 3"`
		Time  time.Duration
		Extra struct {
			Loc string `default:"home"`
		}
	}

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"gua.test",
		"-name", "bob",
		"-age", "30",
		"-alive=false", // a default of true has to be overridable
		"-debug",
		"-slice", "a b c",
		"-time", "3d", // ecp duration syntax
		"-extra.loc", "office",
	}

	c := new(conf)
	set := flag.NewFlagSet("gua", flag.ContinueOnError)
	if err := ParseWithFlagSet(c, set); err != nil {
		t.Fatal(err)
	}

	switch {
	case c.Name != "bob":
		t.Errorf("name: %q", c.Name)
	case c.Age != 30:
		t.Errorf("age: %d", c.Age)
	case c.Alive:
		t.Error("alive should have been turned off by the flag")
	case !c.Debug:
		t.Error("debug should have been turned on by the flag")
	case len(c.Slice) != 3 || c.Slice[2] != "c":
		t.Errorf("slice: %q", c.Slice)
	case len(c.Ints) != 3 || c.Ints[0] != 1:
		t.Errorf("ints: %v", c.Ints)
	case c.Time != 72*time.Hour:
		t.Errorf("time: %v", c.Time)
	case c.Extra.Loc != "office":
		t.Errorf("extra.loc: %q", c.Extra.Loc)
	}
}

// defaults survive when no flag is given
func TestParseDefaults(t *testing.T) {
	type conf struct {
		Name  string `default:"wrfly"`
		Alive bool   `default:"true"`
		Ints  []int  `default:"1 2 3"`
	}

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"gua.test"}

	c := new(conf)
	set := flag.NewFlagSet("gua", flag.ContinueOnError)
	if err := ParseWithFlagSet(c, set); err != nil {
		t.Fatal(err)
	}

	switch {
	case c.Name != "wrfly":
		t.Errorf("name: %q", c.Name)
	case !c.Alive:
		t.Error("alive should have kept its default")
	case len(c.Ints) != 3:
		t.Errorf("ints: %v", c.Ints)
	}
}
