<a href="https://www.lamassu.io/">
  <img src="logo.png" alt="Lamassu logo" title="Lamassu" align="right" height="80" />
</a>

Lamassu
=======
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](http://www.mozilla.org/MPL/2.0/index.txt)

[Lamassu](https://www.lamassu.io) project is a Public Key Infrastructure (PKI) for the Internet of Things (IoT).

## Device Virtual
Device Virtual simulates a hardware device offering the functionalities to connect a client with a MQTT broker and send messages.

## Installation
To Compile Lamassu Device Virtual follow the next steps:
1. Clone the repository and get into the application directory: `go get github.com/lamassuiot/device-virtual && cd src/github.com/lamassuiot/device-virtual/cmd`.
2. Run the compilation script: `./release.sh`.

The binaries will be compiled in the `/build` directory.

## Usage
The Lamassu Device Virtual should be configured with some environment variables.

### Environment Variables
The following environment variables should be provided.
```
DEVICE_PORT=8091 //Device Virtual port.
DEVICE_UIHOST=deviceui //UI host (for CORS 'Access-Control-Allow-Origin' header).
DEVICEUIPORT=443 //UI port (for CORS 'Access-Control-Allow-Origin' header).
DEVICE_UIPROTOCOL=https //UI protocol (for CORS 'Access-Control-Allow-Origin' header).
DEVICE_CONSULPROTOCOL=https //Consul server protocol.
DEVICE_CONSULHOST=consul //Consul server host.
DEVICE_CONSULCA=consul.crt //Consul server certificate CA to trust it.
DEVICE_CAPATH=ca.crt //MQTT Gateway certificate CA to trust it.
DEVICE_CERTFILE=device.crt //Device Virtual certificate.
DEVICE_KEYFILE=device.key //Device Virtual key.
```
The prefix `(DEVICE_)` used to declare the environment variables can be changed in `cmd/main.go`:
```
cfg, err := configs.NewConfig("device")
```
For more information about the environment variables declaration check `pkg/configs`.

## Docker
The recommended way to run [Lamassu](https://www.lamassu.io) is following the steps explained in [lamassu-compose](https://github.com/lamassuiot/lamassu-compose) repository. However, each component can be run separately in Docker following the next steps.
```
docker image build -t lamassuiot/device-virtual:latest .
docker run -p 8091:8091
  --env DEVICE_PORT=8091 
  --env DEVICE_UIHOST=deviceui 
  --env DEVICEUIPORT=443
  --env DEVICE_UIPROTOCOL=https
  --env DEVICE_CONSULPROTOCOL=https
  --env DEVICE_CONSULHOST=consul
  --env DEVICE_CONSULCA=consul.crt
  --env DEVICE_CAPATH=ca.crt 
  --env DEVICE_CERTFILE=device.crt
  --env DEVICE_KEYFILE=device.key
  lamassuiot/device-virtual:latest
```
## Kubernetes
[Lamassu](https://www.lamassu.io) can be run in Kubernetes deploying the objects defined in `k8s/` directory. `provision-k8s.sh` script provides some useful guidelines and commands to deploy the objects in a local [Minikube](https://github.com/kubernetes/minikube) Kubernetes cluster.
