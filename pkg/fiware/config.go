package fiware

// Config provides configurtaion for the request to fiware components
type Config struct {
	fiwareService     string
	fiwareServicePath string
}

// ConfigOpts options for fiware configuration
type ConfigOpts func(*Config)

// NewConfig returns new fiware configuration
func NewConfig(opts ...ConfigOpts) *Config {
	c := &Config{
		fiwareService:     "",
		fiwareServicePath: "/",
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithService sets the fiware-service header for any request made with this configuration
func WithService(service string) ConfigOpts {
	return func(c *Config) {
		c.fiwareService = service
	}
}

// WithServicePath sets the fiware-servicePath header for any request made with this configuration
func WithServicePath(servicePath string) ConfigOpts {
	return func(c *Config) {
		c.fiwareServicePath = servicePath
	}
}

// GetHeader returns the fiware-service and fiware-servicepath as a map[string]string which can be used for header parameter for rusty
func (c *Config) GetHeader() map[string]string {
	m := make(map[string]string)
	m["fiware-service"] = c.fiwareService
	m["fiware-servicePath"] = c.fiwareServicePath
	return m
}
