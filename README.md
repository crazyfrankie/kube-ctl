# Kube-ctl
Kube-ctl 是 kubernetes 的一个可视化系统. 提供给用户一系列方式操作 kubernetes 集群.

## 技术栈
Gin + [kubernetes/client-go](https://github.com/kubernetes/client-go) + GORM + MySQL

## 简介
- Namespace 查询
- Pod 创建、更新、删除、查询、详情
- Node 列表、详情、Node 所包含的 Pods、标签更新、污点更新
- ConfigMap 创建、更新、查询、删除
- TODO

## 项目前端
[kube-ctl-web](https://github.com/crazyfrankie/kube-ctl-web)