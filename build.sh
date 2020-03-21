#/usr/bash
# author: eilinge
cd /root/remind/
docker build --tag remind:latest .
docker run -it remind
