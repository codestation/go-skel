package config

// Config holds all the program arguments and environment variables passed to the application.
// The variables set will depend on the command used.
type Config struct {
	Debug     bool
	Addr      string `mapstructure:"listen"`
	DSN       string
	MasterKey []byte `mapstructure:"master-key"`
	JWTSecret []byte `mapstructure:"jwt-secret"`
}
