# Kubernetes kubectl 命令行总结


这篇文章总结了 Kubernetes 中 命令行的使用。

<!--more-->

{{< admonition info>}}
趁着今年黑色星期五的 Linux CNCF 社区报名有**三五折优惠**，我打算报考 CKAD 认证考试。在学习 Kubernetes 相关知识的同时，我也不忘总结一波有关 `kubectl` 的快捷命令。
{{< /admonition >}}

## 1 Pod 生命周期管理
```bash 
# 命令行创建 nginx pod
kubectl run nginx --image=nginx --restart=Never --dry-run=client -n mynamespace -o yaml > pod.yaml
# 获取所有命名空间的 pod
kubectl get po --all-namespaces
# 运行 busybox 的同时打印 hello world
kubectl run busybox --image=busybox -it --restart=Never -- echo 'hello world'
# 进入容器
kubectl exec -it podName -n nsName /bin/sh
# 列出 pod container 里面的 env variable
kubectl exec podName -- printenv
# 创建 pod 同时为 pod 添加标签
kubectl run nginx1 --image=nginx --restart=Never --labels=app=v1
# 列出所有添加标签的 pods
kubectl get po --show-labels
# 列出标有 app=v2 的 pods
kubectl get po -l app=v2
# 更改 nginx2 的标签为 app=v2
kubectl label po nginx2 app=v2 --overwrite
# 移除 nginx1 标有app的标签
kubectl label po nginx1 app-
```

## 2 暴露 service 或 deployment
```bash
# 输出为yaml文件（推荐）
kubectl expose deployment nginx --port=80 --type=NodePort --target-port=80 --name=web1 -o yaml > web1.yaml
# 暴露端口
kubectl expose deployment nginx -n bigdata --port=80 --type=NodePort
```

## 3 版本更新
```bash 
# 更新 nginx 版本
kubectl set image deployment/nginx nginx=nginx:1.15
# 滚动更新
kubectl rolling-update frontend --image=image:v2
# 扩缩容
kubectl scale deployment nginx --replicas=10
```

## 4 回滚
```bash
# 查看更新过程
kubectl rollout status deployment/nginx --namespace=nsName
# 如果更新成功, 返回值为0 
kubectl rollout status deployment nginx-deployment --watch=false | grep -ic waiting

# 查看变更历史版本信息
kubectl rollout history deployment/nginx
kubectl rollout history deployment/nginx --revision=3 --namespace=nsName

# 终止升级
kubectl rollout pause deployment/nginx --namespace=nsName

# 继续升级
kubectl rollout resume deployment/review-demo --namespace=nsName

# 回滚版本
kubectl rollout undo deployment/nginx --namespace=nsName
kubectl rollout undo deployment/nginx --to-revision=3  --namespace=nsName
```
## 5 探针 livenessProbe 和 readinessProbe
```bash 
kubectl run nginx --image=nginx --restart=Never --dry-run=client -o yaml > pod.yaml
# 然后编辑 pod.yaml 下的(containers -> livenessProbe -> exec -> command -> - 命令) 并保存
kubectl describe pod nginx | grep -i liveness # 检测是否 nginx 存活
```
## 6 Kubernetes configuration 配置
```bash 
# 创建 secret
kubectl create secret generic test-secret --from-literal='username=my-app' --from-literal='password=39528$vdg7Jb'
# 创建 configMap
echo -e "DB_URL=localhost:3306\nDB_USERNAME=postgres" > config.txt
kubectl create cm db-config --from-env-file=config.txt
# 创建 UID = 101
kubectl run nginx --image=nginx --restart=Never --dry-run=client -o yaml > pod.yaml
vi pod.yaml
# spec -> securityContext -> runAsUser -> 101

# requests 和 limits
kubectl run nginx --image=nginx --restart=Never --requests='cpu=100m,memory=256Mi' --limits='cpu=200m,memory=512Mi'
# 创建名字为 myuser 的 serviceAccount 并将它用在 nginx pod 上
kubectl create sa myuser
kubectl run nginx --image=nginx --restart=Never --serviceaccount=myuser -o yaml --dry-run=client > pod.yaml
kubectl apply -f pod.yaml
```

## 7 资源存储和分配
```yaml
# pv 模板
apiVersion: v1
kind: PersistentVolume
metadata:
  name: task-pv-volume
  labels:
    type: local
spec:
  storageClassName: manual
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: "/mnt/data"
--
# pvc 模板
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: task-pv-claim
spec:
  storageClassName: manual
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 3Gi
# 区别在于:
# pv: spec -> capacity -> storage 和 spec -> hostPath -> path
# pvc: spec -> resources -> requests -> storage
```

## 8 表格汇总

| **类型**               | **命令**                                                     | **描述**                           |
| ---------------------- | ------------------------------------------------------------ | ---------------------------------- |
| **基础命令**           | **create**                                                   | 通过文件名或标准输入创建资源       |
| **expose**             | 将一个资源公开为一个新的Service                              |                                    |
| **run**                | 在集群中运行一个特定的镜像                                   |                                    |
| **set**                | 在对象上设置特定的功能                                       |                                    |
| **get**                | 显示一个或多个资源                                           |                                    |
| **explain**            | 文档参考资料                                                 |                                    |
| **edit**               | 使用默认的编辑器编辑资源                                     |                                    |
| **delete**             | 通过文件名、标准输入、资源名称或标签选择器来删除资源         |                                    |
| **部署命令**           | **rollout**                                                  | 管理资源的发布                     |
| **rolling-update**     | 对给定的复制控制器滚动更新                                   |                                    |
| **scale**              | 扩容或缩容Pod、Deployment、ReplicaSet、RC或Job               |                                    |
| **autoscale**          | 创建一个自动选择扩容或缩容并设置Pod数量                      |                                    |
| **集群管理命令**       | **certificate**                                              | 修改证书资源                       |
| **cluster-info**       | 显示集群信息                                                 |                                    |
| **top**                | 显示资源（CPU、Memory、Storage）使用。需要Heapster运行       |                                    |
| **cordon**             | 标记节点不可调度                                             |                                    |
| **uncordon**           | 标记节点可调度                                               |                                    |
| **drain**              | 维护期间排除节点（驱除节点上的应用，准备下线维护）           |                                    |
| **taint**              | 设置污点属性                                                 |                                    |
| **故障诊断和调试命令** | **describe**                                                 | 显示特定资源或资源组的详细信息     |
| **logs**               | 在一个Pod中打印一个容器日志。如果Pod只有一个容器，容器名称是可选的 |                                    |
| **attach**             | 附加到一个运行的容器                                         |                                    |
| **exec**               | 执行命令到容器                                               |                                    |
| **port-forward**       | 转发一个或多个本地端口到一个Pod                              |                                    |
| **proxy**              | 运行一个proxy到Kubernetes API Server                         |                                    |
| **cp**                 | 拷贝文件或目录到容器                                         |                                    |
| **auth**               | 检查授权                                                     |                                    |
| **高级命令**           | **apply**                                                    | 通过文件名或标准输入对资源应用配置 |
| **patch**              | 使用补丁修改、更新资源的字段                                 |                                    |
| **replace**            | 通过文件名或标准输入替换一个资源                             |                                    |
| **convert**            | 不同的API版本之间转换配置文件                                |                                    |
| **设置命令**           | **label**                                                    | 更新资源上的标签                   |
| **annotate**           | 更新资源上的注释                                             |                                    |
| **completion**         | 用于实现kubectl工具自动补全                                  |                                    |
| **其他命令**           | **api-versions**                                             | 打印支持的API版本                  |
| **config**             | 修改kubeconfig文件（用于访问API，比如配置认证信息）          |                                    |
| **help**               | 所有命令帮助                                                 |                                    |
| **plugin**             | 运行一个命令行插件                                           |                                    |
| **version**            | 打印客户端和服务版本信息                                     |                                    |
