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

web:
	cd web

start-frontend:
	${NPMRUN} dev


run-frontend: web start-frontend


frontend-build:
	cd web; ${NPMRUN} build

restart-scheduler: 
	${DCOMPOSE} restart scheduler

restart-api:
	${DCOMPOSE} restart eventgrid-api

.PHONY: all build up down clean
