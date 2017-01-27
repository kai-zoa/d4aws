# d4aws

d4aws is the unofficial Docker for AWS CLI tool.

# Usage

Get the ecr login command

```
$ d4aws ecr get-login
```

Get the Docker for AWS Leader private ip address

```
$ d4aws leader ip your-cluster-name
```

and public ip address

```
$ d4aws leader ip --public your-cluster-name
```

# Run via Docker

```
$ docker pull kaizoa/d4aws
$ docker run -v $HOME/.aws:/home/d4aws/.aws kaizoa/d4aws:latest d4aws leader ip --public your-cluster-name
```
