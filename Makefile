# parameters
REMOTE_HOST=root@10.0.0.249
TARGET_DIR=/opt/monitor-etl
IMAGE_TAG=ghcr.io/ipsums9527/monitor-etl:dev

# Targets
clean:
	docker image prune -f

build: clean
	docker build -f Dockerfile -t $(IMAGE_TAG) .

update:
	ssh $(REMOTE_HOST) "docker-compose -f $(TARGET_DIR)/docker-compose.yml up -d"
	ssh $(REMOTE_HOST) "docker image prune -f"

all: build update
