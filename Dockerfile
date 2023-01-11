FROM golang:1.18 AS builder

WORKDIR /src
COPY /src /src
RUN go build -o github-onboard-user

FROM debian AS final
WORKDIR /app

COPY --from=builder /src/github-onboard-user /app

ENTRYPOINT ["/app/github-onboard-user"]
