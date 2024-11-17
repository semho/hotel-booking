LOCAL_BIN:=$(CURDIR)/bin
PROTO_DIR:=$(CURDIR)/api/proto/v1
DB_PORT ?= 5431
SERVICE ?= auth-service
DB_NAME ?= auth_service
LOCAL_MIGRATION_DIR=$(CURDIR)/$(SERVICE)/deployments/migrations
LOCAL_MIGRATION_DSN="postgres://postgres:postgres@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable"
.PHONY: install-deps
install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.22.1
	go mod tidy

.PHONY: migrate-up
migrate-up:
	@echo "Running migrations UP for service: $(SERVICE)"
	@if [ ! -d "$(LOCAL_MIGRATION_DIR)" ]; then \
		echo "Error: Migration directory not found: $(LOCAL_MIGRATION_DIR)"; \
		exit 1; \
	fi
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) up -v

.PHONY: migrate-down
migrate-down:
	@echo "Running migrations DOWN for service: $(SERVICE)"
	@if [ ! -d "$(LOCAL_MIGRATION_DIR)" ]; then \
		echo "Error: Migration directory not found: $(LOCAL_MIGRATION_DIR)"; \
		exit 1; \
	fi
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) down -v

.PHONY: migrate-status
migrate-status:
	@echo "Checking migration status for service: $(SERVICE)"
	@if [ ! -d "$(LOCAL_MIGRATION_DIR)" ]; then \
		echo "Error: Migration directory not found: $(LOCAL_MIGRATION_DIR)"; \
		exit 1; \
	fi
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) status -v


.PHONY: migrate-create
migrate-create:
	@echo "Creating migration status for service: $(SERVICE)"
	@if [ ! -d "$(LOCAL_MIGRATION_DIR)" ]; then \
		echo "Error: Migration directory not found: $(LOCAL_MIGRATION_DIR)"; \
		exit 1; \
	fi
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) create init sql

.PHONY: vendor-proto
vendor-proto:
	@if [ ! -d vendor.protogen/google ]; then \
		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
		mkdir -p  vendor.protogen/google/ &&\
		mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
		rm -rf vendor.protogen/googleapis ;\
	fi
	@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
		mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
		git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
		mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
		rm -rf vendor.protogen/openapiv2 ;\
	fi

.PHONY: generate-room
generate-room:
	mkdir -p pkg/proto/room_v1
	protoc --proto_path $(PROTO_DIR) --proto_path vendor.protogen \
		--go_out=pkg/proto/room_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=./bin/protoc-gen-go \
		--go-grpc_out=pkg/proto/room_v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc \
		--grpc-gateway_out=pkg/proto/room_v1 --grpc-gateway_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=./bin/protoc-gen-grpc-gateway \
		"$(PROTO_DIR)/room/room.proto"

.PHONY: generate-booking
generate-booking:
	mkdir -p pkg/proto/booking_v1
	protoc --proto_path $(PROTO_DIR) --proto_path vendor.protogen \
		--go_out=pkg/proto/booking_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=./bin/protoc-gen-go \
		--go-grpc_out=pkg/proto/booking_v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc \
		--grpc-gateway_out=pkg/proto/booking_v1 --grpc-gateway_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=./bin/protoc-gen-grpc-gateway \
		"$(PROTO_DIR)/booking/booking.proto"

.PHONY: generate-auth
generate-auth:
	mkdir -p pkg/proto/auth_v1
	protoc --proto_path $(PROTO_DIR) --proto_path vendor.protogen \
		--go_out=pkg/proto/auth_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=./bin/protoc-gen-go \
		--go-grpc_out=pkg/proto/auth_v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc \
		--grpc-gateway_out=pkg/proto/auth_v1 --grpc-gateway_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=./bin/protoc-gen-grpc-gateway \
		"$(PROTO_DIR)/auth/auth.proto"

.PHONY: generate
generate: generate-room generate-booking generate-auth