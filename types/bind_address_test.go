package types

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestBindAddrString(t *testing.T) {
	tests := map[int]string{
		0:     "192.168.1.100",
		1:     "192.168.1.100:1",
		60000: "192.168.1.100:60000",
	}

	for p, expected := range tests {
		bind := BindAddr{
			IP:   []byte{192, 168, 1, 100},
			Port: p,
		}

		if s := bind.String(); s != expected {
			t.Errorf("Incorrect string - expected:%v, got:%v", expected, s)
		}
	}
}

func TestBindAddrResolve(t *testing.T) {
	tests := map[string]BindAddr{
		"192.168.1.100":       BindAddr{IP: []byte{192, 168, 1, 100}, Port: 0},
		"192.168.1.100:0":     BindAddr{IP: []byte{192, 168, 1, 100}, Port: 0},
		"192.168.1.100:10001": BindAddr{IP: []byte{192, 168, 1, 100}, Port: 10001},
		"192.168.1.100:60001": BindAddr{IP: []byte{192, 168, 1, 100}, Port: 60001},
	}

	for s, expected := range tests {
		addr, err := ResolveBindAddr(s)

		if err != nil {
			t.Fatalf("Unexpected error resolving bind address %v (%v)", s, err)
		} else if addr == nil {
			t.Fatalf("Failed to resolve valid bind address %v (%v)", s, addr)
		}

		if !reflect.DeepEqual(*addr, expected) {
			t.Errorf("Incorrectly resolved bind address %v - expected:%v, got:%v", s, expected, *addr)
		}
	}
}

func TestInvalidBindAddrResolve(t *testing.T) {
	s := "192.168.1.100:60000"
	addr, err := ResolveBindAddr(s)

	if err == nil {
		t.Errorf("Expected error resolving bind address %v, got:%v (%v)", s, addr, err)
	}
}

func TestBindAddrSet(t *testing.T) {
	tests := map[string]BindAddr{
		"192.168.1.100":       BindAddr{IP: []byte{192, 168, 1, 100}, Port: 0},
		"192.168.1.100:0":     BindAddr{IP: []byte{192, 168, 1, 100}, Port: 0},
		"192.168.1.100:1":     BindAddr{IP: []byte{192, 168, 1, 100}, Port: 1},
		"192.168.1.100:60001": BindAddr{IP: []byte{192, 168, 1, 100}, Port: 60001},
	}

	for s, expected := range tests {
		addr := BindAddr{}
		err := addr.Set(s)

		if err != nil {
			t.Fatalf("Error 'setting' bind address %v (%v)", s, err)
		}

		if !reflect.DeepEqual(addr, expected) {
			t.Errorf("Incorrect 'bind' address '%v' - expected:%v, got:%v", s, expected, addr)
		}
	}
}

func TestBindAddrMarshalJSON(t *testing.T) {
	tests := map[int]string{
		0:     `"192.168.1.100"`,
		1:     `"192.168.1.100:1"`,
		60000: `"192.168.1.100:60000"`,
	}

	for p, expected := range tests {
		bind := BindAddr{
			IP:   []byte{192, 168, 1, 100},
			Port: p,
		}

		if bytes, err := json.Marshal(bind); err != nil {
			t.Fatalf("Error marshaling BindAddr (%v)", err)
		} else if s := string(bytes); s != expected {
			t.Errorf("Incorrect JSON string - expected:%v, got:%v", expected, s)
		}
	}
}

func TestBindAddrUnmarshalJSON(t *testing.T) {
	tests := map[string]BindAddr{
		`"192.168.1.100"`: BindAddr{
			IP:   []byte{192, 168, 1, 100},
			Port: 0,
		},

		`"192.168.1.100:12345"`: BindAddr{
			IP:   []byte{192, 168, 1, 100},
			Port: 12345,
		},

		`"192.168.1.100:0"`: BindAddr{
			IP:   []byte{192, 168, 1, 100},
			Port: 0,
		},
	}

	for s, expected := range tests {
		bind := BindAddr{}

		if err := json.Unmarshal([]byte(s), &bind); err != nil {
			t.Fatalf("Error unmarshaling BindAddr '%v' (%v)", s, err)
		} else if !reflect.DeepEqual(bind, expected) {
			t.Errorf("Incorrectly unmarshalled bind address '%v'\nexpected:%v\ngot:     %v", s, expected, bind)
		}
	}
}

func TestBindAddrEqual(t *testing.T) {
	tests := []struct {
		bind     BindAddr
		address  Address
		expected bool
	}{
		{
			BindAddr{
				IP:   []byte{192, 168, 1, 100},
				Port: 0,
			},
			Address{
				IP:   []byte{192, 168, 1, 100},
				Port: 0,
			},
			true,
		},
		{
			BindAddr{
				IP:   []byte{192, 168, 1, 100},
				Port: 0,
			},
			Address{
				IP:   []byte{192, 168, 1, 100},
				Port: 12345,
			},
			true,
		},
		{
			BindAddr{
				IP:   []byte{192, 168, 1, 100},
				Port: 0,
			},
			Address{
				IP:   []byte{192, 168, 1, 125},
				Port: 0,
			},
			false,
		},
	}

	for _, test := range tests {
		equal := test.bind.Equal(&test.address)

		if equal != test.expected {
			t.Errorf("Error comparing bind address - expected:%v, got:%v", test.expected, equal)
		}
	}
}

func TestBindAddrClone(t *testing.T) {
	bind := BindAddr{
		IP:   []byte{192, 168, 1, 100},
		Port: 12345,
	}

	expected := BindAddr{
		IP:   []byte{192, 168, 1, 100},
		Port: 12345,
	}

	clone := bind.Clone()

	bind.IP = []byte{192, 168, 1, 125}
	bind.Port = 54321

	if !reflect.DeepEqual(*clone, expected) {
		t.Errorf("Invalid BindAddress clone\nexpected:%#v\ngot:     %#v", expected, clone)
	}

	if reflect.DeepEqual(bind, expected) {
		t.Errorf("Sanity check failed")
	}
}
