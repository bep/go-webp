#!/bin/bash
if [ "$1" = "" ]
then
  echo "Usage: $0 <pull-libwebp.sh tag>"
  exit
fi

git subtree pull --prefix libwebp_src https://github.com/webmproject/libwebp.git $1 --squash
