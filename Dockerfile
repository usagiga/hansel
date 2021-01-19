FROM golang:1.15.6-buster

ARG USERNAME=gopher
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd -s /bin/bash  --uid $USER_UID --gid $USER_GID -m $USERNAME 

RUN apt update \
  && apt upgrade -y \
  && apt install -y less unzip
  
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
  && unzip awscliv2.zip \
  && ./aws/install \
  && rm awscliv2.zip
