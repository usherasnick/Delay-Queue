package redis

import (
	"errors"
	"sync"
	"time"

	"github.com/FZambia/sentinel"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"

	conf "github.com/usherasnick/Delay-Queue/internal/config"
)

type RedisConnPoolSingleton struct {
	db int
	p  *redis.Pool
}

var (
	instance  *RedisConnPoolSingleton
	onceOpen  sync.Once
	onceClose sync.Once
)

func GetOrCreateInstance(cfg *conf.RedisService) *RedisConnPoolSingleton {
	onceOpen.Do(func() {
		sntnl := &sentinel.Sentinel{
			Addrs:      cfg.SentinelEndpoints,
			MasterName: cfg.SentinelMasterName,
			Dial: func(addr string) (redis.Conn, error) {
				conn, err := redis.Dial(
					"tcp",
					addr,
					redis.DialConnectTimeout(time.Duration(cfg.RedisConnectTimeoutMsec)*time.Millisecond),
					redis.DialReadTimeout(time.Duration(cfg.RedisReadTimeoutMsec)*time.Millisecond),
					redis.DialWriteTimeout(time.Duration(cfg.RedisWriteTimeoutMsec)*time.Millisecond),
					redis.DialPassword(cfg.SentinelPassword),
				)
				if err != nil {
					log.Fatal().Err(err).Msg("failed to connect to redis server")
					return nil, err
				}
				return conn, nil
			},
		}
		instance = &RedisConnPoolSingleton{}
		instance.db = cfg.RedisDatabase
		instance.p = &redis.Pool{
			Dial: func() (redis.Conn, error) {
				master, err := sntnl.MasterAddr()
				if err != nil {
					return nil, err
				}
				conn, err := redis.Dial(
					"tcp",
					master,
					redis.DialPassword(cfg.RedisMasterPassword),
				)
				if err != nil {
					log.Fatal().Err(err).Msg("failed to connect to redis server")
					return nil, err
				}
				return conn, nil
			},
			TestOnBorrow: func(conn redis.Conn, t time.Time) error {
				if !sentinel.TestRole(conn, "master") {
					return errors.New("Role check failed")
				} else {
					return nil
				}
			},
			MaxIdle:     cfg.RedisPoolMaxIdleConns,
			MaxActive:   cfg.RedisPoolMaxActiveConns,
			IdleTimeout: time.Duration(cfg.RedisPoolIdleConnTimeoutSec) * time.Second,
			Wait:        true,
		}
	})
	return instance
}

func ReleaseInstance() {
	onceClose.Do(func() {
		if instance != nil {
			instance.p.Close()
		}
	})
}

// ExecCommand 执行redis命令, 完成后自动归还连接.
func (p *RedisConnPoolSingleton) ExecCommand(cmd string, args ...interface{}) (interface{}, error) {
	conn := p.getConn()
	defer conn.Close()
	return conn.Do(cmd, args...)
}

// ExecLuaScript 执行lua脚本, 完成后自动归还连接.
func (p *RedisConnPoolSingleton) ExecLuaScript(src string, keyCount int, keysAndArgs ...interface{}) (interface{}, error) {
	conn := p.getConn()
	defer conn.Close()
	luaScript := redis.NewScript(keyCount, src)
	return luaScript.Do(conn, keysAndArgs...)
}

func (p *RedisConnPoolSingleton) getConn() redis.Conn {
	conn := p.p.Get()
	conn.Do("SELECT", p.db) // nolint
	return conn
}
