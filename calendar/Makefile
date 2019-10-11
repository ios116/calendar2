mockgen:
	mockgen -destination=./internal/calendar/mock_repo.go  -package=calendar -source=./internal/domain/claendar_domain.go EventRepository

gwgen:
	protoc -I/usr/local/include -I.   -I$(GOPATH)/src   -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis   --grpc-gateway_out=logtostderr=true:.   api/api.proto

protogen:
	protoc -I/usr/local/include -I.   -I$(GOPATH)/src   -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis   --go_out=plugins=grpc:. api/api.proto

rpc_impl:
	impl 's *RPCServer' github.com/ios116/calendar/internal/grpcserver.CalendarServer >> internal/grpcserver/impl.go

