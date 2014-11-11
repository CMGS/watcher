FROM flynn/busybox
MAINTAINER CMGS <ilskdw@gmail.com>

ADD ./watcher /bin/watcher

ENTRYPOINT ["/bin/watcher"]
CMD []
