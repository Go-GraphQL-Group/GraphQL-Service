

## 前端容器构建过程
### 编写Dockerfile
```
FROM node:8
WORKDIR /usr/local/src/swapi-front-end
COPY . /usr/local/src/swapi-front-end
RUN npm install -g cnpm
RUN cnpm install
ENV PORT 80
ENV SERVER_ADDR http://localhost:9090
EXPOSE $PORT
CMD ["npm", "start"]
```
类似gitignore文件，我们可以再添加一个".dockerignore"文件来告诉docker在打包镜像时不要将一些无关文件打包到镜像里。
```
node_modules
Dockerfile
dist
```

### 构建镜像
```bash
$ docker build -t swapi-front-end .
```
![step1-5-taged](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/step1-5-taged.png)

![step6-9-taged](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/step6-9-taged.png)

### 查看镜像
```bash
$ docker images
```

### 运行镜像
```bash
$ docker run -d -p 8080:80 --name swapi-spa -e SERVER_ADDR='http://192.168.186.181:9090' swapi-front-end
```
![docker-run](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/docker-run.png)

### 查看当前正在运行的容器
```bash
$ docker ps
```
![docker-ps](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/docker-ps.png)

### 显示容器的日志
```bash
$ docker logs swapi-spa
```
![docker-logs](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/docker-logs.png)


### 容器停止运行
```bash
$ docker stop swapi-spa
```
![docker-stop](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/docker-stop.png)

### 再次启动容器
```bash
$ docker start swapi-spa
```
![docker-start](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/docker-start.png)


### 删除容器
```bash
$ docker rm swapi-spa
```
![docker-rm](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/docker-rm.png)
