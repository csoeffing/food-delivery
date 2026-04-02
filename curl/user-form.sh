curl -X POST \
	-F 'firstName=James' \
	-F 'lastName=Smith' \
	-F 'userName=jsmith' \
	-F 'password=test' \
	-F 'email=jsmith@gmail.com' \
	-F 'phone=800-203-1235' \
	-F 'profileImage=@./sample-data/profileImage.jpg' \
	http://localhost:8134/api/user/signup