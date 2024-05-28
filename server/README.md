# Eventgrid API

## Folder structure
- cmd 
    - api: service which exposes api to manage events and publishes events to scheduler
    - scheduler: service which consumes calendar events and schedules reminder emails using a min-heap
    - messenger: service which sends reminder emails

- docker: dockerfiles for the services

- init-scripts: init.sql file for creating db
- internal
    - commons: has all variables and functions used across application
    - config: has functions related to config
    - handlers: all the http handlers, dtos and middleware functions
    - infrastructure:
        - database: for database connection and utility operations on database
        - email: sends emails
        - mq: message queue for event publishing and consuming
        - transport: has functions for running http server 

    - mocks: contains mocks which can be used for writing tests
    - models: contains db models used for db running operations
    - repository: contains repository interfaces and implementations for accessing data from the database.
    - services: Contains the core business logic and operations, utilizing dependencies like repositories, email, and MQ.
    - workers: contains functions for scheduling events and sending emails


## Architecture Diagram

### Main Architecture 
![Architecture Diagram](./../images/overall-architecture.png)

### Flow Chart for Sending notifications

![Flow Chart for notifications](./../images/notification-worker.png)
