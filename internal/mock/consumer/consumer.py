# -*- coding: utf-8 -*-
import grpc
import kafka
import pprint
import os
import random
import sys
import threading
import time
sys.path.append(os.path.abspath('./api'))

from api import task_delay_queue_service_pb2
from api import task_delay_queue_service_pb2_grpc

'''
ev_epollex_linux.cc:516 Error shutting down fd 43. errno: 9

/* Might be called multiple times */
static void fd_shutdown(grpc_fd* fd, grpc_error* why) {
  if (fd->read_closure.SetShutdown(GRPC_ERROR_REF(why))) {
    if (shutdown(fd->fd, SHUT_RDWR)) {
      if (errno != ENOTCONN) {
        gpr_log(GPR_ERROR, "Error shutting down fd %d. errno: %d",
                grpc_fd_wrapped_fd(fd), errno);
      }
    }
    fd->write_closure.SetShutdown(GRPC_ERROR_REF(why));
    fd->error_closure.SetShutdown(GRPC_ERROR_REF(why));
  }
  GRPC_ERROR_UNREF(why);
}
'''

class GrpcThread(threading.Thread):
    def __init__(self, grpc_channel, task_id):
        threading.Thread.__init__(self)
        self.channel = grpc_channel
        self.task_id = task_id
    
    def run(self):
        time.sleep(random.randrange(1, 11))
        stub = task_delay_queue_service_pb2_grpc.TaskDelayQueueServiceStub(self.channel)
        request = task_delay_queue_service_pb2.FinishTaskRequest(task_id=self.task_id)
        try:
            response = stub.FinishTask(request)
            print(response)
        except grpc.RpcError as e:
            e.details()
            status_code = e.code()
            print("status_code <{}:{}>".format(status_code.name, status_code.value))


if __name__ == "__main__":
    consumer = kafka.KafkaConsumer(
        "inventory_service_line",
        bootstrap_servers=["localhost:9092"],
        auto_offset_reset="earliest",
        group_id="mock_consumer_group_03",
        enable_auto_commit=True,
    )
    with grpc.insecure_channel('localhost:18081', options=(('grpc.so_reuseport', 1),)) as channel:
        for msg in consumer:
            print("[%s:%d:%d]" % (msg.topic, msg.partition, msg.offset))
            value = msg.value
            task = task_delay_queue_service_pb2.Task()
            task.ParseFromString(value)
            pprint.pprint(task)
            t = GrpcThread(channel, task.id)
            t.start()
