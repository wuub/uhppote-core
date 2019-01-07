package messages

import (
	"encoding/binary"
	"fmt"
	"time"
	"uhppote/types"
)

type GetTime struct {
	StartOfMessage byte
	MsgType        byte
	DateTime       types.DateTime
}

func NewGetTime(msg []byte) (*GetTime, error) {
	serialNumber := binary.LittleEndian.Uint32(msg[4:8])
	timestamp := fmt.Sprintf("%04X-%02X-%02X %02X:%02X:%02X", msg[8:10], msg[10:11], msg[11:12], msg[12:13], msg[13:14], msg[14:15])
	datetime, _ := time.ParseInLocation("2006-01-02 15:04:05", timestamp, time.Local)

	return &GetTime{msg[0], msg[1], types.DateTime{
		serialNumber,
		datetime,
	}}, nil
}
