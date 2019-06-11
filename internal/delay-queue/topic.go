package delayqueue

/*
	key -> share the same default key "delay_queue_topic_set"

	using SET
*/

func (dq *DelayQueue) putTopic(key, topic string) error {
	_, err := dq.redisCli.ExecCommand("SADD", key, topic)
	return err
}

func (dq *DelayQueue) listTopic(key string) ([]string, error) {
	v, err := dq.redisCli.ExecCommand("SMEMBERS", key)
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
	topics := make([]string, len(vv))
	for i := 0; i < len(vv); i++ {
		topics[i] = string(vv[i].([]byte))
	}
	return topics, nil
}

func (dq *DelayQueue) hasTopic(key, topic string) (bool, error) {
	v, err := dq.redisCli.ExecCommand("SISMEMBER", key, topic)
	if err != nil {
		return false, err
	}
	if v == nil {
		return false, nil
	}
	return v.(int64) == 1, nil
}

func (dq *DelayQueue) delTopic(key, topic string) error {
	_, err := dq.redisCli.ExecCommand("SREM", key, topic)
	return err
}
