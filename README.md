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