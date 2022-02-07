package streamable

import (
	"encoding/binary"
	"fmt"
	"log"
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

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i) // Field Type
		field := tv.Field(i) // Field Value

		var tag string
		var tagPresent bool
		if tag, tagPresent = tField.Tag.Lookup(tagName); !tagPresent {
			// Continuing because the tag isn't present
			continue
		}

		var err error
		var newVal []byte

		// If optional, should be one byte bool that indicates if its present or not
		// This is the hackiest of hacky ways to check if this is ACTUALLY optional
		// @TODO one day need to actually parse these options out properly
		if strings.Contains(tag, "optional") {
			// Its optional, check if we have actual data
			var presentFlag []byte
			presentFlag, bytes, err = util.ShiftNBytes(1, bytes)
			if presentFlag[0] == boolFalse {
				// Not present in the data, continue
				log.Println("This field was omitted. Skipping...")
				continue
			}
		}

		switch field.Kind() {
		case reflect.Ptr:
			// @TODO
			ptrKind := field.Type().Elem().Kind()
			log.Println("Its a ptr, and its slice type is ", ptrKind)
		case reflect.Uint8:
			// 1 byte
			newVal, bytes, err = util.ShiftNBytes(1, bytes)
			if err != nil {
				return err
			}
			if !field.CanSet() {
				return fmt.Errorf("field %s is not settable", field.String())
			}
			field.SetUint(uint64(newVal[0]))
		case reflect.Uint16:
			// @TODO
			// 2 bytes

		case reflect.Slice:
			// Slice/List is 4 byte prefix (number of items) and then serialization of each item
			// Get 4 byte length prefix
			var length []byte
			length, bytes, err = util.ShiftNBytes(4, bytes)
			numItems := binary.BigEndian.Uint32(length)

			sliceKind := field.Type().Elem().Kind()
			log.Println("Its a slice, and its slice type is ", sliceKind)
			switch sliceKind {
			case reflect.Uint8: // same as byte
				// In this case, numItems == numBytes, because its a uint8
				newVal, bytes, err = util.ShiftNBytes(uint(numItems), bytes)
				if err != nil {
					return err
				}
				if !field.CanSet() {
					return fmt.Errorf("field %s is not settable", field.String())
				}

				sliceReflect := reflect.MakeSlice(field.Type(), 0, 0)
				for _, newValBytes := range newVal {
					rv := reflect.ValueOf(newValBytes)
					sliceReflect = reflect.Append(sliceReflect, rv)
				}
				field.Set(sliceReflect)
			}
		default:
			return fmt.Errorf("unimplemented type")
		}
	}

	return nil
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
			// Its optional, check if we have actual data
			if field.IsNil() {
				// Field is nil, insert `false` byte and continue
				finalBytes = append(finalBytes, boolFalse)
				continue
			}

			finalBytes = append(finalBytes, boolTrue)
		}

		field = reflect.Indirect(field)

		switch field.Kind() {
		case reflect.Ptr:
			panic("How did we get a pointer after calling Indirect?")
		case reflect.Uint8:
			newInt := uint8(field.Uint())
			// Only putting the first one, because it's a uint8
			finalBytes = append(finalBytes, newInt)
		case reflect.Uint16:
			//finalBytes = append(finalBytes, field.Bytes()...)
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
			}
		default:
			return nil, fmt.Errorf("unimplemented type")
		}
	}

	return finalBytes, nil
}
