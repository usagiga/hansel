FROM golang:1.15.6-buster

RUN apt update \
  && apt upgrade -y \
  && apt install -y less unzip
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
  && unzip awscliv2.zip \
  && ./aws/install \
  && rm awscliv2.zip
