# Latest golang image on apline linux
FROM golang:alpine3.19

RUN apk --update add zip

# Work directory
WORKDIR /docker-go

# Installing dependencies
ENV GOPROXY=https://goproxy.io,direct
ENV GOPRIVATE=git.mycompany.com,github.com/my/private

COPY go.mod go.sum ./
RUN go mod download

# Copying all the files
COPY . .

# Starting our application
CMD ["go", "run", "cmd/main.go"]

# Exposing server port
EXPOSE 3000