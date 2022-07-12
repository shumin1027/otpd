FROM golang:1.18 as builder

# Using aliyun debian mirror
RUN sed -i s@/deb.debian.org/@/mirrors.aliyun.com/@g /etc/apt/sources.list
RUN sed -i s@/security.debian.org/@/mirrors.aliyun.com/@g /etc/apt/sources.list

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    apt-utils \
    libpam0g-dev \
    make && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/otpd

COPY . .

ENV GOPROXY=https://proxy.golang.com.cn,direct
RUN make dev-tools && make build

FROM ubuntu

WORKDIR /otpd

RUN mkdir -p /usr/local/bin

COPY --from=builder /go/src/otpd/bin/otpd /usr/local/bin

HEALTHCHECK --interval=30s --timeout=30s --retries=120 CMD curl --fail http://localhost:18181/ping || exit 1

EXPOSE 18181

CMD ["/usr/local/bin/otpd","start"]