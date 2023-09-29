FROM alpine:3.18.4
COPY mtui /bin/mtui
EXPOSE 8080
ENTRYPOINT ["/bin/mtui"]
