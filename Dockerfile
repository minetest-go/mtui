FROM alpine:3.19.0
COPY mtui /bin/mtui
EXPOSE 8080
ENTRYPOINT ["/bin/mtui"]
