package infrastructure

type AlertProducing interface {
	PublicMessage(body []byte) error
}
