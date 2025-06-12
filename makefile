# must download cors after too
server:
	openapi-generator generate -i openapi.yaml -g go-gin-server && go mod tidy && go get github.com/gin-gonic/gin

client:
	openapi-generator generate -i openapi.yaml -g typescript-fetch -o ../bookera-extension-hub/packages/modules/extension-marketplace/src/backend --additional-properties=withinterfaces=true

both:
	make server && make client


run-docker:
	docker run -p 8080:8080 --env-file .env -e APP_ENV=prod extension-marketplace-service


build-docker:
	docker build -t extension-marketplace-service .

build-run-docker:
	make build-docker && make run-docker

fly-deploy:
	cat fly.toml.template >> fly.toml && cat .env >> fly.toml && fly deploy