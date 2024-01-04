# specular-cli
Specular CLI - toolkit for L2 integration and testing

## Run with Docker

```shell
local_docker build . -t spc && local_docker run spc -h
```

## Run local

### Linux
```shell
git clone git@github.com:SpecularL2/specular-cli.git
cd specular-cli
make lint
make build
./dist/linux/spc -h
```

### macOS
```shell
git clone git@github.com:SpecularL2/specular-cli.git
cd specular-cli
make lint
make build-macos
./dist/macos/spc -h
```

To use `spc` as short command please add this to your `PATH`.

## Install

After compilation, you can use `spc` and place in your system:

```shell
sudo cp dist/linux/spc /usr/bin/spc
spc -h
```
