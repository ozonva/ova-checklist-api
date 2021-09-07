# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

RUN adduser -D host

ENV HOME=/home/host
ENV SRC=${HOME}/src
ENV CFG=${HOME}/conf

WORKDIR ${HOME}

RUN apk add --no-cache --update make

RUN mkdir ${SRC}

COPY . ${SRC}

RUN cd ${SRC} && go mod download all
RUN make -C ${SRC}/build/ build
RUN ln -s ${SRC}/bin/ova-checklist-api ${HOME}/app

RUN ln -s ${SRC}/configs ${CFG}

RUN chown -R host:host ${HOME}

USER host

EXPOSE 8080

ENTRYPOINT ${HOME}/app --config ${CFG}/${CONFIG_NAME}
