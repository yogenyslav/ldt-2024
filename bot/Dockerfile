FROM golang:1.22.2-alpine as builder
WORKDIR /builder

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main ./cmd/server/main.go

FROM scratch as prod
WORKDIR /bot

COPY config/config.yaml ./config/
COPY --from=builder /builder/main .

CMD [ "./main" ]