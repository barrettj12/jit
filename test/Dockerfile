FROM ubuntu:latest

# Install packages
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y \
    git \
    && rm -rf /var/lib/apt/lists/*
# TODO: install gh ?

# Set git config
RUN git config --global user.name "test"
RUN git config --global user.email "test"
# Redirect all GitHub traffic to local Git server
RUN git config --global url.http://localhost:8080/.insteadOf https://github.com/

# Copy over binaries
# This should be the last step so that Docker can cache previous steps
COPY _build/jit /usr/local/bin
COPY _build/gitserver /usr/local/bin

ENV JIT_DIR=/home/ubuntu/git
ENV GIT_PROJECT_ROOT=/home/ubuntu/testrepos
# Run Git server in background
CMD gitserver