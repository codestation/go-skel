package model

type HealthCheckResult struct {
	Ping error `json:"ping"`
}

func (h HealthCheckResult) AllOk() bool {
	if h.Ping != nil {
		return false
	}
	return true
}
