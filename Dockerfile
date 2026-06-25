FROM ubuntu:latest
LABEL authors="leval"

ENTRYPOINT ["top", "-b"]