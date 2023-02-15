package webrtc_dev

import (
	"encoding/binary"
)


type PacketType uint8

// RTCP packet types registered with IANA. See: https://www.iana.org/assignments/rtp-parameters/rtp-parameters.xhtml#rtp-parameters-4
const (
	TypeSenderReport              PacketType = 200 // RFC 3550, 6.4.1
	TypeReceiverReport            PacketType = 201 // RFC 3550, 6.4.2
	TypeSourceDescription         PacketType = 202 // RFC 3550, 6.5
	TypeGoodbye                   PacketType = 203 // RFC 3550, 6.6
	TypeApplicationDefined        PacketType = 204 // RFC 3550, 6.7 (unimplemented)
	TypeTransportSpecificFeedback PacketType = 205 // RFC 4585, 6051
	TypePayloadSpecificFeedback   PacketType = 206 // RFC 4585, 6.3
	TypeExtendedReport            PacketType = 207 // RFC 3611
)

/*
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|V=2|P|    RC   |       PT      |             length            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                            Payload                            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 */

type RtcpHeadCommon struct {
	Version			uint8
	Padding			bool
	CountOrFormat		uint8
	PayloadType			PacketType
	PayloadSize			uint16
}

const (
	headLength = 4
	rtpVersion = 2
	maxCount = 31
	versionShift = 6
	paddingShift = 5
)

func (h *RtcpHeadCommon) CheckRtcpHeadPacket() bool {
	if h.CountOrFormat > maxCount {
		return false
	}
	return true
}

func (h *RtcpHeadCommon) Marshal() ([]byte, error) {
	if !h.CheckRtcpHeadPacket() {
		return nil, errInvalidHeader
	}
	headPacket := make([]byte, headLength)
	headPacket[0] = rtpVersion << versionShift
	if h.Padding {
		headPacket[0] |= 1 << paddingShift
	}
	headPacket[0] |= h.CountOrFormat  
	headPacket[1] = uint8(h.PayloadType)
	binary.BigEndian.PutUint16(headPacket[2:], h.PayloadSize)
	return headPacket, nil
}


func (h *RtcpHeadCommon) Unmarshal(data []byte) error {
	if len(data) < headLength {
		return errPacketTooShort
	}
	version := data[0] >>versionShift
	if version != rtpVersion {
		return errBadVersion
	}
	h.Padding = (data[0] >> paddingShift & 0x1f) > 0
	h.PayloadType = PacketType(data[1])
	h.PayloadSize = binary.BigEndian.Uint16(data[2:])
	return nil
}