package api

type HeartbeatAPI struct {
	Request Request
}

func (h HeartbeatAPI) Name() string {
	return "Heartbeat"
}
