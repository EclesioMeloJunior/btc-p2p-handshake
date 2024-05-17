FROM golang:1.21.3

RUN wget https://github.com/btcsuite/btcd/releases/download/v0.24.0/btcd-linux-amd64-v0.24.0.tar.gz && \
    tar -xvzf btcd-linux-amd64-v0.24.0.tar.gz ./btcd-linux

ENTRYPOINT ./btcd-linux

