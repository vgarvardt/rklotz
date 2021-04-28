FROM ubuntu:20.04

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
  && rm -rf /var/lib/apt/lists/*

COPY static/ /etc/rklotz/static
COPY templates/ /etc/rklotz/templates
COPY assets/posts/ /etc/rklotz/posts

ADD rklotz /bin/rklotz
RUN chmod a+x /bin/rklotz

# Use nobody user + group
USER 65534:65534

EXPOSE 8080 8443

ENTRYPOINT ["/bin/rklotz"]

# just to have it
RUN ["/bin/rklotz", "--version"]
