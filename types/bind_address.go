package types

import (
	"encoding/json"
	"fmt"
	"net/netip"
	"regexp"
)

type BindAddr struct {
	netip.AddrPort
}

const BIND_PORT = 0

/*
 * MustParseBindAddr invokes ParseBindAddr and panics on error.
 *
 * It is intended for use in tests with hard-coded strings.
 */
func MustParseBindAddr(s string) BindAddr {
	if addr, err := ParseBindAddr(s); err != nil {
		panic(err)
	} else {
		return addr
	}
}

/*
 * ParseBindAddr parses a string as a bind address.
 *
 * It doesn't do any name resolution i.e.: both the address and the port must be numeric.
 */
func ParseBindAddr(s string) (BindAddr, error) {
	if addr, err := netip.ParseAddrPort(s); err != nil {
		return BindAddr{}, err
	} else {
		return BindAddr{
			addr,
		}, nil
	}
}

/*
 * BindAddrFrom contructs a BindAddr from an address and port.
 */
func BindAddrFrom(addr netip.Addr, port uint16) BindAddr {
	return BindAddr{
		netip.AddrPortFrom(addr, port),
	}
}

/*
 * Formats the bind address and port as an address:port string.
 *
 * Return only the bind address if bind port is the default port (60000).
 */
func (a BindAddr) String() string {
	if a.Port() == BIND_PORT {
		return fmt.Sprintf("%v", a.Addr())
	} else {
		return fmt.Sprintf("%v", a.AddrPort)
	}
}

func (a *BindAddr) Set(v string) error {
	addr, err := ResolveBindAddr(v)
	if err != nil {
		return err
	} else if !addr.IsValid() {
		return fmt.Errorf("invalid bind address '%v'", v)
	}

	*a = addr

	return nil
}

func (a BindAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%v", a))
}

func (a *BindAddr) UnmarshalJSON(bytes []byte) error {
	var s string

	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}

	addr, err := ResolveBindAddr(s)
	if err != nil {
		return err
	}

	*a = addr

	return nil
}

func (a *BindAddr) Equal(addr *Address) bool {
	switch {
	case a == nil && addr == nil:
		return true

	case a != nil && addr != nil:
		p := fmt.Sprintf("%v", a.Addr())
		q := fmt.Sprintf("%v", addr.IP)
		return p == q

	default:
		return false
	}
}

func (a *BindAddr) Clone() *BindAddr {
	if a != nil {
		addr := *a
		return &addr
	}

	return nil
}

func ResolveBindAddr(s string) (BindAddr, error) {
	if matched, err := regexp.MatchString(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}:[0-9]{1,5}`, s); err != nil {
		return BindAddr{}, err
	} else if matched {
		if addr, err := netip.ParseAddrPort(s); err != nil {
			return BindAddr{}, err
		} else if addr.Port() == DEFAULT_PORT {
			return BindAddr{}, fmt.Errorf("%v: invalid 'bind' port (%v)", addr, addr.Port())
		} else {
			return BindAddr{addr}, nil
		}
	}

	if matched, err := regexp.MatchString(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`, s); err != nil {
		return BindAddr{}, err
	} else if matched {
		if addr, err := netip.ParseAddr(s); err != nil {
			return BindAddr{}, err
		} else {
			return BindAddrFrom(addr, BIND_PORT), nil
		}
	}

	return BindAddr{}, fmt.Errorf("%s is not a valid address:port", s)
}
