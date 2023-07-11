protoc-gen-go:
	@echo "Regenerate GO service from .proto"
		@cd langs/go && rm -rf *
		@protoc --go_out=./langs/go --go_opt=paths=source_relative \
        		--go-grpc_out=./langs/go --go-grpc_opt=paths=source_relative \
        		--go-grpc_opt=require_unimplemented_servers=false \
        		service.proto
		@echo "Regenerate GO service completed"
