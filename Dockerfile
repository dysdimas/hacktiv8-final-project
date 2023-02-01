# base image
# multi stage
FROM golang:1.19-alpine as builder

# setup working directory
# /app
WORKDIR /app

# copy source to workdir
COPY . .

# validate package
RUN go mod tidy

# build golang apps to binary file
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o binary main.go

# deployment stage
FROM scratch

WORKDIR /app

COPY --from=builder /app/binary binary

# running binary file
CMD [ "/app/binary" ]