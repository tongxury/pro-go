package rediz

type Config struct {
	// host:port address.
	Addrs    []string
	Type     string // cluster or single
	Password string
	DB       int
}
