package redis

type Config struct {
	Addr    []string
	DB      int
	User    string
	Pass    string
	KeyPath string
	CrtPath string
}
