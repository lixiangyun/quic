package protocol

import "testing"
import "bytes"

type testquicframe struct {
	positiveTest              bool
	leastUnackedDeltaByteSize uint
	data                      []byte
	frame                     QuicFrame
}

var tests_quicframe = []testquicframe{
	// PADDING Frame
	{true, 0,
		[]byte{QUICFRAMETYPE_PADDING},
		QuicFrame{
			frameType:   QUICFRAMETYPE_PADDING,
			frameLength: 0,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_PADDING, 0x00},
		QuicFrame{
			frameType:   QUICFRAMETYPE_PADDING,
			frameLength: 1,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_PADDING, 0x00, 0x00},
		QuicFrame{
			frameType:   QUICFRAMETYPE_PADDING,
			frameLength: 2,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_PADDING, 0x00, 0x00, 0x00},
		QuicFrame{
			frameType:   QUICFRAMETYPE_PADDING,
			frameLength: 3,
		}},
	// PING Frame
	{true, 0, // TEST 4
		[]byte{QUICFRAMETYPE_PING},
		QuicFrame{
			frameType: QUICFRAMETYPE_PING,
		}},
	// BLOCKED Frame
	{false, 0, // not enough data
		[]byte{QUICFRAMETYPE_BLOCKED, 0x12, 0x34, 0x56},
		QuicFrame{
			frameType: QUICFRAMETYPE_BLOCKED,
			streamId:  0x78563412,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_BLOCKED, 0x12, 0x34, 0x56, 0x78},
		QuicFrame{
			frameType: QUICFRAMETYPE_BLOCKED,
			streamId:  0x78563412,
		}},
	// WINDOW_UPDATE Frame
	{false, 0, // not enough data
		[]byte{QUICFRAMETYPE_WINDOW_UPDATE, 0x12, 0x34, 0x56, 0x78, 0x0a, 0x0b, 0x0c, 0x0d, 0xaa, 0xbb, 0xcc},
		QuicFrame{
			frameType:  QUICFRAMETYPE_WINDOW_UPDATE,
			streamId:   0x78563412,
			byteOffset: 0xddccbbaa0d0c0b0a,
		}},
	{true, 0, // TEST 8
		[]byte{QUICFRAMETYPE_WINDOW_UPDATE, 0x12, 0x34, 0x56, 0x78, 0x0a, 0x0b, 0x0c, 0x0d, 0xaa, 0xbb, 0xcc, 0xdd},
		QuicFrame{
			frameType:  QUICFRAMETYPE_WINDOW_UPDATE,
			streamId:   0x78563412,
			byteOffset: 0xddccbbaa0d0c0b0a,
		}},
	// RST_STREAM
	{false, 0, // not enough data
		[]byte{QUICFRAMETYPE_RST_STREAM, 0x12, 0x34, 0x56, 0x78, 0x0a, 0x0b, 0x0c, 0x0d, 0xaa, 0xbb, 0xcc, 0xdd, 0x11, 0x22, 0x33},
		QuicFrame{
			frameType:  QUICFRAMETYPE_RST_STREAM,
			streamId:   0x78563412,
			byteOffset: 0xddccbbaa0d0c0b0a,
			errorCode:  0x44332211,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_RST_STREAM, 0x12, 0x34, 0x56, 0x78, 0x0a, 0x0b, 0x0c, 0x0d, 0xaa, 0xbb, 0xcc, 0xdd, 0x11, 0x22, 0x33, 0x44},
		QuicFrame{
			frameType:  QUICFRAMETYPE_RST_STREAM,
			streamId:   0x78563412,
			byteOffset: 0xddccbbaa0d0c0b0a,
			errorCode:  0x44332211,
		}},
	// CONNECTION_CLOSE Frame
	{false, 0, // not enough data
		[]byte{QUICFRAMETYPE_CONNECTION_CLOSE, 0x11, 0x22, 0x33, 0x44, 0x01, 0x00},
		QuicFrame{
			frameType:   QUICFRAMETYPE_CONNECTION_CLOSE,
			errorCode:   0x44332211,
			frameLength: 1,
			frameData:   []byte{0x1a},
		}},
	{true, 0, // TEST 12
		[]byte{QUICFRAMETYPE_CONNECTION_CLOSE, 0x11, 0x22, 0x33, 0x44, 0x00, 0x00},
		QuicFrame{
			frameType:   QUICFRAMETYPE_CONNECTION_CLOSE,
			errorCode:   0x44332211,
			frameLength: 0,
			frameData:   nil,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_CONNECTION_CLOSE, 0x11, 0x22, 0x33, 0x44, 0x01, 0x00, 0x1a},
		QuicFrame{
			frameType:   QUICFRAMETYPE_CONNECTION_CLOSE,
			errorCode:   0x44332211,
			frameLength: 1,
			frameData:   []byte{0x1a},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_CONNECTION_CLOSE, 0x11, 0x22, 0x33, 0x44, 0x02, 0x00, 0x1a, 0x2b},
		QuicFrame{
			frameType:   QUICFRAMETYPE_CONNECTION_CLOSE,
			errorCode:   0x44332211,
			frameLength: 2,
			frameData:   []byte{0x1a, 0x2b},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_CONNECTION_CLOSE, 0x11, 0x22, 0x33, 0x44, 0x03, 0x00, 0x1a, 0x2b, 0x3c},
		QuicFrame{
			frameType:   QUICFRAMETYPE_CONNECTION_CLOSE,
			errorCode:   0x44332211,
			frameLength: 3,
			frameData:   []byte{0x1a, 0x2b, 0x3c},
		}},
	// GOAWAY Frame
	{false, 0,
		[]byte{QUICFRAMETYPE_GOAWAY, 0x11, 0x22, 0x33, 0x44, 0x12, 0x34, 0x56, 0x78, 0x01, 0x00},
		QuicFrame{
			frameType:   QUICFRAMETYPE_GOAWAY,
			errorCode:   0x44332211,
			streamId:    0x78563412,
			frameLength: 0,
			frameData:   nil,
		}},
	{true, 0, // TEST 17
		[]byte{QUICFRAMETYPE_GOAWAY, 0x11, 0x22, 0x33, 0x44, 0x12, 0x34, 0x56, 0x78, 0x00, 0x00},
		QuicFrame{
			frameType:   QUICFRAMETYPE_GOAWAY,
			errorCode:   0x44332211,
			streamId:    0x78563412,
			frameLength: 0,
			frameData:   nil,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_GOAWAY, 0x11, 0x22, 0x33, 0x44, 0x12, 0x34, 0x56, 0x78, 0x01, 0x00, 0x1a},
		QuicFrame{
			frameType:   QUICFRAMETYPE_GOAWAY,
			errorCode:   0x44332211,
			streamId:    0x78563412,
			frameLength: 1,
			frameData:   []byte{0x1a},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_GOAWAY, 0x11, 0x22, 0x33, 0x44, 0x12, 0x34, 0x56, 0x78, 0x02, 0x00, 0x1a, 0x2b},
		QuicFrame{
			frameType:   QUICFRAMETYPE_GOAWAY,
			errorCode:   0x44332211,
			streamId:    0x78563412,
			frameLength: 2,
			frameData:   []byte{0x1a, 0x2b},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_GOAWAY, 0x11, 0x22, 0x33, 0x44, 0x12, 0x34, 0x56, 0x78, 0x03, 0x00, 0x1a, 0x2b, 0x3c},
		QuicFrame{
			frameType:   QUICFRAMETYPE_GOAWAY,
			errorCode:   0x44332211,
			streamId:    0x78563412,
			frameLength: 3,
			frameData:   []byte{0x1a, 0x2b, 0x3c},
		}},
	// STOP_WAITING
	{false, 1, // not enough data
		[]byte{QUICFRAMETYPE_STOP_WAITING, 0x42},
		QuicFrame{
			frameType:         QUICFRAMETYPE_STOP_WAITING,
			entropyHash:       0x42,
			leastUnackedDelta: 0x0000000000000001,
		}},
	{true, 1, // TEST 22
		[]byte{QUICFRAMETYPE_STOP_WAITING, 0x42, 0x01},
		QuicFrame{
			frameType:                 QUICFRAMETYPE_STOP_WAITING,
			entropyHash:               0x42,
			leastUnackedDeltaByteSize: 1,
			leastUnackedDelta:         0x0000000000000001,
		}},
	{true, 2,
		[]byte{QUICFRAMETYPE_STOP_WAITING, 0x42, 0x01, 0x02},
		QuicFrame{
			frameType:                 QUICFRAMETYPE_STOP_WAITING,
			entropyHash:               0x42,
			leastUnackedDeltaByteSize: 2,
			leastUnackedDelta:         0x0000000000000201,
		}},
	{true, 4,
		[]byte{QUICFRAMETYPE_STOP_WAITING, 0x42, 0x01, 0x02, 0x03, 0x04},
		QuicFrame{
			frameType:                 QUICFRAMETYPE_STOP_WAITING,
			entropyHash:               0x42,
			leastUnackedDeltaByteSize: 4,
			leastUnackedDelta:         0x0000000004030201,
		}},
	{true, 6,
		[]byte{QUICFRAMETYPE_STOP_WAITING, 0x42, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
		QuicFrame{
			frameType:                 QUICFRAMETYPE_STOP_WAITING,
			entropyHash:               0x42,
			leastUnackedDeltaByteSize: 6,
			leastUnackedDelta:         0x0000060504030201,
		}},
	// STREAM Frame
	{true, 0,
		[]byte{QUICFRAMETYPE_STREAM | QUICFLAG_DATALENGTH | QUICFLAG_STREAMID_32bit | QUICFLAG_BYTEOFFSET_64bit, 0x12, 0x34, 0x56, 0x78,
			0x0a, 0x0b, 0x0c, 0x0d, 0xaa, 0xbb, 0xcc, 0xdd,
			0x03, 0x00,
			0x42, 0x17, 0x89},
		QuicFrame{
			frameType:          QUICFRAMETYPE_STREAM,
			flagDataLength:     true,
			streamIdByteSize:   4,
			streamId:           0x78563412,
			byteOffsetByteSize: 8,
			byteOffset:         0xddccbbaa0d0c0b0a,
			frameLength:        3,
			frameData:          []byte{0x42, 0x17, 0x89},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_STREAM | QUICFLAG_DATALENGTH | QUICFLAG_STREAMID_24bit | QUICFLAG_BYTEOFFSET_64bit, 0x12, 0x34, 0x56,
			0x0a, 0x0b, 0x0c, 0x0d, 0xaa, 0xbb, 0xcc, 0xdd,
			0x03, 0x00,
			0x42, 0x17, 0x89},
		QuicFrame{
			frameType:          QUICFRAMETYPE_STREAM,
			flagDataLength:     true,
			streamIdByteSize:   3,
			streamId:           0x563412,
			byteOffsetByteSize: 8,
			byteOffset:         0xddccbbaa0d0c0b0a,
			frameLength:        3,
			frameData:          []byte{0x42, 0x17, 0x89},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_STREAM | QUICFLAG_DATALENGTH | QUICFLAG_STREAMID_16bit | QUICFLAG_BYTEOFFSET_64bit, 0x12, 0x34,
			0x0a, 0x0b, 0x0c, 0x0d, 0xaa, 0xbb, 0xcc, 0xdd,
			0x03, 0x00,
			0x42, 0x17, 0x89},
		QuicFrame{
			frameType:          QUICFRAMETYPE_STREAM,
			flagDataLength:     true,
			streamIdByteSize:   2,
			streamId:           0x3412,
			byteOffsetByteSize: 8,
			byteOffset:         0xddccbbaa0d0c0b0a,
			frameLength:        3,
			frameData:          []byte{0x42, 0x17, 0x89},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_STREAM | QUICFLAG_DATALENGTH | QUICFLAG_STREAMID_8bit | QUICFLAG_BYTEOFFSET_64bit, 0x12,
			0x0a, 0x0b, 0x0c, 0x0d, 0xaa, 0xbb, 0xcc, 0xdd,
			0x03, 0x00,
			0x42, 0x17, 0x89},
		QuicFrame{
			frameType:          QUICFRAMETYPE_STREAM,
			flagDataLength:     true,
			streamIdByteSize:   1,
			streamId:           0x12,
			byteOffsetByteSize: 8,
			byteOffset:         0xddccbbaa0d0c0b0a,
			frameLength:        3,
			frameData:          []byte{0x42, 0x17, 0x89},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_STREAM | QUICFLAG_DATALENGTH | QUICFLAG_STREAMID_16bit | QUICFLAG_BYTEOFFSET_56bit, 0x12, 0x34,
			0x0a, 0x0b, 0x0c, 0x0d, 0xaa, 0xbb, 0xcc,
			0x03, 0x00,
			0x42, 0x17, 0x89},
		QuicFrame{
			frameType:          QUICFRAMETYPE_STREAM,
			flagDataLength:     true,
			streamIdByteSize:   2,
			streamId:           0x3412,
			byteOffsetByteSize: 7,
			byteOffset:         0xccbbaa0d0c0b0a,
			frameLength:        3,
			frameData:          []byte{0x42, 0x17, 0x89},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_STREAM | QUICFLAG_DATALENGTH | QUICFLAG_STREAMID_16bit | QUICFLAG_BYTEOFFSET_48bit, 0x12, 0x34,
			0x0a, 0x0b, 0x0c, 0x0d, 0xaa, 0xbb,
			0x03, 0x00,
			0x42, 0x17, 0x89},
		QuicFrame{
			frameType:          QUICFRAMETYPE_STREAM,
			flagDataLength:     true,
			streamIdByteSize:   2,
			streamId:           0x3412,
			byteOffsetByteSize: 6,
			byteOffset:         0xbbaa0d0c0b0a,
			frameLength:        3,
			frameData:          []byte{0x42, 0x17, 0x89},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_STREAM | QUICFLAG_DATALENGTH | QUICFLAG_STREAMID_16bit | QUICFLAG_BYTEOFFSET_40bit, 0x12, 0x34,
			0x0a, 0x0b, 0x0c, 0x0d, 0xaa,
			0x03, 0x00,
			0x42, 0x17, 0x89},
		QuicFrame{
			frameType:          QUICFRAMETYPE_STREAM,
			flagDataLength:     true,
			streamIdByteSize:   2,
			streamId:           0x3412,
			byteOffsetByteSize: 5,
			byteOffset:         0xaa0d0c0b0a,
			frameLength:        3,
			frameData:          []byte{0x42, 0x17, 0x89},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_STREAM | QUICFLAG_DATALENGTH | QUICFLAG_STREAMID_16bit | QUICFLAG_BYTEOFFSET_32bit, 0x12, 0x34,
			0x0a, 0x0b, 0x0c, 0x0d,
			0x03, 0x00,
			0x42, 0x17, 0x89},
		QuicFrame{
			frameType:          QUICFRAMETYPE_STREAM,
			flagDataLength:     true,
			streamIdByteSize:   2,
			streamId:           0x3412,
			byteOffsetByteSize: 4,
			byteOffset:         0x0d0c0b0a,
			frameLength:        3,
			frameData:          []byte{0x42, 0x17, 0x89},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_STREAM | QUICFLAG_DATALENGTH | QUICFLAG_STREAMID_16bit | QUICFLAG_BYTEOFFSET_24bit, 0x12, 0x34,
			0x0a, 0x0b, 0x0c,
			0x03, 0x00,
			0x42, 0x17, 0x89},
		QuicFrame{
			frameType:          QUICFRAMETYPE_STREAM,
			flagDataLength:     true,
			streamIdByteSize:   2,
			streamId:           0x3412,
			byteOffsetByteSize: 3,
			byteOffset:         0x0c0b0a,
			frameLength:        3,
			frameData:          []byte{0x42, 0x17, 0x89},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_STREAM | QUICFLAG_DATALENGTH | QUICFLAG_STREAMID_16bit | QUICFLAG_BYTEOFFSET_16bit, 0x12, 0x34,
			0x0a, 0x0b,
			0x03, 0x00,
			0x42, 0x17, 0x89},
		QuicFrame{
			frameType:          QUICFRAMETYPE_STREAM,
			flagDataLength:     true,
			streamIdByteSize:   2,
			streamId:           0x3412,
			byteOffsetByteSize: 2,
			byteOffset:         0x0b0a,
			frameLength:        3,
			frameData:          []byte{0x42, 0x17, 0x89},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_STREAM | QUICFLAG_DATALENGTH | QUICFLAG_STREAMID_16bit, 0x12, 0x34,
			0x03, 0x00,
			0x42, 0x17, 0x89},
		QuicFrame{
			frameType:          QUICFRAMETYPE_STREAM,
			flagDataLength:     true,
			streamIdByteSize:   2,
			streamId:           0x3412,
			byteOffsetByteSize: 0,
			byteOffset:         0,
			frameLength:        3,
			frameData:          []byte{0x42, 0x17, 0x89},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_STREAM | QUICFLAG_DATALENGTH | QUICFLAG_STREAMID_16bit | QUICFLAG_BYTEOFFSET_32bit, 0x12, 0x34,
			0x0a, 0x0b, 0x0c, 0x0d,
			0x01, 0x00,
			0x42},
		QuicFrame{
			frameType:          QUICFRAMETYPE_STREAM,
			flagDataLength:     true,
			streamIdByteSize:   2,
			streamId:           0x3412,
			byteOffsetByteSize: 4,
			byteOffset:         0x0d0c0b0a,
			frameLength:        1,
			frameData:          []byte{0x42},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_STREAM | QUICFLAG_FIN | QUICFLAG_DATALENGTH | QUICFLAG_STREAMID_16bit | QUICFLAG_BYTEOFFSET_32bit, 0x12, 0x34,
			0x0a, 0x0b, 0x0c, 0x0d,
			0x01, 0x00,
			0x42},
		QuicFrame{
			frameType:          QUICFRAMETYPE_STREAM,
			flagFIN:            true,
			flagDataLength:     true,
			streamIdByteSize:   2,
			streamId:           0x3412,
			byteOffsetByteSize: 4,
			byteOffset:         0x0d0c0b0a,
			frameLength:        1,
			frameData:          []byte{0x42},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_STREAM | QUICFLAG_FIN | QUICFLAG_STREAMID_16bit | QUICFLAG_BYTEOFFSET_32bit, 0x12, 0x34,
			0x0a, 0x0b, 0x0c, 0x0d},
		QuicFrame{
			frameType:          QUICFRAMETYPE_STREAM,
			flagFIN:            true,
			flagDataLength:     false,
			streamIdByteSize:   2,
			streamId:           0x3412,
			byteOffsetByteSize: 4,
			byteOffset:         0x0d0c0b0a,
			frameLength:        0,
			frameData:          nil,
		}},
	{false, 0, // not enough data
		[]byte{QUICFRAMETYPE_STREAM | QUICFLAG_DATALENGTH | QUICFLAG_STREAMID_16bit | QUICFLAG_BYTEOFFSET_32bit, 0x12, 0x34,
			0x0a, 0x0b, 0x0c, 0x0d,
			0x01, 0x00},
		QuicFrame{
			frameType:          QUICFRAMETYPE_STREAM,
			flagDataLength:     true,
			streamIdByteSize:   2,
			streamId:           0x3412,
			byteOffsetByteSize: 4,
			byteOffset:         0x0d0c0b0a,
			frameLength:        1,
			frameData:          []byte{0x42},
		}},
	// ACK Frame
	{true, 0,
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_LARGESTOBSERVED_48bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_48bit,
			0x42,
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
			0xca, 0xfe,
			0x00},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			entropyHash:                              0x42,
			largestObservedByteSize:                  6,
			largestObserved:                          0x060504030201,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             0,
			missingPacketSequenceNumberDeltaByteSize: 6,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_LARGESTOBSERVED_48bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_16bit,
			0x42,
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
			0xca, 0xfe,
			0x01,
			0x66, 0x0a, 0x0b, 0x0c, 0x0d},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			entropyHash:                              0x42,
			largestObservedByteSize:                  6,
			largestObserved:                          0x060504030201,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             1,
			deltaFromLargestObserved:                 0x66,
			timeSinceLargestObserved:                 0x0d0c0b0a,
			missingPacketSequenceNumberDeltaByteSize: 2,
		}},
	{true, 0, // TEST 43
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_LARGESTOBSERVED_48bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_16bit,
			0x42,
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
			0xca, 0xfe,
			0x02,
			0x66, 0x0a, 0x0b, 0x0c, 0x0d,
			0x67, 0x89, 0x17},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			entropyHash:                              0x42,
			largestObservedByteSize:                  6,
			largestObserved:                          0x060504030201,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             2,
			deltaFromLargestObserved:                 0x66,
			timeSinceLargestObserved:                 0x0d0c0b0a,
			timestampsDeltaLargestObserved:           [255]byte{0, 0x67},
			timestampsTimeSincePrevious:              [255]uint16{0, 0x1789},
			missingPacketSequenceNumberDeltaByteSize: 2,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_LARGESTOBSERVED_48bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_8bit,
			0x42,
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
			0xca, 0xfe,
			0x03,
			0x66, 0x0a, 0x0b, 0x0c, 0x0d,
			0x67, 0x89, 0x17,
			0x68, 0x84, 0x19},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			entropyHash:                              0x42,
			largestObservedByteSize:                  6,
			largestObserved:                          0x060504030201,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             3,
			deltaFromLargestObserved:                 0x66,
			timeSinceLargestObserved:                 0x0d0c0b0a,
			timestampsDeltaLargestObserved:           [255]byte{0, 0x67, 0x68},
			timestampsTimeSincePrevious:              [255]uint16{0, 0x1789, 0x1984},
			missingPacketSequenceNumberDeltaByteSize: 1,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_NACK | QUICFLAG_LARGESTOBSERVED_48bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_8bit,
			0x42,
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
			0xca, 0xfe,
			0x03,
			0x66, 0x0a, 0x0b, 0x0c, 0x0d,
			0x67, 0x89, 0x17,
			0x68, 0x84, 0x19,
			0x01,
			0xaa, 0x55,
			0x00},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			flagNack:                                 true,
			entropyHash:                              0x42,
			largestObservedByteSize:                  6,
			largestObserved:                          0x060504030201,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             3,
			deltaFromLargestObserved:                 0x66,
			timeSinceLargestObserved:                 0x0d0c0b0a,
			timestampsDeltaLargestObserved:           [255]byte{0, 0x67, 0x68},
			timestampsTimeSincePrevious:              [255]uint16{0, 0x1789, 0x1984},
			missingPacketSequenceNumberDeltaByteSize: 1,
			numMissingRanges:                         1,
			missingPacketsSequenceNumberDelta:        [255]QuicPacketSequenceNumber{0x00000000000000aa},
			missingRangeLength:                       [255]byte{0x55},
			numRevived:                               0,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_NACK | QUICFLAG_LARGESTOBSERVED_48bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_16bit,
			0x42,
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
			0xca, 0xfe,
			0x03,
			0x66, 0x0a, 0x0b, 0x0c, 0x0d,
			0x67, 0x89, 0x17,
			0x68, 0x84, 0x19,
			0x02,
			0xaa, 0xbb, 0x55,
			0xcc, 0xdd, 0x44,
			0x00},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			flagNack:                                 true,
			entropyHash:                              0x42,
			largestObservedByteSize:                  6,
			largestObserved:                          0x060504030201,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             3,
			deltaFromLargestObserved:                 0x66,
			timeSinceLargestObserved:                 0x0d0c0b0a,
			timestampsDeltaLargestObserved:           [255]byte{0, 0x67, 0x68},
			timestampsTimeSincePrevious:              [255]uint16{0, 0x1789, 0x1984},
			missingPacketSequenceNumberDeltaByteSize: 2,
			numMissingRanges:                         2,
			missingPacketsSequenceNumberDelta:        [255]QuicPacketSequenceNumber{0x000000000000bbaa, 0x000000000000ddcc},
			missingRangeLength:                       [255]byte{0x55, 0x44},
			numRevived:                               0,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_NACK | QUICFLAG_LARGESTOBSERVED_48bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_32bit,
			0x42,
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
			0xca, 0xfe,
			0x03,
			0x66, 0x0a, 0x0b, 0x0c, 0x0d,
			0x67, 0x89, 0x17,
			0x68, 0x84, 0x19,
			0x02,
			0xaa, 0xbb, 0x0a, 0x0b, 0x55,
			0xcc, 0xdd, 0x0c, 0x0d, 0x44,
			0x00},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			flagNack:                                 true,
			entropyHash:                              0x42,
			largestObservedByteSize:                  6,
			largestObserved:                          0x060504030201,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             3,
			deltaFromLargestObserved:                 0x66,
			timeSinceLargestObserved:                 0x0d0c0b0a,
			timestampsDeltaLargestObserved:           [255]byte{0, 0x67, 0x68},
			timestampsTimeSincePrevious:              [255]uint16{0, 0x1789, 0x1984},
			missingPacketSequenceNumberDeltaByteSize: 4,
			numMissingRanges:                         2,
			missingPacketsSequenceNumberDelta:        [255]QuicPacketSequenceNumber{0x000000000b0abbaa, 0x000000000d0cddcc},
			missingRangeLength:                       [255]byte{0x55, 0x44},
			numRevived:                               0,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_NACK | QUICFLAG_LARGESTOBSERVED_48bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_48bit,
			0x42,
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
			0xca, 0xfe,
			0x03,
			0x66, 0x0a, 0x0b, 0x0c, 0x0d,
			0x67, 0x89, 0x17,
			0x68, 0x84, 0x19,
			0x02,
			0xaa, 0xbb, 0x0a, 0x0b, 0xa0, 0xb0, 0x55,
			0xcc, 0xdd, 0x0c, 0x0d, 0xc0, 0xd0, 0x44,
			0x00},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			flagNack:                                 true,
			entropyHash:                              0x42,
			largestObservedByteSize:                  6,
			largestObserved:                          0x060504030201,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             3,
			deltaFromLargestObserved:                 0x66,
			timeSinceLargestObserved:                 0x0d0c0b0a,
			timestampsDeltaLargestObserved:           [255]byte{0, 0x67, 0x68},
			timestampsTimeSincePrevious:              [255]uint16{0, 0x1789, 0x1984},
			missingPacketSequenceNumberDeltaByteSize: 6,
			numMissingRanges:                         2,
			missingPacketsSequenceNumberDelta:        [255]QuicPacketSequenceNumber{0x0000b0a00b0abbaa, 0x0000d0c00d0cddcc},
			missingRangeLength:                       [255]byte{0x55, 0x44},
			numRevived:                               0,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_NACK | QUICFLAG_LARGESTOBSERVED_48bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_8bit,
			0x42,
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
			0xca, 0xfe,
			0x03,
			0x66, 0x0a, 0x0b, 0x0c, 0x0d,
			0x67, 0x89, 0x17,
			0x68, 0x84, 0x19,
			0x02,
			0xaa, 0x55,
			0xcc, 0x44,
			0x00},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			flagNack:                                 true,
			entropyHash:                              0x42,
			largestObservedByteSize:                  6,
			largestObserved:                          0x060504030201,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             3,
			deltaFromLargestObserved:                 0x66,
			timeSinceLargestObserved:                 0x0d0c0b0a,
			timestampsDeltaLargestObserved:           [255]byte{0, 0x67, 0x68},
			timestampsTimeSincePrevious:              [255]uint16{0, 0x1789, 0x1984},
			missingPacketSequenceNumberDeltaByteSize: 1,
			numMissingRanges:                         2,
			missingPacketsSequenceNumberDelta:        [255]QuicPacketSequenceNumber{0x00000000000000aa, 0x00000000000000cc},
			missingRangeLength:                       [255]byte{0x55, 0x44},
			numRevived:                               0,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_NACK | QUICFLAG_LARGESTOBSERVED_32bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_8bit,
			0x42,
			0x01, 0x02, 0x03, 0x04,
			0xca, 0xfe,
			0x03,
			0x66, 0x0a, 0x0b, 0x0c, 0x0d,
			0x67, 0x89, 0x17,
			0x68, 0x84, 0x19,
			0x02,
			0xaa, 0x55,
			0xcc, 0x44,
			0x00},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			flagNack:                                 true,
			entropyHash:                              0x42,
			largestObservedByteSize:                  4,
			largestObserved:                          0x000004030201,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             3,
			deltaFromLargestObserved:                 0x66,
			timeSinceLargestObserved:                 0x0d0c0b0a,
			timestampsDeltaLargestObserved:           [255]byte{0, 0x67, 0x68},
			timestampsTimeSincePrevious:              [255]uint16{0, 0x1789, 0x1984},
			missingPacketSequenceNumberDeltaByteSize: 1,
			numMissingRanges:                         2,
			missingPacketsSequenceNumberDelta:        [255]QuicPacketSequenceNumber{0x00000000000000aa, 0x00000000000000cc},
			missingRangeLength:                       [255]byte{0x55, 0x44},
			numRevived:                               0,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_NACK | QUICFLAG_LARGESTOBSERVED_16bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_8bit,
			0x42,
			0x01, 0x02,
			0xca, 0xfe,
			0x03,
			0x66, 0x0a, 0x0b, 0x0c, 0x0d,
			0x67, 0x89, 0x17,
			0x68, 0x84, 0x19,
			0x02,
			0xaa, 0x55,
			0xcc, 0x44,
			0x00},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			flagNack:                                 true,
			entropyHash:                              0x42,
			largestObservedByteSize:                  2,
			largestObserved:                          0x000000000201,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             3,
			deltaFromLargestObserved:                 0x66,
			timeSinceLargestObserved:                 0x0d0c0b0a,
			timestampsDeltaLargestObserved:           [255]byte{0, 0x67, 0x68},
			timestampsTimeSincePrevious:              [255]uint16{0, 0x1789, 0x1984},
			missingPacketSequenceNumberDeltaByteSize: 1,
			numMissingRanges:                         2,
			missingPacketsSequenceNumberDelta:        [255]QuicPacketSequenceNumber{0x00000000000000aa, 0x00000000000000cc},
			missingRangeLength:                       [255]byte{0x55, 0x44},
			numRevived:                               0,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_NACK | QUICFLAG_LARGESTOBSERVED_8bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_8bit,
			0x42,
			0x01,
			0xca, 0xfe,
			0x03,
			0x66, 0x0a, 0x0b, 0x0c, 0x0d,
			0x67, 0x89, 0x17,
			0x68, 0x84, 0x19,
			0x02,
			0xaa, 0x55,
			0xcc, 0x44,
			0x00},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			flagNack:                                 true,
			entropyHash:                              0x42,
			largestObservedByteSize:                  1,
			largestObserved:                          0x000000000001,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             3,
			deltaFromLargestObserved:                 0x66,
			timeSinceLargestObserved:                 0x0d0c0b0a,
			timestampsDeltaLargestObserved:           [255]byte{0, 0x67, 0x68},
			timestampsTimeSincePrevious:              [255]uint16{0, 0x1789, 0x1984},
			missingPacketSequenceNumberDeltaByteSize: 1,
			numMissingRanges:                         2,
			missingPacketsSequenceNumberDelta:        [255]QuicPacketSequenceNumber{0x00000000000000aa, 0x00000000000000cc},
			missingRangeLength:                       [255]byte{0x55, 0x44},
			numRevived:                               0,
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_NACK | QUICFLAG_LARGESTOBSERVED_16bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_8bit,
			0x42,
			0x01, 0x02,
			0xca, 0xfe,
			0x03,
			0x66, 0x0a, 0x0b, 0x0c, 0x0d,
			0x67, 0x89, 0x17,
			0x68, 0x84, 0x19,
			0x02,
			0xaa, 0x55,
			0xcc, 0x44,
			0x03,
			0xa1, 0xa2, 0xb1, 0xb2, 0xc1, 0xc2},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			flagNack:                                 true,
			entropyHash:                              0x42,
			largestObservedByteSize:                  2,
			largestObserved:                          0x000000000201,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             3,
			deltaFromLargestObserved:                 0x66,
			timeSinceLargestObserved:                 0x0d0c0b0a,
			timestampsDeltaLargestObserved:           [255]byte{0, 0x67, 0x68},
			timestampsTimeSincePrevious:              [255]uint16{0, 0x1789, 0x1984},
			missingPacketSequenceNumberDeltaByteSize: 1,
			numMissingRanges:                         2,
			missingPacketsSequenceNumberDelta:        [255]QuicPacketSequenceNumber{0x00000000000000aa, 0x00000000000000cc},
			missingRangeLength:                       [255]byte{0x55, 0x44},
			numRevived:                               3,
			revivedPackets:                           [255]QuicPacketSequenceNumber{0x000000000000a2a1, 0x000000000000b2b1, 0x000000000000c2c1},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_NACK | QUICFLAG_LARGESTOBSERVED_16bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_8bit,
			0x42,
			0x01, 0x02,
			0xca, 0xfe,
			0x03,
			0x66, 0x0a, 0x0b, 0x0c, 0x0d,
			0x67, 0x89, 0x17,
			0x68, 0x84, 0x19,
			0x00,
			0x03,
			0xa1, 0xa2, 0xb1, 0xb2, 0xc1, 0xc2},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			flagNack:                                 true,
			entropyHash:                              0x42,
			largestObservedByteSize:                  2,
			largestObserved:                          0x000000000201,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             3,
			deltaFromLargestObserved:                 0x66,
			timeSinceLargestObserved:                 0x0d0c0b0a,
			timestampsDeltaLargestObserved:           [255]byte{0, 0x67, 0x68},
			timestampsTimeSincePrevious:              [255]uint16{0, 0x1789, 0x1984},
			missingPacketSequenceNumberDeltaByteSize: 1,
			numMissingRanges:                         0,
			missingPacketsSequenceNumberDelta:        [255]QuicPacketSequenceNumber{},
			missingRangeLength:                       [255]byte{},
			numRevived:                               3,
			revivedPackets:                           [255]QuicPacketSequenceNumber{0x000000000000a2a1, 0x000000000000b2b1, 0x000000000000c2c1},
		}},
	{true, 0,
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_NACK | QUICFLAG_LARGESTOBSERVED_16bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_8bit,
			0x42,
			0x01, 0x02,
			0xca, 0xfe,
			0x00,
			0x00,
			0x03,
			0xa1, 0xa2, 0xb1, 0xb2, 0xc1, 0xc2},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			flagNack:                                 true,
			entropyHash:                              0x42,
			largestObservedByteSize:                  2,
			largestObserved:                          0x000000000201,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             0,
			missingPacketSequenceNumberDeltaByteSize: 1,
			numMissingRanges:                         0,
			missingPacketsSequenceNumberDelta:        [255]QuicPacketSequenceNumber{},
			missingRangeLength:                       [255]byte{},
			numRevived:                               3,
			revivedPackets:                           [255]QuicPacketSequenceNumber{0x000000000000a2a1, 0x000000000000b2b1, 0x000000000000c2c1},
		}},
	{false, 0, // not enough data
		[]byte{QUICFRAMETYPE_ACK | QUICFLAG_NACK | QUICFLAG_LARGESTOBSERVED_16bit | QUICFLAG_MISSINGPACKETSEQNUMDELTA_8bit,
			0x42,
			0x01, 0x02,
			0xca, 0xfe,
			0x00,
			0x00,
			0x03,
			0xa1, 0xa2, 0xb1, 0xb2, 0xc1},
		QuicFrame{
			frameType:                                QUICFRAMETYPE_ACK,
			flagNack:                                 true,
			entropyHash:                              0x42,
			largestObservedByteSize:                  2,
			largestObserved:                          0x000000000201,
			largestObservedDeltaTime:                 0xfeca,
			numTimestamp:                             0,
			missingPacketSequenceNumberDeltaByteSize: 1,
			numMissingRanges:                         0,
			missingPacketsSequenceNumberDelta:        [255]QuicPacketSequenceNumber{},
			missingRangeLength:                       [255]byte{},
			numRevived:                               3,
			revivedPackets:                           [255]QuicPacketSequenceNumber{0x000000000000a2a1, 0x000000000000b2b1, 0x000000000000c2c1},
		}},
}

func Test_QuicFrame_ParseData(t *testing.T) {
	var f QuicFrame

	for i, v := range tests_quicframe {
		f.SetLeastUnackedDeltaByteSize(v.leastUnackedDeltaByteSize)
		s, err := f.ParseData(v.data)
		if v.positiveTest {
			if err != nil {
				t.Errorf("QuicFrame.ParseData : error %s in test %v with data[%v]%x", err, i, len(v.data), v.data)
			}
			if s != len(v.data) {
				t.Errorf("QuicFrame.ParseData : invalid size %v in test %v with data[%v]%x", s, i, len(v.data), v.data)
			}
			if v.frame.frameType != f.frameType {
				t.Errorf("QuicFrame.ParseData : invalid Frame Type %x in test %v with data[%v]%x", f.frameType, i, len(v.data), v.data)
			}
			if v.frame.streamIdByteSize != f.streamIdByteSize {
				t.Errorf("QuicFrame.ParseData : invalid Stream ID Byte size %v in test %v with data[%v]%x", f.streamIdByteSize, i, len(v.data), v.data)
			}
			if v.frame.streamId != f.streamId {
				t.Errorf("QuicFrame.ParseData : invalid Stream ID %x in test %v with data[%v]%x", f.streamId, i, len(v.data), v.data)
			}
			if v.frame.byteOffsetByteSize != f.byteOffsetByteSize {
				t.Errorf("QuicFrame.ParseData : invalid Byte Offset size value %v in test %v with data[%v]%x", f.byteOffsetByteSize, i, len(v.data), v.data)
			}
			if v.frame.byteOffset != f.byteOffset {
				t.Errorf("QuicFrame.ParseData : invalid Byte Offset value %x in test %v with data[%v]%x", f.byteOffset, i, len(v.data), v.data)
			}
			if v.frame.errorCode != f.errorCode {
				t.Errorf("QuicFrame.ParseData : invalid Error Code %x in test %v with data[%v]%x", f.errorCode, i, len(v.data), v.data)
			}
			if v.frame.frameLength != f.frameLength {
				t.Errorf("QuicFrame.ParseData : invalid frame data length %v in test %v with data[%v]%x", f.frameLength, i, len(v.data), v.data)
			}
			if !bytes.Equal(v.frame.frameData, f.frameData) {
				t.Errorf("QuicFrame.ParseData : invalid frame data in test %v with data[%v]%x", i, len(v.data), v.data)
			}
			if v.frame.entropyHash != f.entropyHash {
				t.Errorf("QuicFrame.ParseData : invalid Entropy Hash %x in test %v with data[%v]%x", f.entropyHash, i, len(v.data), v.data)
			}
			if v.frame.leastUnackedDeltaByteSize != f.leastUnackedDeltaByteSize {
				t.Errorf("QuicFrame.ParseData : invalid Least Unacket Delta Byte Size %v in test %v with data[%v]%x", f.leastUnackedDeltaByteSize, i, len(v.data), v.data)
			}
			if v.frame.leastUnackedDelta != f.leastUnackedDelta {
				t.Errorf("QuicFrame.ParseData : invalid Least Unacket Delta %x in test %v with data[%v]%x", f.leastUnackedDelta, i, len(v.data), v.data)
			}
			if v.frame.flagDataLength != f.flagDataLength {
				t.Errorf("QuicFrame.ParseData : invalid Data Length flag value %v in test %v with data[%v]%x", f.flagDataLength, i, len(v.data), v.data)
			}
			if v.frame.flagFIN != f.flagFIN {
				t.Errorf("QuicFrame.ParseData : invalid FIN flag value %v in test %v with data[%v]%x", f.flagFIN, i, len(v.data), v.data)
			}
			if v.frame.flagNack != f.flagNack {
				t.Errorf("QuicFrame.ParseData : invalid NACK flag value %v in test %v with data[%v]%x", f.flagNack, i, len(v.data), v.data)
			}
			if v.frame.flagTruncated != f.flagTruncated {
				t.Errorf("QuicFrame.ParseData : invalid TRUNCATED flag value %v in test %v with data[%v]%x", f.flagTruncated, i, len(v.data), v.data)
			}
			if v.frame.largestObservedByteSize != f.largestObservedByteSize {
				t.Errorf("QuicFrame.ParseData : invalid Largest Observed size %v in test %v with data[%v]%x", f.largestObservedByteSize, i, len(v.data), v.data)
			}
			if v.frame.largestObserved != f.largestObserved {
				t.Errorf("QuicFrame.ParseData : invalid Largest Observed Sequence Number %x in test %v with data[%v]%x", f.largestObserved, i, len(v.data), v.data)
			}
			if v.frame.numTimestamp != f.numTimestamp {
				t.Errorf("QuicFrame.ParseData : invalid Number of Timestamp %v in test %v with data[%v]%x", f.numTimestamp, i, len(v.data), v.data)
			}
			if v.frame.timeSinceLargestObserved != f.timeSinceLargestObserved {
				t.Errorf("QuicFrame.ParseData : invalid Time Since Largest Observed %x in test %v with data[%v]%x", f.timeSinceLargestObserved, i, len(v.data), v.data)
			}
			if v.frame.deltaFromLargestObserved != f.deltaFromLargestObserved {
				t.Errorf("QuicFrame.ParseData : invalid Delta from Largest Observed %x in test %v with data[%v]%x", f.deltaFromLargestObserved, i, len(v.data), v.data)
			}
			if v.frame.missingPacketSequenceNumberDeltaByteSize != f.missingPacketSequenceNumberDeltaByteSize {
				t.Errorf("QuicFrame.ParseData : invalid Missing Packet Sequence Number Delta size %v in test %v with data[%v]%x", f.missingPacketSequenceNumberDeltaByteSize, i, len(v.data), v.data)
			}
			for j := range v.frame.timestampsDeltaLargestObserved {
				if v.frame.timestampsDeltaLargestObserved[j] != f.timestampsDeltaLargestObserved[j] {
					t.Errorf("QuicFrame.ParseData : invalid Timestamp Delta Largest Observed [%v]%x in test %v with data[%v]%x", j, f.timestampsDeltaLargestObserved[j], i, len(v.data), v.data)
				}
				if v.frame.timestampsTimeSincePrevious[j] != f.timestampsTimeSincePrevious[j] {
					t.Errorf("QuicFrame.ParseData : invalid Timestamp Time Since Previous [%v]%x in test %v with data[%v]%x", j, f.timestampsTimeSincePrevious[j], i, len(v.data), v.data)
				}
				if v.frame.missingPacketsSequenceNumberDelta[j] != f.missingPacketsSequenceNumberDelta[j] {
					t.Errorf("QuicFrame.ParseData : invalid Missing Packet Sequence Number Delta [%v]%x in test %v with data[%v]%x", j, f.missingPacketsSequenceNumberDelta[j], i, len(v.data), v.data)
				}
				if v.frame.missingRangeLength[j] != f.missingRangeLength[j] {
					t.Errorf("QuicFrame.ParseData : invalid Missing Range Length [%v]%x in test %v with data[%v]%x", j, f.missingRangeLength[j], i, len(v.data), v.data)
				}
				if v.frame.revivedPackets[j] != f.revivedPackets[j] {
					t.Errorf("QuicFrame.ParseData : invalid Revived Packet [%v]%x in test %v with data[%v]%x", j, f.revivedPackets[j], i, len(v.data), v.data)
				}
			}
			if v.frame.numMissingRanges != f.numMissingRanges {
				t.Errorf("QuicFrame.ParseData : invalid Number of Missing Ranges %v in test %v with data[%v]%x", f.numMissingRanges, i, len(v.data), v.data)
			}
			if v.frame.numRevived != f.numRevived {
				t.Errorf("QuicFrame.ParseData : invalid Number of Revived Packets %v in test %v with data[%v]%x", f.numRevived, i, len(v.data), v.data)
			}
		} else if err == nil {
			t.Errorf("QuicFrame.ParseData : missing error in test %v with data[%v]%x", i, len(v.data), v.data)
		}
		f.Erase()
	}

}

func Test_QuicFrame_GetSerializedData(t *testing.T) {
	data := make([]byte, 200)
	for i, v := range tests_quicframe {
		if v.positiveTest {
			v.frame.SetLeastUnackedDeltaByteSize(v.leastUnackedDeltaByteSize)
			s, err := v.frame.GetSerializedData(data)
			if err != nil {
				t.Errorf("QuicFrame.GetSerializedData = error %s while serialized data in test n°%v", err, i)
			}
			if s != len(v.data) {
				t.Errorf("QuicFrame.GetSerializedData = invalid serialized size in test n°%v with data[%v]%x", i, s, data[:s])
			}
			if !bytes.Equal(data[:s], v.data) {
				t.Errorf("QuicFrame.GetSerializedData = invalid serialized data %x in test n°%v with data[%v]%x", data[:s], i, s, v.data)
			}
		}
	}

}
