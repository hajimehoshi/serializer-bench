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
	enc.EncodeMapLen(3)
	enc.EncodeString("Int")
	enc.EncodeInt(int64(f.Int))
	enc.EncodeString("String")
	enc.EncodeString(f.String)
	enc.EncodeString("Bars")
	enc.Encode(f.Bars)
	return nil
}

func (f *Foo) DecodeMsgpack(dec *msgpack.Decoder) error {
	l, _ := dec.DecodeMapLen()
	for i := 0; i < l; i++ {
		key, _ := dec.DecodeString()
		switch key {
		case "Int":
			f.Int, _ = dec.DecodeInt()
		case "String":
			f.String, _ = dec.DecodeString()
		case "Bars":
			dec.Decode(&f.Bars)
		}
	}
	return nil
}

type Bar struct {
	Floats  []float64      `json:"floats"`
	Strings []string       `json:"strings"`
	Map     map[string]int `json:"map"`
}

func (b *Bar) EncodeMsgpack(enc *msgpack.Encoder) error {
	enc.EncodeMapLen(3)

	enc.EncodeString("Floats")
	enc.EncodeArrayLen(len(b.Floats))
	for _, v := range b.Floats {
		enc.EncodeFloat64(v)
	}

	enc.EncodeString("Strings")
	enc.EncodeArrayLen(len(b.Strings))
	for _, v := range b.Strings {
		enc.EncodeString(v)
	}

	enc.EncodeString("Map")
	enc.EncodeMapLen(len(b.Map))
	for k, v := range b.Map {
		enc.EncodeString(k)
		enc.EncodeInt(int64(v))
	}

	return nil
}

func (b *Bar) DecodeMsgpack(dec *msgpack.Decoder) error {
	l, _ := dec.DecodeMapLen()
	for i := 0; i < l; i++ {
		key, _ := dec.DecodeString()
		switch key {
		case "Floats":
			ll, _ := dec.DecodeArrayLen()
			b.Floats = make([]float64, ll)
			for i := 0; i < ll; i++ {
				b.Floats[i], _ = dec.DecodeFloat64()
			}
		case "Strings":
			ll, _ := dec.DecodeArrayLen()
			b.Strings = make([]string, ll)
			for i := 0; i < ll; i++ {
				b.Strings[i], _ = dec.DecodeString()
			}
		case "Map":
			ll, _ := dec.DecodeMapLen()
			b.Map = map[string]int{}
			for i := 0; i < ll; i++ {
				k, _ := dec.DecodeString()
				v, _ := dec.DecodeInt()
				b.Map[k] = v
			}
		}
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
