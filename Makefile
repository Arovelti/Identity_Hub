protoc-profile-gen: 
	protoc -I . \
		--go_out ./api --go_opt paths=source_relative \
		--go-grpc_out ./api --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./api \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		proto/profile.proto