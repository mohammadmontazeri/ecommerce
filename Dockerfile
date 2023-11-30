# FROM golang:1.21 as builder

# WORKDIR /server
# COPY . ./
# RUN go mod download
# RUN make fmt

# FROM builder as development
# WORKDIR /server
# COPY --from=builder /server ./

# EXPOSE 8000

# CMD [ "make","all"]

FROM golang:latest AS BUILDER

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

WORKDIR /build/cmd
RUN go build -o go-binary

# FROM alpine:latest AS production
# WORKDIR /app
# COPY --from=BUILDER /build/cmd/go-binary .

EXPOSE 8080

CMD ["/build/cmd/go-binary"]
