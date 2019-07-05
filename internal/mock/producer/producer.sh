#!/bin/bash

echo "/**************************************************/"
echo "list all grpc services"
grpcurl -plaintext localhost:18081 list usherasnick.photon_dance_delay_queue.TaskDelayQueueService
echo "/**************************************************/"
echo "subscribe topic shopping_cart_service_line"
grpcurl -plaintext -d '{"topic": "shopping_cart_service_line"}' localhost:18081 usherasnick.photon_dance_delay_queue.TaskDelayQueueService/SubscribeTopic
echo "/**************************************************/"
echo "subscribe topic order_service_line"
grpcurl -plaintext -d '{"topic": "order_service_line"}' localhost:18081 usherasnick.photon_dance_delay_queue.TaskDelayQueueService/SubscribeTopic
echo "/**************************************************/"
echo "subscribe topic inventory_service_line"
grpcurl -plaintext -d '{"topic": "inventory_service_line"}' localhost:18081 usherasnick.photon_dance_delay_queue.TaskDelayQueueService/SubscribeTopic
echo "/**************************************************/"
echo "add one task"
cat task01.json | grpcurl -plaintext -d @ localhost:18081 usherasnick.photon_dance_delay_queue.TaskDelayQueueService/PushTask
echo "/**************************************************/"
echo "add one task"
cat task02.json | grpcurl -plaintext -d @ localhost:18081 usherasnick.photon_dance_delay_queue.TaskDelayQueueService/PushTask
echo "/**************************************************/"
echo "add one task"
cat task03.json | grpcurl -plaintext -d @ localhost:18081 usherasnick.photon_dance_delay_queue.TaskDelayQueueService/PushTask
echo "/**************************************************/"
echo "add one task"
cat task04.json | grpcurl -plaintext -d @ localhost:18081 usherasnick.photon_dance_delay_queue.TaskDelayQueueService/PushTask
echo "/**************************************************/"
echo "add one task"
cat task05.json | grpcurl -plaintext -d @ localhost:18081 usherasnick.photon_dance_delay_queue.TaskDelayQueueService/PushTask
echo "/**************************************************/"
echo "add one task"
cat task06.json | grpcurl -plaintext -d @ localhost:18081 usherasnick.photon_dance_delay_queue.TaskDelayQueueService/PushTask
echo "/**************************************************/"