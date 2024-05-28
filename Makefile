.ONESHELL: 

DCOMPOSE := docker-compose

NPMRUN := npm run

# Build the Docker environment
build:
	${DCOMPOSE} build

# Run the Docker environment
up:
	${DCOMPOSE} up 

# Stop the Docker environment
down:
	${DCOMPOSE} down


# Clean the Docker environment
clean:
	$(DCOMPOSE) down --rmi all -v

# Restart the scheduler 
restart-scheduler: 
	${DCOMPOSE} restart scheduler

# Restart the API
restart-api:
	${DCOMPOSE} restart eventgrid-api

.PHONY: all build up down clean
