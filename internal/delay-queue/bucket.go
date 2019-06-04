package delayqueue

import (
	"strconv"
)

type BucketItem struct {
	TaskId        string
	TaskTimestamp int64
}

/*
	key -> delay_queue_bucket_${idx}

	using ZSET
*/

func (dq *DelayQueue) pushToBucket(key string, timestamp int64, id string) error {
	_, err := dq.redisCli.ExecCommand("ZADD", key, timestamp, id)
	return err
}

func (dq *DelayQueue) getOneFromBucket(key string) (*BucketItem, error) {
	v, err := dq.redisCli.ExecCommand("ZRANGE", key, 0, 0, "WITHSCORES")
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	vv := v.([]interface{})
	if len(vv) == 0 {
		return nil, nil
	}

	item := BucketItem{}
	item.TaskId = string(vv[0].([]byte))
	timestampStr := string(vv[1].([]byte))
	item.TaskTimestamp, _ = strconv.ParseInt(timestampStr, 10, 64)
	return &item, nil
}

func (dq *DelayQueue) delFromBucket(key string, id string) error {
	_, err := dq.redisCli.ExecCommand("ZREM", key, id)
	return err
}
