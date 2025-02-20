package api

import (
	"log/slog"

	"github.com/ibm/opentalaria/config"
	"github.com/ibm/opentalaria/protocol"
	"github.com/ibm/opentalaria/utils"
)

type ProduceAPI struct {
	Request Request
	Config  *config.Config
}

func (p ProduceAPI) Name() string {
	return "Produce"
}

func (p ProduceAPI) GetRequest() Request {
	return p.Request
}

func (p ProduceAPI) GetHeaderVersion(requestVersion int16) int16 {
	return (&protocol.ProduceResponse{Version: requestVersion}).GetHeaderVersion()
}

// TODO: this is a placeholder function for now. We need to implement a backend that handles cluster topology in order to implement the API correctly and consume the messages.
func (p ProduceAPI) GeneratePayload() ([]byte, error) {
	req := protocol.ProduceRequest{}
	_, err := protocol.VersionedDecode(p.GetRequest().Message, &req, p.GetRequest().Header.RequestApiVersion)
	if err != nil {
		return nil, err
	}

	resp := protocol.ProduceResponse{
		Version: p.GetRequest().Header.RequestApiVersion,
	}

	for _, topic := range req.TopicData {
		topicResponse := protocol.TopicProduceResponse{}
		topicResponse.Version = resp.Version
		topicResponse.Name = topic.Name

		for _, partition := range topic.PartitionData {
			slog.Debug("Received records", "records", partition.Records)

			topicResponse.PartitionResponses = append(topicResponse.PartitionResponses, protocol.PartitionProduceResponse{
				Version:    resp.Version,
				Index:      partition.Index,
				ErrorCode:  int16(utils.ErrNoError),
				BaseOffset: partition.Records.BaseOffset,
				// TODO: this needs to be implemented, see documentation for details
				LogAppendTimeMs: -1,
				LogStartOffset:  0,
				// TODO: Don't forget to handle errors when the protocol is fully implemented
			})
		}

		resp.Responses = append(resp.Responses, topicResponse)
	}

	return protocol.Encode(&resp)
}
