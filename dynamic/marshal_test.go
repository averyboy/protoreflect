package dynamic

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/golang/protobuf/proto"

	"github.com/averyboy/protoreflect/desc"
	"github.com/averyboy/protoreflect/internal/testprotos"
	"github.com/averyboy/protoreflect/internal/testutil"
)

// Shared stuff for marshalling and unmarshalling tests. This is used for the binary format, the text
// format, and the JSON format.

var unaryFieldsPosMsg = &testprotos.UnaryFields{
	I: proto.Int32(1),
	J: proto.Int64(2),
	K: proto.Int32(3),
	L: proto.Int64(4),
	M: proto.Uint32(5),
	N: proto.Uint64(6),
	O: proto.Uint32(7),
	P: proto.Uint64(8),
	Q: proto.Int32(9),
	R: proto.Int64(10),
	S: proto.Float32(11),
	T: proto.Float64(12),
	U: []byte{0, 1, 2, 3, 4, 5, 6, 7},
	V: proto.String("foobar"),
	W: proto.Bool(true),
	X: &testprotos.RepeatedFields{
		I: []int32{3},
		V: []string{"baz"},
	},
	Groupy: &testprotos.UnaryFields_GroupY{
		Ya: proto.String("bedazzle"),
		Yb: proto.Int32(42),
	},
	Z: testprotos.TestEnum_SECOND.Enum(),
}

var unaryFieldsNegMsg = &testprotos.UnaryFields{
	I: proto.Int32(-1),
	J: proto.Int64(-2),
	K: proto.Int32(-3),
	L: proto.Int64(-4),
	M: proto.Uint32(5),
	N: proto.Uint64(6),
	O: proto.Uint32(7),
	P: proto.Uint64(8),
	Q: proto.Int32(-9),
	R: proto.Int64(-10),
	S: proto.Float32(-11),
	T: proto.Float64(-12),
	U: []byte{0, 1, 2, 3, 4, 5, 6, 7},
	V: proto.String("foobar"),
	W: proto.Bool(true),
	X: &testprotos.RepeatedFields{
		I: []int32{-3},
		V: []string{"baz"},
	},
	Groupy: &testprotos.UnaryFields_GroupY{
		Ya: proto.String("bedazzle"),
		Yb: proto.Int32(-42),
	},
	Z: testprotos.TestEnum_SECOND.Enum(),
}

var unaryFieldsPosInfMsg = &testprotos.UnaryFields{
	S: proto.Float32(float32(math.Inf(1))),
	T: proto.Float64(math.Inf(1)),
}

var unaryFieldsNegInfMsg = &testprotos.UnaryFields{
	S: proto.Float32(float32(math.Inf(-1))),
	T: proto.Float64(math.Inf(-1)),
}

var unaryFieldsNanMsg = &testprotos.UnaryFields{
	S: proto.Float32(float32(math.NaN())),
	T: proto.Float64(math.NaN()),
}

var repeatedFieldsMsg = &testprotos.RepeatedFields{
	I: []int32{1, -2, 3},
	J: []int64{-4, 5, -6},
	K: []int32{7, -8, 9},
	L: []int64{-10, 11, -12},
	M: []uint32{13, 14, 15},
	N: []uint64{16, 17, 18},
	O: []uint32{19, 20, 21},
	P: []uint64{22, 23, 24},
	Q: []int32{25, 26, 27},
	R: []int64{28, 29, 30},
	S: []float32{31, 32, 33},
	T: []float64{34, 35, 36},
	U: [][]byte{{0, 1, 2, 3}, {4, 5, 6, 7}, {8, 9, 10, 11}},
	V: []string{"foo", "bar", "baz"},
	W: []bool{true, false, true},
	X: []*testprotos.UnaryFields{
		{I: proto.Int32(-32), V: proto.String("baz")},
		{I: proto.Int32(-64), V: proto.String("bozo")},
	},
	Groupy: []*testprotos.RepeatedFields_GroupY{
		{Ya: proto.String("bedazzle"), Yb: proto.Int32(42)},
		{Ya: proto.String("buzzard"), Yb: proto.Int32(-421)},
	},
	Z: []testprotos.TestEnum{testprotos.TestEnum_SECOND, testprotos.TestEnum_THIRD, testprotos.TestEnum_FIRST},
}

var repeatedFieldsInfNanMsg = &testprotos.RepeatedFields{
	S: []float32{float32(math.Inf(1)), float32(math.Inf(-1)), float32(math.NaN())},
	T: []float64{math.Inf(1), math.Inf(-1), math.NaN()},
}

var repeatedPackedFieldsMsg = &testprotos.RepeatedPackedFields{
	I: []int32{1, -2, 3},
	J: []int64{-4, 5, -6},
	K: []int32{7, -8, 9},
	L: []int64{-10, 11, -12},
	M: []uint32{13, 14, 15},
	N: []uint64{16, 17, 18},
	O: []uint32{19, 20, 21},
	P: []uint64{22, 23, 24},
	Q: []int32{25, 26, 27},
	R: []int64{28, 29, 30},
	S: []float32{31, 32, 33},
	T: []float64{34, 35, 36},
	U: []bool{true, false, true},
	Groupy: []*testprotos.RepeatedPackedFields_GroupY{
		{Yb: []int32{42, 84, 126, 168, 210}},
		{Yb: []int32{-210, -168, -126, -84, -42}},
	},
	V: []testprotos.TestEnum{testprotos.TestEnum_SECOND, testprotos.TestEnum_THIRD, testprotos.TestEnum_FIRST},
}

var repeatedPackedFieldsInfNanMsg = &testprotos.RepeatedPackedFields{
	S: []float32{float32(math.Inf(1)), float32(math.Inf(-1)), float32(math.NaN())},
	T: []float64{math.Inf(1), math.Inf(-1), math.NaN()},
}

var mapKeyFieldsMsg = &testprotos.MapKeyFields{
	I: map[int32]string{1: "foo", -2: "bar", 3: "baz"},
	J: map[int64]string{-4: "foo", 5: "bar", -6: "baz"},
	K: map[int32]string{7: "foo", -8: "bar", 9: "baz"},
	L: map[int64]string{-10: "foo", 11: "bar", -12: "baz"},
	M: map[uint32]string{13: "foo", 14: "bar", 15: "baz"},
	N: map[uint64]string{16: "foo", 17: "bar", 18: "baz"},
	O: map[uint32]string{19: "foo", 20: "bar", 21: "baz"},
	P: map[uint64]string{22: "foo", 23: "bar", 24: "baz"},
	Q: map[int32]string{25: "foo", 26: "bar", 27: "baz"},
	R: map[int64]string{28: "foo", 29: "bar", 30: "baz"},
	S: map[string]string{"a": "foo", "b": "bar", "❤": "baz"},
	T: map[bool]string{true: "foo", false: "bar"},
}

var mapValueFieldsMsg = &testprotos.MapValFields{
	I: map[string]int32{"a": 1, "b": -2, "c": 3},
	J: map[string]int64{"a": -4, "b": 5, "c": -6},
	K: map[string]int32{"a": 7, "b": -8, "c": 9},
	L: map[string]int64{"a": -10, "b": 11, "c": -12},
	M: map[string]uint32{"a": 13, "b": 14, "c": 15},
	N: map[string]uint64{"a": 16, "b": 17, "c": 18},
	O: map[string]uint32{"a": 19, "b": 20, "c": 21},
	P: map[string]uint64{"a": 22, "b": 23, "c": 24},
	Q: map[string]int32{"a": 25, "b": 26, "c": 27},
	R: map[string]int64{"a": 28, "b": 29, "c": 30},
	S: map[string]float32{"a": 31, "b": 32, "c": 33},
	T: map[string]float64{"a": 34, "b": 35, "c": 36},
	U: map[string][]byte{"a": {0, 1, 2, 3}, "b": {4, 5, 6, 7}, "c": {8, 9, 10, 11}},
	V: map[string]string{"a": "foo", "b": "bar", "c": "baz"},
	W: map[string]bool{"a": true, "b": false, "c": true},
	X: map[string]*testprotos.UnaryFields{
		"a": {I: proto.Int32(-32), V: proto.String("baz")},
		"b": {I: proto.Int32(-64), V: proto.String("bozo")},
	},
	Y: map[string]testprotos.TestEnum{"a": testprotos.TestEnum_SECOND, "b": testprotos.TestEnum_THIRD, "c": testprotos.TestEnum_FIRST},
}

var mapValueFieldsInfNanMsg = &testprotos.MapValFields{
	S: map[string]float32{"a": float32(math.Inf(1)), "b": float32(math.Inf(-1)), "c": float32(math.NaN())},
	T: map[string]float64{"a": math.Inf(1), "b": math.Inf(-1), "c": math.NaN()},
}

var mapValueFieldsNilMsg = &testprotos.TestRequest{
	Others: map[string]*testprotos.TestMessage{"a": nil, "b": nil},
}

var mapValueFieldsNilUnknownMsg proto.Message
var mdForUnknownMsg *desc.MessageDescriptor

func init() {
	// NB: can't use desc/builder package because that would cause dependency cycle :(
	fdp := &descriptor.FileDescriptorProto{
		Name:    proto.String("foo.proto"),
		Syntax:  proto.String("proto3"),
		Package: proto.String("example"),
		MessageType: []*descriptor.DescriptorProto{
			{
				Name: proto.String("Message"),
				Field: []*descriptor.FieldDescriptorProto{
					{
						Name:     proto.String("vals"),
						Type:     descriptor.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
						Number:   proto.Int32(1),
						Label:    descriptor.FieldDescriptorProto_LABEL_REPEATED.Enum(),
						TypeName: proto.String(".example.Message.ValsEntry"),
					},
				},
				NestedType: []*descriptor.DescriptorProto{
					{
						Name: proto.String("ValsEntry"),
						Options: &descriptor.MessageOptions{
							MapEntry: proto.Bool(true),
						},
						Field: []*descriptor.FieldDescriptorProto{
							{
								Name:   proto.String("key"),
								Number: proto.Int32(1),
								Type:   descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
								Label:  descriptor.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
							},
							{
								Name:     proto.String("value"),
								Number:   proto.Int32(2),
								Type:     descriptor.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
								Label:    descriptor.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
								TypeName: proto.String(".example.Message"),
							},
						},
					},
				},
			},
		},
	}
	fd, err := desc.CreateFileDescriptor(fdp)
	if err != nil {
		panic(err)
	}
	mdForUnknownMsg = fd.GetMessageTypes()[0]
	mapValueFieldsNilUnknownMsg = &unknownMsg{
		Vals: map[string]*unknownMsg{"a": nil, "b": nil},
	}
}

// This message looks and acts like a proto but is NOT in the registry, so proto.MessageType
// returns nil, which forces us to fallback to a nil *dynamic.Message when representing a
// nil map value in a dynamic message, allowing tests to check that strange edge case.
type unknownMsg struct {
	Vals map[string]*unknownMsg `protobuf:"bytes,1,rep,name=vals,proto3" json:"vals,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *unknownMsg) XXX_MessageName() string {
	return "example.Message"
}
func (m *unknownMsg) Reset() {
	m.Vals = nil
}
func (m *unknownMsg) String() string {
	return fmt.Sprintf("%#v", m.Vals)
}
func (m *unknownMsg) ProtoMessage() {
}
func (m *unknownMsg) GetMessageDescriptor() *desc.MessageDescriptor {
	return mdForUnknownMsg
}

func doTranslationParty(t *testing.T, msg proto.Message,
	marshalPm func(proto.Message) ([]byte, error), unmarshalPm func([]byte, proto.Message) error,
	marshalDm func(*Message) ([]byte, error), unmarshalDm func(*Message, []byte) error,
	includesNaN, compareBytes, outputIsString bool) {

	md, err := desc.LoadMessageDescriptorForMessage(msg)
	testutil.Ok(t, err)
	dm := NewMessage(md)

	b, err := marshalPm(msg)
	testutil.Ok(t, err)
	err = unmarshalDm(dm, b)
	testutil.Ok(t, err, "failed to unmarshal from: %s", b)

	// both techniques to marshal do the same thing
	b2a, err := marshalPm(dm)
	testutil.Ok(t, err)
	b2b, err := marshalDm(dm)
	testutil.Ok(t, err)
	testutil.Eq(t, b2a, b2b)

	// round trip back to proto.Message
	msg2 := reflect.New(reflect.TypeOf(msg).Elem()).Interface().(proto.Message)
	err = unmarshalPm(b2a, msg2)
	testutil.Ok(t, err, "failed to unmarshal from: %s", b2a)

	if !includesNaN {
		// NaN fields are never equal so this would always be false
		testutil.Ceq(t, msg, msg2, eqpm)
	}
	if compareBytes {
		if outputIsString {
			testutil.Eq(t, string(b), string(b2a))
		} else {
			testutil.Eq(t, b, b2a)
		}
	}

	// and back again
	b3, err := marshalPm(msg2)
	testutil.Ok(t, err)
	dm2 := NewMessage(md)
	err = unmarshalDm(dm2, b3)
	testutil.Ok(t, err, "failed to unmarshal from: %s", b3)

	if !includesNaN {
		testutil.Ceq(t, dm, dm2, eqdm)
	}

	// dynamic message -> (bytes) -> dynamic message
	// both techniques to unmarshal are equivalent
	dm3 := NewMessage(md)
	err = unmarshalPm(b2a, dm3)
	testutil.Ok(t, err, "failed to unmarshal from: %s", b2a)
	dm4 := NewMessage(md)
	err = unmarshalDm(dm4, b2a)
	testutil.Ok(t, err, "failed to unmarshal from: %s", b2a)

	if !includesNaN {
		testutil.Ceq(t, dm, dm3, eqdm)
		testutil.Ceq(t, dm, dm4, eqdm)
	}
}
