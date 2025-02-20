package api

import (
	"log/slog"

	"github.com/ibm/opentalaria/config"
	"github.com/ibm/opentalaria/utils"

	"github.com/ibm/opentalaria/protocol"
)

type CreateTopicsAPI struct {
	Request Request
	Config  *config.Config
}

func (m CreateTopicsAPI) Name() string {
	return "CreateTopics"
}

func (m CreateTopicsAPI) GetRequest() Request {
	return m.Request
}

func (m CreateTopicsAPI) GetHeaderVersion(requestVersion int16) int16 {
	return (&protocol.CreateTopicsResponse{Version: requestVersion}).GetHeaderVersion()
}

func (m CreateTopicsAPI) GeneratePayload() ([]byte, error) {
	req := protocol.CreateTopicsRequest{}
	_, err := protocol.VersionedDecode(m.GetRequest().Message, &req, m.GetRequest().Header.RequestApiVersion)

	resp := m.GenerateCreateTopicsResponse(m.GetRequest().Header.RequestApiVersion, req, err)

	return protocol.Encode(resp)
}

func (m CreateTopicsAPI) GenerateCreateTopicsResponse(version int16, req protocol.CreateTopicsRequest, err error) *protocol.CreateTopicsResponse {
	response := protocol.CreateTopicsResponse{}

	response.Version = version
	// TODO: handle throttle time
	response.ThrottleTimeMs = 0

	for _, topic := range req.Topics {
		err := m.Config.Plugin.AddTopic(topic)

		errorCode := int16(utils.ErrNoError)
		if err != nil {
			slog.Error(err.Error())
			errorCode = int16(utils.ErrInvalidRequest)
		}

		response.Topics = append(response.Topics, protocol.CreatableTopicResult{
			Version:   req.Version,
			Name:      topic.Name,
			ErrorCode: errorCode,
		})
	}

	return &response
}
