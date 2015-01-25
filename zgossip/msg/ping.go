package msg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"

	"github.com/zeromq/goczmq"
)

// Client signals liveness
type Ping struct {
	routingId []byte
	version   byte
}

// New creates new Ping message.
func NewPing() *Ping {
	ping := &Ping{}
	return ping
}

// String returns print friendly name.
func (p *Ping) String() string {
	str := "ZGOSSIP_MSG_PING:\n"
	str += fmt.Sprintf("    version = %v\n", p.version)
	return str
}

// Marshal serializes the message.
func (p *Ping) Marshal() ([]byte, error) {
	// Calculate size of serialized data
	bufferSize := 2 + 1 // Signature and message ID

	// version is a 1-byte integer
	bufferSize += 1

	// Now serialize the message
	tmpBuf := make([]byte, bufferSize)
	tmpBuf = tmpBuf[:0]
	buffer := bytes.NewBuffer(tmpBuf)
	binary.Write(buffer, binary.BigEndian, Signature)
	binary.Write(buffer, binary.BigEndian, PingId)

	// version
	value, _ := strconv.ParseUint("1", 10, 1*8)
	binary.Write(buffer, binary.BigEndian, byte(value))

	return buffer.Bytes(), nil
}

// Unmarshal unmarshals the message.
func (p *Ping) Unmarshal(frames ...[]byte) error {
	if frames == nil {
		return errors.New("Can't unmarshal empty message")
	}

	frame := frames[0]
	frames = frames[1:]

	buffer := bytes.NewBuffer(frame)

	// Get and check protocol signature
	var signature uint16
	binary.Read(buffer, binary.BigEndian, &signature)
	if signature != Signature {
		return fmt.Errorf("invalid signature %X != %X", Signature, signature)
	}

	// Get message id and parse per message type
	var id uint8
	binary.Read(buffer, binary.BigEndian, &id)
	if id != PingId {
		return errors.New("malformed Ping message")
	}
	// version
	binary.Read(buffer, binary.BigEndian, &p.version)
	if p.version != 1 {
		return errors.New("malformed version message")
	}

	return nil
}

// Send sends marshaled data through 0mq socket.
func (p *Ping) Send(sock *goczmq.Sock) (err error) {
	frame, err := p.Marshal()
	if err != nil {
		return err
	}

	socType := sock.GetType()
	if err != nil {
		return err
	}

	// If we're sending to a ROUTER, we send the routingId first
	if socType == goczmq.ROUTER {
		err = sock.SendFrame(p.routingId, goczmq.MORE)
		if err != nil {
			return err
		}
	}

	// Now send the data frame
	err = sock.SendFrame(frame, 0)
	if err != nil {
		return err
	}

	return err
}

// RoutingId returns the routingId for this message, routingId should be set
// whenever talking to a ROUTER.
func (p *Ping) RoutingId() []byte {
	return p.routingId
}

// SetRoutingId sets the routingId for this message, routingId should be set
// whenever talking to a ROUTER.
func (p *Ping) SetRoutingId(routingId []byte) {
	p.routingId = routingId
}

// Setversion sets the version.
func (p *Ping) SetVersion(version byte) {
	p.version = version
}

// version returns the version.
func (p *Ping) Version() byte {
	return p.version
}
