#!/usr/bin/env bash
root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

protoExec=$(which "protoc")
if [ -z $protoExec ]; then
    echo 'Please install protoc!'
    exit 0
fi

name=runtime
server=runtimeserver
client=runtimeclient

protos_dir=$root_dir/$name
pb_dir=$root_dir/$name/$server
client_dir=$root_dir/$name/$client
openapi_dir=$root_dir/$name/openapi

mkdir -p $pb_dir
mkdir -p $client_dir
mkdir -p $openapi_dir

echo "generating code"

echo "generating golang stubs..."
cd $protos_dir

# Generate OpenAPI JSON file
protoc -I $protos_dir --openapiv2_out $openapi_dir \
    --openapiv2_opt logtostderr=true \
    --openapiv2_opt=json_names_for_fields=false \
    $protos_dir/*.proto

# go grpc code
protoc -I $protos_dir \
    --go_out $pb_dir --go_opt paths=source_relative \
    --go-grpc_out $pb_dir --go-grpc_opt paths=source_relative \
    $protos_dir/*.proto

# http gw code
protoc -I $protos_dir --grpc-gateway_out $pb_dir \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    $protos_dir/*.proto

# cp golang client code
cp -R $pb_dir/*.go $client_dir

echo "generating golang code success"

echo "done!!!!"

exit 0
