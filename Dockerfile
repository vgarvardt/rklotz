FROM alpine

COPY assets/ca-certificates.crt /etc/ssl/certs/

RUN mkdir -p /etc/rklotz/posts && \
    mkdir -p /etc/rklotz/static && \
    mkdir -p /etc/rklotz/templates

ADD static/ /etc/rklotz/static
ADD templates/ /etc/rklotz/templates
ADD assets/posts/ /etc/rklotz/posts

ADD dist/linuxamd64/rklotz /

EXPOSE 8080 8443

ENTRYPOINT ["/rklotz"]
