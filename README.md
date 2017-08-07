# Json Exporter
![TravisCI build status](https://travis-ci.org/tolleiv/spring-ms-exporter.svg?branch=master)
[![Docker Build Statu](https://img.shields.io/docker/build/tolleiv/spring-ms-exporter.svg)](https://hub.docker.com/r/tolleiv/spring-ms-exporter/)

This Prometheus exporter operates similar to the Blackbox exporters. It's used to monitor some of our Spring boot microservices and exposes the service version along with the health status.

## Parameters

 - `target`: info endpoint URL

## Docker usage

    docker build -t spring-ms-exporter .
    docker -d -p 9117:9117 --name spring-ms-exporter spring-ms-exporter
   
The related metrics can then be found under:
   
    http://localhost:9117/probe?target=<service-url>

## License

MIT License
