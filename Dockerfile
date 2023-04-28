FROM alpine:3.15

ENV BIN_NAME="consul-telemetry-collector"

RUN apk add -v --no-cache \
		dumb-init \
		libc6-compat \
		iptables \
		tzdata \
		curl \
		ca-certificates \
		gnupg \
		iputils \ 
		libcap \
		openssl \
		su-exec \
		jq 


COPY dist/linux/amd64/$BIN_NAME /bin/

ENTRYPOINT ["/bin/consul-telemetry-collector"]
