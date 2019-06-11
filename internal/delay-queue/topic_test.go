package delayqueue

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	conf "github.com/usherasnick/Delay-Queue/internal/config"
	"github.com/usherasnick/Delay-Queue/internal/redis"
)

func TestTopicCURD(t *testing.T) {
	fakeRedisCfg := &conf.RedisService{
		SentinelEndpoints:       []string{"localhost:26379", "localhost:26380", "localhost:26381"},
		SentinelMasterName:      "mymaster",
		SentinelPassword:        "Pwd123!@",
		RedisMasterPassword:     "sOmE_sEcUrE_pAsS",
		RedisPoolMaxIdleConns:   3,
		RedisPoolMaxActiveConns: 64,
		RedisConnectTimeoutMsec: 500,
		RedisReadTimeoutMsec:    500,
		RedisWriteTimeoutMsec:   500,
	}

	dq := &DelayQueue{
		redisCli: redis.GetOrCreateInstance(fakeRedisCfg),
	}

	// 先清理环境
	_, err := dq.redisCli.ExecCommand("DEL", DefaultTopicSetName)
	assert.Empty(t, err)

	fakeTopis := []string{
		"shopping_cart_service_line",
		"order_service_line",
		"inventory_service_line",
	}
	for _, topic := range fakeTopis {
		err := dq.putTopic(DefaultTopicSetName, topic)
		assert.Empty(t, err)
		has, err := dq.hasTopic(DefaultTopicSetName, topic)
		assert.Empty(t, err)
		assert.Equal(t, true, has)
	}

	allTopics, err := dq.listTopic(DefaultTopicSetName)
	assert.Empty(t, err)
	assert.Equal(t, len(fakeTopis), len(allTopics))
	sort.Strings(fakeTopis)
	sort.Strings(allTopics)
	for idx := range fakeTopis {
		assert.Equal(t, fakeTopis[idx], allTopics[idx])
	}

	for _, topic := range fakeTopis {
		err := dq.delTopic(DefaultTopicSetName, topic)
		assert.Empty(t, err)
		has, err := dq.hasTopic(DefaultTopicSetName, topic)
		assert.Empty(t, err)
		assert.Equal(t, false, has)
	}

	// 退出之前, 再次清理环境
	_, err = dq.redisCli.ExecCommand("DEL", DefaultTopicSetName)
	assert.Empty(t, err)

	redis.ReleaseInstance()
}
