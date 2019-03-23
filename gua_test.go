package gua

import (
	"flag"
	"fmt"
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
