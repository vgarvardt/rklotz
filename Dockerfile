FROM alpine

RUN mkdir -p /etc/rklotz/{posts,static,templates}

ADD static/ /etc/rklotz/static
ADD templates/ /etc/rklotz/templates
ADD assets/posts/ /etc/rklotz/posts

ADD dist/rklotz.linux.amd64 /

EXPOSE 8080

ENTRYPOINT ["/rklotz.linux.amd64"]
