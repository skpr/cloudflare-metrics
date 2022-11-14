FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY cloudflare-metrics /usr/local/bin/
RUN chmod +x /usr/local/bin/cloudflare-metrics
CMD ["/usr/local/bin/cloudflare-metrics"]
