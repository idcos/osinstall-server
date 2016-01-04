### 数据中心表-data_center
字段名 | 数据类型 | 备注
----|------|----
id| string| 数据中心ID
name| string| 数据中心名称
showName| string| 数据中心显示名


### 应用表-app
字段名 | 数据类型 | 备注
----|------|----
id| string| 应用ID
name| string| 应用名称
showName| string| 应用显示名


### 部署单元表-unit
字段名 | 数据类型 | 备注
----|------|----
id| string| 部署单元ID
name| string| 部署单元名称
showName| string| 部署单元显示名
confLastChanged| date| 最近一次配置修改时间
confLastApplied| date| 最近一次配置实施时间
execUsers| []string| 可用执行帐号列表
nodeIDs| []string| 主机节点ID列表
appID| string| 应用ID
dataCenterID| string| 数据中心ID

### 主机节点表-node
字段名 | 数据类型 | 备注
----|------|----
id| string| 主机节点ID
name| string| 主机节点名称
showName| string| 主机节点显示名
confLastChanged| date| 最近一次配置修改时间
confLastApplied| date| 最近一次配置实施时间
execUsers| []string| 可用执行帐号列表
stat| enum| 主机节点状态。可选状态包括：`online`、`offline`

### 待执行任务表-task_todos
字段名 | 数据类型 | 备注
----|------|----
id| string| 任务ID
userID| string| 操作人ID
host| string| 操作的目标主机
execUser| string| 命令在目标主机上的执行帐号
cmd| string| 命令内容
startAt| date| 命令执行开始时间
endAt| date| 命令执行结束时间



### 历史命令表-cmd_history
字段名 | 数据类型 | 备注
----|------|----
id| string| 历史命令ID
userID| string| 操作人ID
host| string| 操作的目标主机
execUser| string| 命令在目标主机上的执行帐号
cmd| string| 命令内容
stdout| string| 命令执行后的标准输出
stderr| string| 命令执行后的标准错误输出
startAt| date| 命令执行开始时间
endAt| date| 命令执行结束时间


### 配置信息表-conf
字段名 | 数据类型 | 备注
----|------|----
id| string| 配置记录ID
name| string| 配置信息名称
value| object| 配置信息值
scope| enum| 配置信息的作用范围。可选值包括：`globle`、`role`、`group`、`user`。
scopeID| string| 


### SSH信息表-plugin_ssh
字段名 | 数据类型 | 备注
----|------|----
id| string| 主键ID
nodeID| string| 所属的主机节点ID
host| string| ssh登录主机
user| string| ssh登录用户名
authType| string| 鉴权类型。可选值包括：`passoword`、`key`。
authKey| string| 鉴权所需的密码或者密钥。
port| int| ssh服务端口号

### MCO信息表-plugin_mco
字段名 | 数据类型 | 备注
----|------|----
id| string| 主键ID
dataCenterID| string| 所属数据中心
addrs| []string| mco消息中间件地址列表
nodesQueue| string| 
replyQueue| string| 
certFile| string| 
keyFile| string|
caFile| string| 
userName| string|
userPwd| string|
host| string|
signerKeyFile| string|
signerKeyPwd| string|
master| string|


### 部门表-dept
字段名 | 数据类型 | 备注
----|------|----
id| string| 部门ID
name| string| 部门名称
showName| string| 部门显示名
userIDs| []string| 用户ID列表


### 用户表-user
字段名 | 数据类型 | 备注
----|------|----
id| string| 用户ID
name| string| 用户名称
showName| string| 用户显示名
pwd| string| 用户登录密码
email| string| 电子邮箱
stat| enum| 状态。可选值包括：`enable`、`disable`。
deptID| string| 所属部门ID
groupIDs| []string| 所属用户组ID列表



### 用户组表-group
字段名 | 数据类型 | 备注
----|------|----
id| string| 用户组ID
name| string| 用户组名称
showName| string| 用户组显示名
userIDs| []string| 用户ID列表


### 角色表-role
字段名 | 数据类型 | 备注
----|------|----
id| string| 角色ID
name| string| 角色名称
showName| string| 角色显示名
