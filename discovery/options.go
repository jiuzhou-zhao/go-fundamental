package discovery

type Options struct {
	serviceType string
	serviceName string
	rawKey      string
}

func (options Options) String() string {
	if options.rawKey != "" {
		return options.rawKey
	}
	if options.serviceType == "*" && options.serviceName == "*" {
		return "*"
	}
	return options.serviceType + ":" + options.serviceName + ":*"
}

type funcDiscoveryOption struct {
	f func(*Options)
}

func (fdo *funcDiscoveryOption) apply(do *Options) {
	fdo.f(do)
}

type Option interface {
	apply(*Options)
}

func TypeOption(t string) Option {
	return &funcDiscoveryOption{
		f: func(options *Options) {
			if t == "" {
				t = "*"
			}
			options.serviceType = t
		},
	}
}

func NameOption(t string) Option {
	return &funcDiscoveryOption{
		f: func(options *Options) {
			if t == "" {
				t = "*"
			}
			options.serviceName = t
		},
	}
}

func RawKeyOption(t string) Option {
	return &funcDiscoveryOption{
		f: func(options *Options) {
			options.rawKey = t
		},
	}
}
