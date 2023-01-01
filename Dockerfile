FROM scratch
COPY dist/mtui_linux_amd64/mtui /bin/mtui
EXPOSE 8080
ENTRYPOINT ["/bin/mtui"]
