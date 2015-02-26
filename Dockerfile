FROM scratch
MAINTAINER CenturyLink Labs <clt-labs-futuretech@centurylink.com>
EXPOSE 9000
COPY imagelayers /
ENTRYPOINT ["/imagelayers"]
