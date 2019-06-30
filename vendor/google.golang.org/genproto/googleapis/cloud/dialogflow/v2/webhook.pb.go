// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/dialogflow/v2/webhook.proto

package dialogflow

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	_struct "github.com/golang/protobuf/ptypes/struct"
	_ "google.golang.org/genproto/googleapis/api/annotations"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// The request message for a webhook call.
type WebhookRequest struct {
	// The unique identifier of detectIntent request session.
	// Can be used to identify end-user inside webhook implementation.
	// Format: `projects/<Project ID>/agent/sessions/<Session ID>`, or
	// `projects/<Project ID>/agent/environments/<Environment ID>/users/<User
	// ID>/sessions/<Session ID>`.
	Session string `protobuf:"bytes,4,opt,name=session,proto3" json:"session,omitempty"`
	// The unique identifier of the response. Contains the same value as
	// `[Streaming]DetectIntentResponse.response_id`.
	ResponseId string `protobuf:"bytes,1,opt,name=response_id,json=responseId,proto3" json:"response_id,omitempty"`
	// The result of the conversational query or event processing. Contains the
	// same value as `[Streaming]DetectIntentResponse.query_result`.
	QueryResult *QueryResult `protobuf:"bytes,2,opt,name=query_result,json=queryResult,proto3" json:"query_result,omitempty"`
	// Optional. The contents of the original request that was passed to
	// `[Streaming]DetectIntent` call.
	OriginalDetectIntentRequest *OriginalDetectIntentRequest `protobuf:"bytes,3,opt,name=original_detect_intent_request,json=originalDetectIntentRequest,proto3" json:"original_detect_intent_request,omitempty"`
	XXX_NoUnkeyedLiteral        struct{}                     `json:"-"`
	XXX_unrecognized            []byte                       `json:"-"`
	XXX_sizecache               int32                        `json:"-"`
}

func (m *WebhookRequest) Reset()         { *m = WebhookRequest{} }
func (m *WebhookRequest) String() string { return proto.CompactTextString(m) }
func (*WebhookRequest) ProtoMessage()    {}
func (*WebhookRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2ba880626f278d96, []int{0}
}

func (m *WebhookRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WebhookRequest.Unmarshal(m, b)
}
func (m *WebhookRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WebhookRequest.Marshal(b, m, deterministic)
}
func (m *WebhookRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WebhookRequest.Merge(m, src)
}
func (m *WebhookRequest) XXX_Size() int {
	return xxx_messageInfo_WebhookRequest.Size(m)
}
func (m *WebhookRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WebhookRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WebhookRequest proto.InternalMessageInfo

func (m *WebhookRequest) GetSession() string {
	if m != nil {
		return m.Session
	}
	return ""
}

func (m *WebhookRequest) GetResponseId() string {
	if m != nil {
		return m.ResponseId
	}
	return ""
}

func (m *WebhookRequest) GetQueryResult() *QueryResult {
	if m != nil {
		return m.QueryResult
	}
	return nil
}

func (m *WebhookRequest) GetOriginalDetectIntentRequest() *OriginalDetectIntentRequest {
	if m != nil {
		return m.OriginalDetectIntentRequest
	}
	return nil
}

// The response message for a webhook call.
type WebhookResponse struct {
	// Optional. The text to be shown on the screen. This value is passed directly
	// to `QueryResult.fulfillment_text`.
	FulfillmentText string `protobuf:"bytes,1,opt,name=fulfillment_text,json=fulfillmentText,proto3" json:"fulfillment_text,omitempty"`
	// Optional. The collection of rich messages to present to the user. This
	// value is passed directly to `QueryResult.fulfillment_messages`.
	FulfillmentMessages []*Intent_Message `protobuf:"bytes,2,rep,name=fulfillment_messages,json=fulfillmentMessages,proto3" json:"fulfillment_messages,omitempty"`
	// Optional. This value is passed directly to `QueryResult.webhook_source`.
	Source string `protobuf:"bytes,3,opt,name=source,proto3" json:"source,omitempty"`
	// Optional. This value is passed directly to `QueryResult.webhook_payload`.
	// See the related `fulfillment_messages[i].payload field`, which may be used
	// as an alternative to this field.
	//
	// This field can be used for Actions on Google responses.
	// It should have a structure similar to the JSON message shown here. For more
	// information, see
	// [Actions on Google Webhook
	// Format](https://developers.google.com/actions/dialogflow/webhook)
	// <pre>{
	//   "google": {
	//     "expectUserResponse": true,
	//     "richResponse": {
	//       "items": [
	//         {
	//           "simpleResponse": {
	//             "textToSpeech": "this is a simple response"
	//           }
	//         }
	//       ]
	//     }
	//   }
	// }</pre>
	Payload *_struct.Struct `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`
	// Optional. The collection of output contexts. This value is passed directly
	// to `QueryResult.output_contexts`.
	OutputContexts []*Context `protobuf:"bytes,5,rep,name=output_contexts,json=outputContexts,proto3" json:"output_contexts,omitempty"`
	// Optional. Makes the platform immediately invoke another `DetectIntent` call
	// internally with the specified event as input.
	// When this field is set, Dialogflow ignores the `fulfillment_text`,
	// `fulfillment_messages`, and `payload` fields.
	FollowupEventInput   *EventInput `protobuf:"bytes,6,opt,name=followup_event_input,json=followupEventInput,proto3" json:"followup_event_input,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *WebhookResponse) Reset()         { *m = WebhookResponse{} }
func (m *WebhookResponse) String() string { return proto.CompactTextString(m) }
func (*WebhookResponse) ProtoMessage()    {}
func (*WebhookResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2ba880626f278d96, []int{1}
}

func (m *WebhookResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WebhookResponse.Unmarshal(m, b)
}
func (m *WebhookResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WebhookResponse.Marshal(b, m, deterministic)
}
func (m *WebhookResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WebhookResponse.Merge(m, src)
}
func (m *WebhookResponse) XXX_Size() int {
	return xxx_messageInfo_WebhookResponse.Size(m)
}
func (m *WebhookResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_WebhookResponse.DiscardUnknown(m)
}

var xxx_messageInfo_WebhookResponse proto.InternalMessageInfo

func (m *WebhookResponse) GetFulfillmentText() string {
	if m != nil {
		return m.FulfillmentText
	}
	return ""
}

func (m *WebhookResponse) GetFulfillmentMessages() []*Intent_Message {
	if m != nil {
		return m.FulfillmentMessages
	}
	return nil
}

func (m *WebhookResponse) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *WebhookResponse) GetPayload() *_struct.Struct {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *WebhookResponse) GetOutputContexts() []*Context {
	if m != nil {
		return m.OutputContexts
	}
	return nil
}

func (m *WebhookResponse) GetFollowupEventInput() *EventInput {
	if m != nil {
		return m.FollowupEventInput
	}
	return nil
}

// Represents the contents of the original request that was passed to
// the `[Streaming]DetectIntent` call.
type OriginalDetectIntentRequest struct {
	// The source of this request, e.g., `google`, `facebook`, `slack`. It is set
	// by Dialogflow-owned servers.
	Source string `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	// Optional. The version of the protocol used for this request.
	// This field is AoG-specific.
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	// Optional. This field is set to the value of the `QueryParameters.payload`
	// field passed in the request. Some integrations that query a Dialogflow
	// agent may provide additional information in the payload.
	//
	// In particular for the Telephony Gateway this field has the form:
	// <pre>{
	//  "telephony": {
	//    "caller_id": "+18558363987"
	//  }
	// }</pre>
	// Note: The caller ID field (`caller_id`) will be redacted for Standard
	// Edition agents and populated with the caller ID in [E.164
	// format](https://en.wikipedia.org/wiki/E.164) for Enterprise Edition agents.
	Payload              *_struct.Struct `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *OriginalDetectIntentRequest) Reset()         { *m = OriginalDetectIntentRequest{} }
func (m *OriginalDetectIntentRequest) String() string { return proto.CompactTextString(m) }
func (*OriginalDetectIntentRequest) ProtoMessage()    {}
func (*OriginalDetectIntentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2ba880626f278d96, []int{2}
}

func (m *OriginalDetectIntentRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OriginalDetectIntentRequest.Unmarshal(m, b)
}
func (m *OriginalDetectIntentRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OriginalDetectIntentRequest.Marshal(b, m, deterministic)
}
func (m *OriginalDetectIntentRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OriginalDetectIntentRequest.Merge(m, src)
}
func (m *OriginalDetectIntentRequest) XXX_Size() int {
	return xxx_messageInfo_OriginalDetectIntentRequest.Size(m)
}
func (m *OriginalDetectIntentRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OriginalDetectIntentRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OriginalDetectIntentRequest proto.InternalMessageInfo

func (m *OriginalDetectIntentRequest) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *OriginalDetectIntentRequest) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *OriginalDetectIntentRequest) GetPayload() *_struct.Struct {
	if m != nil {
		return m.Payload
	}
	return nil
}

func init() {
	proto.RegisterType((*WebhookRequest)(nil), "google.cloud.dialogflow.v2.WebhookRequest")
	proto.RegisterType((*WebhookResponse)(nil), "google.cloud.dialogflow.v2.WebhookResponse")
	proto.RegisterType((*OriginalDetectIntentRequest)(nil), "google.cloud.dialogflow.v2.OriginalDetectIntentRequest")
}

func init() {
	proto.RegisterFile("google/cloud/dialogflow/v2/webhook.proto", fileDescriptor_2ba880626f278d96)
}

var fileDescriptor_2ba880626f278d96 = []byte{
	// 567 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0x41, 0x8f, 0xd3, 0x3c,
	0x10, 0x86, 0x95, 0xf4, 0xfb, 0x5a, 0xad, 0xbb, 0xda, 0x22, 0xb3, 0x82, 0xa8, 0x8b, 0x96, 0xaa,
	0x48, 0x6c, 0xe1, 0x90, 0x88, 0x80, 0xc4, 0x81, 0xdb, 0x6e, 0x01, 0x15, 0x81, 0x58, 0x02, 0x02,
	0x84, 0x84, 0xa2, 0x34, 0x71, 0x83, 0x85, 0xeb, 0x49, 0x63, 0xbb, 0xdd, 0x4a, 0x9c, 0xf8, 0x01,
	0xdc, 0x38, 0x71, 0xe3, 0xc8, 0x2f, 0xe4, 0x88, 0x62, 0x3b, 0xb4, 0x20, 0x36, 0x70, 0x1c, 0xcf,
	0x33, 0xef, 0xcc, 0xbc, 0x99, 0xa0, 0x51, 0x0e, 0x90, 0x33, 0x12, 0xa4, 0x0c, 0x54, 0x16, 0x64,
	0x34, 0x61, 0x90, 0xcf, 0x18, 0xac, 0x82, 0x65, 0x18, 0xac, 0xc8, 0xf4, 0x1d, 0xc0, 0x7b, 0xbf,
	0x28, 0x41, 0x02, 0xee, 0x1b, 0xd2, 0xd7, 0xa4, 0xbf, 0x21, 0xfd, 0x65, 0xd8, 0x6f, 0x52, 0x49,
	0x81, 0x4b, 0x72, 0x26, 0x8d, 0x4a, 0xff, 0xa8, 0x81, 0xa4, 0x5c, 0x12, 0x5e, 0x83, 0x4d, 0x92,
	0x82, 0x08, 0x41, 0x81, 0x5b, 0xf2, 0xce, 0xdf, 0xc9, 0x98, 0x70, 0x49, 0xe5, 0x3a, 0x96, 0xeb,
	0x82, 0xd8, 0xaa, 0x2b, 0xb6, 0x4a, 0x47, 0x53, 0x35, 0x0b, 0x84, 0x2c, 0x55, 0x2a, 0x7f, 0xcb,
	0x26, 0x05, 0x0d, 0x12, 0xce, 0x41, 0x26, 0x92, 0x02, 0x17, 0x26, 0x3b, 0xfc, 0xec, 0xa2, 0xbd,
	0x57, 0xc6, 0x9c, 0x88, 0x2c, 0x14, 0x11, 0x12, 0x7b, 0xa8, 0x63, 0x7b, 0x79, 0xff, 0x0d, 0x9c,
	0xd1, 0x4e, 0x54, 0x87, 0xf8, 0x2a, 0xea, 0x96, 0x44, 0x14, 0xc0, 0x05, 0x89, 0x69, 0xe6, 0x39,
	0x3a, 0x8b, 0xea, 0xa7, 0x49, 0x86, 0x1f, 0xa1, 0xdd, 0x85, 0x22, 0xe5, 0x3a, 0x2e, 0x89, 0x50,
	0x4c, 0x7a, 0xee, 0xc0, 0x19, 0x75, 0xc3, 0x23, 0xff, 0x7c, 0xbf, 0xfd, 0x67, 0x15, 0x1f, 0x69,
	0x3c, 0xea, 0x2e, 0x36, 0x01, 0xfe, 0x80, 0x0e, 0xa1, 0xa4, 0x39, 0xe5, 0x09, 0x8b, 0x33, 0x22,
	0x49, 0x2a, 0x63, 0xe3, 0x6a, 0x5c, 0x9a, 0x41, 0xbd, 0x96, 0x56, 0xbf, 0xdb, 0xa4, 0xfe, 0xd4,
	0x2a, 0x8c, 0xb5, 0xc0, 0x44, 0xd7, 0xdb, 0x3d, 0xa3, 0x03, 0x38, 0x3f, 0x39, 0xfc, 0xd4, 0x42,
	0xbd, 0x9f, 0xbe, 0x98, 0xfd, 0xf0, 0x0d, 0x74, 0x61, 0xa6, 0xd8, 0x8c, 0x32, 0x36, 0xaf, 0xc6,
	0xa8, 0x4e, 0xc1, 0x7a, 0xd0, 0xdb, 0x7a, 0x7f, 0x41, 0xce, 0x24, 0x7e, 0x8b, 0xf6, 0xb7, 0xd1,
	0x39, 0x11, 0x22, 0xc9, 0x89, 0xf0, 0xdc, 0x41, 0x6b, 0xd4, 0x0d, 0x6f, 0x36, 0x8d, 0x6c, 0xe6,
	0xf0, 0x9f, 0x98, 0x92, 0xe8, 0xe2, 0x96, 0x8e, 0x7d, 0x13, 0xf8, 0x12, 0x6a, 0x0b, 0x50, 0x65,
	0x4a, 0xb4, 0x07, 0x3b, 0x91, 0x8d, 0xf0, 0x2d, 0xd4, 0x29, 0x92, 0x35, 0x83, 0x24, 0xd3, 0x9f,
	0xae, 0x1b, 0x5e, 0xae, 0x3b, 0xd5, 0xb7, 0xe1, 0x3f, 0xd7, 0xb7, 0x11, 0xd5, 0x1c, 0x7e, 0x8c,
	0x7a, 0xa0, 0x64, 0xa1, 0x64, 0x6c, 0xaf, 0x5b, 0x78, 0xff, 0xeb, 0x21, 0xaf, 0x35, 0x0d, 0x79,
	0x62, 0xd8, 0x68, 0xcf, 0xd4, 0xda, 0x50, 0xe0, 0xd7, 0x68, 0x7f, 0x06, 0x8c, 0xc1, 0x4a, 0x15,
	0x31, 0x59, 0x56, 0xab, 0x53, 0x5e, 0x28, 0xe9, 0xb5, 0xf5, 0x34, 0xd7, 0x9b, 0x24, 0xef, 0x57,
	0xf8, 0xa4, 0xa2, 0x23, 0x5c, 0x6b, 0x6c, 0xde, 0x86, 0x1f, 0x1d, 0x74, 0xd0, 0xf0, 0x35, 0xb7,
	0x2c, 0x71, 0x7e, 0xb1, 0xc4, 0x43, 0x9d, 0x25, 0x29, 0xf5, 0x35, 0xbb, 0xe6, 0x9a, 0x6d, 0xb8,
	0x6d, 0x56, 0xeb, 0xdf, 0xcc, 0x3a, 0xfe, 0xe2, 0xa0, 0xc3, 0x14, 0xe6, 0x0d, 0x6b, 0x1c, 0xef,
	0xda, 0xab, 0x39, 0xad, 0x34, 0x4e, 0x9d, 0x37, 0x63, 0xcb, 0xe6, 0xc0, 0x12, 0x9e, 0xfb, 0x50,
	0xe6, 0x41, 0x4e, 0xb8, 0xee, 0x10, 0x98, 0x54, 0x52, 0x50, 0xf1, 0xa7, 0x3f, 0xfe, 0xde, 0x26,
	0xfa, 0xee, 0x38, 0x5f, 0x5d, 0x77, 0xfc, 0xe0, 0x9b, 0xdb, 0x7f, 0x68, 0xe4, 0x4e, 0x74, 0xeb,
	0xf1, 0xa6, 0xf5, 0xcb, 0x70, 0xda, 0xd6, 0xaa, 0xb7, 0x7f, 0x04, 0x00, 0x00, 0xff, 0xff, 0x78,
	0xbf, 0xb5, 0xab, 0x09, 0x05, 0x00, 0x00,
}
