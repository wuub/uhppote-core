package uhppote

import (
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/uhppoted/uhppote-core/types"
)

func TestGetCardByID(t *testing.T) {
	message := []byte{
		0x17, 0x5a, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xd4, 0x88, 0x5d, 0x00, 0x20, 0x20, 0x01, 0x01,
		0x20, 0x20, 0x12, 0x31, 0x01, 0x00, 0x1d, 0x01, 0x31, 0xd4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	expected := types.Card{
		CardNumber: 6129876,
		From:       types.ToDate(2020, time.January, 1),
		To:         types.ToDate(2020, time.December, 31),
		Doors: map[uint8]uint8{
			1: 1,
			2: 0,
			3: 29,
			4: 1,
		},
		PIN: 54321,
	}

	u := uhppote{
		driver: &stub{
			broadcastTo: func(addr *net.UDPAddr, request []byte, handler func([]byte) bool) ([]byte, error) {
				return message, nil
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

func TestGetCardByIDWithInvalidFromDate(t *testing.T) {
	message := []byte{
		0x17, 0x5a, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xd4, 0x88, 0x5d, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x20, 0x23, 0x12, 0x31, 0x01, 0x00, 0x1d, 0x01, 0x31, 0xd4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	expected := types.Card{
		CardNumber: 6129876,
		From:       types.Date{},
		To:         types.ToDate(2023, time.December, 31),
		Doors: map[uint8]uint8{
			1: 1,
			2: 0,
			3: 29,
			4: 1,
		},
		PIN: 54321,
	}

	u := uhppote{
		driver: &stub{
			broadcastTo: func(addr *net.UDPAddr, request []byte, handler func([]byte) bool) ([]byte, error) {
				return message, nil
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

func TestGetCardByIDWithInvalidToDate(t *testing.T) {
	message := []byte{
		0x17, 0x5a, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xd4, 0x88, 0x5d, 0x00, 0x20, 0x23, 0x01, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x1d, 0x01, 0x31, 0xd4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	expected := types.Card{
		CardNumber: 6129876,
		From:       types.ToDate(2023, time.January, 1),
		To:         types.Date{},
		Doors: map[uint8]uint8{
			1: 1,
			2: 0,
			3: 29,
			4: 1,
		},
		PIN: 54321,
	}

	u := uhppote{
		driver: &stub{
			broadcastTo: func(addr *net.UDPAddr, request []byte, handler func([]byte) bool) ([]byte, error) {
				return message, nil
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
		t.Errorf("Invalid response:\nexpected:%v\ngot:     %v", expected, *card)
	}
}

func TestGetCardByIdWithInvalidDeviceID(t *testing.T) {
	u := uhppote{}

	_, err := u.GetCardByID(0, 8165535)
	if err == nil {
		t.Fatalf("Expected 'Invalid device ID' error, got %v", err)
	}
}

func TestGetCardByIndex(t *testing.T) {
	message := []byte{
		0x17, 0x5c, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xd4, 0x88, 0x5d, 0x00, 0x20, 0x20, 0x01, 0x01,
		0x20, 0x20, 0x12, 0x31, 0x01, 0x00, 0x1d, 0x01, 0x31, 0xd4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	expected := types.Card{
		CardNumber: 6129876,
		From:       types.ToDate(2020, time.January, 1),
		To:         types.ToDate(2020, time.December, 31),
		Doors: map[uint8]uint8{
			1: 1,
			2: 0,
			3: 29,
			4: 1,
		},
		PIN: 54321,
	}

	u := uhppote{
		driver: &stub{
			send: func(addr *net.UDPAddr, request []byte, handler func([]byte) bool) error {
				handler(message)
				return nil
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

func TestGetCardByIndexWithInvalidFromDate(t *testing.T) {
	message := []byte{
		0x17, 0x5c, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xd4, 0x88, 0x5d, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x20, 0x23, 0x12, 0x31, 0x01, 0x00, 0x1d, 0x01, 0x31, 0xd4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	expected := types.Card{
		CardNumber: 6129876,
		From:       types.Date{},
		To:         types.ToDate(2023, time.December, 31),
		Doors: map[uint8]uint8{
			1: 1,
			2: 0,
			3: 29,
			4: 1,
		},
		PIN: 54321,
	}

	u := uhppote{
		driver: &stub{
			send: func(addr *net.UDPAddr, request []byte, handler func([]byte) bool) error {
				handler(message)
				return nil
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

func TestGetCardByIndexWithInvalidToDate(t *testing.T) {
	message := []byte{
		0x17, 0x5c, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xd4, 0x88, 0x5d, 0x00, 0x20, 0x23, 0x01, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x1d, 0x01, 0x31, 0xd4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	expected := types.Card{
		CardNumber: 6129876,
		From:       types.ToDate(2023, time.January, 1),
		To:         types.Date{},
		Doors: map[uint8]uint8{
			1: 1,
			2: 0,
			3: 29,
			4: 1,
		},
		PIN: 54321,
	}

	u := uhppote{
		driver: &stub{
			send: func(addr *net.UDPAddr, request []byte, handler func([]byte) bool) error {
				handler(message)
				return nil
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

	u := uhppote{
		driver: &stub{
			send: func(addr *net.UDPAddr, request []byte, handler func([]byte) bool) error {
				handler(message)
				return nil
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

	u := uhppote{
		driver: &stub{
			send: func(addr *net.UDPAddr, request []byte, handler func([]byte) bool) error {
				handler(message)
				return nil
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

func TestGetCardByIndexWithInvalidDeviceID(t *testing.T) {
	u := uhppote{}

	_, err := u.GetCardByIndex(0, 17)
	if err == nil {
		t.Fatalf("Expected 'Invalid device ID' error, got %v", err)
	}
}
