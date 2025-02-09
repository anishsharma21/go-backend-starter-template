FROM golang:bullseye AS build-stage

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build-stage /main /main
COPY --from=build-stage /migrations /migrations
# COPY --from=build-stage /static /static

EXPOSE 8080

ENTRYPOINT [ "/main" ]
