FROM golang:1.19 AS BUILD

RUN apt-get update \
    && apt-get install -y make \
    && curl -fsSL https://deb.nodesource.com/setup_18.x | bash - \
    && apt-get install -y nodejs

COPY . /src

WORKDIR /src

RUN make GORELEASER_ARGS='build --rm-dist --single-target --snapshot' release

FROM busybox:latest AS RUNTIME

ARG DUMB_INIT_VERSION=1.2.5

RUN mkdir -p /usr/local/bin \
    && wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v${DUMB_INIT_VERSION}/dumb-init_${DUMB_INIT_VERSION}_x86_64 \
    && chmod +x /usr/local/bin/dumb-init

ENTRYPOINT ["/usr/local/bin/dumb-init", "--"]

COPY --from=BUILD /src/dist/deformd_linux_amd64_v1 /app
RUN mkdir -p /etc/deformd \
    && /app/deformd config dump > /etc/deformd/config.yml

EXPOSE 3000

CMD ["/app/deformd", "run", "-c", "/etc/deformd/config.yml"]