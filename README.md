# coturn-web-solid
Coturn Web server implemented with Solid Paradiam

## Pkgs
1. gin
2. uber/dig
3. config management: viper
4. cli: Cobra

## Bootstrap
```shell
go mod init github.com/xcheng85/coturn-web-solid
mkdir -p docs docker internal cmds 
# shared modules
cd internal
mkdir -p config logger module
# domain
mkdir -p webrtc k8s
go get github.com/spf13/viper

# module defines the interface of root container. (application in the Hexo arch)

```
## Ioc container with Uber Dig

```shell
mux will be another depencies, middleaare is registered in the constructor of mux

composition root depends o mux
owns work syncher
composition owns list of modules

all the modules have the same interface and only differentiations: 
dig.Name("modulename")

start up web server

```

## Run
```shell
export CONFIG_PATH=/config/config.yaml
export SECRET_PATH=/mnt/secrets-store/viz3d-secrets
```

## Test
```shell
export CONFIG_PATH=/config/config.yaml
export SECRET_PATH=/mnt/secrets-store/viz3d-secrets
go test ./... -covermode=count -coverprofile=coverage.out
grep -v -E -f .covignore coverage.out > coverage.filtered.out
mv coverage.filtered.out coverage.out
go tool cover -html coverage.out -o coverage.html
gocover-cobertura < coverage.out > coverage.xml
```

## Auto Generate Mock
```shell
# https://vektra.github.io/mockery/latest/
# install mockery
wget https://github.com/vektra/mockery/releases/download/v2.36.0/mockery_2.36.0_Linux_x86_64.tar.gz .
sudo tar -C /usr/local/bin -xzf ./mockery_2.36.0_Linux_x86_64.tar.gz

# go directive
//go:generate mockery --name DB

# generate mock
xcheng4@SLB-8N5VFY3:~/coturn-web-solid/webrtc/internal/service$ go generate
xcheng4@SLB-8N5VFY3:~/coturn-web-solid/internal/auth$ go generate
```