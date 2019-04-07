package uhppote

import (
	"reflect"
	"testing"
	"time"
	"uhppote/encoding"
)

func TestMarshalGetCardByIndexRequest(t *testing.T) {
	expected := []byte{
		0x17, 0x5C, 0x00, 0x00, 0x2D, 0x55, 0x39, 0x19, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	request := GetCardByIndexRequest{
		MsgType:      0x5C,
		SerialNumber: 423187757,
		Index:        4,
	}

	m, err := uhppote.Marshal(request)

	if err != nil {
		t.Errorf("Marshal returned unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Marshall returned Invalid byte array:\nExpected:\n%s\nReturned:\n%s", print(expected), print(m))
		return
	}
}

func TestUnmarshalGetCardResponse(t *testing.T) {
	message := []byte{
		0x17, 0x5C, 0x00, 0x00, 0x2D, 0x55, 0x39, 0x19, 0xAC, 0xE8, 0x5D, 0x00, 0x20, 0x19, 0x02, 0x03,
		0x20, 0x19, 0x12, 0x29, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply := GetCardResponse{}

	err := uhppote.Unmarshal(message, &reply)

	if err != nil {
		t.Errorf("Unmarshal returned error from valid message: %v\n", err)
	}

	if reply.MsgType != 0x5C {
		t.Errorf("Incorrect 'message type' - expected:%02X, got:%02x\n", 0x5C, reply.MsgType)
	}

	if reply.SerialNumber != 423187757 {
		t.Errorf("Incorrect 'serial number' from valid message: %v\n", reply.SerialNumber)
	}

	if reply.CardNumber != 6154412 {
		t.Errorf("Incorrect 'card number' from valid message: %v\n", reply.CardNumber)
	}

	from, _ := time.ParseInLocation("2006-01-02", "2019-02-03", time.Local)
	if reply.From.Date != from {
		t.Errorf("Incorrect 'from date' - expected:%s, got:%s\n", from.Format("2006-01-02"), reply.From)
	}

	to, _ := time.ParseInLocation("2006-01-02", "2019-12-29", time.Local)
	if reply.To.Date != to {
		t.Errorf("Incorrect 'to date' - expected:%s, got:%s\n", to.Format("2006-01-02"), reply.To)
	}

	if reply.Door1 != false {
		t.Errorf("Incorrect 'door 1' - expected:%v, got:%v\n", false, reply.Door1)
	}

	if reply.Door2 != false {
		t.Errorf("Incorrect 'door 2' - expected:%v, got:%v\n", false, reply.Door2)
	}

	if reply.Door3 != true {
		t.Errorf("Incorrect 'door 3' - expected:%v, got:%v\n", true, reply.Door3)
	}

	if reply.Door4 != true {
		t.Errorf("Incorrect 'door 4' - expected:%v, got:%v\n", true, reply.Door4)
	}
}
