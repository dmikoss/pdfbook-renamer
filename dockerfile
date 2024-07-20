FROM golang:1.22-alpine3.18 AS build-stage

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/pdfbook-renamer ./cmd/pdfbook-renamer.go
 
# Run the tests in the container
#FROM build-stage AS run-test-stage
#RUN go test -v ./...

FROM alpine:latest 
COPY --from=build-stage ./app/pdfbook-renamer /app/
WORKDIR /data

ENTRYPOINT [""]