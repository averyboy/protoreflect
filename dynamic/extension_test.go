package dynamic

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/averyboy/protoreflect/desc"

	"github.com/averyboy/protoreflect/internal/testprotos"
	"github.com/averyboy/protoreflect/internal/testutil"
)

func TestSetExtension(t *testing.T) {
	extd, err := desc.LoadFieldDescriptorForExtension(testprotos.E_TestMessage_NestedMessage_AnotherNestedMessage_Flags)
	testutil.Ok(t, err)

	// with dynamic message
	dm := NewMessage(extd.GetOwner())
	err = SetExtension(dm, extd, []bool{true, false, true})
	testutil.Ok(t, err)
	testutil.Eq(t, []bool{true, false, true}, dm.GetField(extd))

	// with non-dynamic message
	var msg testprotos.AnotherTestMessage
	err = SetExtension(&msg, extd, []bool{true, false, true})
	testutil.Ok(t, err)
	val, err := proto.GetExtension(&msg, testprotos.E_TestMessage_NestedMessage_AnotherNestedMessage_Flags)
	testutil.Ok(t, err)
	testutil.Eq(t, []bool{true, false, true}, val)

	// fails with wrong value type
	err = SetExtension(&msg, extd, "foo")
	testutil.Require(t, err != nil)

	// fails if you use wrong type of message
	var msg2 testprotos.TestMessage
	err = SetExtension(&msg2, extd, []bool{true, false, true})
	testutil.Require(t, err != nil)

}
