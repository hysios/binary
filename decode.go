package binary

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"time"
)

var decoder = &Decoder{}

type Decoder struct {
}

func (dec *Decoder) Decode(b []byte, val interface{}) (int, error) {
	if len(b) == 0 {
		return 0, ErrBufTooSmall
	}

	typ := Type(b[0])
	switch typ {
	case Uint8, Uint16, Uint32, Uint64, Uint, Int8, Int16,
		Int32, Int64, Int, Bool, Float32, Float64, Timestamp, Duration,
		String, Bytes:
		return dec.decode(b, val)
	case ArrayInt:
		var (
			l      = ByteOrder.Uint16(b[1:])
			v      = reflect.ValueOf(val)
			nv     = make([]int, l)
			offset = 3
			count  = offset
		)
		v = reflect.Indirect(v)
		for i := 0; i < int(l); i++ {
			n, err := dec.decode(b[offset+i*9:], &nv[i])
			if err != nil {
				return 0, err
			}
			count += n
		}

		v.Set(reflect.ValueOf(nv))
		return count, nil
	case ArrayUint:
		var (
			l      = ByteOrder.Uint16(b[1:])
			v      = reflect.ValueOf(val)
			nv     = make([]uint, l)
			offset = 3
			count  = offset
		)
		v = reflect.Indirect(v)
		for i := 0; i < int(l); i++ {
			n, err := dec.decode(b[offset+i*9:], &nv[i])
			if err != nil {
				return 0, err
			}
			count += n
		}

		v.Set(reflect.ValueOf(nv))
		return count, nil
	case ArrayFloat32:
		var (
			l      = ByteOrder.Uint16(b[1:])
			v      = reflect.ValueOf(val)
			nv     = make([]float32, l)
			offset = 3
			count  = offset
		)
		v = reflect.Indirect(v)
		for i := 0; i < int(l); i++ {
			n, err := dec.decode(b[offset+i*5:], &nv[i])
			if err != nil {
				return 0, err
			}
			count += n
		}

		v.Set(reflect.ValueOf(nv))
		return count, nil
	case ArrayFloat:
		var (
			l      = ByteOrder.Uint16(b[1:])
			v      = reflect.ValueOf(val)
			nv     = make([]float64, l)
			offset = 3
			count  = offset
		)
		v = reflect.Indirect(v)
		for i := 0; i < int(l); i++ {
			n, err := dec.decode(b[offset+i*9:], &nv[i])
			if err != nil {
				return 0, err
			}
			count += n
		}

		v.Set(reflect.ValueOf(nv))
		return count, nil
	case ArrayString:
		var (
			l      = ByteOrder.Uint16(b[1:])
			v      = reflect.ValueOf(val)
			nv     = make([]string, l)
			offset = 3
			count  = offset
		)
		v = reflect.Indirect(v)
		for i := 0; i < int(l); i++ {
			n, err := dec.decode(b[offset:], &nv[i])
			if err != nil {
				return 0, err
			}
			offset += bytes.Index(b[offset:], []byte{0}) + 1
			count += n
		}

		v.Set(reflect.ValueOf(nv))
		return count, nil
	case Map:
		var (
			l      = ByteOrder.Uint32(b[1:])
			v      = reflect.ValueOf(val)
			offset = 5
			count  = offset
		)
		v = reflect.Indirect(v)
		t := v.Type()

		for i := 0; i < int(l); i++ {
			if Type(b[offset]) != MapKey {
				return 0, ErrInvalidMapKey
			}

			offset++
			key := reflect.New(t.Key())
			n, err := dec.decode(b[offset:], key.Interface())
			if err != nil {
				return 0, err
			}

			offset += n
			if Type(b[offset]) != MapValue {
				return 0, ErrInvalidMapKey
			}
			offset++
			val := reflect.New(t.Elem()).Elem()
			n, err = dec.decodeVal(b[offset:], val)
			if err != nil {
				return 0, err
			}
			count += n
			offset += n

			v.SetMapIndex(key.Elem(), val)
		}

		return count, nil
	case Struct:
		var (
			l      = ByteOrder.Uint32(b[1:])
			v      = reflect.ValueOf(val)
			offset = 5
			count  = offset
		)
		v = reflect.Indirect(v)
		t := v.Type()
		for i := 0; i < int(l); i++ {
			if Type(b[offset]) != StructField {
				return 0, ErrInvalidStructField
			}

			offset++
			var key string
			idx := bytes.Index(b[offset:], []byte{0})
			if idx <= 0 {
				return 0, ErrInvalidStructField
			}

			key = string(b[offset : offset+idx])

			offset += idx + 1
			if Type(b[offset]) != StructValue {
				return 0, ErrInvalidStructValue
			}

			ft, ok := t.FieldByName(key)
			if !ok {
				return 0, fmt.Errorf("missing struct field %s", key)
			}
			offset++
			val := reflect.New(ft.Type).Elem()
			n, err := dec.decodeVal(b[offset:], val)
			if err != nil {
				return 0, err
			}
			count += n
			offset += n

			v.FieldByName(key).Set(val)
			// offset++
			// count++
			// v.SetMapIndex(key.Elem(), val)
		}
		return count, nil
	default:
		return 0, errors.New("invalid codec of decode")
	}
}

func (dec *Decoder) DecodeElement(b []byte, val interface{}) (bool, error) {
	if len(b) < 1 {
		return false, ErrBufTooSmall
	}

	switch b[0] {
	case byte(ElementValue):
		if _, err := dec.decode(b[1:], val); err != nil {
			return false, err
		}
		return false, nil
	case byte(ElementRef):
		if _, err := dec.decode(b[1:], val); err != nil {
			return true, err
		}
		return true, nil
	default:
		return false, ErrInvalidElementType
	}
}

func (dec *Decoder) decode(b []byte, val interface{}) (int, error) {
	if len(b) == 0 {
		return 0, ErrBufTooSmall
	}

	v := reflect.ValueOf(val)
	v = reflect.Indirect(v)
	if !v.CanSet() {
		return 0, errors.New("decode to value must can be set")
	}

	return dec.decodeVal(b, v)
}

func (dec *Decoder) decodeVal(b []byte, v reflect.Value) (int, error) {
	typ := Type(b[0])
	switch typ {
	case Uint8:
		if len(b) < 2 {
			return 0, ErrBufTooSmall
		}
		// v.SetUint(uint64(b[1]))
		dec.setVal(v, uint64(b[1]))
		return 2, nil
	case Uint16:
		if len(b) < 3 {
			return 0, ErrBufTooSmall
		}
		dec.setVal(v, uint64(ByteOrder.Uint16(b[1:])))
		// v.SetUint(uint64(ByteOrder.Uint16(b[1:])))
		return 3, nil
	case Uint32:
		if len(b) < 5 {
			return 0, ErrBufTooSmall
		}
		// v.SetUint(uint64(ByteOrder.Uint32(b[1:])))
		dec.setVal(v, uint64(ByteOrder.Uint32(b[1:])))
		return 5, nil
	case Uint64:
		if len(b) < 9 {
			return 0, ErrBufTooSmall
		}
		dec.setVal(v, ByteOrder.Uint64(b[1:]))
		// v.SetUint(ByteOrder.Uint64(b[1:]))
		return 9, nil
	case Int8:
		if len(b) < 2 {
			return 0, ErrBufTooSmall
		}
		v.SetInt(int64(b[1]))
		return 2, nil
	case Int16:
		if len(b) < 3 {
			return 0, ErrBufTooSmall
		}
		v.SetInt(int64(ByteOrder.Uint16(b[1:])))
		return 3, nil
	case Int32:
		if len(b) < 5 {
			return 0, ErrBufTooSmall
		}
		v.SetInt(int64(ByteOrder.Uint32(b[1:])))
		return 5, nil
	case Int64:
		if len(b) < 9 {
			return 0, ErrBufTooSmall
		}
		v.SetInt(int64(ByteOrder.Uint64(b[1:])))
		return 9, nil
	case Int:
		if len(b) < 9 {
			return 0, ErrBufTooSmall
		}
		// v.SetInt(int64(ByteOrder.Uint64(b[1:])))
		dec.setVal(v, int(ByteOrder.Uint64(b[1:])))
		return 9, nil
	case Uint:
		if len(b) < 9 {
			return 0, ErrBufTooSmall
		}
		// v.SetUint(ByteOrder.Uint64(b[1:]))
		dec.setVal(v, ByteOrder.Uint64(b[1:]))
		return 9, nil
	case Float32:
		var (
			buf = bytes.NewBuffer(b[1:])
			val float32
		)
		if len(b) < 5 {
			return 0, ErrBufTooSmall
		}
		if err := binary.Read(buf, ByteOrder, &val); err != nil {
			return 0, err
		}
		v.Set(reflect.ValueOf(val))
		return 5, nil
	case Float64:
		var (
			buf = bytes.NewBuffer(b[1:])
			val float64
		)
		if len(b) < 9 {
			return 0, ErrBufTooSmall
		}
		if err := binary.Read(buf, ByteOrder, &val); err != nil {
			return 0, err
		}
		v.Set(reflect.ValueOf(val))
		return 9, nil
	case Bool:
		var (
			buf = bytes.NewBuffer(b[1:])
			val bool
		)
		if len(b) < 2 {
			return 0, ErrBufTooSmall
		}
		if err := binary.Read(buf, ByteOrder, &val); err != nil {
			return 0, err
		}
		v.Set(reflect.ValueOf(val))
		return 2, nil
	case String:
		p := bytes.IndexByte(b[1:], 0)
		if p < 0 {
			return 0, ErrNonStringTailZero
		}
		// v.SetString(string(b[1 : p+1]))
		v.Set(reflect.ValueOf(string(b[1 : p+1])))
		return p + 2, nil
	case Bytes:
		l := ByteOrder.Uint32(b[1:])
		// v.SetBytes(b[5 : l+5])
		dec.setVal(v, b[5:l+5])
		return int(l) + 5, nil
	case Timestamp:
		var (
			buf = bytes.NewBuffer(b[1:])
			dt  int64
		)
		if err := binary.Read(buf, ByteOrder, &dt); err != nil {
			return 0, err
		}
		t := time.Unix(dt/1000000000, dt%1000000000).UTC()

		v.Set(reflect.ValueOf(t))
		return 9, nil
	case Duration:
		var (
			buf = bytes.NewBuffer(b[1:])
			dt  int64
		)
		if err := binary.Read(buf, ByteOrder, &dt); err != nil {
			return 0, err
		}
		// v.SetInt(dt)
		dec.setVal(v, dt)
		return 9, nil
	default:
		return 0, errors.New("invalid codec of decode")
	}
}

func (dec *Decoder) setVal(v reflect.Value, val interface{}) {
	if v.Kind() == reflect.Interface {
		v.Set(reflect.ValueOf(val))
	} else {
		switch x := val.(type) {
		case uint64:
			v.SetUint(x)
		case int64:
			v.SetInt(x)
		case int:
			v.Set(reflect.ValueOf(x))
		case float64:
			v.SetFloat(x)
		case float32:
			v.SetFloat(float64(x))
		case []byte:
			v.SetBytes(x)
		default:
			panic("invalid type")
		}
	}
}

func Decode(b []byte, val interface{}) (int, error) {
	return decoder.Decode(b, val)
}

func DecodeTo(b []byte) (interface{}, error) {
	var v = new(interface{})
	_, err := Decode(b, v)
	return reflect.ValueOf(v).Elem().Interface(), err
}

func DecodeElement(b []byte, val interface{}) (bool, error) {
	return decoder.DecodeElement(b, val)
}

func DecodeElementTo(b []byte) (interface{}, bool, error) {
	var v = new(interface{})
	ref, err := decoder.DecodeElement(b, v)
	return *v, ref, err
}
