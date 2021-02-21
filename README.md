# Janus

Janus 是罗马神话中开端、大门、选择、过渡、时间、对偶、道路、门框和结尾的神。他通常被描述成有前后两张面孔，展望着过去和未来；也有描述成四方四个面孔的。(引自维基百科)

Janus类似于gingate的网关服务，支持自定义修改请求参数、接口请求次数限制等功能的 proxy 服务。

   ## 结构
   
   ```
   
   │———— app                           
   │      │—— main                     # 目录下有项目入口
   │
   |
   │———— controllers
   │       │———— middlewares           # 中间件
   │       │      │── client
   │       │      │── log
   │       │      │── reqid  
   │       │
   │       │──—— proxy                  # 代理层 
   │       │      │—— proxy             # 代理层的主要逻辑
   │       │
   │
   │———— config  # 有关 janus 和 proxy 的配置  
   │                    
   │     
   
   ```
  
## 开发

janus 使用 `go mod` 管理依赖，把代码放在任何位置,

```
git@github.com:ViolaTangxl/janus.git

make
```
即可


## janus 配置

在 ``/config/proxy_XXXX.yml`` 中， 可以通过简单配置，转发前端请求到对应后端服务
```
# 以 name 为 github 的 proxy 为例
 
# 默认proxy请求路由前添加 /api/proxy
# name 必填 该组 proxy 名字 
 
proxy_entries: 
    - name: "github" 
      # target 必填 代理后端地址
      target: "https://github.com"
      matches:
          # path必填 具体 api 路由 支持 /* 、* 匹配
        - path: "/ViolaTangxl/janus/show_partial"
          # method必填  请求类型 支持*匹配
          method: GET 
          params:
            # rename:  选填，重命名，如果 rename 为空，则使用session_key作为参数名，如果 session_key 为空，则使用 rename 字段作为额外附加的自定义 param
            # custom_value: 选填，自定义的参数值，如果session中没有存某个值，支持自定义, session和custom_value均有值时，sessionValue 优先于 custom_value
            # location: 必填，指定参数位置,目前支持: url_param, body, header, url_path
            - location: "url_param"
              rename: "partial"
              custom_value: "tree/recently_touched_branches_list" 
            
 
# 通过上述配置，可以通过请求 
localhost:11009/api/proxy/github/ViolaTangxl/janus/show_partial 
就可以正常访问 github 的 https://github.com/ViolaTangxl/janus/show_partial?partial=tree%2Frecently_touched_branches_list
接口，request param  中默认添加 partial: tree/recently_touched_branches_list

``` 


