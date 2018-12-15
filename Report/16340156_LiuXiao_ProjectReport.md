# 服务计算 - 5 | GraphQL简单web服务与客户端开发
---
## 概述
利用 web 客户端调用远端服务是服务开发本实验的重要内容。其中，要点建立 API First 的开发理念，实现前后端分离，使得团队协作变得更有效率。
### 任务目标
1. 选择合适的 API 风格，实现从接口或资源（领域）建模，到 API 设计的过程
2. 使用 API 工具，编制 API 描述文件，编译生成服务器、客户端原型
3. 使用 Github 建立一个组织，通过 API 文档，实现 客户端项目 与 RESTful 服务项目同步开发
4. 使用 API 设计工具提供 Mock 服务，两个团队独立测试 API
5. 使用 travis 测试相关模块
### 开发环境选取
1. API使用[GraphQL](http://graphql.cn/learn/)规范进行设计
2. 客户端使用[Vue框架](https://cn.vuejs.org/index.html)
3. 服务器使用GraphQL官方提供的生成基于 graphql 的服务器的库[GQLGen](https://github.com/graphql-go/graphql)进行开发。
4. 数据库使用[BoltDB](https://github.com/boltdb/bolt)实现
### [GITHUB传送门](https://github.com/Go-GraphQL-Group)
---
## 项目实现
项目参照星球大战API [SWAPI](https://swapi.co/documentation)编写
### [API设计](https://github.com/Ernie1/GraphQL/blob/master/summary_16340286.md)
### [客户端实现](https://blog.csdn.net/linshk_ver18/article/details/85015778)
### [数据库实现](https://liu-yt.github.io/2018/12/14/%E6%9C%8D%E5%8A%A1%E8%AE%A1%E7%AE%97-6-BoltDB%E5%AD%A6%E4%B9%A0%E4%B8%8E%E7%AE%80%E5%8D%95%E5%89%96%E6%9E%90/?tdsourcetag=s_pctim_aiomsg)
### 服务器实现
服务器实现其实并不难，但是理解GraphQL和GQLGEN这两个预备工作比较困难，需要大量阅读。
#### GraphQL
GraphQL 是一个用于 API 的查询语言，是一个使用基于类型系统来执行查询的服务端运行时（类型系统由你的数据定义）。GraphQL 并没有和任何特定数据库或者存储引擎绑定，而是依靠你现有的代码和数据支撑。
与Restful相比，GraphQL不会由复杂的URL，请求的Json按照规范被放在数据中。由于有完备的规范，使用GraphQL构建服务器时不需要自行对每个请求进行解析，可以使用现成的框架，如[GQLGen](https://github.com/graphql-go/graphql)，按规范编写Schema后即可生成相应的解析函数，最终只需要自己编写resolve中的查询函数即可。无需对每个数据规定复杂的URL，大大简化了开发流程。
> [GraphQL官网](http://graphql.cn/learn/)
> [GraphQL核心概念](https://blog.csdn.net/liuyh73/article/details/85010148?tdsourcetag=s_pctim_aiomsg)
> [GQLGEN样例](https://gqlgen.com/getting-started/)

GraphQL只是一个规范，具体使用时必须自行实现解析。这里可以用各种开源库来简化开发流程。
#### GQLGEN
使用GQLGEN首先应该编写`schema.graphql`文件，其中按照GraphQL规范，定义了所有结构的内容，以及查询的方法，在这个项目中没有用到客户端更新数据，所以没有使用Mutation。GQLGEN这个组件会根据schema生成对应的请求路径解析和请求中GraphQL规则的查询的解析，并且使用者只需要实现每个请求的处理函数即可，简化了开发流程。
* `type Query`中定义了所有的查询查询方法，在这个类型中的查询函数会被GQLGEN自动实现解析，并在`resolver.go`文件中新建空白查询函数，而我们的任务就是编写该文件中的函数，返回对应的数据。
 ```graphql
 """
 The query root, from which multiple types of requests can be made.
 """
 type Query {
     """
     Look up a specific people by its ID.
     """
     people(
         """
         The ID of the entity.
         """
         id: ID!
     ): People

     """
     Look up a specific film by its ID.
     """
     film(
         """
         The ID of the entity.
         """
         id: ID!
     ): Film
     
     """
     Look up a specific starship by its ID.
     """
     starship(
         """
         The ID of the entity.
         """
         id: ID!
     ): Starship
     
     """
     Look up a specific vehicle by its ID.
     """
     vehicle(
         """
         The ID of the entity.
         """
         id: ID!
     ): Vehicle
     
     """
     Look up a specific specie by its ID.
     """
     specie(
         """
         The ID of the entity.
         """
         id: ID!
     ): Specie
     
     """
     Look up a specific planet by its ID.
     """
     planet(
         """
         The ID of the entity.
         """
         id: ID!
     ): Planet

     """
     Browse people entities.
     """
     peoples (
         """
         The number of entities in the connection.
         """
         first: Int

         """
         The connection follows by.
         """
         after: ID
     ): PeopleConnection!

     """
     Browse film entities.
     """
     films (
         """
         The number of entities in the connection.
         """
         first: Int

         """
         The connection follows by.
         """
         after: ID
     ): FilmConnection!

     """
     Browse starship entities.
     """
     starships (
         """
         The number of entities in the connection.
         """
         first: Int

         """
         The connection follows by.
         """
         after: ID
     ): StarshipConnection!

     """
     Browse vehicle entities.
     """
     vehicles (
         """
         The number of entities in the connection.
         """
         first: Int

         """
         The connection follows by.
         """
         after: ID
     ): VehicleConnection!

     """
     Browse specie entities.
     """
     species (
         """
         The number of entities in the connection.
         """
         first: Int

         """
         The connection follows by.
         """
         after: ID
     ): SpecieConnection!

     """
     Browse planet entities.
     """
     planets (
         """
         The number of entities in the connection.
         """
         first: Int

         """
         The connection follows by.
         """
         after: ID
     ): PlanetConnection!

     """
     Search for people entities matching the given query.
     """
     peopleSearch (
         """
         The search field for name, in Lucene search syntax.
         """
         search: String!
         
         """
         The number of entities in the connection.
         """
         first: Int

         """
         The connection follows by.
         """
         after: ID
     ): PeopleConnection

     """
     Search for film entities matching the given query.
     """
     filmsSearch (
         """
         The search field for title, in Lucene search syntax.
         """
         search: String!
         
         """
         The number of entities in the connection.
         """
         first: Int

         """
         The connection follows by.
         """
         after: ID
     ): FilmConnection

     """
     Search for starship entities matching the given query.
     """
     starshipsSearch (
         """
         The search field for name or model, in Lucene search syntax.
         """
         search: String!
         
         """
         The number of entities in the connection.
         """
         first: Int

         """
         The connection follows by.
         """
         after: ID
     ): StarshipConnection

     """
     Search for vehicle entities matching the given query.
     """
     vehiclesSearch (
         """
         The search field for name or model, in Lucene search syntax.
         """
         search: String!
         
         """
         The number of entities in the connection.
         """
         first: Int

         """
         The connection follows by.
         """
         after: ID
     ): VehicleConnection

     """
     Search for specie entities matching the given query.
     """
     speciesSearch (
         """
         The search field for name, in Lucene search syntax.
         """
         search: String!
         
         """
         The number of entities in the connection.
         """
         first: Int

         """
         The connection follows by.
         """
         after: ID
     ): SpecieConnection

     """
     Search for planet entities matching the given query.
     """
     planetsSearch (
         """
         The search field for name, in Lucene search syntax.
         """
         search: String!
         
         """
         The number of entities in the connection.
         """
         first: Int

         """
         The connection follows by.
         """
         after: ID
     ): PlanetConnection
 }

 ```
* 其它部分按照GraphQL规范编写即可，具体可以查看项目中的`schema.graphql`文件。这里的`schema.graphql`有些臃肿，可以通过实现共同属性的`interface`来减少定义的工作量。
* 具体设计参阅API文档
#### 解析函数的编写
1. 对于普通的通过ID查询的函数，直接通过数据库提供的方法查询对应ID的对象。
    ```go
    func (r *queryResolver) People(ctx context.Context, id string) (*People, error) {
        err, people := GetPeopleByID(id, nil)
        checkErr(err)
        return people, err
    }
    ```
2. **分页查询**则需要解析需要的元素数量，起始位置即`after`游标在数据库中的位置，是否有前后页及当前页开始和结束位置元素的游标，用于客户端在需要的时候获取前后页。
   ```go
    func (r *queryResolver) Peoples(ctx context.Context, first *int, after *string) (PeopleConnection, error) {
        from := -1
        if after != nil {
            b, err := base64.StdEncoding.DecodeString(*after)
            if err != nil {
                return PeopleConnection{}, err
            }
            i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
            if err != nil {
                return PeopleConnection{}, err
            }
            from = i
        }
        count := 0
        startID := ""
        hasPreviousPage := true
        hasNextPage := true
        // 获取edges
        edges := []PeopleEdge{}
        db, err := bolt.Open("./data/data.db", 0600, nil)
        CheckErr(err)
        defer db.Close()
        db.View(func(tx *bolt.Tx) error {
            c := tx.Bucket([]byte(peopleBucket)).Cursor()

            // 判断是否还有前向页
            k, v := c.First()
            if from == -1 || strconv.Itoa(from) == string(k) {
                startID = string(k)
                hasPreviousPage = false
            }

            if from == -1 {
                for k, _ := c.First(); k != nil; k, _ = c.Next() {
                    _, people := GetPeopleByID(string(k), db)
                    edges = append(edges, PeopleEdge{
                        Node:   people,
                        Cursor: encodeCursor(string(k)),
                    })
                    count++
                    if count == *first {
                        break
                    }
                }
            } else {
                for k, _ := c.First(); k != nil; k, _ = c.Next() {
                    if strconv.Itoa(from) == string(k) {
                        k, _ = c.Next()
                        startID = string(k)
                    }
                    if startID != "" {
                        _, people := GetPeopleByID(string(k), db)
                        edges = append(edges, PeopleEdge{
                            Node:   people,
                            Cursor: encodeCursor(string(k)),
                        })
                        count++
                        if count == *first {
                            break
                        }
                    }
                }
            }

            k, v = c.Next()
            if k == nil && v == nil {
                hasNextPage = false
            }
            return nil
        })
        if count == 0 {
            return PeopleConnection{}, nil
        }
        // 获取pageInfo
        pageInfo := PageInfo{
            HasPreviousPage: hasPreviousPage,
            HasNextPage:     hasNextPage,
            StartCursor:     encodeCursor(startID),
            EndCursor:       encodeCursor(edges[count-1].Node.ID),
        }

        return PeopleConnection{
            PageInfo:   pageInfo,
            Edges:      edges,
            TotalCount: count,
        }, nil
    }
   ```
3. 其次是**基于相关字段的分页查询**，与普通分页查询类似，只是多了一个查询字段的字符串来限定，获取对应的页。
    ```go
    func (r *queryResolver) PeopleSearch(ctx context.Context, search string, first *int, after *string) (*PeopleConnection, error) {

        if strings.HasPrefix(search, "Name:") {
            search = strings.TrimPrefix(search, "Name:")
        } else {
            return &PeopleConnection{}, errors.New("Search content must be ' Name:<People's Name you want to get> ' ")
        }
        from := -1
        if after != nil {
            b, err := base64.StdEncoding.DecodeString(*after)
            if err != nil {
                return &PeopleConnection{}, err
            }
            i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
            if err != nil {
                return &PeopleConnection{}, err
            }
            from = i
        }
        count := 0
        hasPreviousPage := false
        hasNextPage := false
        // 获取edges
        edges := []PeopleEdge{}
        db, err := bolt.Open("./data/data.db", 0600, nil)
        CheckErr(err)
        defer db.Close()
        db.View(func(tx *bolt.Tx) error {
            c := tx.Bucket([]byte(peopleBucket)).Cursor()
            k, _ := c.First()
            // 判断是否还有前向页
            if from != -1 {
                for k != nil {
                    _, people := GetPeopleByID(string(k), db)
                    if people.Name == search {
                        hasPreviousPage = true
                    }
                    if strconv.Itoa(from) == string(k) {
                        k, _ = c.Next()
                        break
                    }
                    k, _ = c.Next()
                }
            }

            // 添加edge
            for k != nil {
                _, people := GetPeopleByID(string(k), db)
                if people.Name == search {
                    edges = append(edges, PeopleEdge{
                        Node:   people,
                        Cursor: encodeCursor(string(k)),
                    })
                    count++
                }
                k, _ = c.Next()
                if first != nil && count == *first {
                    break
                }
            }

            // 判断是否还有后向页
            for k != nil {
                _, people := GetPeopleByID(string(k), db)
                if people.Name == search {
                    hasNextPage = true
                    break
                }
                k, _ = c.Next()
            }
            return nil
        })
        if count == 0 {
            return &PeopleConnection{}, nil
        }
        // 获取pageInfo
        pageInfo := PageInfo{
            StartCursor:     encodeCursor(edges[0].Node.ID),
            EndCursor:       encodeCursor(edges[count-1].Node.ID),
            HasPreviousPage: hasPreviousPage,
            HasNextPage:     hasNextPage,
        }

        return &PeopleConnection{
            PageInfo:   pageInfo,
            Edges:      edges,
            TotalCount: count,
        }, nil

    }
    ```
其他的查询函数实现和上述People方法的实现基本相同。

#### JWT 产生 token 实现用户认证
* **基于 Token 的身份验证方法** 
* 使用基于 Token 的身份验证方法，在服务端不需要存储用户的登录记录。大概的流程是这样的： 
  1. 客户端使用用户名跟密码请求登录 
  2. 服务端收到请求，去验证用户名与密码 
  3. 验证成功后，服务端会签发一个 Token，再把这个 Token 发送给客户端 
  4. 客户端收到 Token 以后可以把它存储起来，比如放在 Cookie 里或者 Local Storage 中
  5. 客户端每次向服务端请求资源的时候需要带着服务端签发的 Token 
  6. 服务端收到请求，然后去验证客户端请求里面带着的 Token，如果验证成功，就向客户端返回请求的数据

* **JSON Web Token**
根据官网的定义，JSON Web Token（以下简称 JWT）是一套开放的标准（RFC 7519），它定义了一套简洁（compact）且 URL 安全（URL-safe）的方案，以安全地在客户端和服务器之间传输 JSON 格式的信息。
  * 优点
    1. 体积小（一串字符串）。因而传输速度快
    2. 传输方式多样。可以通过 HTTP 头部（推荐）/URL/POST 参数等方式传输
    3. 严谨的结构化。它自身（在 payload 中）就包含了所有与用户相关的验证消息，如用户可访问路由、访问有效期等信息，服务器无需再去连接数据库验证信息的有效性，并且 payload 支持应用定制
    4. 支持跨域验证，多应用于单点登录
        > 单点登录（Single Sign On）：在多个应用系统中，用户只需登陆一次，就可以访问所有相互信任的应用
  * 实现
    1. 这里没有实现账号密码的数据库，只是在路由处理`router.go`中比对固定密码，实现登出。
    2. 每次请求到达服务器时，服务器中间件判断是否是登陆请求。如果不是登陆请求，则获取保存在请求头部的Token进行比对，相同则调用正常的HttpHandler，错误则返回错误回复。
    ```go
    func TokenMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.RequestURI[1:] != "login" {
                /*
                    // token位于Authorization中，用此方法
                    token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
                        return []byte(SecretKey), nil
                    })
                */
                tokenStr := ""
                for k, v := range r.Header {
                    if strings.ToUpper(k) == TokenName {
                        tokenStr = v[0]
                        break
                    }
                }
                validToken := false
                for _, token := range tokens {
                    if token.SW_TOKEN == tokenStr {
                        validToken = true
                    }
                }
                if validToken {
                    ctx := context.WithValue(r.Context(), TokenName, tokenStr)
                    next.ServeHTTP(w, r.WithContext(ctx))
                } else {
                    w.WriteHeader(http.StatusUnauthorized)
                    w.Write([]byte("Unauthorized access to this resource"))
                    //fmt.Fprint(w, "Unauthorized access to this resource")
                }
            } else {
                next.ServeHTTP(w, r)
            }
        })
    }
    ```
    3. `jwt.go`中实现了创建Token和比对Token的方法，由于只有一个Token，所以直接比对即可。
    ```go
    var tokens []Token

    const TokenName = "SW-TOKEN"
    const Issuer = "Go-GraphQL-Group"
    const SecretKey = "StarWars"

    type Token struct {
        SW_TOKEN string `json:"SW-TOKEN"`
    }

    type jwtCustomClaims struct {
        jwt.StandardClaims

        Admin bool `json:"admin"`
    }

    func CreateToken(secretKey []byte, issuer string, isAdmin bool) (token Token, err error) {
        claims := &jwtCustomClaims{
            jwt.StandardClaims{
                ExpiresAt: int64(time.Now().Add(time.Hour * 1).Unix()),
                Issuer:    issuer,
            },
            isAdmin,
        }

        tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
        token = Token{
            tokenStr,
        }
        return
    }

    func ParseToken(tokenStr string, secretKey []byte) (claims jwt.Claims, err error) {
        var token *jwt.Token
        token, err = jwt.Parse(tokenStr, func(*jwt.Token) (interface{}, error) {
            return secretKey, nil
        })
        claims = token.Claims
        return
    }
    ```