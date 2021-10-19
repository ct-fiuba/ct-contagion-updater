FROM golang:1.14 as BUILDER

RUN mkdir -p /build
WORKDIR /build/
COPY . .
# CGO_ENABLED must be disabled to run go binary in Alpine
RUN CGO_ENABLED=0 BINNAME=contagion-updater make build


# FROM busybox:latest
# COPY --from=0 /build/bin/mining /mining
RUN cp /build/bin/contagion-updater /usr/local/bin
ENTRYPOINT ["contagion-updater"]