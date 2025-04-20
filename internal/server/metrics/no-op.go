package metrics

type NoOpClient struct{}

func NewNoOpClient() NoOpClient {
	return NoOpClient{}
}

func (NoOpClient) Publish(code int, ip, page, lang string) error { return nil }
