package webrtc_dev


const  kRtpVersion uint8 = 2
const kMinRtpPacketSize  = 12
const kMinRtcpPacketSize  = 4

// 72--76 reserved for RTCP conflict avoidance [RFC1889]
func isRtcpPacketByPayLoad(payloadType uint8) bool {
	if payloadType <= 76 && payloadType >= 72 {
		return true
	} else {
		return false
	}
}

func isVersionOk(packet []uint8) bool {
	return packet[0] >> 6 == kRtpVersion
}

func IsRtcpPacket(packet []uint8) bool {
	return len(packet) >= kMinRtcpPacketSize && isVersionOk(packet) && isRtcpPacketByPayLoad(packet[1] & 0x7f)
}