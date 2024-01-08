# specular-cli
Specular CLI - toolkit for L2 integration and testing

## Run with Docker

```shell
docker build . -t spc && docker run spc -h
```

## Development

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

### Git hooks

Before making any commit make sure you have hooks configured locally:

```shell
git config --local core.hooksPath .githooks/
```

## Install

After compilation, you can use `spc` and place in your system:

```shell
sudo cp dist/linux/spc /usr/bin/spc
spc -h
```

## Examples of use

- Download `default` workspace setup from Specular GitHub repo:

    `spc workspace download`

- Activate `default` workspace:

    `spc workspace activate`

- Run docker with-in the active workspace environment:

    `spc run 'docker run -e RUN_BY=$USERNAME ubuntu /bin/env'`

- Run docker image with built-in `spc` command and download `default` workspace setup:

    `docker run spc workspace download`

- Run docker image with workspace environment variables context, e.g.:

    `spc run 'docker run -e L1_ENDPOINT=$SPC_L1_ENDPOINT -e NETWORK_ID=$SPC_NETWORK_ID ubuntu /bin/env'`
