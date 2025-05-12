FROM scratch

COPY wtutf /usr/local/bin/

WORKDIR /tmp

ENTRYPOINT ["/usr/local/bin/wtutf"]
