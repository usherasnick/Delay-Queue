#!/bin/bash

# ****************************** #
# generate grpc apis for golang 
#
# [third-party lib dependencies]
# * libprotoc 3.12.3
#     download protoc-3.12.3-linux-x86_64.zip
# * protoc-gen-go v1.25.0 
#     go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0
#     go install google.golang.org/protobuf/cmd/protoc-gen-go
# * protoc-gen-go-grpc v1.0.1
#     go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.0.1
#     go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
# * protoc-gen-grpc-gateway v2.1.0
# * protoc-gen-openapiv2 v2.1.0
#     git clone https://github.com/grpc-ecosystem/grpc-gateway.git
#     cd grpc-gateway && git checkout v2.1.0 && git switch -c v2.1.0
#     go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
#     go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
#
# More info: https://blog.golang.org/protobuf-apiv2
# ****************************** #

if [ -z ${GOMODULEPATH} ]; then
	echo "no GOMODULEPATH provided!!!"
	exit 1
fi

cd ${GOMODULEPATH}
for i in $(ls ${GOMODULEPATH}/github.com/usherasnick/Delay-Queue/pb/*.proto); do
	fn=github.com/usherasnick/Delay-Queue/pb/$(basename "$i")
	echo "compile" $fn
	# generate the messages
	protoc -I/usr/local/include -I . --go_out=. "$fn"
	# generate the services
	protoc -I/usr/local/include -I . --go-grpc_out=. "$fn"
done