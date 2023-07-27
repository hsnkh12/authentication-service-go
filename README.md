# authentication-service-go
Authentication service in Golang, built to be used in microservies architecture app that requires authentication.

## Protocols 
- HTTP - Rest API
- gRPC 


## Setup

- Clone the repo
```
git clone https://github.com/hsnkh12/authentication-service-go.git
```
- Setup a mysql container
```
docker run --name CONTAINER_NAME -d -p localhost:HOST_PORT:CONTAINER_PORT -e ROOT_PASSWORD=YOUR_PASSWORD mysql
```
- Login into the container and create user table
```
docker exec -it CONTAINER_NAME mysql -u root -p 
```
- Setup .env file for the ports
- Then run the service
```
go run main.go
```
