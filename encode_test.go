package binary

import (
	"reflect"
	"testing"
	"time"
)

func testEncode(val interface{}, typ Type, b []byte, err bool) testEnc {
	return testEnc{
		name:    typ.String(),
		args:    val,
		want:    concat(byte(typ), b),
		wantErr: err,
	}
}

type testEnc struct {
	name    string
	args    interface{}
	want    []byte
	wantErr bool
}

func TestEncoder_Encode(t *testing.T) {
	var tests = []testEnc{
		testEncode(uint8(1), Uint8, []byte{1}, false),
		testEncode(uint16(2), Uint16, []byte{2, 0}, false),
		testEncode(uint32(3), Uint32, []byte{3, 0, 0, 0}, false),
		testEncode(uint64(4), Uint64, []byte{4, 0, 0, 0, 0, 0, 0, 0}, false),
		testEncode(uint(23), Uint, []byte{23, 0, 0, 0, 0, 0, 0, 0}, false),
		testEncode(int(5), Int, []byte{5, 0, 0, 0, 0, 0, 0, 0}, false),
		testEncode(int8(6), Int8, []byte{6}, false),
		testEncode(int16(7), Int16, []byte{7, 0}, false),
		testEncode(int32(8), Int32, []byte{8, 0, 0, 0}, false),
		testEncode(int64(9), Int64, []byte{9, 0, 0, 0, 0, 0, 0, 0}, false),
		testEncode(float32(10.3), Float32, []byte{205, 204, 36, 65}, false),
		testEncode(float64(12.5), Float64, []byte{0, 0, 0, 0, 0, 0, 41, 64}, false),
		testEncode(float64(0.123456789), Float64, []byte{95, 99, 57, 55, 221, 154, 191, 63}, false),
		testEncode(time.Second*10, Duration, []byte{00, 0xE4, 0x0B, 0x54, 0x02, 0x00, 0x00, 0x00}, false),
		testEncode(true, Bool, []byte{1}, false),
		testEncode(false, Bool, []byte{0}, false),
		testEncode("hello", String, []byte{'h', 'e', 'l', 'l', 'o', 0}, false),
		testEncode("世界，你好", String, []byte{0xE4, 0xB8, 0x96, 0xE7, 0x95, 0x8C, 0xEF, 0xBC, 0x8C, 0xE4, 0xBD, 0xA0, 0xE5, 0xA5, 0xBD, 00}, false),
		testEncode(time.Date(2020, 1, 1, 12, 59, 59, 0, time.UTC), Timestamp, []byte{0x00, 0x56, 0x4F, 0xF7, 0xC5, 0xC4, 0xE5, 0x15}, false),
		testEncode([]byte("hello world"), Bytes, []byte{0x0b, 0, 0, 0, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'}, false),
		testEncode([]int{1, 2, 3, 4}, ArrayInt, []byte{0x04, 0,
			0xF7, 1, 0, 0, 0, 0, 0, 0, 0,
			0xF7, 2, 0, 0, 0, 0, 0, 0, 0,
			0xF7, 3, 0, 0, 0, 0, 0, 0, 0,
			0xF7, 4, 0, 0, 0, 0, 0, 0, 0,
		}, false),
		testEncode([]bool{true, false, true, false}, ArrayBool, []byte{0x04, 0, 0xF2, 01, 0xF2, 00, 0xF2, 0x01, 0xF2, 00}, false),
		testEncode([]float32{1.23, 3.45, 6.789}, ArrayFloat32, []byte{0x03, 0,
			0xF5, 0xA4, 0x70, 0x9D, 0x3F,
			0xF5, 0xCD, 0xCC, 0x5C, 0x40,
			0xF5, 0x7D, 0x3F, 0xD9, 0x40,
		}, false),
		testEncode([]float64{1.23, 3.45, 6.789}, ArrayFloat, []byte{0x03, 0,
			0xF4, 0xAE, 0x47, 0xE1, 0x7A, 0x14, 0xAE, 0xF3, 0x3F,
			0xF4, 0x9A, 0x99, 0x99, 0x99, 0x99, 0x99, 0x0B, 0x40,
			0xF4, 0x0E, 0x2D, 0xB2, 0x9D, 0xEF, 0x27, 0x1B, 0x40,
		}, false),
		testEncode([]string{"hello", "world"}, ArrayString, []byte{0x02, 0x00, 0xF3, 'h', 'e', 'l', 'l', 'o', 0x00, 0xF3, 'w', 'o', 'r', 'l', 'd', 0x00}, false),
		testEncode(map[string]interface{}{
			"hello": "world",
			"世界":    "你好",
		}, Map, []byte{2, 0, 0, 0, 0xe2, 0xF3, 'h', 'e', 'l', 'l', 'o', 0x00, 0xE1, 0xF3, 'w', 'o', 'r', 'l', 'd', 0x00,
			0xe2, 0xf3, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c, 0x00, 0xe1, 0xf3, 0xe4, 0xbd, 0xa0, 0xe5, 0xa5, 0xbd, 0x00,
		}, false),
		testEncode(struct {
			Name   string
			Price  float64
			Off    bool
			Amount int
		}{
			Name:   "Apple",
			Price:  13.5,
			Off:    true,
			Amount: 500,
		}, Struct, []byte{0x04, 0x00, 0x00, 0x00, 0xE5, 0x4E, 0x61, 0x6D, 0x65, 0x00, 0xE4, 0xF3, 0x41, 0x70, 0x70, 0x6C, 0x65, 0x00, 0xE5, 0x50, 0x72, 0x69, 0x63, 0x65, 0x00, 0xE4, 0xF4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2B, 0x40, 0xE5, 0x4F, 0x66, 0x66, 0x00, 0xE4, 0xF2, 0x01, 0xE5, 0x41, 0x6D, 0x6F, 0x75, 0x6E, 0x74, 0x00, 0xE4, 0xF7, 0xF4, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, false),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &Encoder{}
			got, err := enc.Encode(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s Encoder.Encode() = %v, want %v", tt.name, got, tt.want)
			}
			t.Logf("%s % X", tt.name, got)
		})
	}
}

func TestEncoder_encode(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		enc     *Encoder
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := &Encoder{}
			got, err := enc.encode(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encoder.encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
