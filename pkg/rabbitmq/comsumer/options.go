package consumer

type Option func(*consumer)

func ExChangeName(exchangeName string) Option {
	return func(c *consumer) {
		c.exchangeName = exchangeName
	}
}
