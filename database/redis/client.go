package redis

import (
	"context"
	"crypto/tls"
	"github.com/go-redis/redis/v8"
)

func Init(c *Config) (redis.UniversalClient, error) {

	var tlsConf *tls.Config
	var err error

	if c.CrtPath != "" && c.KeyPath != "" {
		tlsConf = &tls.Config{}
		tlsConf.Certificates = make([]tls.Certificate, 1)
		tlsConf.Certificates[0], err = tls.LoadX509KeyPair(c.CrtPath, c.KeyPath)
		if err != nil {
			return nil, err
		}
	}
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:      c.Addr,
		DB:         c.DB,
		Username:   c.User,
		Password:   c.Pass,
		MaxRetries: 3,
		TLSConfig:  tlsConf,
	})
	err = client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return client, nil
}
