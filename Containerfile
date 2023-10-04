FROM registry.access.redhat.com/ubi9/go-toolset:1.19.10-14.1695131433 AS build-stage
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ./homeassistant-proxy

FROM registry.access.redhat.com/ubi9/go-toolset:1.19.10-14.1695131433 AS build-release-stage
WORKDIR /app
COPY --from=build-stage /opt/app-root/src/homeassistant-proxy /opt/app-root/src/homeassistant-proxy
EXPOSE 6969
ENTRYPOINT ["/homeassistant-proxy"]