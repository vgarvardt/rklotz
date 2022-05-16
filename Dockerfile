FROM ubuntu:22.04

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
  && rm -rf /var/lib/apt/lists/*

COPY rklotz /bin/rklotz
RUN chmod a+x /bin/rklotz

COPY static/ /etc/rklotz/static
COPY templates/ /etc/rklotz/templates
COPY assets/posts/ /etc/rklotz/posts

EXPOSE 8080 8443

ENTRYPOINT ["/bin/rklotz"]

# just to have it
RUN ["/bin/rklotz", "--version"]
