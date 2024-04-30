package uhppote

import (
	"fmt"

	"github.com/uhppoted/uhppote-core/messages"
	"github.com/uhppoted/uhppote-core/types"
)

func (u *uhppote) GetCards(deviceID uint32) (uint32, error) {
	if deviceID == 0 {
		return 0, fmt.Errorf("invalid device ID (%v)", deviceID)
	}

	request := messages.GetCardsRequest{
		SerialNumber: types.SerialNumber(deviceID),
	}

	if reply, err := u.sendTo(deviceID, request, messages.GetCardsResponse{}); err != nil {
		return 0, err
	} else {
		return reply.(messages.GetCardsResponse).Records, nil
	}
}
