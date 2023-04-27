// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bson

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsonrw/bsonrwtest"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

func TestBasicDecode(t *testing.T) {
	t.Parallel()

	for _, tc := range unmarshalingTestCases() {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := reflect.New(tc.sType).Elem()
			vr := bsonrw.NewBSONDocumentReader(tc.data)
			reg := DefaultRegistry
			decoder, err := reg.LookupDecoder(reflect.TypeOf(got))
			noerr(t, err)
			err = decoder.DecodeValue(bsoncodec.DecodeContext{Registry: reg}, vr, got)
			noerr(t, err)
			assert.Equal(t, tc.want, got.Addr().Interface(), "Results do not match.")
		})
	}
}

func TestDecoderv2(t *testing.T) {
	t.Parallel()

	t.Run("Decode", func(t *testing.T) {
		t.Parallel()

		for _, tc := range unmarshalingTestCases() {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got := reflect.New(tc.sType).Interface()
				vr := bsonrw.NewBSONDocumentReader(tc.data)
				dec, err := NewDecoderWithContext(bsoncodec.DecodeContext{Registry: DefaultRegistry}, vr)
				noerr(t, err)
				err = dec.Decode(got)
				noerr(t, err)
				assert.Equal(t, tc.want, got, "Results do not match.")
			})
		}
		t.Run("lookup error", func(t *testing.T) {
			t.Parallel()

			type certainlydoesntexistelsewhereihope func(string, string) string
			// Avoid unused code lint error.
			_ = certainlydoesntexistelsewhereihope(func(string, string) string { return "" })

			cdeih := func(string, string) string { return "certainlydoesntexistelsewhereihope" }
			dec, err := NewDecoder(bsonrw.NewBSONDocumentReader([]byte{}))
			noerr(t, err)
			want := bsoncodec.ErrNoDecoder{Type: reflect.TypeOf(cdeih)}
			got := dec.Decode(&cdeih)
			assert.Equal(t, want, got, "Received unexpected error.")
		})
		t.Run("Unmarshaler", func(t *testing.T) {
			t.Parallel()

			testCases := []struct {
				name    string
				err     error
				vr      bsonrw.ValueReader
				invoked bool
			}{
				{
					"error",
					errors.New("Unmarshaler error"),
					&bsonrwtest.ValueReaderWriter{BSONType: bsontype.EmbeddedDocument, Err: bsonrw.ErrEOD, ErrAfter: bsonrwtest.ReadElement},
					true,
				},
				{
					"copy error",
					errors.New("copy error"),
					&bsonrwtest.ValueReaderWriter{Err: errors.New("copy error"), ErrAfter: bsonrwtest.ReadDocument},
					false,
				},
				{
					"success",
					nil,
					&bsonrwtest.ValueReaderWriter{BSONType: bsontype.EmbeddedDocument, Err: bsonrw.ErrEOD, ErrAfter: bsonrwtest.ReadElement},
					true,
				},
			}

			for _, tc := range testCases {
				tc := tc

				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()

					unmarshaler := &testUnmarshaler{err: tc.err}
					dec, err := NewDecoder(tc.vr)
					noerr(t, err)
					got := dec.Decode(unmarshaler)
					want := tc.err
					if !compareErrors(got, want) {
						t.Errorf("Did not receive expected error. got %v; want %v", got, want)
					}
					if unmarshaler.invoked != tc.invoked {
						if tc.invoked {
							t.Error("Expected to have UnmarshalBSON invoked, but it wasn't.")
						} else {
							t.Error("Expected UnmarshalBSON to not be invoked, but it was.")
						}
					}
				})
			}

			t.Run("Unmarshaler/success bsonrw.ValueReader", func(t *testing.T) {
				t.Parallel()

				want := bsoncore.BuildDocument(nil, bsoncore.AppendDoubleElement(nil, "pi", 3.14159))
				unmarshaler := &testUnmarshaler{}
				vr := bsonrw.NewBSONDocumentReader(want)
				dec, err := NewDecoder(vr)
				noerr(t, err)
				err = dec.Decode(unmarshaler)
				noerr(t, err)
				got := unmarshaler.data
				if !bytes.Equal(got, want) {
					t.Errorf("Did not unmarshal properly. got %v; want %v", got, want)
				}
			})
		})
	})
	t.Run("NewDecoder", func(t *testing.T) {
		t.Parallel()

		t.Run("error", func(t *testing.T) {
			t.Parallel()

			_, got := NewDecoder(nil)
			want := errors.New("cannot create a new Decoder with a nil ValueReader")
			if !cmp.Equal(got, want, cmp.Comparer(compareErrors)) {
				t.Errorf("Was expecting error but got different error. got %v; want %v", got, want)
			}
		})
		t.Run("success", func(t *testing.T) {
			t.Parallel()

			got, err := NewDecoder(bsonrw.NewBSONDocumentReader([]byte{}))
			noerr(t, err)
			if got == nil {
				t.Errorf("Was expecting a non-nil Decoder, but got <nil>")
			}
		})
	})
	t.Run("NewDecoderWithContext", func(t *testing.T) {
		t.Parallel()

		t.Run("errors", func(t *testing.T) {
			t.Parallel()

			dc := bsoncodec.DecodeContext{Registry: DefaultRegistry}
			_, got := NewDecoderWithContext(dc, nil)
			want := errors.New("cannot create a new Decoder with a nil ValueReader")
			if !cmp.Equal(got, want, cmp.Comparer(compareErrors)) {
				t.Errorf("Was expecting error but got different error. got %v; want %v", got, want)
			}
		})
		t.Run("success", func(t *testing.T) {
			t.Parallel()

			got, err := NewDecoderWithContext(bsoncodec.DecodeContext{}, bsonrw.NewBSONDocumentReader([]byte{}))
			noerr(t, err)
			if got == nil {
				t.Errorf("Was expecting a non-nil Decoder, but got <nil>")
			}
			dc := bsoncodec.DecodeContext{Registry: DefaultRegistry}
			got, err = NewDecoderWithContext(dc, bsonrw.NewBSONDocumentReader([]byte{}))
			noerr(t, err)
			if got == nil {
				t.Errorf("Was expecting a non-nil Decoder, but got <nil>")
			}
		})
	})
	t.Run("Decode doesn't zero struct", func(t *testing.T) {
		t.Parallel()

		type foo struct {
			Item  string
			Qty   int
			Bonus int
		}
		var got foo
		got.Item = "apple"
		got.Bonus = 2
		data := docToBytes(D{{"item", "canvas"}, {"qty", 4}})
		vr := bsonrw.NewBSONDocumentReader(data)
		dec, err := NewDecoder(vr)
		noerr(t, err)
		err = dec.Decode(&got)
		noerr(t, err)
		want := foo{Item: "canvas", Qty: 4, Bonus: 2}
		assert.Equal(t, want, got, "Results do not match.")
	})
	t.Run("Reset", func(t *testing.T) {
		t.Parallel()

		vr1, vr2 := bsonrw.NewBSONDocumentReader([]byte{}), bsonrw.NewBSONDocumentReader([]byte{})
		dc := bsoncodec.DecodeContext{Registry: DefaultRegistry}
		dec, err := NewDecoderWithContext(dc, vr1)
		noerr(t, err)
		if dec.vr != vr1 {
			t.Errorf("Decoder should use the value reader provided. got %v; want %v", dec.vr, vr1)
		}
		err = dec.Reset(vr2)
		noerr(t, err)
		if dec.vr != vr2 {
			t.Errorf("Decoder should use the value reader provided. got %v; want %v", dec.vr, vr2)
		}
	})
	t.Run("SetContext", func(t *testing.T) {
		t.Parallel()

		dc1 := bsoncodec.DecodeContext{Registry: DefaultRegistry}
		dc2 := bsoncodec.DecodeContext{Registry: NewRegistryBuilder().Build()}
		dec, err := NewDecoderWithContext(dc1, bsonrw.NewBSONDocumentReader([]byte{}))
		noerr(t, err)
		if !reflect.DeepEqual(dec.dc, dc1) {
			t.Errorf("Decoder should use the Registry provided. got %v; want %v", dec.dc, dc1)
		}
		err = dec.SetContext(dc2)
		noerr(t, err)
		if !reflect.DeepEqual(dec.dc, dc2) {
			t.Errorf("Decoder should use the Registry provided. got %v; want %v", dec.dc, dc2)
		}
	})
	t.Run("SetRegistry", func(t *testing.T) {
		t.Parallel()

		r1, r2 := DefaultRegistry, NewRegistryBuilder().Build()
		dc1 := bsoncodec.DecodeContext{Registry: r1}
		dc2 := bsoncodec.DecodeContext{Registry: r2}
		dec, err := NewDecoder(bsonrw.NewBSONDocumentReader([]byte{}))
		noerr(t, err)
		if !reflect.DeepEqual(dec.dc, dc1) {
			t.Errorf("Decoder should use the Registry provided. got %v; want %v", dec.dc, dc1)
		}
		err = dec.SetRegistry(r2)
		noerr(t, err)
		if !reflect.DeepEqual(dec.dc, dc2) {
			t.Errorf("Decoder should use the Registry provided. got %v; want %v", dec.dc, dc2)
		}
	})
	t.Run("DecodeToNil", func(t *testing.T) {
		t.Parallel()

		data := docToBytes(D{{"item", "canvas"}, {"qty", 4}})
		vr := bsonrw.NewBSONDocumentReader(data)
		dec, err := NewDecoder(vr)
		noerr(t, err)

		var got *D
		err = dec.Decode(got)
		if err != ErrDecodeToNil {
			t.Fatalf("Decode error mismatch; expected %v, got %v", ErrDecodeToNil, err)
		}
	})
	t.Run("DefaultDocuemntD embedded map as empty interface", func(t *testing.T) {
		t.Parallel()

		type someMap map[string]interface{}

		in := make(someMap)
		in["foo"] = map[string]interface{}{"bar": "baz"}

		bytes, err := Marshal(in)
		if err != nil {
			t.Fatal(err)
		}

		var bsonOut someMap
		dec, err := NewDecoder(bsonrw.NewBSONDocumentReader(bytes))
		if err != nil {
			t.Fatal(err)
		}
		dec.DefaultDocumentM()
		if err := dec.Decode(&bsonOut); err != nil {
			t.Fatal(err)
		}

		// Ensure that interface{}-typed top-level data is converted to the document type.
		bsonOutType := reflect.TypeOf(bsonOut)
		inType := reflect.TypeOf(in)
		assert.Equal(t, inType, bsonOutType, "expected %v to equal %v", inType.String(), bsonOutType.String())

		// Ensure that the embedded type is a primitive map.
		mType := reflect.TypeOf(primitive.M{})
		bsonFooOutType := reflect.TypeOf(bsonOut["foo"])
		assert.Equal(t, mType, bsonFooOutType, "expected %v to equal %v", mType.String(), bsonFooOutType.String())
	})
	t.Run("DefaultDocuemntD for decoding into interface{} alias", func(t *testing.T) {
		t.Parallel()

		var in interface{} = map[string]interface{}{"bar": "baz"}

		bytes, err := Marshal(in)
		if err != nil {
			t.Fatal(err)
		}

		var bsonOut interface{}
		dec, err := NewDecoder(bsonrw.NewBSONDocumentReader(bytes))
		if err != nil {
			t.Fatal(err)
		}
		dec.DefaultDocumentD()
		if err := dec.Decode(&bsonOut); err != nil {
			t.Fatal(err)
		}

		// Ensure that interface{}-typed top-level data is converted to the document type.
		dType := reflect.TypeOf(primitive.D{})
		bsonOutType := reflect.TypeOf(bsonOut)
		assert.Equal(t, dType, bsonOutType,
			"expected %v to equal %v", dType.String(), bsonOutType.String())
	})
	t.Run("DefaultDocuemntD for decoding into non-interface{} alias", func(t *testing.T) {
		t.Parallel()

		var in interface{} = map[string]interface{}{"bar": "baz"}

		bytes, err := Marshal(in)
		if err != nil {
			t.Fatal(err)
		}

		var bsonOut struct{}
		dec, err := NewDecoder(bsonrw.NewBSONDocumentReader(bytes))
		if err != nil {
			t.Fatal(err)
		}
		dec.DefaultDocumentD()
		if err := dec.Decode(&bsonOut); err != nil {
			t.Fatal(err)
		}

		// Ensure that typed top-level data is not converted to the document type.
		dType := reflect.TypeOf(primitive.D{})
		bsonOutType := reflect.TypeOf(bsonOut)
		assert.NotEqual(t, dType, bsonOutType,
			"expected %v to not equal %v", dType.String(), bsonOutType.String())
	})
	t.Run("DefaultDocumentD for deep struct values", func(t *testing.T) {
		t.Parallel()

		type emb struct {
			Foo map[int]interface{} `bson:"foo"`
		}

		objID := primitive.NewObjectID()

		in := emb{
			Foo: map[int]interface{}{
				1: map[string]interface{}{"bar": "baz"},
				2: map[int]interface{}{
					3: map[string]interface{}{"bar": "baz"},
				},
				4: map[primitive.ObjectID]interface{}{
					objID: map[string]interface{}{"bar": "baz"},
				},
			},
		}

		bytes, err := Marshal(in)
		if err != nil {
			t.Fatal(err)
		}

		dec, err := NewDecoder(bsonrw.NewBSONDocumentReader(bytes))
		if err != nil {
			t.Fatal(err)
		}

		dec.DefaultDocumentD()

		var out emb
		if err := dec.Decode(&out); err != nil {
			t.Fatal(err)
		}

		mType := reflect.TypeOf(primitive.M{})
		bsonOutType := reflect.TypeOf(out)
		assert.NotEqual(t, mType, bsonOutType,
			"expected %v to not equal %v", mType.String(), bsonOutType.String())

		want := emb{
			Foo: map[int]interface{}{
				1: primitive.D{{Key: "bar", Value: "baz"}},
				2: primitive.D{{Key: "3", Value: primitive.D{{Key: "bar", Value: "baz"}}}},
				4: primitive.D{{Key: objID.Hex(), Value: primitive.D{{Key: "bar", Value: "baz"}}}},
			},
		}

		assert.Equal(t, want, out, "expected %v, got %v", want, out)
	})
}

type testUnmarshaler struct {
	invoked bool
	err     error
	data    []byte
}

func (tu *testUnmarshaler) UnmarshalBSON(d []byte) error {
	tu.invoked = true
	tu.data = d
	return tu.err
}
