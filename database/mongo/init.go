package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Config struct {
	Host            string
	Port            string
	DbName          string
	Username        string
	Password        string
	PasswordSet     bool
	MaxPoolSize     int64
	MinPoolSize     int64
	ConnectTimeout  time.Duration
	MaxConnIdleTime time.Duration
}

func Init(conf *Config) *mongo.Client {
	if conf.Username != "" || conf.Password != "" {
		conf.PasswordSet = true
	}

	var mgOptions = new(options.ClientOptions)
	mgOptions = mgOptions.SetHosts([]string{conf.Host + ":" + conf.Port})
	mgOptions = mgOptions.SetConnectTimeout(conf.ConnectTimeout)
	mgOptions = mgOptions.SetMaxConnIdleTime(conf.MaxConnIdleTime)
	mgOptions = mgOptions.SetMaxPoolSize(uint64(conf.MaxPoolSize))
	mgOptions = mgOptions.SetMinPoolSize(uint64(conf.MinPoolSize))

	if conf.PasswordSet {
		mgOptions = mgOptions.SetAuth(options.Credential{
			Username:    conf.Username,
			Password:    conf.Password,
			PasswordSet: conf.PasswordSet,
		})
	}

	ctx, c := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, mgOptions)
	c()
	if err != nil {
		panic(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	return client
}
