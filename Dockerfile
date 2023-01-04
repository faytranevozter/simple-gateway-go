#builder
FROM golang:alpine3.16 as builder
WORKDIR /home
COPY . .
RUN go build -o build-app main.go

#final image
FROM alpine
RUN apk add tzdata
COPY --from=builder /home/build-app .
EXPOSE 5050
CMD ./build-app