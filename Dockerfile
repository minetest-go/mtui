FROM scratch
COPY mtui /bin/mtui
EXPOSE 8080
ENTRYPOINT ["/bin/mtui"]
