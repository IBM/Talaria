package api

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"

	"github.com/ibm/opentalaria/protocol"
)

type API interface {
	Name() string
	GeneratePayload() ([]byte, error)
	GetRequest() Request
	GetHeaderVersion(requestVersion int16) int16
}

type Request struct {
	Header  protocol.RequestHeader
	Message []byte
	Conn    net.Conn
}

func HandleResponse(api API) error {
	payload := make([]byte, 0)

	resHeader := protocol.ResponseHeader{
		Version:       api.GetHeaderVersion(api.GetRequest().Header.RequestApiVersion),
		CorrelationID: api.GetRequest().Header.CorrelationID,
	}

	resHeaderBytes, err := protocol.Encode(&resHeader)
	if err != nil {
		return err
	}
	// TODO: calculate the payload size before merging the header with the message payload, to avoid the append operation
	payload = append(payload, resHeaderBytes...)

	msg, err := api.GeneratePayload()
	if err != nil {
		return err
	}

	payload = append(payload, msg...)

	// prepend payload size to the final byte array that will be sent back via the wire
	result := make([]byte, 0)
	result = binary.BigEndian.AppendUint32(result, uint32(len(payload)))
	result = append(result, payload...)

	slog.Debug(fmt.Sprintf("writing %d bytes", len(result)), "api", api.Name())

	_, err = api.GetRequest().Conn.Write(result)
	return err
}
