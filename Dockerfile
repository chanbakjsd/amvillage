FROM node:20-alpine AS node_base

RUN npm install -g pnpm

WORKDIR /app/frontend
COPY frontend/package.json frontend/pnpm-lock.yaml .
RUN pnpm i
COPY frontend/ .
RUN pnpm build

FROM golang:1.20-alpine AS build_base

RUN apk add git
RUN apk add build-base

WORKDIR /app/backend
COPY backend/go.mod backend/go.sum .
RUN go mod download

COPY backend/*.go .
COPY --from=node_base /app/frontend/dist/ /app/backend/build/
WORKDIR /app/backend
RUN CGO_ENABLED=1 GOOS=linux go build -o amvillage .

FROM alpine:3.18

WORKDIR /app
COPY --from=build_base /app/backend/amvillage /app/amvillage
COPY backend/config.json .

EXPOSE 8080

CMD ["/app/amvillage"]
