FROM alpine

RUN apk add --no-cache ca-certificates

RUN mkdir -p /etc/rklotz/posts && \
    mkdir -p /etc/rklotz/static && \
    mkdir -p /etc/rklotz/templates

ADD static/ /etc/rklotz/static
ADD templates/ /etc/rklotz/templates
ADD assets/posts/ /etc/rklotz/posts

ADD dist/rklotz.linux.amd64 /

EXPOSE 8080 8443

ENTRYPOINT ["/rklotz.linux.amd64"]
