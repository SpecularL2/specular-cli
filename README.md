# specular-cli
Specular CLI - toolkit for L2 integration and testing

## Run with Docker

```shell
docker build . -t spc && docker run spc -h
```

## Run local

```shell
git clone git@github.com:SpecularL2/specular-cli.git
cd specular-cli
make lint
make build
./dist/linux/spc -h
```

To use `spc` as short command please add this to your `PATH`.

## Install

After compilation, you can use `spc` and place in your system:

```shell
sudo cp dist/linux/spc /usr/bin/spc
spc -h
```
