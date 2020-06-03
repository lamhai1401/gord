package redis

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Config redis
type Config struct {
	Addrs []string
	Pwd   string
	DB    int
}

// Client redis
type Client struct {
	cluster     *redis.ClusterClient
	single      *redis.Client
	clusterNode bool
	mutex       sync.Mutex
}

// NewClient connect redis client
func NewClient(c *Config) *Client {
	if c.Addrs == nil || len(c.Addrs) == 0 {
		return nil
	}

	client := &Client{}
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	if len(c.Addrs) == 1 {
		client.single = redis.NewClient(&redis.Options{
			Addr:         c.Addrs[0],
			Password:     c.Pwd,
			DB:           c.DB,
			DialTimeout:  3 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		})

		if err := client.single.Ping(ctx).Err(); err != nil {
			log.Println(err.Error())
			return nil
		}

		client.single.Do(ctx, "CONFIG", "SET", "notify-keyspace-events", "AKE")
		client.clusterNode = false
		return client
	}

	client.cluster = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        c.Addrs,
		Password:     c.Pwd,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	if err := client.cluster.Ping(ctx).Err(); err != nil {
		log.Println(err.Error())
		return nil
	}

	client.cluster.Do(ctx, "CONFIG", "SET", "notify-keyspace-events", "AKE")
	client.clusterNode = true
	return client
}

// Set linter
func (c *Client) Set(k, v string, t time.Duration) error {
	ctx, cancel := context.WithTimeout(context.TODO(), t)
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
		cancel()

	}()
	if c.clusterNode {
		return c.cluster.Set(ctx, k, v, t).Err()
	}

	return c.single.Set(ctx, k, v, t).Err()
}

// Get linter
func (c *Client) Get(k string) string {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
		cancel()

	}()

	if c.clusterNode {
		return c.cluster.Get(ctx, k).Val()
	}

	return c.single.Get(ctx, k).Val()
}
