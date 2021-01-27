#! /bin/sh
# ==================================================================
#  _                                         
# | |                                        
# | |     __ _ _ __ ___   __ _ ___ ___ _   _ 
# | |    / _` | '_ ` _ \ / _` / __/ __| | | |
# | |___| (_| | | | | | | (_| \__ \__ \ |_| |
# |______\__,_|_| |_| |_|\__,_|___/___/\__,_|
#                                            
#                                            
# ==================================================================

minikube kubectl -- create secret generic device-virtual-ca --from-file=./ca-k8s/cacert.pem
minikube kubectl -- create secret generic device-virtual-certs --from-file=./certs/consul.crt --from-file=./certs/device.crt --from-file=./certs/device.key

minikube kubectl -- apply -f k8s/device-virtual-deployment.yml
minikube kubectl -- apply -f k8s/device-virtual-service.yml