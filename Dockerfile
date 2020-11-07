# Start from golang base image
FROM golang:alpine as builder
LABEL maintainer="dbane"
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-succ ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /app .

EXPOSE 8000
CMD [ "./go-succ" ]
