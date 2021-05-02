package uhppote

import (
	"github.com/uhppoted/uhppote-core/messages"
	"github.com/uhppoted/uhppote-core/types"
)

func (u *UHPPOTE) DeleteCards(serialNumber uint32) (bool, error) {
	driver := iuhppote(u)
	if u.driver != nil {
		driver = u.driver
	}

	request := messages.DeleteCardsRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		MagicWord:    0x55aaaa55,
	}

	reply := messages.DeleteCardsResponse{}

	err := driver.Send(serialNumber, request, &reply)
	if err != nil {
		return false, err
	}

	return reply.Succeeded, nil
}
