package uhppote

import (
	"reflect"
	"testing"

	codec "github.com/uhppoted/uhppote-core/encoding/UTO311-L0x"
	"github.com/uhppoted/uhppote-core/types"
)

func TestGetCardByID(t *testing.T) {
	message := []byte{
		0x17, 0x5a, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xd4, 0x88, 0x5d, 0x00, 0x20, 0x20, 0x01, 0x01,
		0x20, 0x20, 0x12, 0x31, 0x01, 0x00, 29, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	expected := types.Card{
		CardNumber: 6129876,
		From:       date("2020-01-01"),
		To:         date("2020-12-31"),
		Doors: map[uint8]types.Permission{
			1: true,
			2: false,
			3: uint8(29),
			4: true,
		},
	}

	u := UHPPOTE{
		driver: &mock{
			send: func(deviceID uint32, request, response interface{}) error {
				return codec.Unmarshal(message, response)
			},
		},
	}

	card, err := u.GetCardByID(423187757, 6129876)
	if err != nil {
		t.Fatalf("Unexpected error returned from GetCardByIndex (%v)", err)
	}

	if card == nil {
		t.Fatalf("Expected response from GetCardByIndex, got:%v", card)
	}

	if !reflect.DeepEqual(*card, expected) {
		t.Errorf("Invalid response:\nexpected:%#v\ngot:     %#v", expected, *card)
	}
}

func TestGetCardByIndex(t *testing.T) {
	message := []byte{
		0x17, 0x5c, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xd4, 0x88, 0x5d, 0x00, 0x20, 0x20, 0x01, 0x01,
		0x20, 0x20, 0x12, 0x31, 0x01, 0x00, 29, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	expected := types.Card{
		CardNumber: 6129876,
		From:       date("2020-01-01"),
		To:         date("2020-12-31"),
		Doors: map[uint8]types.Permission{
			1: true,
			2: false,
			3: uint8(29),
			4: true,
		},
	}

	u := UHPPOTE{
		driver: &mock{
			send: func(deviceID uint32, request, response interface{}) error {
				return codec.Unmarshal(message, response)
			},
		},
	}

	card, err := u.GetCardByIndex(423187757, 67)
	if err != nil {
		t.Fatalf("Unexpected error returned from GetCardByIndex (%v)", err)
	}

	if card == nil {
		t.Fatalf("Expected response from GetCardByIndex, got:%v", card)
	}

	if !reflect.DeepEqual(*card, expected) {
		t.Errorf("Invalid response:\nexpected:%#v\ngot:     %#v", expected, *card)
	}
}

func TestGetCardByIndexWithCardNotFound(t *testing.T) {
	message := []byte{
		0x17, 0x5c, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	u := UHPPOTE{
		driver: &mock{
			send: func(deviceID uint32, request, response interface{}) error {
				return codec.Unmarshal(message, response)
			},
		},
	}

	card, err := u.GetCardByIndex(423187757, 67)
	if err != nil {
		t.Fatalf("Unexpected error returned from GetCardByIndex (%v)", err)
	}

	if card != nil {
		t.Fatalf("Expected <nil> from GetCardByIndex, got:%v", *card)
	}
}

func TestGetCardByIndexWithCardDeleted(t *testing.T) {
	message := []byte{
		0x17, 0x5c, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	u := UHPPOTE{
		driver: &mock{
			send: func(deviceID uint32, request, response interface{}) error {
				return codec.Unmarshal(message, response)
			},
		},
	}

	card, err := u.GetCardByIndex(423187757, 67)
	if err != nil {
		t.Fatalf("Unexpected error returned from GetCardByIndex (%v)", err)
	}

	if card != nil {
		t.Fatalf("Expected <nil> from GetCardByIndex, got:%v", *card)
	}
}
