FROM alpine:3.16
ENTRYPOINT ["sneasler"]
COPY sneasler /usr/local/bin
