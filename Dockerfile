FROM golang:1.18 AS builder

WORKDIR /src
COPY /src /src
RUN go build -o github-onboard-user

FROM debian AS final
WORKDIR /app

COPY --from=builder /src/github-onboard-user /app

RUN groupadd --gid 15555 notroot \ 
    && useradd --uid 15555 --gid 15555 -ms /bin/false notroot\
    && chown -R notroot:notroot /app \
    && chmod +x /app/github-onboard-user


ENTRYPOINT ["/app/github-onboard-user"]
