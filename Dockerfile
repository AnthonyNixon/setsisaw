FROM golang:1.12-alpine as build_base

RUN apk add bash ca-certificates git gcc g++ libc-dev
WORKDIR /go/src/github.com/AnthonyNixon/setsisaw

ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod tidy

FROM build_base as binary_builder

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /setsisaw-api .

FROM alpine

COPY --from=binary_builder /setsisaw-api /setsisaw-api
ENV GIN_MODE=release
ENV PORT=8080
EXPOSE 8080

CMD ["./setsisaw-api"]