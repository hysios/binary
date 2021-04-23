package binary

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"time"
)

var (
	encoder   = &Encoder{}
	ByteOrder = binary.LittleEndian
)

//go:generate Stringer -type Type
type Type uint8

const (
	Uint8 Type = 0xFF - iota
	Uint16
	Uint32
	Uint64
	Int8
	Int16
	Int32
	Int64
	Int
	Uint
	Float32
	Float64
	String
	Bool
	Duration
	Rune
	Timestamp
	Bytes
	Array
	ArrayInt
	ArrayUint
	ArrayFloat
	ArrayFloat32
	ArrayString
	ArrayBool
	Struct
	StructField
	StructValue
	Map
	MapKey
	MapValue
	ElementValue
	ElementRef
)

type Encoder struct {
}

func (enc *Encoder) Encode(val interface{}) ([]byte, error) {
	switch x := val.(type) {
	case uint8, uint16, uint32, uint64, int8, int16, int32,
		int64, int, uint, float32, float64, bool, string, time.Time, time.Duration:
		return enc.encode(x)
	case []byte:
		var buf bytes.Buffer
		buf.WriteByte(byte(Bytes))
		binary.Write(&buf, ByteOrder, uint32(len(x)))
		buf.Write(x)
		return buf.Bytes(), nil
	case []int:
		var buf bytes.Buffer
		buf.WriteByte(byte(ArrayInt))
		binary.Write(&buf, ByteOrder, uint16(len(x)))
		for _, m := range x {
			mb, err := enc.encode(m)
			if err != nil {
				return nil, err
			}
			buf.Write(mb)
		}
		return buf.Bytes(), nil
	case []uint:
		var buf bytes.Buffer
		buf.WriteByte(byte(ArrayUint))
		binary.Write(&buf, ByteOrder, uint16(len(x)))
		for _, m := range x {
			mb, err := enc.encode(m)
			if err != nil {
				return nil, err
			}
			buf.Write(mb)
		}
		return buf.Bytes(), nil
	case []float32:
		var buf bytes.Buffer
		buf.WriteByte(byte(ArrayFloat32))
		binary.Write(&buf, ByteOrder, uint16(len(x)))
		for _, m := range x {
			mb, err := enc.encode(m)
			if err != nil {
				return nil, err
			}
			buf.Write(mb)
		}
		return buf.Bytes(), nil
	case []float64:
		var buf bytes.Buffer
		buf.WriteByte(byte(ArrayFloat))
		binary.Write(&buf, ByteOrder, uint16(len(x)))
		for _, m := range x {
			mb, err := enc.encode(m)
			if err != nil {
				return nil, err
			}
			buf.Write(mb)
		}
		return buf.Bytes(), nil
	case []bool:
		var buf bytes.Buffer
		buf.WriteByte(byte(ArrayBool))
		binary.Write(&buf, ByteOrder, uint16(len(x)))
		for _, m := range x {
			mb, err := enc.encode(m)
			if err != nil {
				return nil, err
			}
			buf.Write(mb)
		}
		return buf.Bytes(), nil
	case []string:
		var buf bytes.Buffer
		buf.WriteByte(byte(ArrayString))
		binary.Write(&buf, ByteOrder, uint16(len(x)))
		for _, m := range x {
			mb, err := enc.encode(m)
			if err != nil {
				return nil, err
			}
			buf.Write(mb)
		}
		return buf.Bytes(), nil
	case []interface{}:
		var buf bytes.Buffer
		buf.WriteByte(byte(Array))
		binary.Write(&buf, ByteOrder, uint16(len(x)))
		for _, m := range x {
			mb, err := enc.encode(m)
			if err != nil {
				return nil, err
			}
			buf.Write(mb)
		}
		return buf.Bytes(), nil
	default:
		v := reflect.ValueOf(x)
		switch v.Kind() {
		case reflect.Struct:
			var buf bytes.Buffer
			buf.WriteByte(byte(Struct))
			binary.Write(&buf, ByteOrder, uint32(v.NumField()))
			t := v.Type()
			for i := 0; i < t.NumField(); i++ {
				var (
					ft = t.Field(i)
					fv = v.Field(i)
				)
				buf.WriteByte(byte(StructField))
				buf.WriteString(ft.Name)
				buf.WriteByte(0)
				fb, err := enc.encode(fv.Interface())
				if err != nil {
					return nil, err
				}
				buf.WriteByte(byte(StructValue))
				buf.Write(fb)
			}
			return buf.Bytes(), nil
		case reflect.Map:
			var buf bytes.Buffer
			buf.WriteByte(byte(Map))
			binary.Write(&buf, ByteOrder, uint32(len(v.MapKeys())))
			iter := v.MapRange()
			for iter.Next() {
				mk := iter.Key()
				mv := iter.Value()
				buf.WriteByte(byte(MapKey))
				kb, err := enc.encode(mk.Interface())
				if err != nil {
					return nil, err
				}
				buf.Write(kb)
				buf.WriteByte(byte(MapValue))
				vb, err := enc.encode(mv.Interface())
				if err != nil {
					return nil, err
				}
				buf.Write(vb)
			}
			return buf.Bytes(), nil
		}

		return nil, errors.New("invalid type")
	}
}

func (enc *Encoder) encode(val interface{}) ([]byte, error) {
	var b = make([]byte, 9)
	switch x := val.(type) {
	case uint8:
		b[0] = uint8(Uint8)
		b[1] = x
		return b[:2], nil
	case uint16:
		b[0] = uint8(Uint16)
		ByteOrder.PutUint16(b[1:], x)
		return b[:3], nil
	case uint32:
		b[0] = uint8(Uint32)
		ByteOrder.PutUint32(b[1:], x)
		return b[:5], nil
	case uint64:
		b[0] = uint8(Uint64)
		ByteOrder.PutUint64(b[1:], x)
		return b[:9], nil
	case int8:
		b[0] = uint8(Int8)
		b[1] = uint8(x)
		return b[:2], nil
	case int16:
		b[0] = uint8(Int16)
		ByteOrder.PutUint16(b[1:], uint16(x))
		return b[:3], nil
	case int32:
		b[0] = uint8(Int32)
		ByteOrder.PutUint32(b[1:], uint32(x))
		return b[:5], nil
	case int64:
		b[0] = uint8(Int64)
		ByteOrder.PutUint64(b[1:], uint64(x))
		return b[:9], nil
	case int:
		b[0] = uint8(Int)
		ByteOrder.PutUint64(b[1:], uint64(x))
		return b[:9], nil
	case uint:
		b[0] = uint8(Uint)
		ByteOrder.PutUint64(b[1:], uint64(x))
		return b[:9], nil
	case float32:
		var buf bytes.Buffer
		buf.WriteByte(byte(Float32))
		if err := binary.Write(&buf, ByteOrder, x); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	case float64:
		var buf bytes.Buffer
		buf.WriteByte(byte(Float64))
		if err := binary.Write(&buf, ByteOrder, x); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	case bool:
		var buf bytes.Buffer
		buf.WriteByte(byte(Bool))
		if err := binary.Write(&buf, ByteOrder, x); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	case string:
		var buf bytes.Buffer
		buf.WriteByte(byte(String))
		buf.WriteString(x)
		buf.WriteByte(0)
		return buf.Bytes(), nil
	case time.Time:
		var buf bytes.Buffer
		buf.WriteByte(byte(Timestamp))
		if err := binary.Write(&buf, ByteOrder, x.UnixNano()); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	case time.Duration:
		var buf bytes.Buffer
		buf.WriteByte(byte(Duration))
		if err := binary.Write(&buf, ByteOrder, int64(x)); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	default:
		return nil, errors.New("invalid type")
	}
}

func (enc *Encoder) EncodeIndex(idx int, val interface{}) ([]byte, error) {
	switch x := val.(type) {
	case uint8, uint16, uint32, uint64, uint,
		int8, int16, int32, int64, int,
		string, bool, float32, float64,
		time.Duration, time.Time:
		b, err := enc.encode(x)
		if err != nil {
			return nil, err
		}
		b = concat(byte(ElementValue), b)
		return b, nil
	default:
		return nil, ErrMustScalarType
	}
}

func (enc *Encoder) EncodeRef(idx int, refId string) ([]byte, error) {
	b, err := enc.encode(refId)
	if err != nil {
		return nil, err
	}
	b = concat(byte(ElementRef), b)
	return b, nil
}

func Encode(val interface{}) ([]byte, error) {
	return encoder.Encode(val)
}

func EncodeIndex(idx int, val interface{}) ([]byte, error) {
	return encoder.EncodeIndex(idx, val)
}

func EncodeRef(idx int, refId string) ([]byte, error) {
	return encoder.EncodeRef(idx, refId)
}

func concat(t byte, b []byte) []byte {
	var o = make([]byte, len(b)+1)
	o[0] = t
	copy(o[1:], b)
	return o
}
