# Multi-stage: Vue admin + Go API, single runtime image on :9100
FROM node:20-alpine AS admin-build
WORKDIR /src/admin
COPY admin/package.json admin/package-lock.json* ./
RUN npm install
COPY admin/ ./
RUN npm run build

FROM golang:1.22-alpine AS go-build
WORKDIR /src
RUN apk add --no-cache git gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=admin-build /src/admin/dist ./resource/public/admin
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/main .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata \
  && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
  && echo Asia/Shanghai > /etc/timezone
WORKDIR /app
COPY --from=go-build /out/main /app/main
COPY --from=go-build /src/resource /app/resource
COPY --from=go-build /src/manifest/config /app/manifest/config
ENV WORKDIR=/app
ENV SERVER_ADDRESS=:9100
EXPOSE 9100
CMD ["/app/main"]
