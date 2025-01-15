PROTO_NAMES=(
    "admin"
    "chat"
    "common"
)

for name in "${PROTO_NAMES[@]}"; do
 protoc --go_out=./${name} --go_opt=module=github.com/openimsdk/chat/pkg/protocol/${name} ${name}/${name}.proto
  if [ $? -ne 0 ]; then
      echo "error processing ${name}.proto (go_out)"
      exit $?
  fi
done

# generate go-grpc

for name in "${PROTO_NAMES[@]}"; do
 protoc --go-grpc_out=./${name} --go-grpc_opt=module=github.com/openimsdk/chat/pkg/protocol/${name} ${name}/${name}.proto
  if [ $? -ne 0 ]; then
      echo "error processing ${name}.proto (go-grpc_out)"
      exit $?
  fi
done

if [ "$(uname -s)" == "Darwin" ]; then
    find . -type f -name '*.pb.go' -exec sed -i '' 's/,omitempty"`/\"\`/g' {} +
else
    find . -type f -name '*.pb.go' -exec sed -i 's/,omitempty"`/\"\`/g' {} +
fi