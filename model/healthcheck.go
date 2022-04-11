package model

type HealthCheckResult struct {
	Ping error `json:"ping"`
}

func (h HealthCheckResult) AllOk() bool {
	return h.Ping == nil
}
