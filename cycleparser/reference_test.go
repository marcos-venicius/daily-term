// Copyright (c) 2019 go-extras

package cycleparser_test

import (
	"reflect"
	"testing"

	"github.com/marcos-venicius/daily-term/cycleparser"
)

type resolverTest struct {
	in             *cycleparser.Value
	out            *cycleparser.Value
	hasUnresolved  bool
	unresolvedRefs []uint64
}

func resolverTests() []resolverTest {
	result := make([]resolverTest, 0)

	result = append(result, resolverTest{
		in: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Bool,
				Value: true,
			},
		},
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Bool,
				Value: true,
			},
		},
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	result = append(result, resolverTest{
		in: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Map,
			Value: map[string]any{
				"test": &cycleparser.Value{
					Refid: 2,
					Kind:  cycleparser.Bool,
					Value: true,
				},
			},
		},
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Map,
			Value: map[string]any{
				"test": &cycleparser.Value{
					Refid: 2,
					Kind:  cycleparser.Bool,
					Value: true,
				},
			},
		},
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	result = append(result, resolverTest{
		in: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Slice,
			Value: []any{
				&cycleparser.Value{
					Refid: 2,
					Kind:  cycleparser.Bool,
					Value: true,
				},
			},
		},
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Slice,
			Value: []any{
				&cycleparser.Value{
					Refid: 2,
					Kind:  cycleparser.Bool,
					Value: true,
				},
			},
		},
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	result = append(result, resolverTest{
		in: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: nil,
		},
		out: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: nil,
		},
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	res1 := &cycleparser.Value{
		Refid: 1,
		Kind:  cycleparser.Ptr,
		Value: &cycleparser.Value{
			Refid: 2,
			Kind:  cycleparser.Ref,
			Value: true,
		},
	}
	res1.Value.(*cycleparser.Value).Value = &cycleparser.Reference{
		Refid: 1,
		Value: res1,
	}
	result = append(result, resolverTest{
		in: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Ref,
				Value: uint64(1),
			},
		},
		out:            res1,
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	res2 := &cycleparser.Value{
		Refid: 1,
		Kind:  cycleparser.Ptr,
		Value: &cycleparser.Value{
			Refid: 2,
			Kind:  cycleparser.Struct,
			Value: map[string]*cycleparser.Value{
				"Name": {
					Refid: 3,
					Kind:  cycleparser.Ptr,
					Value: &cycleparser.Value{
						Refid: 4,
						Kind:  cycleparser.String,
						Value: "Mike",
					},
				},
				"Children": {
					Refid: 5,
					Kind:  cycleparser.Slice,
					Value: []*cycleparser.Value{
						{
							Refid: 6,
							Kind:  cycleparser.Ref,
							Value: nil,
						},
					},
				},
			},
		},
	}
	res2.Value.(*cycleparser.Value).Value.(map[string]*cycleparser.Value)["Children"].Value.([]*cycleparser.Value)[0].Value = &cycleparser.Reference{
		Refid: 1,
		Value: res2,
	}
	result = append(result, resolverTest{
		in: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Struct,
				Value: map[string]*cycleparser.Value{
					"Name": {
						Refid: 3,
						Kind:  cycleparser.Ptr,
						Value: &cycleparser.Value{
							Refid: 4,
							Kind:  cycleparser.String,
							Value: "Mike",
						},
					},
					"Children": {
						Refid: 5,
						Kind:  cycleparser.Slice,
						Value: []*cycleparser.Value{
							{
								Refid: 6,
								Kind:  cycleparser.Ref,
								Value: uint64(1),
							},
						},
					},
				},
			},
		},
		out:            res2,
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	result = append(result, resolverTest{
		in: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Struct,
				Value: map[string]*cycleparser.Value{
					"Name": {
						Refid: 3,
						Kind:  cycleparser.Ptr,
						Value: &cycleparser.Value{
							Refid: 4,
							Kind:  cycleparser.String,
							Value: "Mike",
						},
					},
					"Children": {
						Refid: 5,
						Kind:  cycleparser.Slice,
						Value: []*cycleparser.Value{
							{
								Refid: 6,
								Kind:  cycleparser.Ref,
								Value: uint64(9),
							},
						},
					},
				},
			},
		},
		hasUnresolved:  true,
		unresolvedRefs: []uint64{9},
	})

	res3 := &cycleparser.Value{
		Refid: 1,
		Kind:  cycleparser.Ptr,
		Value: &cycleparser.Value{
			Refid: 2,
			Kind:  cycleparser.Struct,
			Value: map[string]*cycleparser.Value{
				"Sibling": {
					Refid: 3,
					Kind:  cycleparser.Ref,
					Value: uint64(10),
				},
				"Name": {
					Refid: 4,
					Kind:  cycleparser.String,
					Value: "Mike",
				},
				"Parent": {
					Refid: 5,
					Kind:  cycleparser.Ptr,
					Value: &cycleparser.Value{
						Refid: 6,
						Kind:  cycleparser.Struct,
						Value: map[string]*cycleparser.Value{
							"Name": {
								Refid: 7,
								Kind:  cycleparser.String,
								Value: "Frank",
							},
							"Children": {
								Refid: 8,
								Kind:  cycleparser.Slice,
								Value: []*cycleparser.Value{
									{
										Refid: 9,
										Kind:  cycleparser.Ref,
										Value: uint64(1),
									},
									{
										Refid: 10,
										Kind:  cycleparser.Ptr,
										Value: &cycleparser.Value{
											Refid: 11,
											Kind:  cycleparser.Struct,
											Value: map[string]*cycleparser.Value{
												"Name": {
													Refid: 12,
													Kind:  cycleparser.String,
													Value: "Zak",
												},
												"Sibling": {
													Refid: 13,
													Kind:  cycleparser.Ref,
													Value: uint64(1),
												},
												"Parent": {
													Refid: 14,
													Kind:  cycleparser.Ref,
													Value: uint64(5),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	res3.Value.(*cycleparser.Value).
		Value.(map[string]*cycleparser.Value)["Parent"].
		Value.(*cycleparser.Value).
		Value.(map[string]*cycleparser.Value)["Children"].
		Value.([]*cycleparser.Value)[0].
		Value = &cycleparser.Reference{
		Refid: 1,
		Value: res3,
	}
	res3.Value.(*cycleparser.Value).
		Value.(map[string]*cycleparser.Value)["Parent"].
		Value.(*cycleparser.Value).
		Value.(map[string]*cycleparser.Value)["Children"].
		Value.([]*cycleparser.Value)[1].
		Value.(*cycleparser.Value).
		Value.(map[string]*cycleparser.Value)["Parent"].Value = &cycleparser.Reference{
		Refid: 5,
		Value: res3.Value.(*cycleparser.Value).
			Value.(map[string]*cycleparser.Value)["Parent"],
	}
	res3.Value.(*cycleparser.Value).Value.(map[string]*cycleparser.Value)["Parent"].
		Value.(*cycleparser.Value).Value.(map[string]*cycleparser.Value)["Children"].Value.([]*cycleparser.Value)[1].
		Value.(*cycleparser.Value).Value.(map[string]*cycleparser.Value)["Sibling"].Value = &cycleparser.Reference{
		Refid: 1,
		Value: res3,
	}
	res3.Value.(*cycleparser.Value).Value.(map[string]*cycleparser.Value)["Sibling"].Value = &cycleparser.Reference{
		Refid: 10,
		Value: res3.Value.(*cycleparser.Value).
			Value.(map[string]*cycleparser.Value)["Parent"].
			Value.(*cycleparser.Value).
			Value.(map[string]*cycleparser.Value)["Children"].
			Value.([]*cycleparser.Value)[1],
	}
	result = append(result, resolverTest{
		in: &cycleparser.Value{
			Refid: 1,
			Kind:  cycleparser.Ptr,
			Value: &cycleparser.Value{
				Refid: 2,
				Kind:  cycleparser.Struct,
				Value: map[string]*cycleparser.Value{
					"Sibling": {
						Refid: 3,
						Kind:  cycleparser.Ref,
						Value: uint64(10),
					},
					"Name": {
						Refid: 4,
						Kind:  cycleparser.String,
						Value: "Mike",
					},
					"Parent": {
						Refid: 5,
						Kind:  cycleparser.Ptr,
						Value: &cycleparser.Value{
							Refid: 6,
							Kind:  cycleparser.Struct,
							Value: map[string]*cycleparser.Value{
								"Name": {
									Refid: 7,
									Kind:  cycleparser.String,
									Value: "Frank",
								},
								"Children": {
									Refid: 8,
									Kind:  cycleparser.Slice,
									Value: []*cycleparser.Value{
										{
											Refid: 9,
											Kind:  cycleparser.Ref,
											Value: uint64(1),
										},
										{
											Refid: 10,
											Kind:  cycleparser.Ptr,
											Value: &cycleparser.Value{
												Refid: 11,
												Kind:  cycleparser.Struct,
												Value: map[string]*cycleparser.Value{
													"Name": {
														Refid: 12,
														Kind:  cycleparser.String,
														Value: "Zak",
													},
													"Sibling": {
														Refid: 13,
														Kind:  cycleparser.Ref,
														Value: uint64(1),
													},
													"Parent": {
														Refid: 14,
														Kind:  cycleparser.Ref,
														Value: uint64(5),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		out:            res3,
		hasUnresolved:  false,
		unresolvedRefs: []uint64{},
	})

	return result
}

func TestResolver_Resolve(t *testing.T) {
	for i, arg := range resolverTests() {
		r := cycleparser.NewResolver(arg.in)
		err := r.Resolve()
		if err != nil {
			t.Errorf("#%d: Resolver.Resolve() returned an error: %s", i, err.Error())
			continue
		}
		hasUnresolved := r.HasUnresolved()
		unresolvedRefs := r.Unresolved()
		if hasUnresolved != arg.hasUnresolved {
			t.Errorf("#%d: Resolver.HasUnresolved mismatch\nhave: %v\nwant: %v", i, hasUnresolved, arg.hasUnresolved)
		}
		if !reflect.DeepEqual(unresolvedRefs, arg.unresolvedRefs) {
			t.Errorf("#%d: Resolver.Unresolved mismatch\nhave: %#+v\nwant: %#+v", i, unresolvedRefs, arg.unresolvedRefs)
		}
		// don't check if arg.out is nil
		if arg.out != nil && !reflect.DeepEqual(arg.in, arg.out) {
			t.Errorf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, arg.in, arg.out)
		}
	}
}
