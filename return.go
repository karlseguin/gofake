package gofake

import (
	"errors"
)

type Return struct {
	Values []interface{}
}

func (r *Return) String(index int, d string) string {
	return r.Type(index, d).(string)
}

func (r *Return) Bool(index int, d bool) bool {
	return r.Type(index, d).(bool)
}

func (r *Return) Byte(index int, d byte) byte {
	return r.Type(index, d).(byte)
}

func (r *Return) Rune(index int, d rune) rune {
	return r.Type(index, d).(rune)
}

func (r *Return) Error(index int, d error) error {
	value := r.At(index, d)
	if value == nil {
		return nil
	}
	switch typed := value.(type) {
		case string:
			return errors.New(typed)
		default:
			return value.(error)
	}
}

func (r *Return) Int(index int, d int) int {
	return r.Type(index, d).(int)
}

func (r *Return) Int8(index int, d int8) int8 {
	return r.Type(index, d).(int8)
}

func (r *Return) Int16(index int, d int16) int16 {
	return r.Type(index, d).(int16)
}

func (r *Return) Int32(index int, d int32) int32 {
	return r.Type(index, d).(int32)
}

func (r *Return) Int64(index int, d int64) int64 {
	return r.Type(index, d).(int64)
}

func (r *Return) UInt(index int, d uint) uint {
	return r.Type(index, d).(uint)
}

func (r *Return) UInt8(index int, d uint8) uint8 {
	return r.Type(index, d).(uint8)
}

func (r *Return) UInt16(index int, d uint16) uint16 {
	return r.Type(index, d).(uint16)
}

func (r *Return) UInt32(index int, d uint32) uint32 {
	return r.Type(index, d).(uint32)
}

func (r *Return) UInt64(index int, d uint64) uint64 {
	return r.Type(index, d).(uint64)
}

func (r *Return) Float64(index int, d float64) float64 {
	return r.Type(index, d).(float64)
}

func (r *Return) Float32(index int, d float32) float32 {
	return r.Type(index, d).(float32)
}

func (r *Return) Type(index int, d interface{}) interface{} {
	if len(r.Values) > index  && r.Values[index] != nil {
		return r.Values[index]
	}
	return d
}

func (r *Return) At(index int, d interface{}) interface{} {
	if len(r.Values) > index  {
		return r.Values[index]
	}
	return d
}
