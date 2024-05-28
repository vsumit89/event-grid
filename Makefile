.ONESHELL: 

DCOMPOSE := docker-compose

NPMRUN := npm run

build:
	${DCOMPOSE} build


up:
	${DCOMPOSE} up 


down:
	${DCOMPOSE} down


# Clean the Docker environment
clean:
	$(DCOMPOSE) down --rmi all -v


restart-scheduler: 
	${DCOMPOSE} restart scheduler

restart-api:
	${DCOMPOSE} restart eventgrid-api

.PHONY: all build up down clean
