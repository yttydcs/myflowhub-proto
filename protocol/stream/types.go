package stream

// 本文件定义 Proto 中 `stream` 子协议的稳定类型与载荷结构。

import "encoding/json"

const SubProtoStream uint8 = 8

const (
	KindCtrl byte = 0x01
	KindData byte = 0x02
	KindAck  byte = 0x03
)

const (
	ActionAnnounce             = "announce"
	ActionAnnounceResp         = "announce_resp"
	ActionWithdraw             = "withdraw"
	ActionWithdrawResp         = "withdraw_resp"
	ActionListSources          = "list_sources"
	ActionListSourcesResp      = "list_sources_resp"
	ActionGetSource            = "get_source"
	ActionGetSourceResp        = "get_source_resp"
	ActionAnnounceConsumer     = "announce_consumer"
	ActionAnnounceConsumerResp = "announce_consumer_resp"
	ActionWithdrawConsumer     = "withdraw_consumer"
	ActionWithdrawConsumerResp = "withdraw_consumer_resp"
	ActionListConsumers        = "list_consumers"
	ActionListConsumersResp    = "list_consumers_resp"
	ActionGetConsumer          = "get_consumer"
	ActionGetConsumerResp      = "get_consumer_resp"
	ActionSubscribe            = "subscribe"
	ActionSubscribeResp        = "subscribe_resp"
	ActionUnsubscribe          = "unsubscribe"
	ActionUnsubscribeResp      = "unsubscribe_resp"
	ActionConnect              = "connect"
	ActionConnectResp          = "connect_resp"
	ActionDisconnect           = "disconnect"
	ActionDisconnectResp       = "disconnect_resp"
	ActionSignal               = "signal"
	ActionSignalResp           = "signal_resp"
)

const (
	PermStreamPublish   = "stream.publish"
	PermStreamConsume   = "stream.consume"
	PermStreamSubscribe = "stream.subscribe"
	PermStreamConnect   = "stream.connect"
)

const (
	StreamKindMusic  = "music"
	StreamKindVideo  = "video"
	StreamKindText   = "text"
	StreamKindCustom = "custom"
)

const (
	ModeLive    = "live"
	ModeBounded = "bounded"
)

const (
	UnitModeFrame = "frame"
	UnitModeChunk = "chunk"
)

const (
	SignalOpPause           = "pause"
	SignalOpResume          = "resume"
	SignalOpMetadataUpdate  = "metadata_update"
	SignalOpKeyframeRequest = "keyframe_request"
	SignalOpCustom          = "custom"
)

const (
	HeaderVersionV1 uint8 = 1

	DataFlagEOS           uint8 = 1
	DataFlagKeyframe      uint8 = 2
	DataFlagConfig        uint8 = 4
	DataFlagDiscontinuity uint8 = 8
)

type Message struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type SourceDescriptor struct {
	SourceID    string          `json:"source_id"`
	Producer    uint32          `json:"producer"`
	Name        string          `json:"name,omitempty"`
	Kind        string          `json:"kind"`
	ContentType string          `json:"content_type,omitempty"`
	Mode        string          `json:"mode,omitempty"`
	UnitMode    string          `json:"unit_mode,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
}

type ConsumerDescriptor struct {
	ConsumerID  string          `json:"consumer_id"`
	Consumer    uint32          `json:"consumer"`
	Name        string          `json:"name,omitempty"`
	Kind        string          `json:"kind"`
	ContentType string          `json:"content_type,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
}

type AnnounceReq struct {
	ReqID  string           `json:"req_id"`
	Source SourceDescriptor `json:"source"`
}

type AnnounceResp struct {
	ReqID  string            `json:"req_id"`
	Code   int               `json:"code"`
	Msg    string            `json:"msg,omitempty"`
	Source *SourceDescriptor `json:"source,omitempty"`
}

type WithdrawReq struct {
	ReqID    string `json:"req_id"`
	SourceID string `json:"source_id"`
}

type WithdrawResp struct {
	ReqID    string `json:"req_id"`
	Code     int    `json:"code"`
	Msg      string `json:"msg,omitempty"`
	SourceID string `json:"source_id,omitempty"`
}

type ListSourcesReq struct {
	ReqID    string `json:"req_id"`
	Producer uint32 `json:"producer"`
	Kind     string `json:"kind,omitempty"`
	Tag      string `json:"tag,omitempty"`
}

type ListSourcesResp struct {
	ReqID    string             `json:"req_id"`
	Code     int                `json:"code"`
	Msg      string             `json:"msg,omitempty"`
	Producer uint32             `json:"producer,omitempty"`
	Sources  []SourceDescriptor `json:"sources,omitempty"`
}

type GetSourceReq struct {
	ReqID    string `json:"req_id"`
	Producer uint32 `json:"producer"`
	SourceID string `json:"source_id"`
}

type GetSourceResp struct {
	ReqID  string            `json:"req_id"`
	Code   int               `json:"code"`
	Msg    string            `json:"msg,omitempty"`
	Source *SourceDescriptor `json:"source,omitempty"`
}

type AnnounceConsumerReq struct {
	ReqID            string             `json:"req_id"`
	ConsumerEndpoint ConsumerDescriptor `json:"consumer_endpoint"`
}

type AnnounceConsumerResp struct {
	ReqID            string              `json:"req_id"`
	Code             int                 `json:"code"`
	Msg              string              `json:"msg,omitempty"`
	ConsumerEndpoint *ConsumerDescriptor `json:"consumer_endpoint,omitempty"`
}

type WithdrawConsumerReq struct {
	ReqID      string `json:"req_id"`
	ConsumerID string `json:"consumer_id"`
}

type WithdrawConsumerResp struct {
	ReqID      string `json:"req_id"`
	Code       int    `json:"code"`
	Msg        string `json:"msg,omitempty"`
	ConsumerID string `json:"consumer_id,omitempty"`
}

type ListConsumersReq struct {
	ReqID    string `json:"req_id"`
	Consumer uint32 `json:"consumer"`
	Kind     string `json:"kind,omitempty"`
	Tag      string `json:"tag,omitempty"`
}

type ListConsumersResp struct {
	ReqID             string               `json:"req_id"`
	Code              int                  `json:"code"`
	Msg               string               `json:"msg,omitempty"`
	Consumer          uint32               `json:"consumer,omitempty"`
	ConsumerEndpoints []ConsumerDescriptor `json:"consumer_endpoints,omitempty"`
}

type GetConsumerReq struct {
	ReqID      string `json:"req_id"`
	Consumer   uint32 `json:"consumer"`
	ConsumerID string `json:"consumer_id"`
}

type GetConsumerResp struct {
	ReqID            string              `json:"req_id"`
	Code             int                 `json:"code"`
	Msg              string              `json:"msg,omitempty"`
	ConsumerEndpoint *ConsumerDescriptor `json:"consumer_endpoint,omitempty"`
}

type SubscribeReq struct {
	ReqID         string `json:"req_id"`
	Producer      uint32 `json:"producer"`
	SourceID      string `json:"source_id"`
	ConsumerID    string `json:"consumer_id"`
	ResumeFrom    uint64 `json:"resume_from,omitempty"`
	WindowBytes   uint32 `json:"window_bytes,omitempty"`
	AckIntervalMs uint32 `json:"ack_interval_ms,omitempty"`
}

type SubscribeResp struct {
	ReqID            string              `json:"req_id"`
	Code             int                 `json:"code"`
	Msg              string              `json:"msg,omitempty"`
	Accept           bool                `json:"accept,omitempty"`
	Source           *SourceDescriptor   `json:"source,omitempty"`
	ConsumerEndpoint *ConsumerDescriptor `json:"consumer_endpoint,omitempty"`
	DeliveryID       string              `json:"delivery_id,omitempty"`
	Producer         uint32              `json:"producer,omitempty"`
	Consumer         uint32              `json:"consumer,omitempty"`
	ConsumerID       string              `json:"consumer_id,omitempty"`
	StartPosition    uint64              `json:"start_position,omitempty"`
	WindowBytes      uint32              `json:"window_bytes,omitempty"`
	AckIntervalMs    uint32              `json:"ack_interval_ms,omitempty"`
}

type UnsubscribeReq struct {
	ReqID      string `json:"req_id"`
	DeliveryID string `json:"delivery_id"`
	Reason     string `json:"reason,omitempty"`
}

type UnsubscribeResp struct {
	ReqID      string `json:"req_id"`
	Code       int    `json:"code"`
	Msg        string `json:"msg,omitempty"`
	DeliveryID string `json:"delivery_id,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

type ConnectReq struct {
	ReqID         string `json:"req_id"`
	Producer      uint32 `json:"producer"`
	SourceID      string `json:"source_id"`
	Consumer      uint32 `json:"consumer"`
	ConsumerID    string `json:"consumer_id"`
	ResumeFrom    uint64 `json:"resume_from,omitempty"`
	WindowBytes   uint32 `json:"window_bytes,omitempty"`
	AckIntervalMs uint32 `json:"ack_interval_ms,omitempty"`
}

type ConnectResp struct {
	ReqID            string              `json:"req_id"`
	Code             int                 `json:"code"`
	Msg              string              `json:"msg,omitempty"`
	Accept           bool                `json:"accept,omitempty"`
	Source           *SourceDescriptor   `json:"source,omitempty"`
	ConsumerEndpoint *ConsumerDescriptor `json:"consumer_endpoint,omitempty"`
	DeliveryID       string              `json:"delivery_id,omitempty"`
	Producer         uint32              `json:"producer,omitempty"`
	Consumer         uint32              `json:"consumer,omitempty"`
	ConsumerID       string              `json:"consumer_id,omitempty"`
	StartPosition    uint64              `json:"start_position,omitempty"`
	WindowBytes      uint32              `json:"window_bytes,omitempty"`
	AckIntervalMs    uint32              `json:"ack_interval_ms,omitempty"`
}

type DisconnectReq struct {
	ReqID      string `json:"req_id"`
	DeliveryID string `json:"delivery_id"`
	Reason     string `json:"reason,omitempty"`
}

type DisconnectResp struct {
	ReqID      string `json:"req_id"`
	Code       int    `json:"code"`
	Msg        string `json:"msg,omitempty"`
	DeliveryID string `json:"delivery_id,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

type SignalReq struct {
	ReqID      string          `json:"req_id"`
	DeliveryID string          `json:"delivery_id"`
	Op         string          `json:"op"`
	Data       json.RawMessage `json:"data,omitempty"`
}

type SignalResp struct {
	ReqID      string `json:"req_id"`
	Code       int    `json:"code"`
	Msg        string `json:"msg,omitempty"`
	DeliveryID string `json:"delivery_id,omitempty"`
	Op         string `json:"op,omitempty"`
}

type StreamDataHeaderV1 struct {
	Ver        uint8
	Flags      uint8
	DeliveryID [16]byte
	Position   uint64
	PTSMs      uint64
}

type StreamAckHeaderV1 struct {
	Ver         uint8
	Flags       uint8
	DeliveryID  [16]byte
	Position    uint64
	CreditUnits uint32
	Reserved    uint32
}
