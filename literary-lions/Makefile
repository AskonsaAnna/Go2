


IMAGE_NAME := lions


CONTAINER_PORT := 8080
HOST_PORT := 8080


DEFAULT_VOLUME_PATH := $(CURDIR)/pkg/database


VOLUME_PATH := $(shell [ -d $(DEFAULT_VOLUME_PATH) ] && echo $(DEFAULT_VOLUME_PATH) || echo $(CURDIR)/pkg/database)

# Build the Docker image
build:
	docker build -t $(IMAGE_NAME) .

# Run the Docker container
run:
	docker run -d -p $(HOST_PORT):$(CONTAINER_PORT) -v $(VOLUME_PATH):/app/pkg/database $(IMAGE_NAME)

# Combined target to build and run
all: build run
