// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: event.proto

package event

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EventType int32

const (
	EventType_UNKNOWN EventType = 0
	EventType_CREATED EventType = 1
	EventType_REMOVED EventType = 2
	EventType_UPDATED EventType = 3
)

// Enum value maps for EventType.
var (
	EventType_name = map[int32]string{
		0: "UNKNOWN",
		1: "CREATED",
		2: "REMOVED",
		3: "UPDATED",
	}
	EventType_value = map[string]int32{
		"UNKNOWN": 0,
		"CREATED": 1,
		"REMOVED": 2,
		"UPDATED": 3,
	}
)

func (x EventType) Enum() *EventType {
	p := new(EventType)
	*p = x
	return p
}

func (x EventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EventType) Descriptor() protoreflect.EnumDescriptor {
	return file_event_proto_enumTypes[0].Descriptor()
}

func (EventType) Type() protoreflect.EnumType {
	return &file_event_proto_enumTypes[0]
}

func (x EventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EventType.Descriptor instead.
func (EventType) EnumDescriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{0}
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId      uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ChecklistId string `protobuf:"bytes,2,opt,name=checklist_id,json=checklistId,proto3" json:"checklist_id,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Event) GetChecklistId() string {
	if x != nil {
		return x.ChecklistId
	}
	return ""
}

var File_event_proto protoreflect.FileDescriptor

var file_event_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x6f,
	0x7a, 0x6f, 0x6e, 0x76, 0x61, 0x2e, 0x6f, 0x76, 0x61, 0x2e, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x6c,
	0x69, 0x73, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x22, 0x43, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x65,
	0x63, 0x6b, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x6c, 0x69, 0x73, 0x74, 0x49, 0x64, 0x2a, 0x3f, 0x0a, 0x09,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b,
	0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45,
	0x44, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x45, 0x4d, 0x4f, 0x56, 0x45, 0x44, 0x10, 0x02,
	0x12, 0x0b, 0x0a, 0x07, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x44, 0x10, 0x03, 0x42, 0x2f, 0x5a,
	0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x7a, 0x6f, 0x6e,
	0x76, 0x61, 0x2f, 0x6f, 0x76, 0x61, 0x2d, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x6c, 0x69, 0x73, 0x74,
	0x2d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_event_proto_rawDescOnce sync.Once
	file_event_proto_rawDescData = file_event_proto_rawDesc
)

func file_event_proto_rawDescGZIP() []byte {
	file_event_proto_rawDescOnce.Do(func() {
		file_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_event_proto_rawDescData)
	})
	return file_event_proto_rawDescData
}

var file_event_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_event_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_event_proto_goTypes = []interface{}{
	(EventType)(0), // 0: ozonva.ova.checklist.api.EventType
	(*Event)(nil),  // 1: ozonva.ova.checklist.api.Event
}
var file_event_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_event_proto_init() }
func file_event_proto_init() {
	if File_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_event_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_event_proto_goTypes,
		DependencyIndexes: file_event_proto_depIdxs,
		EnumInfos:         file_event_proto_enumTypes,
		MessageInfos:      file_event_proto_msgTypes,
	}.Build()
	File_event_proto = out.File
	file_event_proto_rawDesc = nil
	file_event_proto_goTypes = nil
	file_event_proto_depIdxs = nil
}
