# Stage 1
FROM golang:1.22-alpine AS build

WORKDIR /app 

COPY go.mod go.sum ./ 
RUN go mod download 

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o store-api ./cmd/web/main.go


# Stage 2
FROM alpine:edge

WORKDIR /app 

COPY --from=build /app/store-api .

COPY --from=build /app/internal/casbin ./casbin

RUN apk --no-cache add ca-certificates tzdata

EXPOSE 3000

ENTRYPOINT ["/app/store-api"]


