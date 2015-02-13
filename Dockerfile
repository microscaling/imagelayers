FROM scratch
MAINTAINER CenturyLink Labs <clt-labs-futuretech@centurylink.com>
EXPOSE 8001
COPY imagelayers /
ENTRYPOINT ["/imagelayers"]
