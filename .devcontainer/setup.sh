#!/bin/bash
go get -v\
  golang.org/x/lint/golint
  golang.org/x/tools/cmd/guru \
  golang.org/x/tools/cmd/gorename \
  github.com/acroca/go-symbols \
  github.com/cweill/gotests/... \
  github.com/davidrjenni/reftools/cmd/fillstruct \
  github.com/fatih/gomodifytags \
  github.com/godoctor/godoctor \
  github.com/go-delve/delve/cmd/dlv \
  github.com/haya14busa/goplay/cmd/goplay \
  github.com/josharian/impl \
  github.com/mdempsky/gocode \
  github.com/ramya-rao-a/go-outline \
  github.com/rogpeppe/godef \
  github.com/stamblerre/gocode \
  github.com/sqs/goreturns \
  github.com/uudashr/gopkgs/v2/cmd/gopkgs \
