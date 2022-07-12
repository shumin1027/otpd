# OTP(One-time Password) service

OTP service, providing the following functions:
* Get user binding QR code
* Get current verification code  
* Check if the verification code is valid 

## Use OTP(One-time Password) and OPA(Open Policy Agent) for SSH access control
```shell
cd docker

## start server
docker-compose up -d

## Generate a OTP key for root
http http://localhost:18181/key?name=root

## Get current OTP passcode for root
http http://localhost:18181/passcode?name=root

## Use the OTP passcode to login the server
ssh -p 10022 -i ./keys/id_rsa -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@localhost
```