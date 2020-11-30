// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.7.1
// source: currency.proto

package gen

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Currency struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title     string               `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Value     float64              `protobuf:"fixed64,2,opt,name=value,proto3" json:"value,omitempty"`
	UpdatedAt *timestamp.Timestamp `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	Base      string               `protobuf:"bytes,4,opt,name=base,proto3" json:"base,omitempty"`
}

func (x *Currency) Reset() {
	*x = Currency{}
	if protoimpl.UnsafeEnabled {
		mi := &file_currency_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Currency) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Currency) ProtoMessage() {}

func (x *Currency) ProtoReflect() protoreflect.Message {
	mi := &file_currency_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Currency.ProtoReflect.Descriptor instead.
func (*Currency) Descriptor() ([]byte, []int) {
	return file_currency_proto_rawDescGZIP(), []int{0}
}

func (x *Currency) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Currency) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *Currency) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Currency) GetBase() string {
	if x != nil {
		return x.Base
	}
	return ""
}

type Currencies struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rates []*Currency `protobuf:"bytes,1,rep,name=rates,proto3" json:"rates,omitempty"`
}

func (x *Currencies) Reset() {
	*x = Currencies{}
	if protoimpl.UnsafeEnabled {
		mi := &file_currency_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Currencies) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Currencies) ProtoMessage() {}

func (x *Currencies) ProtoReflect() protoreflect.Message {
	mi := &file_currency_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Currencies.ProtoReflect.Descriptor instead.
func (*Currencies) Descriptor() ([]byte, []int) {
	return file_currency_proto_rawDescGZIP(), []int{1}
}

func (x *Currencies) GetRates() []*Currency {
	if x != nil {
		return x.Rates
	}
	return nil
}

type CurrencyTitle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
}

func (x *CurrencyTitle) Reset() {
	*x = CurrencyTitle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_currency_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CurrencyTitle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CurrencyTitle) ProtoMessage() {}

func (x *CurrencyTitle) ProtoReflect() protoreflect.Message {
	mi := &file_currency_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CurrencyTitle.ProtoReflect.Descriptor instead.
func (*CurrencyTitle) Descriptor() ([]byte, []int) {
	return file_currency_proto_rawDescGZIP(), []int{2}
}

func (x *CurrencyTitle) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

type InitialDayCurrency struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title string  `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Value float64 `protobuf:"fixed64,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *InitialDayCurrency) Reset() {
	*x = InitialDayCurrency{}
	if protoimpl.UnsafeEnabled {
		mi := &file_currency_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitialDayCurrency) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitialDayCurrency) ProtoMessage() {}

func (x *InitialDayCurrency) ProtoReflect() protoreflect.Message {
	mi := &file_currency_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitialDayCurrency.ProtoReflect.Descriptor instead.
func (*InitialDayCurrency) Descriptor() ([]byte, []int) {
	return file_currency_proto_rawDescGZIP(), []int{3}
}

func (x *InitialDayCurrency) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *InitialDayCurrency) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type InitialDayCurrencies struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Currencies []*InitialDayCurrency `protobuf:"bytes,1,rep,name=currencies,proto3" json:"currencies,omitempty"`
}

func (x *InitialDayCurrencies) Reset() {
	*x = InitialDayCurrencies{}
	if protoimpl.UnsafeEnabled {
		mi := &file_currency_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitialDayCurrencies) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitialDayCurrencies) ProtoMessage() {}

func (x *InitialDayCurrencies) ProtoReflect() protoreflect.Message {
	mi := &file_currency_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitialDayCurrencies.ProtoReflect.Descriptor instead.
func (*InitialDayCurrencies) Descriptor() ([]byte, []int) {
	return file_currency_proto_rawDescGZIP(), []int{4}
}

func (x *InitialDayCurrencies) GetCurrencies() []*InitialDayCurrency {
	if x != nil {
		return x.Currencies
	}
	return nil
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_currency_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_currency_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_currency_proto_rawDescGZIP(), []int{5}
}

var File_currency_proto protoreflect.FileDescriptor

var file_currency_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x85, 0x01, 0x0a, 0x08,
	0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x62, 0x61, 0x73, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62,
	0x61, 0x73, 0x65, 0x22, 0x36, 0x0a, 0x0a, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x69, 0x65,
	0x73, 0x12, 0x28, 0x0a, 0x05, 0x72, 0x61, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x2e, 0x43, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x63, 0x79, 0x52, 0x05, 0x72, 0x61, 0x74, 0x65, 0x73, 0x22, 0x25, 0x0a, 0x0d, 0x43,
	0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x22, 0x40, 0x0a, 0x12, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x44, 0x61, 0x79,
	0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x22, 0x54, 0x0a, 0x14, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x44,
	0x61, 0x79, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x12, 0x3c, 0x0a, 0x0a,
	0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x2e, 0x49, 0x6e, 0x69, 0x74,
	0x69, 0x61, 0x6c, 0x44, 0x61, 0x79, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x52, 0x0a,
	0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x32, 0x8c, 0x02, 0x0a, 0x0f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x33, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x52, 0x61,
	0x74, 0x65, 0x73, 0x12, 0x0f, 0x2e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x14, 0x2e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x2e,
	0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x22, 0x00, 0x12, 0x3a, 0x0a, 0x07,
	0x47, 0x65, 0x74, 0x52, 0x61, 0x74, 0x65, 0x12, 0x17, 0x2e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x63, 0x79, 0x2e, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x54, 0x69, 0x74, 0x6c, 0x65,
	0x1a, 0x14, 0x2e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x2e, 0x43, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x22, 0x00, 0x12, 0x3c, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x4c,
	0x61, 0x73, 0x74, 0x52, 0x61, 0x74, 0x65, 0x12, 0x17, 0x2e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x63, 0x79, 0x2e, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x54, 0x69, 0x74, 0x6c, 0x65,
	0x1a, 0x12, 0x2e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x2e, 0x43, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x63, 0x79, 0x22, 0x00, 0x12, 0x4a, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x69,
	0x74, 0x69, 0x61, 0x6c, 0x44, 0x61, 0x79, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12,
	0x0f, 0x2e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x1e, 0x2e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x2e, 0x49, 0x6e, 0x69, 0x74,
	0x69, 0x61, 0x6c, 0x44, 0x61, 0x79, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73,
	0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x67, 0x65, 0x6e, 0x3b, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_currency_proto_rawDescOnce sync.Once
	file_currency_proto_rawDescData = file_currency_proto_rawDesc
)

func file_currency_proto_rawDescGZIP() []byte {
	file_currency_proto_rawDescOnce.Do(func() {
		file_currency_proto_rawDescData = protoimpl.X.CompressGZIP(file_currency_proto_rawDescData)
	})
	return file_currency_proto_rawDescData
}

var file_currency_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_currency_proto_goTypes = []interface{}{
	(*Currency)(nil),             // 0: currency.Currency
	(*Currencies)(nil),           // 1: currency.Currencies
	(*CurrencyTitle)(nil),        // 2: currency.CurrencyTitle
	(*InitialDayCurrency)(nil),   // 3: currency.InitialDayCurrency
	(*InitialDayCurrencies)(nil), // 4: currency.InitialDayCurrencies
	(*Empty)(nil),                // 5: currency.Empty
	(*timestamp.Timestamp)(nil),  // 6: google.protobuf.Timestamp
}
var file_currency_proto_depIdxs = []int32{
	6, // 0: currency.Currency.updated_at:type_name -> google.protobuf.Timestamp
	0, // 1: currency.Currencies.rates:type_name -> currency.Currency
	3, // 2: currency.InitialDayCurrencies.currencies:type_name -> currency.InitialDayCurrency
	5, // 3: currency.currencyService.GetRates:input_type -> currency.Empty
	2, // 4: currency.currencyService.GetRate:input_type -> currency.CurrencyTitle
	2, // 5: currency.currencyService.GetLastRate:input_type -> currency.CurrencyTitle
	5, // 6: currency.currencyService.GetInitialDayCurrency:input_type -> currency.Empty
	1, // 7: currency.currencyService.GetRates:output_type -> currency.Currencies
	1, // 8: currency.currencyService.GetRate:output_type -> currency.Currencies
	0, // 9: currency.currencyService.GetLastRate:output_type -> currency.Currency
	4, // 10: currency.currencyService.GetInitialDayCurrency:output_type -> currency.InitialDayCurrencies
	7, // [7:11] is the sub-list for method output_type
	3, // [3:7] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_currency_proto_init() }
func file_currency_proto_init() {
	if File_currency_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_currency_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Currency); i {
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
		file_currency_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Currencies); i {
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
		file_currency_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CurrencyTitle); i {
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
		file_currency_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitialDayCurrency); i {
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
		file_currency_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitialDayCurrencies); i {
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
		file_currency_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
			RawDescriptor: file_currency_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_currency_proto_goTypes,
		DependencyIndexes: file_currency_proto_depIdxs,
		MessageInfos:      file_currency_proto_msgTypes,
	}.Build()
	File_currency_proto = out.File
	file_currency_proto_rawDesc = nil
	file_currency_proto_goTypes = nil
	file_currency_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CurrencyServiceClient is the client API for CurrencyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CurrencyServiceClient interface {
	GetRates(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Currencies, error)
	GetRate(ctx context.Context, in *CurrencyTitle, opts ...grpc.CallOption) (*Currencies, error)
	GetLastRate(ctx context.Context, in *CurrencyTitle, opts ...grpc.CallOption) (*Currency, error)
	GetInitialDayCurrency(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*InitialDayCurrencies, error)
}

type currencyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCurrencyServiceClient(cc grpc.ClientConnInterface) CurrencyServiceClient {
	return &currencyServiceClient{cc}
}

func (c *currencyServiceClient) GetRates(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Currencies, error) {
	out := new(Currencies)
	err := c.cc.Invoke(ctx, "/currency.currencyService/GetRates", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *currencyServiceClient) GetRate(ctx context.Context, in *CurrencyTitle, opts ...grpc.CallOption) (*Currencies, error) {
	out := new(Currencies)
	err := c.cc.Invoke(ctx, "/currency.currencyService/GetRate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *currencyServiceClient) GetLastRate(ctx context.Context, in *CurrencyTitle, opts ...grpc.CallOption) (*Currency, error) {
	out := new(Currency)
	err := c.cc.Invoke(ctx, "/currency.currencyService/GetLastRate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *currencyServiceClient) GetInitialDayCurrency(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*InitialDayCurrencies, error) {
	out := new(InitialDayCurrencies)
	err := c.cc.Invoke(ctx, "/currency.currencyService/GetInitialDayCurrency", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CurrencyServiceServer is the server API for CurrencyService service.
type CurrencyServiceServer interface {
	GetRates(context.Context, *Empty) (*Currencies, error)
	GetRate(context.Context, *CurrencyTitle) (*Currencies, error)
	GetLastRate(context.Context, *CurrencyTitle) (*Currency, error)
	GetInitialDayCurrency(context.Context, *Empty) (*InitialDayCurrencies, error)
}

// UnimplementedCurrencyServiceServer can be embedded to have forward compatible implementations.
type UnimplementedCurrencyServiceServer struct {
}

func (*UnimplementedCurrencyServiceServer) GetRates(context.Context, *Empty) (*Currencies, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRates not implemented")
}
func (*UnimplementedCurrencyServiceServer) GetRate(context.Context, *CurrencyTitle) (*Currencies, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRate not implemented")
}
func (*UnimplementedCurrencyServiceServer) GetLastRate(context.Context, *CurrencyTitle) (*Currency, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLastRate not implemented")
}
func (*UnimplementedCurrencyServiceServer) GetInitialDayCurrency(context.Context, *Empty) (*InitialDayCurrencies, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInitialDayCurrency not implemented")
}

func RegisterCurrencyServiceServer(s *grpc.Server, srv CurrencyServiceServer) {
	s.RegisterService(&_CurrencyService_serviceDesc, srv)
}

func _CurrencyService_GetRates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CurrencyServiceServer).GetRates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/currency.currencyService/GetRates",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CurrencyServiceServer).GetRates(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CurrencyService_GetRate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CurrencyTitle)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CurrencyServiceServer).GetRate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/currency.currencyService/GetRate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CurrencyServiceServer).GetRate(ctx, req.(*CurrencyTitle))
	}
	return interceptor(ctx, in, info, handler)
}

func _CurrencyService_GetLastRate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CurrencyTitle)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CurrencyServiceServer).GetLastRate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/currency.currencyService/GetLastRate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CurrencyServiceServer).GetLastRate(ctx, req.(*CurrencyTitle))
	}
	return interceptor(ctx, in, info, handler)
}

func _CurrencyService_GetInitialDayCurrency_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CurrencyServiceServer).GetInitialDayCurrency(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/currency.currencyService/GetInitialDayCurrency",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CurrencyServiceServer).GetInitialDayCurrency(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _CurrencyService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "currency.currencyService",
	HandlerType: (*CurrencyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRates",
			Handler:    _CurrencyService_GetRates_Handler,
		},
		{
			MethodName: "GetRate",
			Handler:    _CurrencyService_GetRate_Handler,
		},
		{
			MethodName: "GetLastRate",
			Handler:    _CurrencyService_GetLastRate_Handler,
		},
		{
			MethodName: "GetInitialDayCurrency",
			Handler:    _CurrencyService_GetInitialDayCurrency_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "currency.proto",
}
