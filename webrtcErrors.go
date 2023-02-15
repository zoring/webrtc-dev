package webrtc_dev

import "errors"

var (
	errInvalidHeader            = errors.New("rtcp: invalid header")
	errPacketTooShort           = errors.New("rtcp: packet too short")
	errBadVersion               = errors.New("rtcp: invalid packet version")
)
