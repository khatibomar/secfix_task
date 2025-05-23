SOCKET_FILE := /tmp/osquery.socket.path
SOCKET_PATH := $(if $(wildcard $(SOCKET_FILE)),$(shell cat $(SOCKET_FILE)),/tmp/osquery.$(shell whoami).$(shell bash -c 'echo $$RANDOM').em)
OSQUERYD_BIN := /opt/osquery/lib/osquery.app/Contents/MacOS/osqueryd
CONFIG_PATH := /var/osquery/osquery.conf
LOG_DIR := /tmp/osquery_logs_$(shell whoami)
PID_FILE := /tmp/osquery_temp.pid
BUILD_DIR := ./bin
APP_NAME := osquery
API_NAME := api
UI_NAME := ui

.PHONY: deamon-run deamon-stop deamon-status deamon-setup deamon-cleanup gen build-app app build-api api docker-up docker-down db-up ui build-ui

# Deamon
deamon-setup:
	sudo cp /var/osquery/osquery.example.conf /var/osquery/osquery.conf
	sudo cp /var/osquery/io.osquery.agent.plist /Library/LaunchDaemons
	sudo launchctl load /Library/LaunchDaemons/io.osquery.agent.plist

deamon-run: deamon-setup
	@mkdir -p $(LOG_DIR)
	@echo "$(SOCKET_PATH)" > $(SOCKET_FILE)
	@echo "Starting temporary osqueryd..."
	@echo "  - Socket: $(SOCKET_PATH)"
	@echo "  - Logs: $(LOG_DIR)"
	@$(OSQUERYD_BIN) \
		--database_path=/tmp/osquery_temp.db \
		--logger_path="$(LOG_DIR)" \
		--extensions_socket="$(SOCKET_PATH)" \
		--config_path="$(CONFIG_PATH)" \
		--disable_database \
		--ephemeral \
		--disable_audit \
		--disable_events & echo $$! > $(PID_FILE)
	@for i in $$(seq 1 10); do \
		if [ -S "$(SOCKET_PATH)" ]; then \
			echo -e "\n✅ Ready! Use this socket path in your Go app:"; \
			echo "$(SOCKET_PATH)"; \
			exit 0; \
		fi; \
		sleep 1; \
		echo -n "."; \
	done; \
	echo -e "\n❌ Failed to create socket!"; \
	$(MAKE) stop; \
	exit 1

deamon-stop:
	@if [ -f $(PID_FILE) ]; then \
		kill $$(cat $(PID_FILE)) 2>/dev/null || true; \
		rm -f $(PID_FILE); \
	fi
	@rm -f $(SOCKET_PATH)
	@rm -f $(SOCKET_FILE)
	@echo "Cleaned up temporary osqueryd"

deamon-status:
	@if [ -f $(PID_FILE) ] && kill -0 $$(cat $(PID_FILE)) 2>/dev/null; then \
		echo "osqueryd is running with PID $$(cat $(PID_FILE))"; \
		echo "Socket: $(SOCKET_PATH)"; \
	else \
		echo "osqueryd is not running"; \
	fi

# Apps
gen:
	sqlc generate

build-app: gen
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/osquery

build-api: gen 
	@echo "Building api..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(API_NAME) ./cmd/api

build-ui: 
	@echo "Building ui..."
	@mkdir -p $(BUILD_DIR)
	cd cmd/ui && go build -o ../../$(BUILD_DIR)/$(UI_NAME)

ui: build-ui
	$(BUILD_DIR)/$(UI_NAME)

api: build-api
	SECFIX_CONNECTION_STRING="postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable" $(BUILD_DIR)/$(API_NAME) $(ARGS) 

app: build-app
	SECFIX_CONNECTION_STRING="postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable" $(BUILD_DIR)/$(APP_NAME) -socket-path=$(SOCKET_PATH) $(ARGS)

# Docker
docker-up:
	@echo "Building docker image..."
	docker compose up -d 

docker-down:
	@echo "Stopping docker containers..."
	docker compose down

# database
db-up:
	@echo "Running migrations..."
	DATABASE_URL="postgres://postgres@127.0.0.1:5430/postgres?sslmode=disable" dbmate -d ./data/sql/migrations up

curl:
	curl localhost:4000/v1/latest_data
