FROM alpine:3.20.1
COPY mtui /bin/mtui
EXPOSE 8080
ENTRYPOINT ["/bin/mtui"]
