# e-resep-be
Receipt Medicine Delivery Services

## Stack / Technologies
- **Golang** _(using echo framework)_
- **PostgreSQL**
- **Docker & Docker Compose** _(for deployment)_

## Prerequisites
1. must have installed golang version 1.20 or higher  
2. installed docker, docker-compose, 
3. [golang migrate](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md)

## SETUP
### A. Local Mode _(must installed postgres version 14 or higher)_
  1. first, copy `.env.example -> .env`, then adjust database configuration for local
  2. next download all module by run command :
  ```
  go mod download
  ```
  3. execute this command to run the API :
  ```
  go run . local
  ``` 
### B. Using Docker Compose
  1. first, copy `.env.example -> .env`, then adjust database configuration for local
  2. execute this command to run all services docker composes:
  ```
  docker compose up -d --build --force-recreate
  ```