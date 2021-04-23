package binary

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testDec struct {
	name    string
	b       []byte
	args    interface{}
	want    interface{}
	wantErr bool
}

func testDecode(typ Type, b []byte, val interface{}, want interface{}, err bool) testDec {

	return testDec{
		name:    typ.String(),
		b:       concat(byte(typ), b),
		args:    val,
		want:    want,
		wantErr: err,
	}
}

func TestDecoder_Decode(t *testing.T) {
	bs := make([]byte, 0)
	is := make([]int, 0)
	uis := make([]uint, 0)
	fs := make([]float32, 0)
	ffs := make([]float64, 0)
	ss := make([]string, 0)
	mm := make(map[string]interface{})

	type testStruct struct {
		Name   string
		Price  float64
		Off    bool
		Amount int
	}

	var tests = []testDec{
		testDecode(Uint8, []byte{1}, new(uint8), uint8(1), false),
		testDecode(Uint16, []byte{2, 0}, new(uint16), uint16(2), false),
		testDecode(Uint32, []byte{3, 0, 0, 0}, new(uint32), uint32(3), false),
		testDecode(Uint64, []byte{4, 0, 0, 0, 0, 0, 0, 0}, new(uint64), uint64(4), false),
		testDecode(Int8, []byte{255}, new(int8), int8(-1), false),
		testDecode(Int16, []byte{254, 255}, new(int16), int16(-2), false),
		testDecode(Int32, []byte{253, 255, 255, 255}, new(int32), int32(-3), false),
		testDecode(Int64, []byte{252, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, new(int64), int64(-4), false),
		testDecode(Duration, []byte{0x00, 0x5E, 0xD0, 0xB2, 0, 0, 0, 0}, new(time.Duration), 3*time.Second, false),
		testDecode(Float32, []byte{205, 204, 36, 65}, new(float32), float32(10.3), false),
		testDecode(Float64, []byte{0, 0, 0, 0, 0, 0, 41, 64}, new(float64), float64(12.5), false),
		testDecode(Float64, []byte{95, 99, 57, 55, 221, 154, 191, 63}, new(float64), float64(0.123456789), false),
		testDecode(Bool, []byte{1}, new(bool), true, false),
		testDecode(Bool, []byte{0}, new(bool), false, false),
		testDecode(String, []byte{'h', 'e', 'l', 'l', 'o', 0}, new(string), "hello", false),
		testDecode(String, []byte{0xE4, 0xB8, 0x96, 0xE7, 0x95, 0x8C, 0xEF, 0xBC, 0x8C, 0xE4, 0xBD, 0xA0, 0xE5, 0xA5, 0xBD, 00}, new(string), "世界，你好", false),
		testDecode(Timestamp, []byte{0xD2, 0x5A, 0x4F, 0xF7, 0xC5, 0xC4, 0xE5, 0x15}, new(time.Time), time.Date(2020, 1, 1, 12, 59, 59, 1234, time.UTC), false),
		testDecode(Bytes, []byte{0x0b, 0, 0, 0, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'}, &bs, []byte("hello world"), false),
		testDecode(ArrayInt, []byte{0x04, 0,
			0xF7, 1, 0, 0, 0, 0, 0, 0, 0,
			0xF7, 2, 0, 0, 0, 0, 0, 0, 0,
			0xF7, 3, 0, 0, 0, 0, 0, 0, 0,
			0xF7, 4, 0, 0, 0, 0, 0, 0, 0,
		}, &is, []int{1, 2, 3, 4}, false),
		testDecode(ArrayUint, []byte{0x04, 0,
			0xF6, 1, 0, 0, 0, 0, 0, 0, 0,
			0xF6, 2, 0, 0, 0, 0, 0, 0, 0,
			0xF6, 3, 0, 0, 0, 0, 0, 0, 0,
			0xF6, 4, 0, 0, 0, 0, 0, 0, 0,
		}, &uis, []uint{1, 2, 3, 4}, false),
		testDecode(ArrayFloat32, []byte{0x03, 0,
			0xF5, 0xA4, 0x70, 0x9D, 0x3F,
			0xF5, 0xCD, 0xCC, 0x5C, 0x40,
			0xF5, 0x7D, 0x3F, 0xD9, 0x40,
		}, &fs, []float32{1.23, 3.45, 6.789}, false),
		testDecode(ArrayFloat, []byte{0x03, 0,
			0xF4, 0xAE, 0x47, 0xE1, 0x7A, 0x14, 0xAE, 0xF3, 0x3F,
			0xF4, 0x9A, 0x99, 0x99, 0x99, 0x99, 0x99, 0x0B, 0x40,
			0xF4, 0x0E, 0x2D, 0xB2, 0x9D, 0xEF, 0x27, 0x1B, 0x40,
		}, &ffs, []float64{1.23, 3.45, 6.789}, false),
		testDecode(ArrayString, []byte{0x02, 0x00, 0xF3, 'h', 'e', 'l', 'l', 'o', 0x00, 0xF3, 'w', 'o', 'r', 'l', 'd', 0x00}, &ss, []string{"hello", "world"}, false),
		testDecode(Map, []byte{4, 0, 0, 0, 0xe2, 0xF3, 'h', 'e', 'l', 'l', 'o', 0x00, 0xE1, 0xF3, 'w', 'o', 'r', 'l', 'd', 0x00,
			0xe2, 0xf3, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c, 0x00, 0xe1, 0xf3, 0xe4, 0xbd, 0xa0, 0xe5, 0xa5, 0xbd, 0x00,
			0xe2, 0xf3, 'c', 'o', 'u', 'n', 't', 0, 0xe1, 0xf7, 0x1, 0, 0, 0, 0, 0, 0, 0,
			0xe2, 0xf3, 'p', 'r', 'i', 'c', 'e', 0, 0xe1, 0xf4, 0, 0, 0, 0, 0, 0, 41, 64,
		}, &mm, map[string]interface{}{
			"hello": "world",
			"世界":    "你好",
			"count": 1,
			"price": 12.5,
		}, false),
		testDecode(Struct, []byte{0x04, 0x00, 0x00, 0x00, 0xE5, 0x4E, 0x61, 0x6D, 0x65, 0x00, 0xE4, 0xF3, 0x41, 0x70, 0x70, 0x6C, 0x65, 0x00, 0xE5, 0x50, 0x72, 0x69, 0x63, 0x65, 0x00, 0xE4, 0xF4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2B, 0x40, 0xE5, 0x4F, 0x66, 0x66, 0x00, 0xE4, 0xF2, 0x01, 0xE5, 0x41, 0x6D, 0x6F, 0x75, 0x6E, 0x74, 0x00, 0xE4, 0xF7, 0xF4, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, new(testStruct), testStruct{
			Name:   "Apple",
			Price:  13.5,
			Off:    true,
			Amount: 500,
		}, false),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := &Decoder{}
			if _, err := dec.Decode(tt.b, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("%s Decoder.Decode() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}

			val := reflect.ValueOf(tt.args)
			assert.Equal(t, val.Elem().Interface(), tt.want)
			t.Logf("got val %v", tt.args)
		})
	}
}

type testDecto struct {
	name    string
	b       []byte
	want    interface{}
	wantErr bool
}

func testDecodeTo(typ Type, b []byte, val interface{}, wantErr bool) testDecto {
	return testDecto{
		b:       concat(byte(typ), b),
		want:    val,
		wantErr: wantErr,
	}
}

func TestDecodeTo(t *testing.T) {

	tests := []struct {
		name    string
		b       []byte
		want    interface{}
		wantErr bool
	}{
		testDecodeTo(Uint8, []byte{1}, uint8(1), false),
		testDecodeTo(Uint16, []byte{2, 0}, uint16(2), false),
		testDecodeTo(Uint32, []byte{4, 0, 0, 0}, uint32(4), false),
		testDecodeTo(Uint64, []byte{8, 0, 0, 0, 0, 0, 0, 0}, uint64(8), false),
		testDecodeTo(String, []byte{'h', 'e', 'l', 'l', 'o', 0}, "hello", false),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeTo(tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
