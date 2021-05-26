package uhppote

import (
	"github.com/uhppoted/uhppote-core/messages"
	"github.com/uhppoted/uhppote-core/types"
)

func (u *uhppote) GetCards(deviceID uint32) (uint32, error) {
	request := messages.GetCardsRequest{
		SerialNumber: types.SerialNumber(deviceID),
	}

	reply := messages.GetCardsResponse{}

	err := u.send(deviceID, request, &reply)
	if err != nil {
		return 0, err
	}

	return reply.Records, nil
}
