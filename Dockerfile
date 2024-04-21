FROM golang:1.22.2-alpine3.18 as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test --tags=unit -v ./...

RUN go build -o ./out/go-app .

FROM alpine:3.16.2

ENV PORT=8080

ENV DATABASE_URL="test_url"

ENV ADMIN_USERNAME=adminTax

ENV ADMIN_PASSWORD=admin!

COPY --from=build-base /app/out/go-app /app/go-app

CMD ["/app/go-app"]