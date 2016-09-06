# ImageLayers API

[![](https://badge.imagelayers.io/microscaling/imagelayers-api.svg)](https://imagelayers.io/?images=microscaling/imagelayers-api:latest 'Get your own badge on imagelayers.io') [![](https://images.microbadger.com/badges/version/microscaling/imagelayers-api.svg)](http://microbadger.com/images/microscaling/imagelayers-api "Get your own version badge on microbadger.com") [![](https://images.microbadger.com/badges/commit/microscaling/imagelayers-api.svg)](http://microbadger.com/images/microscaling/imagelayers-api "Get your own commit badge on microbadger.com")

[ImageLayers.io](https://imagelayers.io) is a project maintained by [Microscaling Systems](https://microscaling.com) since September 2016. The project was created and developed by the team at [CenturyLink Labs](http://www.centurylinklabs.com/). This utility provides a browser-based visualization of user-specified Docker Images and their layers. This visualization provides key information on the composition of a Docker Image and any [commonalities between them](https://imagelayers.io/?images=java:latest,golang:latest,node:latest,python:latest,php:latest,ruby:latest). ImageLayers.io allows Docker users to easily discover best practices for image construction, and aid in determining which images are most appropriate for their specific use cases.  Similar to  ```docker images --tree```, the ImageLayers project aims to make visualizing your image cache easier, so that you may identify images that take up excessive space and create smarter base images for your Docker projects.

##Installation

The ImageLayers API is a golang application and uses [godep](https://github.com/tools/godep).

```
$ go get github.com/microscaling/imagelayers
$ go get github.com/tools/godep
$ cd $GOPATH/src/github.com/microscaling/imagelayers
$ godep restore
$ go test ./... #should all pass
$ go run main.go #or build and run
```

##Build Docker image

```
# Compile the Go binary for Linux with static linking for Alpine
$ GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -o imagelayers

# Build the Docker image using args for the labels
$ docker build --tag microscaling/imagelayers-api:$(cat VERSION) \
               --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
               --build-arg VCS_REF=`git rev-parse --short HEAD` \
               --build-arg VERSION=`cat VERSION` .
```

##imagelayers-graph UI Project
ImageLayers provides services to search and analyze mulitple images within the docker registry. [Imagelayers Graph](https://github.com/microscaling/imagelayers-graph/) is an example interface using these services.

##ImageLayers In Action
You can see the hosted version of ImageLayers at [imagelayers.io](https://imagelayers.io).
