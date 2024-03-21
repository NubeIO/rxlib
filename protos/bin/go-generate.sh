#!/usr/bin/env bash
root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

protoExec=$(which "protoc")
if [ -z $protoExec ]; then
    echo 'Please install protoc!'
    exit 0
fi

name=runtimebase
combined_dir=runtime

protos_dir=$root_dir/$name
combined_dir=$root_dir/$name/$combined_dir
openapi_dir=$root_dir/$name/openapi
rust_out=$root_dir/rust

mkdir -p $combined_dir
mkdir -p $openapi_dir
mkdir -p $rust_out

echo "generating code"

echo "generating golang stubs..."
cd $protos_dir

# Generate OpenAPI JSON file
protoc -I $protos_dir --openapiv2_out $openapi_dir \
    --openapiv2_opt logtostderr=true \
    --openapiv2_opt=json_names_for_fields=false \
    $protos_dir/*.proto

# go grpc code (both server and client)
protoc -I $protos_dir \
    --go_out $combined_dir --go_opt paths=source_relative \
    --go-grpc_out $combined_dir --go-grpc_opt paths=source_relative \
    $protos_dir/*.proto

# http gw code
protoc -I $protos_dir --grpc-gateway_out $combined_dir \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    $protos_dir/*.proto

# rust code
protoc -I $protos_dir --rust_out $rust_out \
    $protos_dir/*.proto

echo "generating golang and rust code success"

echo "done!!!!"

exit 0
