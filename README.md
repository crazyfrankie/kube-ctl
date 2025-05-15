# Kube-ctl
Kube-ctl 是 kubernetes 的一个可视化系统. 提供给用户一系列方式操作 kubernetes 集群

演示地址: http://www.crazyfrank.top:8080

本项目是一个严格意义上的 CRUD 项目, 但侧重点在于学习 kubernetes 整个系统的组成. 

阅读[官方文档](https://kubernetes.io/docs)、阅读 [k8s](https://github.com/kubernetes/kubernetes) 源码是本项目的核心

## 技术栈
Gin + [kubernetes/client-go](https://github.com/kubernetes/client-go)

## 简介
- [x] Namespace 查询
- [x] Pod 创建、更新、删除、查询（详情和列表）
- [x] Node 列表、详情、Node 所包含的 Pods、标签更新、污点更新
- [x] ConfigMap 创建、更新、删除、查询（详情和列表）
- [x] Secret 创建、更新、删除、查询（详情和列表）
- [x] PersistentVolume 创建、查询、删除
- [x] PersistentVolumeClaim 创建、查询、删除
- [x] StorageClass 创建、查询、删除
- [x] Pod 支持多种存储卷: EmptyDir、ConfigMap、Secret、HostPath、DownwardAPI、PersistentVolume 
- [x] Service 创建、更新、删除、查询（详情和列表）
- [x] Ingress 创建、更新、删除、查询（详情和列表）
  - 注： Ingress controller 在本系统中作为了系统内置资源，如果在使用 Ingress 之前没有编写 IngressClass 资源的配置文件去创建 Ingress Controller, 请先创建
- [x] IngressRoute 创建、更新、删除、查询
- [x] Deployment 创建、更新、删除、查询（详情和列表）
- [x] DaemonSet 创建、更新、删除、查询（详情和列表）
- [x] StatefulSet 创建、更新、删除、查询（详情和列表）
- [x] Job 创建、更新、删除、查询（详情和列表）
- [x] CronJob 创建、更新、删除、查询（详情和列表）
- [x] ServiceAccount 创建、更新、删除、查询（详情和列表）
- [x] Role | ClusterRole 创建、更新、删除、查询（详情和列表）
- [x] RoleBinding | ClusterRoleBinding 创建、更新、删除、查询（详情和列表）

## 启动
### v1:
项目在开发和测试阶段, 均以集群外的方式访问, 需按照下面的步骤进行启动
1. 在项目根目录下新建 `.kube` 文件夹, 将集群中 `control-plane` 机器上的 `~/.kube/config` 复制到文件夹下
2. 参照 `conf/test/example.yaml` 修改对应配置
### v2:
服务内置发现机制
- 动态判断是集群外访问还是集群内访问
- 若为集群外, 走默认路径, 用户仍需参照 v1 进行配置
- 若为集群内, 走集群内访问机制

## 项目前端
[kube-ctl-web](https://github.com/crazyfrankie/kube-ctl-web)
