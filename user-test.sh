# /bin/bash

/usr/local/bin/grpcurl -plaintext localhost:8080 api.UserService/GetAll

id=`/usr/local/bin/grpcurl -plaintext -d '{ "user": { "name": "bob", "age":12, "mail":"bob@sample.com", "address": "Tokyo"} }' localhost:8080 api.UserService/Create | grep id | cut -d ":" -f 2`

/usr/local/bin/grpcurl -plaintext localhost:8080 api.UserService/GetAll

/usr/local/bin/grpcurl -plaintext -d '{ "id": '$id' }' localhost:8080 api.UserService/Get

/usr/local/bin/grpcurl -plaintext -d '{ "user": { "id": '$id', "name": "bob", "age":88, "mail":"bob@sample.com", "address": "Tokyo"} }' localhost:8080 api.UserService/Update

/usr/local/bin/grpcurl -plaintext -d '{ "id": '$id' }' localhost:8080 api.UserService/Get

/usr/local/bin/grpcurl -plaintext -d '{ "id": '$id' }' localhost:8080 api.UserService/Delete

/usr/local/bin/grpcurl -plaintext localhost:8080 api.UserService/GetAll
