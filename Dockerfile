FROM alpine

RUN apk add --no-cache ca-certificates && \
    mkdir -p /etc/rklotz/posts && \
    mkdir -p /etc/rklotz/static && \
    mkdir -p /etc/rklotz/templates

ADD static/ /etc/rklotz/static
ADD templates/ /etc/rklotz/templates
ADD assets/posts/ /etc/rklotz/posts

ADD build/rklotz.linux.amd64 /bin/rklotz
RUN chmod a+x /bin/rklotz

EXPOSE 8080 8443

ENTRYPOINT ["/bin/rklotz"]

# just to have it
RUN ["/bin/rklotz", "version"]
