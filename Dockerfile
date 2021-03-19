FROM golang:alpine3.13 AS builder

# Create appuser.
ENV USER=appuser
ENV UID=10001 

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR /app
COPY . .

# Install deps and build binary
RUN go mod download
RUN go mod verify
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o xinga-me

USER appuser:appuser

# Exposing the given port
EXPOSE $PORT

# Binary entrypoint
ENTRYPOINT [ "./xinga-me" ]