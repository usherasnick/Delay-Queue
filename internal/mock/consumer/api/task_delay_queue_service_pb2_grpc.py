# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import task_delay_queue_service_pb2 as task__delay__queue__service__pb2


class TaskDelayQueueServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.PushTask = channel.unary_unary(
                '/usherasnick.photon_dance_delay_queue.TaskDelayQueueService/PushTask',
                request_serializer=task__delay__queue__service__pb2.PushTaskRequest.SerializeToString,
                response_deserializer=task__delay__queue__service__pb2.PushTaskResponse.FromString,
                )
        self.FinishTask = channel.unary_unary(
                '/usherasnick.photon_dance_delay_queue.TaskDelayQueueService/FinishTask',
                request_serializer=task__delay__queue__service__pb2.FinishTaskRequest.SerializeToString,
                response_deserializer=task__delay__queue__service__pb2.FinishTaskResponse.FromString,
                )
        self.CheckTask = channel.unary_unary(
                '/usherasnick.photon_dance_delay_queue.TaskDelayQueueService/CheckTask',
                request_serializer=task__delay__queue__service__pb2.CheckTaskRequest.SerializeToString,
                response_deserializer=task__delay__queue__service__pb2.CheckTaskResponse.FromString,
                )
        self.SubscribeTopic = channel.unary_unary(
                '/usherasnick.photon_dance_delay_queue.TaskDelayQueueService/SubscribeTopic',
                request_serializer=task__delay__queue__service__pb2.SubscribeTopicRequest.SerializeToString,
                response_deserializer=task__delay__queue__service__pb2.SubscribeTopicResponse.FromString,
                )
        self.UnsubscribeTopic = channel.unary_unary(
                '/usherasnick.photon_dance_delay_queue.TaskDelayQueueService/UnsubscribeTopic',
                request_serializer=task__delay__queue__service__pb2.UnsubscribeTopicRequest.SerializeToString,
                response_deserializer=task__delay__queue__service__pb2.UnsubscribeTopicResponse.FromString,
                )


class TaskDelayQueueServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def PushTask(self, request, context):
        """for task producers 
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def FinishTask(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def CheckTask(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SubscribeTopic(self, request, context):
        """for task comsumers 
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def UnsubscribeTopic(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_TaskDelayQueueServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'PushTask': grpc.unary_unary_rpc_method_handler(
                    servicer.PushTask,
                    request_deserializer=task__delay__queue__service__pb2.PushTaskRequest.FromString,
                    response_serializer=task__delay__queue__service__pb2.PushTaskResponse.SerializeToString,
            ),
            'FinishTask': grpc.unary_unary_rpc_method_handler(
                    servicer.FinishTask,
                    request_deserializer=task__delay__queue__service__pb2.FinishTaskRequest.FromString,
                    response_serializer=task__delay__queue__service__pb2.FinishTaskResponse.SerializeToString,
            ),
            'CheckTask': grpc.unary_unary_rpc_method_handler(
                    servicer.CheckTask,
                    request_deserializer=task__delay__queue__service__pb2.CheckTaskRequest.FromString,
                    response_serializer=task__delay__queue__service__pb2.CheckTaskResponse.SerializeToString,
            ),
            'SubscribeTopic': grpc.unary_unary_rpc_method_handler(
                    servicer.SubscribeTopic,
                    request_deserializer=task__delay__queue__service__pb2.SubscribeTopicRequest.FromString,
                    response_serializer=task__delay__queue__service__pb2.SubscribeTopicResponse.SerializeToString,
            ),
            'UnsubscribeTopic': grpc.unary_unary_rpc_method_handler(
                    servicer.UnsubscribeTopic,
                    request_deserializer=task__delay__queue__service__pb2.UnsubscribeTopicRequest.FromString,
                    response_serializer=task__delay__queue__service__pb2.UnsubscribeTopicResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'usherasnick.photon_dance_delay_queue.TaskDelayQueueService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class TaskDelayQueueService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def PushTask(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/usherasnick.photon_dance_delay_queue.TaskDelayQueueService/PushTask',
            task__delay__queue__service__pb2.PushTaskRequest.SerializeToString,
            task__delay__queue__service__pb2.PushTaskResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def FinishTask(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/usherasnick.photon_dance_delay_queue.TaskDelayQueueService/FinishTask',
            task__delay__queue__service__pb2.FinishTaskRequest.SerializeToString,
            task__delay__queue__service__pb2.FinishTaskResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def CheckTask(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/usherasnick.photon_dance_delay_queue.TaskDelayQueueService/CheckTask',
            task__delay__queue__service__pb2.CheckTaskRequest.SerializeToString,
            task__delay__queue__service__pb2.CheckTaskResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def SubscribeTopic(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/usherasnick.photon_dance_delay_queue.TaskDelayQueueService/SubscribeTopic',
            task__delay__queue__service__pb2.SubscribeTopicRequest.SerializeToString,
            task__delay__queue__service__pb2.SubscribeTopicResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def UnsubscribeTopic(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/usherasnick.photon_dance_delay_queue.TaskDelayQueueService/UnsubscribeTopic',
            task__delay__queue__service__pb2.UnsubscribeTopicRequest.SerializeToString,
            task__delay__queue__service__pb2.UnsubscribeTopicResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)