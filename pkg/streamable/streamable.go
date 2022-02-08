package streamable

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"

	"github.com/cmmarslender/go-chia-lib/pkg/util"
)

const (
	// Name of the struct tag used to identify the streamable properties
	tagName = "streamable"

	// ProtocolVersion Current supported Protocol Version
	// Not all of this is supported, but this was the current version at the time
	// This library was started
	ProtocolVersion string = "0.0.33"

	// Bytes that indicate bool yes or no when serialized
	boolFalse uint8 = 0
	boolTrue  uint8 = 1
)

// Unmarshal unmarshals a streamable type based on struct tags
// Struct order is extremely important in this decoding. Ensure the order/types are identical
// on both sides of the stream
// Ugly, but.. it works? So we can make it pretty later...
func Unmarshal(bytes []byte, v interface{}) error {
	tv := reflect.ValueOf(v)
	if tv.Kind() != reflect.Ptr || tv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	// Gets rid of the pointer
	tv = reflect.Indirect(tv)

	// Get the actual type
	t := tv.Type()

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("streamable can't unmarshal into non-struct type")
	}

	bytes, err := unmarshalStruct(bytes, t, tv)
	return err
}

func unmarshalStruct(bytes []byte, t reflect.Type, tv reflect.Value) ([]byte, error) {
	var err error

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i) // Struct Field
		fieldValue := tv.Field(i) // Field Value
		fieldType := fieldValue.Type()

		bytes, err = unmarshalField(bytes, fieldType, fieldValue, structField)
		if err != nil {
			return bytes, err
		}
	}

	return bytes, nil
}

func unmarshalSlice(bytes []byte, t reflect.Type, v reflect.Value) ([]byte, error) {
	var err error
	var newVal []byte

	// Slice/List is 4 byte prefix (number of items) and then serialization of each item
	// Get 4 byte length prefix
	var length []byte
	length, bytes, err = util.ShiftNBytes(4, bytes)
	numItems := binary.BigEndian.Uint32(length)

	sliceKind := t.Elem().Kind()
	switch sliceKind {
	case reflect.Uint8: // same as byte
		// In this case, numItems == numBytes, because its a uint8
		newVal, bytes, err = util.ShiftNBytes(uint(numItems), bytes)
		if err != nil {
			return bytes, err
		}
		if !v.CanSet() {
			return bytes, fmt.Errorf("field %s is not settable", v.String())
		}

		sliceReflect := reflect.MakeSlice(v.Type(), 0, 0)
		for _, newValBytes := range newVal {
			sliceReflect = reflect.Append(sliceReflect, reflect.ValueOf(newValBytes))
		}
		v.Set(sliceReflect)
	case reflect.Struct:
		// Recursion, I guess
		sliceReflect := reflect.MakeSlice(v.Type(), 0, 0)
		for j := uint32(0); j < numItems; j++ {
			newValue := reflect.Indirect(reflect.New(v.Type().Elem()))
			bytes, err = unmarshalStruct(bytes, t.Elem(), newValue)
			sliceReflect = reflect.Append(sliceReflect, newValue)
		}
		v.Set(sliceReflect)
	default:
		return bytes, fmt.Errorf("encountered type inside slice that is not implemented")
	}

	return bytes, nil
}

func unmarshalField(bytes []byte, fieldType reflect.Type, fieldValue reflect.Value, structField reflect.StructField) ([]byte, error) {
	var tag string
	var tagPresent bool
	if tag, tagPresent = structField.Tag.Lookup(tagName); !tagPresent {
		// Continuing because the tag isn't present
		return bytes, nil
	}

	var err error
	var newVal []byte

	// If optional, should be one byte bool that indicates if its present or not
	// This is the hackiest of hacky ways to check if this is ACTUALLY optional
	// @TODO one day need to actually parse these options out properly
	if strings.Contains(tag, "optional") {
		if fieldValue.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("optional fields must be pointer types")
		}

		// Its optional, check if we have actual data
		var presentFlag []byte
		presentFlag, bytes, err = util.ShiftNBytes(1, bytes)
		if presentFlag[0] == boolFalse {
			// Not present in the data, continue
			return bytes, nil
		}
	}

	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()

		// Need to init the field to something non-nil before using it
		fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
		fieldValue = fieldValue.Elem()
	}

	switch kind := fieldType.Kind(); kind {
	case reflect.Uint8:
		newVal, bytes, err = util.ShiftNBytes(1, bytes)
		if err != nil {
			return bytes, err
		}
		if !fieldValue.CanSet() {
			return bytes, fmt.Errorf("field %s is not settable", fieldValue.String())
		}
		fieldValue.SetUint(uint64(util.BytesToUint8(newVal)))
	case reflect.Uint16:
		newVal, bytes, err = util.ShiftNBytes(2, bytes)
		if err != nil {
			return bytes, err
		}
		if !fieldValue.CanSet() {
			return bytes, fmt.Errorf("field %s is not settable", fieldValue.String())
		}
		newInt := util.BytesToUint16(newVal)
		fieldValue.SetUint(uint64(newInt))
	case reflect.Slice:
		bytes, err = unmarshalSlice(bytes, fieldType, fieldValue)
		if err != nil {
			return bytes, err
		}
	case reflect.String:
		// 4 byte size prefix, then []byte which can be converted to utf-8 string
		// Get 4 byte length prefix
		var length []byte
		length, bytes, err = util.ShiftNBytes(4, bytes)
		numBytes := binary.BigEndian.Uint32(length)

		var strBytes []byte
		strBytes, bytes, err = util.ShiftNBytes(uint(numBytes), bytes)
		fieldValue.SetString(string(strBytes))
	default:
		return bytes, fmt.Errorf("unimplemented type %s", fieldValue.Kind())
	}

	return bytes, nil
}

// Marshal marshals the item into the streamable byte format
func Marshal(v interface{}) ([]byte, error) {
	// Doesn't matter if a pointer or not for marshalling, so
	// we just call this and let it deal with ptr or not ptr
	tv := reflect.Indirect(reflect.ValueOf(v))

	// Get the actual type
	t := tv.Type()

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("streamable can't marshal a non-struct type")
	}

	// This will become the final encoded data
	var finalBytes []byte

	// Iterate over all available fields in the type and encode to bytes
	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i) // Field Type
		field := tv.Field(i) // Field Value

		var tag string
		var tagPresent bool
		if tag, tagPresent = tField.Tag.Lookup(tagName); !tagPresent {
			// Continuing because the tag isn't present
			continue
		}

		// If optional, the type MUST be a pointer type
		// nil pointer will be assumed to be not present, and we'll insert 0x00 and move on
		// Anything other than nil pointer we'll insert 0x01 and encode the value
		// This is the hackiest of hacky ways to check if this is ACTUALLY optional
		// @TODO one day need to actually parse these options out properly
		if strings.Contains(tag, "optional") {
			if field.Kind() != reflect.Ptr {
				return nil, fmt.Errorf("optional fields must be pointer types")
			}

			// Its optional, check if we have actual data
			if field.IsNil() {
				// Field is nil, insert `false` byte and continue
				finalBytes = append(finalBytes, boolFalse)
				continue
			}

			finalBytes = append(finalBytes, boolTrue)
		}

		// If field is still a pointer, get rid of that now that we're past the optional checking
		field = reflect.Indirect(field)

		switch field.Kind() {
		case reflect.Ptr:
			panic("How did we get a pointer after calling Indirect?")
		case reflect.Uint8:
			newInt := uint8(field.Uint())
			finalBytes = append(finalBytes, newInt)
		case reflect.Uint16:
			newInt := uint16(field.Uint())
			finalBytes = append(finalBytes, util.Uint16ToBytes(newInt)...)
		case reflect.Slice:
			// Slice/List is 4 byte prefix (number of items) and then serialization of each item
			// Get 4 byte length prefix
			numItems := uint32(field.Len())
			finalBytes = append(finalBytes, util.Uint32ToBytes(numItems)...)

			sliceKind := field.Type().Elem().Kind()
			switch sliceKind {
			case reflect.Uint8: // same as byte
				// This is the easy case - already a slice of bytes
				finalBytes = append(finalBytes, field.Bytes()...)
			case reflect.Struct:
				// @TODO again this is super hacky and manual and only works for Capability right now
				for j := 0; j < field.Len(); j++ {
					currentStruct := field.Index(j)
					capabilityType := uint16(currentStruct.Field(0).Uint())
					valueStrBytes := []byte(currentStruct.Field(1).String())

					finalBytes = append(finalBytes, util.Uint16ToBytes(capabilityType)...)

					// String stuff...
					numBytes := uint32(len(valueStrBytes))
					finalBytes = append(finalBytes, util.Uint32ToBytes(numBytes)...)

					finalBytes = append(finalBytes, valueStrBytes...)
				}
			}
		case reflect.String:
			// Strings get converted to []byte with a 4 byte size prefix
			strBytes := []byte(field.String())
			numBytes := uint32(len(strBytes))
			finalBytes = append(finalBytes, util.Uint32ToBytes(numBytes)...)

			finalBytes = append(finalBytes, strBytes...)
		default:
			return nil, fmt.Errorf("unimplemented type %s", field.Kind())
		}
	}

	return finalBytes, nil
}
