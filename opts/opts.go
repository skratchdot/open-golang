package opts

// Option configures how a target is opened.
type Option func(builder *ConfigBuilder)

// Config holds the resolved options used to open a target.
type Config struct {
	appName string
	wait    bool
}

// AppName returns the name of the application used to open the target.
func (cfg Config) AppName() string {
	return cfg.appName
}

// Wait reports whether opening the target waits for completion.
func (cfg Config) Wait() bool {
	return cfg.wait
}

// Builder returns a builder initialized with cfg's values.
func (cfg Config) Builder() ConfigBuilder {
	return ConfigBuilder(cfg)
}

// ConfigBuilder incrementally builds a [Config].
type ConfigBuilder struct {
	appName string
	wait    bool
}

// AppName sets the application used to open the target and returns the
// updated builder.
func (builder ConfigBuilder) AppName(name string) ConfigBuilder {
	builder.appName = name
	return builder
}

// Wait sets whether opening the target waits for completion and returns the
// updated builder.
func (builder ConfigBuilder) Wait(wait bool) ConfigBuilder {
	builder.wait = wait
	return builder
}

// Build returns the [Config] represented by builder.
func (builder ConfigBuilder) Build() Config {
	return Config(builder)
}

// WithApp specifies the application that will
// be used to open the value.
func WithApp(appName string) Option {
	return func(builder *ConfigBuilder) {
		*builder = builder.AppName(appName)
	}
}

// Wait will indicate that we need to wait for the
// command to complete.
func Wait() Option {
	return func(builder *ConfigBuilder) {
		*builder = builder.Wait(true)
	}
}

// Init applies opts to builder and returns the resulting [Config].
func Init(builder ConfigBuilder, opts ...Option) Config {
	for _, opt := range opts {
		opt(&builder)
	}
	return builder.Build()
}
