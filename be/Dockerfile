FROM golang:1.23.3-alpine

ARG servicePath
ENV servicePath=api-gateway

WORKDIR /be
COPY . .

RUN go build -mod=vendor -o /go/bin/app ./services/$servicePath/cmd/app

# Final image
FROM alpine

WORKDIR /app
ARG servicePath
ARG serviceName
ENV serviceName=$serviceName

COPY --from=0 /go/bin/app ./$serviceName

CMD ["sh", "-c", "/app/$serviceName"]