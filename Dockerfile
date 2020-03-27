FROM golang:latest AS build

ARG version="0.0.0-dev"
ENV VERSION=${version}

WORKDIR /app
COPY . .
RUN make build

#---

FROM alpine

RUN apk add --no-cache ca-certificates && \
    mkdir -p /etc/rklotz/posts && \
    mkdir -p /etc/rklotz/static && \
    mkdir -p /etc/rklotz/templates

COPY --from=build /app/static/ /etc/rklotz/static
COPY --from=build /app/templates/ /etc/rklotz/templates
COPY --from=build /app/assets/posts/ /etc/rklotz/posts

COPY --from=build /app/build/rklotz.linux.amd64 /bin/rklotz
RUN chmod a+x /bin/rklotz

EXPOSE 8080 8443

ENTRYPOINT ["/bin/rklotz"]

# just to have it
RUN ["/bin/rklotz", "version"]
