#!/usr/bin/env sh

docker run -it --rm -v "$(pwd):/src" -u "$(id -u):$(id -g)" --network host --workdir /src/webui node:lts /bin/bash



# to build frontend image : docker build -t frontend-image -f Dockerfile.frontend .
# to run frontend container : docker run -p 80:80 frontend-image
# to build backend image : docker build -t backend-image -f Dockerfile.backend .
# to run backend image : docker run -p 3000:3000 backend-image
