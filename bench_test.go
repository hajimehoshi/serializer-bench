// Copyright 2017 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main_test

import (
	"encoding/json"
	"testing"

	"github.com/vmihailenco/msgpack"
)

type Foo struct {
	Int    int    `json:"int"`
	String string `json:"string"`
	Bars   []*Bar `json:"bars"`
}

func (f *Foo) EncodeMsgpack(enc *msgpack.Encoder) error {
	if err := enc.EncodeInt(int64(f.Int)); err != nil {
		return err
	}
	if err := enc.EncodeString(f.String); err != nil {
		return err
	}
	if err := enc.Encode(f.Bars); err != nil {
		return err
	}
	return nil
}

func (f *Foo) DecodeMsgpack(dec *msgpack.Decoder) error {
	var err error
	f.Int, err = dec.DecodeInt()
	if err != nil {
		return err
	}
	f.String, err = dec.DecodeString()
	if err != nil {
		return err
	}
	if err := dec.Decode(&f.Bars); err != nil {
		return err
	}
	return nil
}

type Bar struct {
	Floats  []float64      `json:"floats"`
	Strings []string       `json:"strings"`
	Map     map[string]int `json:"map"`
}

func (b *Bar) EncodeMsgpack(enc *msgpack.Encoder) error {
	if err := enc.EncodeArrayLen(len(b.Floats)); err != nil {
		return err
	}
	for _, v := range b.Floats {
		if err := enc.EncodeFloat64(v); err != nil {
			return err
		}
	}
	if err := enc.EncodeArrayLen(len(b.Strings)); err != nil {
		return err
	}
	for _, v := range b.Strings {
		if err := enc.EncodeString(v); err != nil {
			return err
		}
	}
	if err := enc.EncodeMapLen(len(b.Map)); err != nil {
		return err
	}
	for k, v := range b.Map {
		if err := enc.EncodeString(k); err != nil {
			return err
		}
		if err := enc.EncodeInt(int64(v)); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bar) DecodeMsgpack(dec *msgpack.Decoder) error {
	var err error
	l := 0

	l, err = dec.DecodeArrayLen()
	if err != nil {
		return err
	}
	b.Floats = make([]float64, l)
	for i := 0; i < l; i++ {
		b.Floats[i], err = dec.DecodeFloat64()
		if err != nil {
			return err
		}
	}

	l, err = dec.DecodeArrayLen()
	if err != nil {
		return err
	}
	b.Strings = make([]string, l)
	for i := 0; i < l; i++ {
		b.Strings[i], err = dec.DecodeString()
		if err != nil {
			return err
		}
	}

	l, err = dec.DecodeMapLen()
	b.Map = map[string]int{}
	for i := 0; i < l; i++ {
		k, err := dec.DecodeString()
		if err != nil {
			return err
		}
		v, err := dec.DecodeInt()
		if err != nil {
			return err
		}
		b.Map[k] = v
	}
	return nil
}

func input() *Foo {
	return &Foo{
		Int:    123,
		String: "hello",
		Bars: []*Bar{
			&Bar{
				Floats:  []float64{4, 5, 6},
				Strings: []string{"world", "こんにちは", "世界"},
				Map: map[string]int{
					"one":   1,
					"two":   2,
					"three": 3,
				},
			},
			&Bar{
				Floats:  []float64{3.14159, 1.41421},
				Strings: []string{"", "\n", "\r"},
				Map: map[string]int{
					"digit": 1234567890,
				},
			},
		},
	}
}

func BenchmarkJSON(b *testing.B) {
	f := input()
	var ff *Foo
	for i := 0; i < b.N; i++ {
		bin, err := json.Marshal(f)
		if err != nil {
			b.Fatal(err)
		}
		if err := json.Unmarshal(bin, &ff); err != nil {
			b.Fatal(err)
		}
		f = ff
	}
	want := 1234567890
	got := f.Bars[1].Map["digit"]
	if want != got {
		b.Fatalf("want %d but %d", want, got)
	}
}

func BenchmarkMsgpack(b *testing.B) {
	f := input()
	var ff *Foo
	for i := 0; i < b.N; i++ {
		bin, err := msgpack.Marshal(f)
		if err != nil {
			b.Fatal(err)
		}
		if err := msgpack.Unmarshal(bin, &ff); err != nil {
			b.Fatal(err)
		}
		f = ff
	}
	want := 1234567890
	got := f.Bars[1].Map["digit"]
	if want != got {
		b.Fatalf("want %d but %d", want, got)
	}
}
