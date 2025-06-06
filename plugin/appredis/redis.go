package appredis

import (
	"context"
	"flag"
	"fmt"

	"github.com/200Lab-Education/go-sdk/logger"
	"github.com/redis/go-redis/v9"
)

var (
	defaultRedisName      = "DefaultRedis"
	defaultRedisMaxActive = 0 // 0 is unlimited max active connection
	defaultRedisMaxIdle   = 10
)

type RedisDBOpt struct {
	Prefix    string
	RedisUri  string
	MaxActive int
	MaxIdle   int
}

type redisDB struct {
	name   string
	client *redis.Client
	logger logger.Logger
	*RedisDBOpt
}

func NewRedisDB(name, flagPrefix string) *redisDB {
	return &redisDB{
			name: name,
			RedisDBOpt: &RedisDBOpt{
					Prefix:    flagPrefix,
					MaxActive: defaultRedisMaxActive,
					MaxIdle:   defaultRedisMaxIdle, 
			},
	}
}

func (r *redisDB) GetPrefix() string {
	return r.Prefix
}

func (r *redisDB) IsDisabled() bool {
	return r.RedisUri == ""
}

// Rgisters command-line flags that you can use when running the program.
// ./app --redis-main-uri="redis://example.com:6380/1" --redis-main-pool-max-active=20
func (r *redisDB) InitFlags() {
	prefix := r.Prefix
	if r.Prefix != "" {
			prefix += "-"
	}

	flag.StringVar(
			&r.RedisUri, 
			prefix+"-uri", 
			"redis://localhost:6379", 
			"(For go-redis) Redis connection-string. Ex: redis://ocalhost:6379/0",
	)
	flag.IntVar(
			&r.MaxActive, 
			prefix+"pool-max-active", 
			defaultRedisMaxActive, 
			"(For go-redis) Override redis pool MaxActive",
	)
	flag.IntVar(
			&r.MaxIdle, 
			prefix+"pool-max-idle", 
			defaultRedisMaxIdle, 
			"(For go-redis) Override redis pool MaxIdle",
	)
}


/* 
* Establishes a connection to the Redis server 
*/

func (r *redisDB) Configure() error {
	if r.IsDisabled() {
			return nil
	}

	// Initializes the logger specific to this Redis instance 
	r.logger = logger.GetCurrent().GetLogger(r.name)
	r.logger.Info("Connecting to Redis at ", r.RedisUri, "...")

	// Takes the Redis connection string and parse, then return a Redis options 
	// that contain parsed connection details (addr, password, db number,...)
	opt, err := redis.ParseURL(r.RedisUri)
	if err != nil {
			r.logger.Error("Cannot parse Redis URI: ", err.Error())
			return fmt.Errorf("failed to parse Redis URI: %w", err)
	}

	// Sets the pool size and minimum idle connections for the Redis client.
	opt.PoolSize = r.MaxActive
	opt.MinIdleConns = r.MaxIdle 

	// Initializes a new Redis client with the given options.
	client := redis.NewClient(opt)

	// Ping to test Redis connection
	if err := client.Ping(context.Background()).Err(); err != nil {
			r.logger.Error("Cannot connect to Redis:", err.Error())
			return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Connect successfully, assign client to redisDB
	r.client = client
	r.logger.Info("Redis connection established successfully")
	return nil
}

func (r *redisDB) Name() string {
	return r.name
}

func (r *redisDB) Get() interface{} {
	return r.client
}

func (r *redisDB) Run() error {
	return r.Configure()
}

func (r *redisDB) Stop() <-chan bool {
	// stops Redis connection
	if r.client != nil {
			if err := r.client.Close(); err != nil {
					r.logger.Info("cannot close ", r.name, " error:", err)
			}
	}

	c := make(chan bool, 1)
	go func() { 
			c <- true 
	}()
	return c
}