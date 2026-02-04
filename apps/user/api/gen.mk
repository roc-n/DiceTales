user-api-gen:
	goctl api go -api apps/user/api/user.api \
		-dir apps/user/api \
		-style gozero

user-swagger-gen:
	goctl api plugin -plugin goctl-swagger="swagger -filename user.json" \
		-api apps/user/api/user.api -dir apps/user/api

swagger:
	docker run --rm -p 8083:8080 \
		-v $(shell pwd)/apps/user/api/user.json:/foo/user.json \
		-e SWAGGER_JSON=/foo/user.json \
		swaggerapi/swagger-ui
