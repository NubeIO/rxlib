#!/usr/bin/env bash
set -e

root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

protoExec=$(which "protoc")
if [ -z "$protoExec" ]; then
    echo 'Please install protoc!'
    exit 0
fi

name=runtimebase
combined_dir=runtime

protos_dir="$root_dir/$name"
combined_dir="$root_dir/$name/$combined_dir"
openapi_dir="$root_dir/$name/openapi"

mkdir -p "$combined_dir"
mkdir -p "$openapi_dir"

echo "generating code"

echo "generating golang stubs..."
cd "$protos_dir"

# List .proto files to check if they exist
echo "Listing .proto files in $protos_dir"
ls -la "$protos_dir"/*.proto

# Ensure the proto include path for google/api/annotations.proto
PROTOC_INCLUDE="-I$protos_dir -I/usr/local/include -I$root_dir/third_party"

# Generate OpenAPI JSON file
echo "protoc $PROTOC_INCLUDE --experimental_allow_proto3_optional --openapiv2_out $openapi_dir --openapiv2_opt logtostderr=true --openapiv2_opt=json_names_for_fields=false $protos_dir/*.proto"
protoc $PROTOC_INCLUDE --experimental_allow_proto3_optional --openapiv2_out "$openapi_dir" \
    --openapiv2_opt logtostderr=true \
    --openapiv2_opt=json_names_for_fields=false \
    "$protos_dir"/*.proto

# go grpc code (both server and client)
echo "protoc $PROTOC_INCLUDE --experimental_allow_proto3_optional --go_out $combined_dir --go_opt paths=source_relative --go-grpc_out $combined_dir --go-grpc_opt paths=source_relative $protos_dir/*.proto"
protoc $PROTOC_INCLUDE --experimental_allow_proto3_optional \
    --go_out "$combined_dir" --go_opt paths=source_relative \
    --go-grpc_out "$combined_dir" --go-grpc_opt paths=source_relative \
    "$protos_dir"/*.proto

# http gw code
echo "protoc $PROTOC_INCLUDE --experimental_allow_proto3_optional --grpc-gateway_out $combined_dir --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative $protos_dir/*.proto"
protoc $PROTOC_INCLUDE --experimental_allow_proto3_optional --grpc-gateway_out "$combined_dir" \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    "$protos_dir"/*.proto

echo "generating success"

echo "done!!!!"

exit 0
