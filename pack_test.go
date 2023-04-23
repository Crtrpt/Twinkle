package twinkle

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBit(t *testing.T) {
	data := 0b1000_0000
	b := BitSet(byte(data), 0)
	assert.Equal(t, b, byte(0b1000_0001))
	c := BitClear(byte(b), 0)
	assert.Equal(t, c, byte(data))
	assert.Equal(t, BitGet(byte(data), 7), true)
	assert.Equal(t, BitGet(byte(data), 6), false)
}

func TestPackAndUnpack(t *testing.T) {
	role := 0
	ip := net.IP([]byte{127, 0, 0, 1})
	port := 9001
	payload := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	enc := UDPForwardPacket(role, ip, port, payload)

	dec_role, dec_ip, dec_port, dec_paload, err := UDPForwardUnPacket(enc)
	if err != nil {
		t.Fatalf("error %v", err)
		return
	}

	assert.Equal(t, role, dec_role)
	assert.Equal(t, ip, dec_ip)
	assert.Equal(t, port, dec_port)
	assert.Equal(t, payload, dec_paload)
}

func TestIPV6PackAndUnpack(t *testing.T) {
	role := 0
	ip := net.IP([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 127, 0, 0, 1})
	port := 9001
	payload := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	enc := UDPForwardPacket(role, ip, port, payload)

	dec_role, dec_ip, dec_port, dec_paload, err := UDPForwardUnPacket(enc)
	if err != nil {
		t.Fatalf("error %v", err)
		return
	}

	assert.Equal(t, role, dec_role)
	assert.Equal(t, ip, dec_ip)
	assert.Equal(t, port, dec_port)
	assert.Equal(t, payload, dec_paload)
}
