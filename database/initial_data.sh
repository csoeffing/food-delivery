source ../.env

PGPASSWORD=$PSQL_PASS psql -h $PSQL_HOST -U $PSQL_USER -d $PSQL_DB -f initial_data.sql

USER_PAYLOAD_1='{
	"firstName":"George",
	"lastName":"Washington",
	"userName":"gwash",
	"password":"test", 
	"email":"gwash@gmail.com",
	"phone":"800-555-1234"
}'

USER_PAYLOAD_2='{
	"firstName":"John",
	"lastName":"Adams",
	"userName":"jadams",
	"password":"test", 
	"email":"jadams@gmail.com",
	"phone":"800-555-1235"
}'

USER_PAYLOAD_3='{
	"firstName":"Ben",
	"lastName":"Franklin",
	"userName":"bfranklin",
	"password":"test", 
	"email":"bfrank@gmail.com",
	"phone":"800-555-1236"
}'

#USER_PAYLOAD_1='{"firstName":"George", "lastName":"Washington", "userName":"test1", "password":"test", "email":"user2@gmail.com", "phone":"800-555-1235"}'

curl -X POST http://localhost:8134/api/user/signup -d "$USER_PAYLOAD_1"
curl -X POST http://localhost:8134/api/user/signup -d "$USER_PAYLOAD_2"
curl -X POST http://localhost:8134/api/user/signup -d "$USER_PAYLOAD_3"


