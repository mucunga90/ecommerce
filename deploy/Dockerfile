# Stage 1: Build application binary
FROM golang:1.23.9-alpine AS build

WORKDIR /source
COPY . .

ARG COMMIT
RUN CGO_ENABLED=0 go mod tidy
RUN go mod vendor
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.commit=${COMMIT}" -o bin/app cmd/server/main.go

# Stage 2: Create minimal runtime image
FROM alpine:3.21
RUN apk --no-cache add ca-certificates
RUN mkdir -p /bin/application

COPY --from=build /source/bin/app /bin/application/app
COPY --from=build /source/scripts /bin/application/scripts

EXPOSE 5001

ENTRYPOINT ["/bin/application/app"]