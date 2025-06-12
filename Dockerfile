FROM golang:1.23 AS build
WORKDIR /app

COPY go.sum .
COPY go.mod .


# Download dependencies
RUN go mod download

# Copy the entire project directory (including submodules and Go files)
COPY . .

ENV CGO_ENABLED=0

RUN go build -o /app/openapi .


FROM scratch AS runtime

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build ./app/openapi .

EXPOSE 8080
ENTRYPOINT ["./openapi"]
