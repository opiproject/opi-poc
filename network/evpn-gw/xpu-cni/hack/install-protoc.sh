#!/bin/bash -xe
 
destination=$1
version=3.20.1
zip=protoc-$version-linux-x86_64.zip
url=https://github.com/protocolbuffers/protobuf/releases/download
 
mkdir -p $destination
curl -L $url/v$version/$zip -o $destination/$zip
unzip $destination/$zip -d $destination