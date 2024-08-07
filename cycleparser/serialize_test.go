// Copyright (c) 2019 go-extras

package cycleparser_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"unsafe"

	"github.com/marcos-venicius/daily-term/cycleparser"
)

type parentSerT struct {
	Name     string
	Children []*childSerT
}

type childSerT struct {
	Name   string
	Parent *parentSerT
}

type selfRefT struct {
	Name string
	Self *selfRefT
}

type interfaceST struct {
	Value        any
	NoSerialize1 bool `json:"-"`
	NoSerialize2 bool `json:"_"`
}

type valueTest struct {
	in  any
	out *cycleparser.Value
	err any
}

func valueTests() []valueTest {
	result := make([]valueTest, 0)

	result = append(result, valueTest{
		in: true,
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Bool,
				Value: true,
			},
		},
	})

	result = append(result, valueTest{
		in: nil,
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: nil,
		},
	})

	result = append(result, valueTest{
		in: "test",
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.String,
				Value: "test",
			},
		},
	})

	s := "test"
	result = append(result, valueTest{
		in: &s,
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.String,
				Value: "test",
			},
		},
	})

	result = append(result, valueTest{
		in: &interfaceST{Value: "test"},
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Struct,
				Value: map[string]*cycleparser.Value{
					"Value": {
						Refid: 3,
						Kind:  cycleparser.String,
						Value: "test",
					},
				},
			},
		},
	})

	result = append(result, valueTest{
		in: any("test"),
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.String,
				Value: "test",
			},
		},
	})

	result = append(result, valueTest{
		in: int(47),
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Int,
				Value: int(47),
			},
		},
	})
	result = append(result, valueTest{
		in: int8(47),
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Int8,
				Value: int8(47),
			},
		},
	})
	result = append(result, valueTest{
		in: int16(47),
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Int16,
				Value: int16(47),
			},
		},
	})
	result = append(result, valueTest{
		in: int32(47),
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Int32,
				Value: int32(47),
			},
		},
	})
	result = append(result, valueTest{
		in: int64(47),
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Int64,
				Value: int64(47),
			},
		},
	})

	result = append(result, valueTest{
		in: uint(47),
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Uint,
				Value: uint(47),
			},
		},
	})
	result = append(result, valueTest{
		in: uint8(47),
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Uint8,
				Value: uint8(47),
			},
		},
	})
	result = append(result, valueTest{
		in: byte(47),
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Uint8,
				Value: uint8(47),
			},
		},
	})
	result = append(result, valueTest{
		in: uint16(47),
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Uint16,
				Value: uint16(47),
			},
		},
	})
	result = append(result, valueTest{
		in: uint32(47),
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Uint32,
				Value: uint32(47),
			},
		},
	})
	result = append(result, valueTest{
		in: float32(47.47),
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Float32,
				Value: float32(47.47),
			},
		},
	})
	result = append(result, valueTest{
		in: float64(47.47),
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Float64,
				Value: float64(47.47),
			},
		},
	})

	result = append(result, valueTest{
		in: &parentSerT{
			Name: "Patrik",
		},
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Struct,
				Value: map[string]*cycleparser.Value{
					"Name": {
						Refid: 3,
						Kind:  cycleparser.String,
						Value: "Patrik",
					},
					"Children": {
						Refid: 4,
						Kind:  cycleparser.Slice,
						Value: []*cycleparser.Value{},
					},
				},
			},
		},
	})

	p1 := &parentSerT{
		Name:     "Patrik",
		Children: nil,
	}
	c1 := &childSerT{
		Name:   "Valentine",
		Parent: p1,
	}
	p1.Children = append(p1.Children, c1)
	result = append(result, valueTest{
		in: p1,
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Struct,
				Value: map[string]*cycleparser.Value{
					"Name": {
						Refid: 3,
						Kind:  cycleparser.String,
						Value: "Patrik",
					},
					"Children": {
						Refid: 4,
						Kind:  cycleparser.Slice,
						Value: []*cycleparser.Value{
							{
								Refid: 5,
								Kind:  cycleparser.Ptr,
								Value: &cycleparser.Value{
									Refid: 6,
									Kind:  cycleparser.Struct,
									Value: map[string]*cycleparser.Value{
										"Name": {
											Refid: 7,
											Kind:  cycleparser.String,
											Value: "Valentine",
										},
										"Parent": {
											Refid: 8,
											Kind:  cycleparser.Ref,
											Value: uint64(1),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	selfRef := &selfRefT{
		Name: "Klark",
	}
	selfRef.Self = selfRef
	result = append(result, valueTest{
		in: selfRef,
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Struct,
				Value: map[string]*cycleparser.Value{
					"Name": {
						Refid: 3,
						Kind:  cycleparser.String,
						Value: "Klark",
					},
					"Self": {
						Refid: 4,
						Kind:  cycleparser.Ref,
						Value: uint64(1),
					},
				},
			},
		},
	})

	result = append(result, valueTest{
		in: map[string]any{
			"Id": uint64(1),
		},
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Map,
				Value: map[string]*cycleparser.Value{
					"Id": {
						Refid: 3,
						Kind:  cycleparser.Uint64,
						Value: uint64(1),
					},
				},
			},
		},
	})

	result = append(result, valueTest{
		in:  uintptr(1),
		err: &cycleparser.InvalidMapperKindError{Kind: "uintptr"},
	})

	result = append(result, valueTest{
		in:  [4]int{1, 2, 3, 4},
		err: &cycleparser.InvalidMapperKindError{Kind: cycleparser.Array},
	})

	result = append(result, valueTest{
		in:  make(chan any),
		err: &cycleparser.InvalidMapperKindError{Kind: "chan"},
	})

	dummy1 := true
	result = append(result, valueTest{
		in:  unsafe.Pointer(&dummy1),
		err: &cycleparser.InvalidMapperKindError{Kind: "unsafe.Pointer"},
	})

	dummy2 := &interfaceST{
		Value:        uintptr(1),
		NoSerialize1: true,
		NoSerialize2: true,
	}
	result = append(result, valueTest{
		in:  &dummy2,
		err: &cycleparser.InvalidMapperKindError{Kind: "uintptr"},
	})

	return result
}

func TestToValue(t *testing.T) {
	for i, arg := range valueTests() {
		v, err := cycleparser.ToValue(arg.in)
		if err != nil {
			if !reflect.DeepEqual(err, arg.err) {
				t.Errorf("#%d: %#+v", i, err)
			}
			// otherwise the error is expected
		} else if !reflect.DeepEqual(v, arg.out) {
			x, err := json.Marshal(v)
			if err != nil {
				t.Fatalf("#%d: %#+v", i, err)
			}
			y, err := json.Marshal(arg.out)
			if err != nil {
				t.Fatalf("#%d: %#+v", i, err)
			}
			t.Errorf("#%d: mismatch\nhave: %v\nwant: %v", i, string(x), string(y))
		}
	}
}
