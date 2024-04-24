package uhppote

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/netip"
	"os"
	"time"

	codec "github.com/uhppoted/uhppote-core/encoding/UTO311-L0x"
	"github.com/uhppoted/uhppote-core/types"
)

const (
	VERSION string = "v0.8.9" // common version number used across uhppoted ecosystem
)

type driver interface {
	Broadcast([]byte, *net.UDPAddr) ([][]byte, error)
	Send([]byte, *net.UDPAddr, func([]byte) bool) error
	SendTo(*net.UDPAddr, []byte) ([]byte, error)
	Listen(chan any, chan any, func([]byte)) error
}

type uhppote struct {
	bindAddr      *net.UDPAddr
	broadcastAddr *net.UDPAddr
	listenAddr    *net.UDPAddr
	devices       map[uint32]Device
	debug         bool
	driver        driver
}

func NewUHPPOTE(
	bindAddr types.BindAddr,
	broadcastAddr types.BroadcastAddr,
	listenAddr types.ListenAddr,
	timeout time.Duration,
	devices []Device,
	debug bool) IUHPPOTE {
	bind := net.UDPAddr(bindAddr)
	broadcast := net.UDPAddr(broadcastAddr)
	listen := net.UDPAddr(listenAddr)

	uhppote := uhppote{
		bindAddr:      &bind,
		broadcastAddr: &broadcast,
		listenAddr:    &listen,
		devices:       map[uint32]Device{},
		driver: &udp{
			bindAddr:   bind,
			listenAddr: listen,
			timeout:    timeout,
			debug:      debug,
		},
		debug: debug,
	}

	for _, device := range devices {
		uhppote.devices[device.DeviceID] = device.Clone()
	}

	return &uhppote
}

func (u *uhppote) DeviceList() map[uint32]Device {
	list := map[uint32]Device{}
	if u != nil {
		for k, v := range u.devices {
			list[k] = v
		}
	}

	return list
}

func (u *uhppote) BroadcastAddr() *net.UDPAddr {
	if u != nil {
		return u.broadcastAddr
	}

	return nil
}

func (u *uhppote) ListenAddr() *net.UDPAddr {
	if u != nil {
		return u.listenAddr
	}

	return nil
}

func (u *uhppote) broadcast(request, reply any) ([]any, error) {
	return u.sendTo(0, request, reply)
}

// Sends a UDP message to a specific device but anticipates replies from more than one device
// because it may fall back to the broadcast address if the device ID has no configured IP
// address.
func (u *uhppote) sendTo(serialNumber uint32, request, reply any) ([]any, error) {
	m, err := codec.Marshal(request)
	if err != nil {
		return nil, err
	}

	f := func() ([][]byte, error) {
		if serialNumber == 0 {
			return u.udpBroadcastTo(m)
		} else if controller, ok := u.devices[serialNumber]; !ok {
			return u.udpBroadcastTo(m)
		} else if controller.Address == nil {
			return u.udpBroadcastTo(m)
		} else if controller.Protocol == "tcp" {
			return u.tcpSendTo(*controller.Address, m)
		} else {
			return u.udpSendTo(*controller.Address, m)
		}
	}

	if responses, err := f(); err != nil {
		return nil, err
	} else {
		replies := []any{}

		for _, bytes := range responses {
			// ... discard invalid replies
			if len(bytes) != 64 {
				u.debugf(" ... receive error", fmt.Errorf("invalid message length - expected:%v, got:%v", 64, len(bytes)))
				continue
			}

			// ... discard replies without a valid device ID
			if deviceID := binary.LittleEndian.Uint32(bytes[4:8]); serialNumber != 0 && deviceID != serialNumber {
				u.debugf(" ... receive error", fmt.Errorf("invalid device ID - expected:%v, got:%v", serialNumber, deviceID))
				continue
			}

			// ... discard unparseable replies
			v, err := codec.UnmarshalAs(bytes, reply)
			if err != nil {
				u.debugf(" ... receive error", err)
				continue
			}

			replies = append(replies, v)
		}

		return replies, nil
	}
}

// Broadcasts a UDP request and returns all received replies.
func (u *uhppote) udpBroadcastTo(request []byte) ([][]byte, error) {
	dest := u.broadcastAddress()

	return u.driver.Broadcast(request, dest)
}

// Sends a UDP message to a specific device but anticipates replies from more than one device
// because the controller address may be a broadcast address (unlikely but possible).
func (u *uhppote) udpSendTo(address netip.AddrPort, request []byte) ([][]byte, error) {
	dest := net.UDPAddrFromAddrPort(address)

	if response, err := u.driver.SendTo(dest, request); err != nil {
		return nil, err
	} else {
		return [][]byte{
			response,
		}, nil
	}
}

/* Sends the request as a TCP message and wraps the reply (if any) as slice of byte arrays.
 *
 */
func (u *uhppote) tcpSendTo(address netip.AddrPort, request []byte) ([][]byte, error) {
	// dest := net.UDPAddrFromAddrPort(address)
	//
	// if response, err := u.driver.SendTo(dest, request); err != nil {
	// 	return nil, err
	// } else {
	// 	return [][]byte{
	// 		response,
	// 	}, nil
	// }

	return nil, fmt.Errorf("NOT IMPLEMENTED")
}

func (u *uhppote) send(serialNumber uint32, request, reply any) error {
	dest := u.broadcastAddress()
	if device, ok := u.devices[serialNumber]; ok {
		if device.Address != nil {
			dest = net.UDPAddrFromAddrPort(*device.Address)
		}
	}

	m, err := codec.Marshal(request)
	if err != nil {
		return err
	}

	var handler func([]byte) bool

	if reply != nil {
		handler = func(bytes []byte) bool {
			// ... discard invalid replies
			if len(bytes) != 64 {
				u.debugf(" ... receive error", fmt.Errorf("invalid message length - expected:%v, got:%v", 64, len(bytes)))
				return false
			}

			// ... discard replies without a valid device ID
			if deviceID := binary.LittleEndian.Uint32(bytes[4:8]); deviceID != serialNumber {
				u.debugf(" ... receive error", fmt.Errorf("invalid device ID - expected:%v, got:%v", serialNumber, deviceID))
				return false
			}

			// .. discard unparseable messages
			if err := codec.Unmarshal(bytes, reply); err != nil {
				u.debugf(" ... receive error", err)
				return false
			}

			return true
		}
	}

	if err := u.driver.Send(m, dest, handler); err != nil {
		return err
	}

	return nil
}

func (u *uhppote) listen(p chan *event, q chan os.Signal, listener Listener) error {
	handler := func(bytes []byte) {
		// ... discard invalid events
		if len(bytes) != 64 {
			listener.OnError(fmt.Errorf("invalid message length - expected:%v, got:%v", 64, len(bytes)))
			return
		}

		// ... discard replies without a valid device ID
		if deviceID := binary.LittleEndian.Uint32(bytes[4:8]); deviceID == 0 {
			listener.OnError(fmt.Errorf("invalid device ID (%v)", deviceID))
			return
		}

		// .. discard unparseable messages
		e := event{}
		if err := codec.Unmarshal(bytes, &e); err != nil {
			listener.OnError(err)
			return
		}

		p <- &e
	}

	signal := make(chan any)
	closed := make(chan any)

	err := u.driver.Listen(signal, closed, handler)
	if err != nil {
		return err
	}

	listener.OnConnected()

	<-q
	close(signal)
	<-closed

	return nil
}

func (u *uhppote) broadcastAddress() *net.UDPAddr {
	if u.broadcastAddr != nil {
		return u.broadcastAddr
	}

	addr := net.UDPAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: 60000,
		Zone: "",
	}

	copy(addr.IP, net.IPv4bcast)

	return &addr
}

func (u *uhppote) debugf(msg string, err error) {
	if u.debug {
		if err != nil {
			fmt.Printf("%v: %v\n", msg, err)
		} else {
			fmt.Printf("%v\n", msg)
		}
	}
}
