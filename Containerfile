FROM scratch

COPY wtutf /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/wtutf"]
