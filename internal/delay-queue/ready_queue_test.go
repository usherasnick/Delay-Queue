package delayqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"

	conf "github.com/usherasnick/Delay-Queue/internal/config"
	"github.com/usherasnick/Delay-Queue/internal/redis"
)

func TestReadyQueueCURD(t *testing.T) {
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
	redisCli := redis.GetOrCreateInstance(fakeRedisCfg)
	rq := NewReadyQueue()

	// 先清理环境
	_, err := redisCli.ExecCommand("DEL", DefaultReadyQueueName)
	assert.Empty(t, err)

	fakeTaskIds := []string{
		"ff74da2f-20c1-45c4-9570-a01b192f1c9d",
		"3fa03ec8-9d75-4002-914b-e7b79f25c323",
		"f58d644e-a2fc-44ce-851c-7530390cfce",
	}
	for _, taskId := range fakeTaskIds {
		err = rq.PushToReadyQueue(redisCli, DefaultReadyQueueName, taskId)
		assert.Empty(t, err)
	}

	taskId, err := rq.BlockPopFromReadyQueue(redisCli, DefaultReadyQueueName, 1)
	assert.Empty(t, err)
	assert.Equal(t, fakeTaskIds[0], taskId)
	taskId, err = rq.BlockPopFromReadyQueue(redisCli, DefaultReadyQueueName, 1)
	assert.Empty(t, err)
	assert.Equal(t, fakeTaskIds[1], taskId)
	taskId, err = rq.BlockPopFromReadyQueue(redisCli, DefaultReadyQueueName, 1)
	assert.Empty(t, err)
	assert.Equal(t, fakeTaskIds[2], taskId)

	// 退出之前, 再次清理环境
	_, err = redisCli.ExecCommand("DEL", DefaultReadyQueueName)
	assert.Empty(t, err)

	redis.ReleaseInstance()
}
