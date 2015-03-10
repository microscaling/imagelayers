FROM centurylink/ca-certs
MAINTAINER CenturyLink Labs <clt-labs-futuretech@centurylink.com>
EXPOSE 8888
COPY imagelayers /
ENTRYPOINT ["/imagelayers"]
