package local

import "github.com/ryt-io/ryt-v2/message"

type stubOutboundMessage struct {
	bypassThrottling      bool
	op                    message.Op
	bytes                 []byte
	bytesSavedCompression int
}

func (s stubOutboundMessage) BypassThrottling() bool      { return s.bypassThrottling }
func (s stubOutboundMessage) Op() message.Op              { return s.op }
func (s stubOutboundMessage) Bytes() []byte               { return s.bytes }
func (s stubOutboundMessage) BytesSavedCompression() int { return s.bytesSavedCompression }

func newOutboundMessage(op message.Op, payload []byte, bypassThrottling bool) message.OutboundMessage {
	return stubOutboundMessage{
		bypassThrottling:      bypassThrottling,
		op:                    op,
		bytes:                 payload,
		bytesSavedCompression: 0,
	}
}
