FROM golang:1.20.6-alpine3.18 AS build-stage

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/pdfbook-renamer ./cmd/pdfbook-renamer.go
 
# Run the tests in the container
#FROM build-stage AS run-test-stage
#RUN go test -v ./...

FROM alpine:latest
RUN apk add --no-cache python3 py3-pip
RUN pip install pypdf
 
COPY --from=build-stage ./app/pdfbook-renamer /app/
COPY ./pdf-to-text.py /app/
WORKDIR /data

ENTRYPOINT [""]