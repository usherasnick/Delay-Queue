package delayqueue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	conf "github.com/usherasnick/Delay-Queue/internal/config"
	"github.com/usherasnick/Delay-Queue/internal/redis"
)

func TestBucketCURD(t *testing.T) {
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
	_, err := dq.redisCli.ExecCommand("DEL", DefaultBucketName)
	assert.Empty(t, err)

	fakeTaskIds := []string{
		"ff74da2f-20c1-45c4-9570-a01b192f1c9d",
		"3fa03ec8-9d75-4002-914b-e7b79f25c323",
		"f58d644e-a2fc-44ce-851c-7530390cfce",
	}
	for _, taskId := range fakeTaskIds {
		err = dq.pushToBucket(DefaultBucketName, time.Now().Unix(), taskId)
		assert.Empty(t, err)
		time.Sleep(time.Second)
	}

	next := 0
	for {
		item, err := dq.getOneFromBucket(DefaultBucketName)
		assert.Empty(t, err)
		if item == nil {
			break
		}
		assert.Equal(t, fakeTaskIds[next], item.TaskId)

		err = dq.delFromBucket(DefaultBucketName, item.TaskId)
		assert.Empty(t, err)
		next++
	}
	assert.Equal(t, 3, next)

	// 退出之前, 再次清理环境
	_, err = dq.redisCli.ExecCommand("DEL", DefaultBucketName)
	assert.Empty(t, err)

	redis.ReleaseInstance()
}
