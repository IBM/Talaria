package api

import (
	"opentalaria/protocol"
)

type APIVersionsAPI struct {
	Request Request
}

func (a APIVersionsAPI) Name() string {
	return "API Versions"
}

func (a APIVersionsAPI) GetRequest() Request {
	return a.Request
}

func (a APIVersionsAPI) GeneratePayload() ([]byte, error) {
	// handle response
	apiVersionRequest := protocol.ApiVersionsRequest{}
	_, err := protocol.VersionedDecode(a.Request.Message, &apiVersionRequest, a.Request.Header.RequestApiVersion)
	if err != nil {
		return nil, err
	}

	response := NewAPIVersionsResponse(a.GetRequest().Header.RequestApiVersion)
	return protocol.Encode(response)
}

func (a APIVersionsAPI) GetHeaderVersion(requestVersion int16) int16 {
	return (&protocol.ApiVersionsResponse{Version: requestVersion}).GetHeaderVersion()
}

func getAPIVersions() []protocol.ApiVersion {
	return []protocol.ApiVersion{
		{ApiKey: (&protocol.ApiVersionsRequest{}).GetKey(), MinVersion: 0, MaxVersion: 3},
		{ApiKey: (&protocol.MetadataRequest{}).GetKey(), MinVersion: 0, MaxVersion: 8},
		{ApiKey: (&protocol.ProduceRequest{}).GetKey(), MinVersion: 0, MaxVersion: 8},
		{ApiKey: (&protocol.HeartbeatRequest{}).GetKey(), MinVersion: 0, MaxVersion: 1},
		// {APIKey: FetchKey, MinVersion: 0, MaxVersion: 3},
		// {APIKey: OffsetsKey, MinVersion: 0, MaxVersion: 2},
		// {APIKey: LeaderAndISRKey, MinVersion: 0, MaxVersion: 1},
		// {APIKey: StopReplicaKey, MinVersion: 0, MaxVersion: 0},
		// {APIKey: FindCoordinatorKey, MinVersion: 0, MaxVersion: 1},
		// {APIKey: JoinGroupKey, MinVersion: 0, MaxVersion: 1},
		// {APIKey: HeartbeatKey, MinVersion: 0, MaxVersion: 1},
		// {APIKey: LeaveGroupKey, MinVersion: 0, MaxVersion: 1},
		// {APIKey: SyncGroupKey, MinVersion: 0, MaxVersion: 1},
		// {APIKey: DescribeGroupsKey, MinVersion: 0, MaxVersion: 1},
		// {APIKey: ListGroupsKey, MinVersion: 0, MaxVersion: 1},
		// {APIKey: CreateTopicsKey, MinVersion: 0, MaxVersion: 1},
		// {APIKey: DeleteTopicsKey, MinVersion: 0, MaxVersion: 1},
	}
}

func NewAPIVersionsResponse(version int16) *protocol.ApiVersionsResponse {
	return &protocol.ApiVersionsResponse{
		Version:        version,
		ErrorCode:      0,
		ApiKeys:        getAPIVersions(),
		ThrottleTimeMs: 0,
	}
}
