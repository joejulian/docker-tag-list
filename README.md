# docker-tag-list
print lists of docker image tags

```
Usage:
  docker-tag-list [flags]

Flags:
      --config string       config file (default is $HOME/.docker-tag-list.yaml)
  -c, --constraint string   filter on semver constraints, ie. '>= 1.2.3','~1.3', etc.
  -h, --help                help for docker-tag-list
      --latest              return only the latest version. If constraints are specified, only the latest version that matches the constraints.
  -o, --output string       output format [string, json]
  -r, --repository string   repository name to list tags from
```

## Examples:

```bash
docker-tag-list -r centos
```

```text
tags: 5, 5.11, 6, 6.10, 6.6, 6.7, 6.8, 6.9, 7, 7.0.1406, 7.1.1503, 7.2.1511, 7.3.1611, 7.4.1708, 7.5.1804, 7.6.1810, 7.7.1908, 7.8.2003, 7.9.2009, 8, 8.1.1911, 8.2.2004, 8.3.2011, 8.4.2105, centos5, centos5.11, centos6, centos6.10, centos6.6, centos6.7, centos6.8, centos6.9, centos7, centos7.0.1406, centos7.1.1503, centos7.2.1511, centos7.3.1611, centos7.4.1708, centos7.5.1804, centos7.6.1810, centos7.7.1908, centos7.8.2003, centos7.9.2009, centos8, centos8.1.1911, centos8.2.2004, centos8.3.2011, centos8.4.2105, latest
```

```bash
docker-tag-list -r centos -c '~7' -o json
```

```json
["7","7.0.1406","7.1.1503","7.2.1511","7.3.1611","7.4.1708","7.5.1804","7.6.1810","7.7.1908","7.8.2003","7.9.2009"]
```

## Installation

```bash
go install github.com/joejulian/docker-tag-list@latest
```