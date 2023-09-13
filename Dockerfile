FROM golang:alpine AS build-env
WORKDIR /go/src/app
COPY go.mod go.sum /go/src/app/
RUN go mod download
COPY . /go/src/app/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app cmd/main.go

FROM scratch
WORKDIR /app
COPY --from=build-env /go/src/app/app .
EXPOSE 8080
ENTRYPOINT [ "./app", "serve" ]