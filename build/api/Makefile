PROTO_DIR=../../api/samples_storage/protobuf-spec
SWAGGER_DIR=../../api/samples_storage/swagger
OUT_DIR=../../internal/services/samples_storage/api

all:
	mkdir -p ${OUT_DIR}
	mkdir -p ${SWAGGER_DIR}

	protoc -I/usr/local/include \
			-I${GOPATH}/src \
			-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
			--proto_path=${PROTO_DIR} \
			--go_out=plugins=grpc:${OUT_DIR} \
			messages.proto service.proto

	protoc -I/usr/local/include \
			-I${GOPATH}/src \
			-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
			--proto_path=${PROTO_DIR} \
			--grpc-gateway_out=logtostderr=true:${OUT_DIR} \
			messages.proto service.proto

	protoc -I/usr/local/include \
			-I${GOPATH}/src \
			-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
			--proto_path=${PROTO_DIR} \
			--swagger_out=logtostderr=true:${SWAGGER_DIR} \
			messages.proto service.proto

clean:
	rm -rf ${OUT_DIR}
	rm -rf ${SWAGGER_DIR}
