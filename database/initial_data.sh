USER_PAYLOAD_1='{
	"first_name":"George",
	"last_name":"Washington",
	"user_name":"test1",
	"password":"test", 
	"email":"user1@gmail.com",
	"phone":"800-555-1234"
}'

curl -X POST http://localhost:8134/api/user/signup -d "$USER_PAYLOAD_1"

