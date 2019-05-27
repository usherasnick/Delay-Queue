package config

type DelayQueueService struct {
	GRPCEndpoint string        `json:"grpc_endpoint"` // grpc服务地址
	RedisService *RedisService `json:"redis_backend"`
	KafkaService *KafkaService `json:"kafka_backend"`
}

type RedisService struct {
	SentinelEndpoints           []string `json:"sentinel_endpoints"`
	SentinelMasterName          string   `json:"sentinel_master_name"`
	SentinelPassword            string   `json:"sentinel_password"`
	RedisDatabase               int      `json:"redis_database"`
	RedisMasterPassword         string   `json:"redis_master_password"`
	RedisPoolMaxIdleConns       int      `json:"redis_pool_max_idle_conns"`        // 连接池最大空闲连接数
	RedisPoolMaxActiveConns     int      `json:"redis_pool_max_active_conns"`      // 连接池最大激活连接数
	RedisPoolIdleConnTimeoutSec int      `json:"redis_pool_idle_conn_timeout_sec"` // 连接池空闲连接最长存活时间
	RedisConnectTimeoutMsec     int      `json:"redis_connect_timeout_msec"`       // 连接超时
	RedisReadTimeoutMsec        int      `json:"redis_read_timeout_msec"`          // 读取超时
	RedisWriteTimeoutMsec       int      `json:"redis_write_timeout_msec"`         // 写入超时
}

type KafkaService struct {
	KafkaBrokers []string `json:"kafka_brokers"`
}
