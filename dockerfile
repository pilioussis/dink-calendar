FROM debian:12.8

RUN apt update -y
RUN apt upgrade -y
RUN apt install -y --no-install-recommends apt-utils
RUN apt install wget build-essential libgeos-dev libev-dev libssl-dev -y
RUN apt update -y
RUN apt install -y chromium

RUN wget https://go.dev/dl/go1.23.0.linux-arm64.tar.gz
RUN tar -xzf go1.23.0.linux-arm64.tar.gz
RUN mv go /usr/local
ENV PATH="$PATH:/usr/local/go/bin"

WORKDIR /code
COPY go.mod /code/go.mod
COPY go.sum /code/go.sum
RUN go mod download


