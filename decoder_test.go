package amf

import (
	"bytes"
	"testing"
	"reflect"
	"github.com/marcuswu/amf/amf3"
)

func TestReadAMFPacket(t *testing.T) {
	buf := bytes.NewReader([]byte{0x00, 0x03, 0x00, 0x00, 0x00, 0x01, 0x00, 0x12, 0x4C, 0x6F, 0x67, 0x69, 0x6E, 0x2E,
	0x70, 0x72, 0x6F, 0x63, 0x65, 0x73, 0x73, 0x4C, 0x6F, 0x67, 0x69, 0x6E, 0x00, 0x02, 0x2F, 0x30, 0x00, 0x00, 0x00,
	0x01, 0x11, 0x09, 0x03, 0x01, 0x09, 0x0F, 0x01, 0x04, 0x36, 0x04, 0x01, 0x06, 0x33, 0x38, 0x31, 0x62, 0x32, 0x62,
	0x34, 0x64, 0x65, 0x62, 0x65, 0x37, 0x35, 0x61, 0x34, 0x64, 0x2D, 0x35, 0x36, 0x39, 0x33, 0x39, 0x36, 0x34, 0x34,
	0x33, 0x06, 0x01, 0x06, 0x1F, 0x33, 0x35, 0x39, 0x31, 0x32, 0x35, 0x30, 0x35, 0x31, 0x35, 0x36, 0x31, 0x32, 0x37,
	0x34, 0x06, 0x1F, 0x38, 0x31, 0x62, 0x32, 0x62, 0x34, 0x64, 0x65, 0x62, 0x65, 0x37, 0x35, 0x61, 0x34, 0x64, 0x06,
	0x01})
	decoder := NewDecoder(buf)
	got, err := decoder.Decode()
	if err != nil {
		t.Errorf("test for %s error: %s", "foo", err)
	} else {
		if got.Version != 3 {
			t.Errorf("expected Version 3, got %v", got.Version)
		}
		if len(got.Headers) != 0 {
			t.Errorf("expected 0 Headers, got %v", len(got.Headers))
		}

		if len(got.Messages) != 1 {
			t.Errorf("expected 1 message, got %v", len(got.Messages))
		}

		m := got.Messages[0]
		//now check message contents -- should contain an array of length 1 containing an array of length 7
		if _, ok := m.Data.(*amf3.ArrayType); !ok {
			t.Errorf("expected amf3 ArrayType message Data, got %s", reflect.TypeOf(m.Data).Kind())
		}

		var arrayValue *amf3.ArrayType = m.Data.(*amf3.ArrayType)

		if len(arrayValue.Dense) != 1 {
			t.Errorf("expected message contents to be an array of length 1, got length %d", len(arrayValue.Dense))
		}

		var internalArray *amf3.ArrayType
		var ok bool
		internalArray, ok = arrayValue.Dense[0].(*amf3.ArrayType)

		if !ok {
			t.Errorf("expected amf3 ArrayType message Data, got %s", reflect.TypeOf(internalArray).Kind())
		}

		if len(internalArray.Dense) != 7 {
			t.Errorf("expected internal array to be of length 7, got length %d", len(internalArray.Dense))
		}

		//var intVal uint32
		var intVal amf3.IntegerType
		intVal, ok = internalArray.Dense[0].(amf3.IntegerType)

		if !ok || intVal != 54 {
			t.Error("expected first value of internal array to be 54, got ", internalArray.Dense[0])
		}

		intVal, ok = internalArray.Dense[1].(amf3.IntegerType)
		if !ok || intVal != 1 {
			t.Errorf("expected first value of internal array to be 1, got %v", internalArray.Dense[1])
		}

		var stringVal amf3.StringType
		stringVal, ok = internalArray.Dense[2].(amf3.StringType)
		if !ok || stringVal != "81b2b4debe75a4d-569396443" {
			t.Errorf("expected first value of internal array to be 81b2b4debe75a4d-569396443, got %v", internalArray.Dense[2])
		}

		stringVal, ok = internalArray.Dense[3].(amf3.StringType)
		if !ok || stringVal != "" {
			t.Errorf("expected first value of internal array to be an empty string, got %v", internalArray.Dense[3])
		}

		stringVal, ok = internalArray.Dense[4].(amf3.StringType)
		if !ok || stringVal != "359125051561274" {
			t.Errorf("expected first value of internal array to be 359125051561274, got %v", internalArray.Dense[4])
		}

		stringVal, ok = internalArray.Dense[5].(amf3.StringType)
		if !ok || stringVal != "81b2b4debe75a4d" {
			t.Errorf("expected first value of internal array to be 81b2b4debe75a4d, got %v", internalArray.Dense[5])
		}

		stringVal, ok = internalArray.Dense[6].(amf3.StringType)
		if !ok || stringVal != "" {
			t.Errorf("expected first value of internal array to be an empty string, got %v", internalArray.Dense[6])
		}
	}
}
