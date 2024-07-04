COMPOSEV2 := $(shell docker compose version 2> /dev/null)

ifdef COMPOSEV2
    COMMAND_DOCKER=docker compose
else
    COMMAND_DOCKER=docker-compose
endif

build:
	$(COMMAND_DOCKER) up -d --build

up: 
	$(COMMAND_DOCKER) up -d

down:
	$(COMMAND_DOCKER) down

test:
	go test ./...