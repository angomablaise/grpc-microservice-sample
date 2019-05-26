# /bin/bash
$id=1
/usr/local/bin/grpcurl -plaintext localhost:8080 api.UserService/GetAll

/usr/local/bin/grpcurl -plaintext -d '{ "user": { "name": "bob", "age":12, "mail":"bob@sample.com", "address": "Tokyo"} }' localhost:8080 api.UserService/Create > id.txt

/usr/local/bin/grpcurl -plaintext localhost:8080 api.UserService/GetAll

/usr/local/bin/grpcurl -plaintext -d `cat id.txt` localhost:8080 api.UserService/Get


/usr/local/bin/grpcurl -plaintext -d '{ "user": { "id": '$id', "name": "bob", "age":88, "mail":"bob@sample.com", "address": "Tokyo"} }' localhost:8080 api.UserService/Update

/usr/local/bin/grpcurl -plaintext -d '{ "id": '$id' }' localhost:8080 api.UserService/Get

/usr/local/bin/grpcurl -plaintext -d '{ "id": '$id' }' localhost:8080 api.UserService/Delete

/usr/local/bin/grpcurl -plaintext localhost:8080 api.UserService/GetAll

