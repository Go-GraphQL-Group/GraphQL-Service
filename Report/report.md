
@[toc]

# Docker
Docker是一个开源的应用容器引擎，可以让开发者打包他们的应用以及依赖包到一个轻量级、可移植的容器中，然后发布到任何流行的 Linux 机器上，也可以实现虚拟化。
## docker安装（CentOS）
-	查看CentOS系统内核版本，需要高于3.10
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181224222736674.png)
如果内核版本过低，可以参考此[链接](https://blog.csdn.net/kikajack/article/details/79396793)进行升级。
-	安装docker
	-	移除旧版本
	```bash
	$ sudo yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-selinux \
                  docker-engine-selinux \
                  docker-engine
	```
	-	安装必要的系统工具
	```bash
	$ sudo yum install -y yum-utils device-mapper-persistent-data lvm2
	```
	-	添加软件源信息
	```bash
	$ sudo yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
	```
	-	更新yum缓存
	```bash
	$ sudo yum makecache fast
	```
	-	安装Docker-ce
	```bash
	$ sudo yum -y install docker-ce
	```
	-	启动Docker后台服务
	```bash
	$ sudo systemctl start docker
	```
	-	运行hello-word
	由于本地没有hello-world镜像，所以从仓库中下载该镜像并在容器中运行。![在这里插入图片描述](https://img-blog.csdnimg.cn/20181224223829426.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
	-	镜像加速
	在`/etc/docker/daemon.json`文件(如果没有，自行创建)中添加：
	```json
	{
	    "registry-mirrors": ["http://hub-mirror.c.163.com"]
	}
	```
## docker基本操作
-	查看镜像`$ docker images`
-	查看容器`$ docker ps`，如果想要查看已经关闭的镜像可以加上`-a`参数
-	删除镜像`$ docker rmi <image-name/id>`，强行删除加上`-f`参数
-	删除容器`$ docker rm <container-name/id>`，强行删除加上`-f`参数
-	构建镜像`$ docker build -t <image-name> .`，`.`表示当前路径，可以使用具体路径代替
-	运行容器`$ docker run -p <本机端口>:<docker容器内部端口> --name <container-name> -d <image-name>`；-d参数表示在后台运行容器，并返回容器ID
-	交互模式运行容器`$ docker exec -it <container-name> /bin/bash`；`-i`以交互模式运行容器、`-t`为容器重新分配一个伪输入终端
-	标记本地镜像，将其归入某一仓库`$ docker tag <image-name[:tag]> <registryhost/username/image-name[:tag]>`
-	查看docker日志`$ docker logs <container-name>`
## docker实战（构建client前端镜像）

**编写Dockerfile**
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

**构建镜像**
```bash
$ docker build -t swapi-front-end .
```
![step1-5-taged](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/step1-5-taged.png)

![step6-9-taged](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/step6-9-taged.png)

**查看镜像**
```bash
$ docker images
```

**运行镜像**
```bash
$ docker run -d -p 8080:80 --name swapi-spa -e SERVER_ADDR='http://192.168.186.181:9090' swapi-front-end
```
![docker-run](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/docker-run.png)

**查看当前正在运行的容器**
```bash
$ docker ps
```
![docker-ps](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/docker-ps.png)

**显示容器的日志**
```bash
$ docker logs swapi-spa
```
![docker-logs](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/docker-logs.png)


**容器停止运行**
```bash
$ docker stop swapi-spa
```
![docker-stop](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/docker-stop.png)

**再次启动容器**
```bash
$ docker start swapi-spa
```
![docker-start](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/docker-start.png)

**删除容器**
```bash
$ docker rm swapi-spa
```
![docker-rm](https://github.com/Go-GraphQL-Group/front-end/blob/master/documents/img/docker-rm.png)

## docker实战（构建MySQL数据库镜像）
-	**首先，准备好数据库所需的数据（sql文件）**
```sql
-- init.sql --
drop database if exists starwars;
create DATABASE starwars;
use starwars;
create table if not exists people (
	_id INT NOT NULL AUTO_INCREMENT,
    ID char(100),
    Name char(100),
    Heigth char(100),
    Mass char(100),
    Hair_color char(100),
    Skin_color char(100),
    Eye_color char(100),
    Birth_year char(100),
    Gender char(100),
    Homeworld char(100),
    Films char(200),
    Species char(200),
    Vehicles char(200),
    Starships char(200),
    primary key(_id)
); 
……
……
```
-	**之后，书写`Dockerfile`文件**
以mysql5.7为基础镜像构建我们所需要的镜像，首先设置`MYSQL_ALLOW_EMPTY_PASSWORD`便于我们对数据库进行数据导入的操作，然后将文件(相对/绝对路径)拷贝到容器中，最后执行脚本写入数据。
```dockerfile
	FROM mysql:5.7
	# no password
	ENV MYSQL_ALLOW_EMPTY_PASSWORD yes
	
	# put the file to container
	COPY setup.sh /mysql/setup.sh
	COPY data/mysql/priviledges.sql /mysql/priviledges.sql
	COPY sql/init.sql /mysql/init.sql
	# command
	CMD ["sh", "/mysql/setup.sh"]
```
具体Dockerfile语法参见[链接](https://docs.docker.com/v17.09/engine/reference/builder/#parser-directives)
- **编写setup.sh脚本**
```sh
#!/bin/bash
set -e

#查看mysql服务的状态，方便调试，这条语句可以删除
echo `service mysql status`

echo '1.启动mysql....'
#启动mysql
service mysql start
sleep 3
echo `service mysql status`

echo '2.开始导入数据....'
#导入数据
mysql < /mysql/init.sql
echo '3.导入数据完毕....'

sleep 3
echo `service mysql status`

#重新设置mysql密码
echo '4.开始修改密码....'
mysql < /mysql/priviledges.sql
echo '5.修改密码完毕....'

#sleep 3
echo `service mysql status`
echo 'mysql容器启动完毕,且数据导入成功'

tail -f /dev/null
```
-	**MySQL权限设置priviledges.sql**
```sql
use mysql;
select host, user from user;
-- 因为mysql版本是5.7，因此新建用户为如下命令：
create user starwars identified by 'starwars';
-- 将starwars数据库的权限授权给创建的starwars用户，密码为starwars：
grant all on starwars.* to starwars@'%' identified by 'starwars' with grant option;
-- 这一条命令一定要有：
flush privileges;
```
-	**创建镜像**
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181224232513645.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181224232729684.png)
由于未使用`--name`参数指定容器名称，所以随机生成了一个名字`brave_agnesi`
-	**启动容器**
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181224232821531.png)
使用`docker logs`命令查看日志：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181224233104944.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
-	验证数据库镜像是否真正拥有数据
1.进入容器
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181224233307864.png)
2.登录mysql，使用我们在`priviledges.sql`文件中的用户名密码
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181224233429793.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
尝试查询数据：`select * from film;`
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181224233541548.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
由上图可知，我们成功插入了数据。
## docker实战（构建server后台镜像）
-	**首先，定义Dockerfile**
```dockerfile
#源镜像
FROM golang:latest
#设置工作目录
WORKDIR $GOPATH/src/github.com/Go-GraphQL-Group/GraphQL-Service
#将服务器的go工程代码加入到docker容器中
ADD . $GOPATH/src/github.com/Go-GraphQL-Group/GraphQL-Service
#go构建可执行文件
RUN go get github.com/Go-GraphQL-Group/GraphQL-Service
RUN go build .
# 设置 PORT 环境变量
ENV PORT 9090
#暴露端口
EXPOSE 9090
#最终运行docker的命令
ENTRYPOINT  ["./GraphQL-Service"]
```
-	**之后，生成镜像**
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225000134495.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225000224562.png)
-	**启动容器**
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225000423555.png)
-	**测试访问后台服务**
测试所用用户名密码`admin-password`为代码内置，并没有从数据库中获取，所以可成功登录。
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225000632222.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
## docker进阶
此时，我们已经得到了三个镜像（starwars_client(重命名得到)、starwars_db、starwars_server），每个镜像都可以在各自的容器中正常运行，但是三个容器相互独立，前端容器无法获取后端服务，后端也无法从数据库中获取服务。例如：
- 前端容器想要获取后台服务（不使用SERVER_ADDR指定服务器地址）：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225001729285.png)
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225001714178.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
- 后台想要从数据库获取数据：测试此部分，我们前端可以不使用docker容器启动，而是使用`npm start`在项目根目录下启动，此时前端可以访问后台服务，端口为`8080`。
我们可以成功登陆：，因为用户名密码只有一组（`admin-password`）直接写入到后台代码，并没有存储到数据库中。
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225005103438.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
但是，我们后台无法获取mysql数据库中的数据，返回数据为空：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225011015452.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)

此时，我们便需要运用到容器间通信的内容。
### docker容器间通信
同主机下容器间通信，此篇文章讲解的比较详细，[传送门](https://blog.csdn.net/dream_broken/article/details/52414560)。关于容器间通信的方法，可以参见此文章，[传送门](https://www.cnblogs.com/CloudMan6/p/7096731.html)。我们小组在实现的过程中，使用`--link`参数完成通信。
**starwars_server与starwars_db通信具体实现如下：**
- 首先，删除之前创建的容器
![在这里插入图片描述](https://img-blog.csdnimg.cn/2018122500374138.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
- 之后，使用`starwars_db`mysql镜像启动容器，并命名为`starwars_db`：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225004006314.png)
- 然后，修改后台代码中链接数据库host配置：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225004123633.png)
- 然后重新编译生成新的镜像`starwars_server`，然后使用该镜像启动容器，命名为`starwars_server`，并指定`--link`参数：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225004353408.png)
其中`--link`参数格式为`<container name>:<container alias>`。而我们修改的配置文件`host`即为`container alias`
- 后台服务尝试获取数据（依然使用`npm start`启动的前端）,由下图可知，成功获取数据。
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225005624394.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)

**starwars_client与starwars_server通信具体实现如下：**
- 删除之前的前端容器和镜像
- 修改前端Dockerfile
将`SERVER_ADDR`修改为`http://starwars_server:9090`
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225005833776.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
- 重新生成镜像，并在运行`starwars_client`容器时加上`--link`指定关联的容器（`starwars_server:starwars_server`）：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225010052743.png)
- 测试`starwars_client`访问后台服务，`80端口`
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225010256537.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)

此时，我们已经完成了三个容器之间的通信，并实现了该项目的功能。但是，这样的操作过于繁琐，而且当容器过多时会更加难以部署。所以，我们需要用到`docker-compose`来部署应用。
### docker-compose部署应用
Compose 项目是 Docker 官方的开源项目，负责实现对 Docker 容器集群的快速编排。Compose 定位是 「定义和运行多个 Docker 容器的应用（Defining and running multi-container Docker applications）」。
Compose 恰好满足了这样的需求。它允许用户通过一个单独的 docker-compose.yml 模板文件（YAML 格式）来定义一组相关联的应用容器为一个项目（project）。

Compose 中有两个重要的概念：

服务 (service)：一个应用的容器，实际上可以包括若干运行相同镜像的容器实例。

项目 (project)：由一组关联的应用容器组成的一个完整业务单元，在 docker-compose.yml 文件中定义。
**docker-compose安装**

- **首先，尝试编写docker-compose.yml**
```yml
version: "3.3"
services:
  database:
    image: starwars_db:latest
    container_name: starwars_db
    restart: always
    ports: 
      - "3306:3306"
  server:
    depends_on:
      - database
    image: starwars_server:latest
    container_name: starwars_server
    restart: always
    ports: 
      - "9090:9090"
    external_links:
      - starwars_db:starwars_db
  client:
    depends_on:
      - database
      - server
    image: starwars_client:latest
    container_name: starwars_client
    restart: always
    ports:
      - "80:80"
    external_links:
      - starwars_server:starwars_server
```
在上述配置中，我们定义了三个服务：database、server、client，分别对应于我们所创建的三个镜像：starwars_db，starwars_server，starwars_client。
`depends_on：`
由于server服务依赖于database，所以我们使用`depends_on`来管理依赖，表示server在database启动后在启动。client依赖也类似处理。
`container_name：`
此属性即为我们创建容器时的容器名称，用于`external_links`标识。
`external_links：`
此属性即为我们之前使用的`--link`，管理容器间通信。
- 删除之前所创建的容器
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225103551347.png)
- 然后在`docker-compose.yml`所在目录运行`docker-compose up -d`。
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225103530674.png)
- 测试应用
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225103818460.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
大功告成。但是，我们依然还有一个问题，所有的镜像我们都存储在本地，那么如何让其他人也使用我们的镜像构建服务呢？
### Docker-Hub使用
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225104026686.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
Docker Hub是一个镜像仓库，方便我们构建管理镜像。
- **注册\登录账号**
可参考[链接](https://blog.csdn.net/liuyh73/article/details/84181436)
- **docker push**
首先需要将我们的镜像重命名，可以加上版本号：`docker tag`
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225104711602.png)
之前的镜像可删可不删，然后使用`docker push liuyh73/starwars_XXX`推送到自己的`Docker Hub`仓库中。如果提示没有登录，则可以使用`docker login`进行登录。
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225104958291.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
- **修改`docker-compose.yml`**
```yml
# 修改image
version: "3.3"
services:
  database:
    image: liuyh73/starwars_db:latest
    container_name: starwars_db
    restart: always
    ports: 
      - "3306:3306"
  server:
    depends_on:
      - database
    image: liuyh73/starwars_server:latest
    container_name: starwars_server
    restart: always
    ports: 
      - "9090:9090"
    external_links:
      - starwars_db:starwars_db
  client:
    depends_on:
      - database
      - server
    image: liuyh73/starwars_client:latest
    container_name: starwars_client
    restart: always
    ports:
      - "80:80"
    external_links:
      - starwars_server:starwars_server
```
- 删除本地镜像`liuyh73/starwars_XXX`，
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225105757207.png)
- 然后执行`docker-compose up -d`进行安装运行。
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225112210123.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
![在这里插入图片描述](https://img-blog.csdnimg.cn/2018122511242031.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)
- 测试
![在这里插入图片描述](https://img-blog.csdnimg.cn/20181225112612900.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2xpdXloNzM=,size_16,color_FFFFFF,t_70)

完整代码详见[Github](https://github.com/Go-GraphQL-Group)