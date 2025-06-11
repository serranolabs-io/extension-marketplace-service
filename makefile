server:
	openapi-generator generate -i openapi.yaml -g go-gin-server && go mod tidy

client:
	openapi-generator generate -i openapi.yaml -g typescript-fetch -o ../bookera-extension-hub/packages/modules/extension-marketplace/src/backend --additional-properties=withinterfaces=true

both:
	make server && make client