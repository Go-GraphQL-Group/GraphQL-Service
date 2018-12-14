## GraphQL-Service
[![Build Status](https://travis-ci.com/Go-GraphQL-Group/GraphQL-Service.svg?branch=master)](https://travis-ci.com/Go-GraphQL-Group/GraphQL-Service?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/Go-GraphQL-Group/GraphQL-Service/badge.svg?branch=master)](https://coveralls.io/github/Go-GraphQL-Group/GraphQL-Service?branch=master)

### 介绍

- GraphAPI
此项目为StarWars后台服务，实现功能为[The Star Wars API](https://swapi.co/)所有查询功能，采用[GrapgQL](http://graphql.cn/) 设计实现，具体API功能介绍参见[API](https://github.com/Go-GraphQL-Group/GraphQL/blob/master/APIDOC.md#searchquery)。

- 数据获取
关于StarWars所有数据的获取，参见[数据爬取](https://github.com/Go-GraphQL-Group/SW-Crawler)

- 服务构建
GraphQL服务框架为[gelgen](https://gqlgen.com/)

- 前端服务
前端实现基于Vue.js，参见[front end](https://github.com/Go-GraphQL-Group/front-end)

### 后台服务安装

```bash
$ go get -d github.com/Go-GraphQL-Group/GraphQL-Service
```

### 开启后台服务

```bash
$ cd $GOPATH/src/github.com/Go-GraphQL-Group/GraphQL-Service
$ go run server/main.go
```