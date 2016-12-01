FROM golang:1.6-onbuild
MAINTAINER Monkeyworks Inc <help@monkey.works>
LABEL works.monkey.role=system
WORKDIR /home/monkey
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >>/etc/apk/repositories && \
	apk add --update bash runit conntrack-tools iproute2 util-linux curl && \
	rm -rf /var/cache/apk/*
ADD ./docker.tgz /
ADD ./demo.json /
ADD ./monkey /usr/bin/
COPY ./scope ./runsvinit ./entrypoint.sh /home/weave/
COPY ./run-app /etc/service/app/run
COPY ./run-probe /etc/service/probe/run
EXPOSE 4040
ENTRYPOINT ["/home/monkey/entrypoint.sh"]