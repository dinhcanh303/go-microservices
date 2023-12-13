package publisher

type Option func(*publisher)

func ExChangeName(exchangeName string) Option {
	return func(p *publisher) {
		p.exchangeName = exchangeName
	}
}

func BindingKey(bindingKey string) Option {
	return func(p *publisher) {
		p.bindingKey = bindingKey
	}
}

func MessageTypeName(messageTypeName string) Option {
	return func(p *publisher) {
		p.messageTypeName = messageTypeName
	}
}
