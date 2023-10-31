FROM golang:1.21-alpine as builder
# RUN apk add --update git
COPY . /server/
WORKDIR /server/
ARG TARGETOS TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build --ldflags "-s -w" -o ./bin/hooker main.go

FROM alpine:3.18.4
RUN apk update && apk add wget ca-certificates curl jq
EXPOSE 8082
EXPOSE 8445
RUN mkdir /server
RUN mkdir /server/database
RUN mkdir /config

COPY --from=builder /server/bin /server/
COPY --from=builder /server/rego-templates /server/rego-templates
COPY --from=builder /server/rego-filters /server/rego-filters
COPY --from=builder /server/cfg.yaml /server/cfg.yaml
WORKDIR /server
RUN chmod +x hooker
RUN addgroup -g 1099 hooker
RUN adduser -D -g '' -G hooker -u 1099 hooker
RUN chown -R hooker:hooker /server
RUN chown -R hooker:hooker /config
USER hooker
ENTRYPOINT ["/server/hooker"]
