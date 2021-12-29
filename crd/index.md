# Kubernetes CRD Operator 开发系列


这篇文章总结了 Kubernetes 中 CRD 的知识体系和其中的底层实现等等。

<!--more-->

## 1 初识 CRD 和 Operator 框架

### 1.Custom Resource Define

`Custom Resource Define` 简称 CRD，是 Kubernetes（v1.7+）为提高可扩展性，让开发者去自定义资源的一种方式。CRD 资源可以动态注册到集群中，注册完毕后，用户可以通过 kubectl 来创建访问这个自定义的资源对象，类似于操作 Pod 一样。不过需要注意的是 CRD 仅仅是资源的定义而已，需要一个对应的控制器去监听 CRD 的各种事件来添加自定义的业务逻辑。

### 2.YAML 文件格式

我们可以定义一个如下所示的 CRD 资源清单文件：

```yaml
# crd-demo.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  # name 必须匹配下面的spec字段：<plural>.<group>  
  name: crontabs.stable.example.com
spec:
  # group 名用于 REST API 中的定义：/apis/<group>/<version>  
  group: stable.example.com
  # 列出自定义资源的所有 API 版本  
  versions:
  - name: v1beta1  # 版本名称，比如 v1、v2beta1 等等    
    served: true  # 是否开启通过 REST APIs 访问 `/apis/<group>/<version>/...`    
    storage: true # 必须将一个且只有一个版本标记为存储版本    
    schema:  # 定义自定义对象的声明规范      
      openAPIV3Schema:
        description: Define CronTab YAML Spec
        type: object
        properties:
          spec:
            type: object
            properties:
              cronSpec:
                type: string
              image:
                type: string
              replicas:
                type: integer
  # 定义作用范围：Namespaced（命名空间级别）或者 Cluster（整个集群）  
  scope: Namespaced
  names:
    # kind 是 singular 的一个驼峰形式定义，在资源清单中会使用    
    kind: CronTab
    # plural 名字用于 REST API 中的定义：/apis/<group>/<version>/<plural>    
    plural: crontabs
    # singular 名称用于 CLI 操作或显示的一个别名    
    singular: crontab
    # shortNames 相当于缩写形式    
    shortNames:
    - ct
```
我们在创建资源的时候，肯定不是任由我们随意去编写 YAML 文件的，当我们把上面的 CRD 文件提交给 Kubernetes 之后，Kubernetes 会对我们提交的声明文件进行校验，从定义可以看出 CRD 是基于 [OpenAPI v3 schem](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.0.md#schemaObject) 进行规范的。

使用 `kubectl` 来创建这个 CRD 资源清单：

```Shell
[root@master ~]# kubectl apply -f crd-demo.yaml
customresourcedefinition.apiextensions.k8s.io/crontabs.stable.example.com created
```
然后我们就可以使用这个 API 端点来创建和管理自定义的对象，这些对象的类型就是上面创建的 CRD 对象规范中的 `CronTab`。
现在在 Kubernetes 集群中我们就多了一种新的资源叫 `crontabs.stable.example.com`，我们就可以使用它来定义一个 `CronTab` 资源对象了，这个自定义资源对象里面可以包含的字段我们在定义的时候通过 `schema` 进行了规范，比如现在我们来创建一个如下所示的资源清单：

```yaml
# crd-crontab-demo.yaml
apiVersion: "stable.example.com/v1beta1"
kind: CronTab
metadata:
  name: my-new-cron-object
spec:
  cronSpec: "* * * * */5"
  image: my-awesome-cron-image
```
直接创建这个对象：
```shell
[root@master ~]# kubectl apply -f crd-crontab-demo.yaml
crontab.stable.example.com/my-new-cron-object created
````
我们就可以用 `kubectl` 来管理我们这里创建 `CronTab` 对象了，比如：

{{< admonition type=example >}}
```shell
[root@master ~]# kubectl get ct  # 简写
NAME                 AGE
my-new-cron-object   42s
[root@master ~]# kubectl get crontab
NAME                 AGE
my-new-cron-object   88s
```
{{< /admonition >}}

### 3.kubebuilder 脚手架使用

创建一个目录，然后在里面运行 `kubebuilder init` 命令，初始化一个新项目。

```shell
$ cd go/src/github.com/peterliao96
$ mkdir builder-demo
$ cd builder-demo
# 开启 go modules
$ export GO111MODULE=on
$ export GOPROXY=https://goproxy.cn
$ kubebuilder init --domain ydzs.io --owner peterliao96 --repo github.com/peterliao96/builder-demo
# 创建一个新的 API（组/版本）为 “webapp/v1”，并在上面创建新的 Kind(CRD) “Guestbook”
$ kubebuilder create api --group webapp --version v1 --kind Guestbook
```
上面的命令会创建文件 `api/v1/guestbook_types.go`，该文件中定义相关 API ，而针对于这一类型 (CRD) 的业务逻辑生成在 `controller/guestbook_controller.go` 文件中。
我们可以根据自己的需求去修改资源对象的定义结构，修改 `api/v1/guestbook_types.go` 文件：

```go
// api/v1/guestbook_types.go
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GuestbookSpec defines the desired state of Guestbook
type GuestbookSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Guestbook. Edit guestbook_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// GuestbookStatus defines the observed state of Guestbook
type GuestbookStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Guestbook is the Schema for the guestbooks API
type Guestbook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GuestbookSpec   `json:"spec,omitempty"`
	Status GuestbookStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GuestbookList contains a list of Guestbook
type GuestbookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Guestbook `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Guestbook{}, &GuestbookList{})
}
```
本部分摘录于 Kubernetes 开发课文档中的 [CRD介绍](https://www.notion.so/CRD-133b377c4b5e408ea2ede4d40bb31108) 和 [kubebuilder 介绍](https://www.notion.so/kubebuilder-3cbdc53796044a939a58eaaba67f415b)。

## 2 client-go 实现原理

我们自定义的资源 CRD 创建完成以后，其实如果没有一个 controller 运行它的话，它并没有起到任何用处。于是，我们需要通过创建 controller 来运行 CRD。
要想了解 controller 的实现原理和方式，我们就需要了解下 `client-go` 这个库的实现，Kubernetes 部分代码也是基于这个库实现的，也包含了开发自定义控制器时可以使用的各种机制。
下图显示了 `client-go` 中的各个组件是如何公众的以及与我们要编写的自定义控制器代码的交互入口：

![controller interaction](/client-go-controller-interaction.jpeg "图1:client-go 实现流程图")

### 1. 组件介绍

`client-go` 组件：
- `Reflector`：通过 Kubernetes API 监控 Kubernetes 的资源类型 采用 List/Watch 机制, 可以 Watch 任何资源包括 CRD 添加 object 对象到 FIFO 队列，然后 Informer 会从队列里面取数据
- `Informer`：controller 机制的基础，循环处理 object 对象 从 Reflector 取出数据，然后将数据给到 Indexer 去缓存，提供对象事件的 `handler` 接口，只要给 Informer 添加 `ResourceEventHandler` 实例的回调函数，去实现 `OnAdd(obj interface{})`、 `OnUpdate(oldObj, newObj interface{})` 和 `OnDelete(obj interface{})` 这三个方法，就可以处理好资源的创建、更新和删除操作了。
- `Indexer`：提供 object 对象的索引，是线程安全的，缓存对象信息

`controller` 组件：
- `Informer reference`: controller 需要创建合适的 Informer 才能通过 Informer reference 操作资源对象
- `Indexer reference`: controller 创建 Indexer reference 然后去利用索引做相关处理
- `Resource Event Handlers`：Informer 会回调这些 handlers
- `Work queue`: Resource Event Handlers 被回调后将 key 写到工作队列，这里的 key 相当于事件通知，后面根据取出事件后，做后续的处理
- `Process Item`：从工作队列中取出 key 后进行后续处理，具体处理可以通过 Indexer reference controller 可以直接创建上述两个引用对象去处理，也可以采用工厂模式，官方都有相关示例

### 2. 自定义 Controller 的控制流

![controller workflow](/client-go-controller-workflow.png "图2：controller 工作流")


如上图所示主要有两个部分，一个是发生在 `SharedIndexInformer` 中，另外一个是在自定义控制器中。

- `Reflector` 通过 Kubernetes APIServer 执行对象（比如 Pod）的 `ListAndWatch` 查询，记录和对象相关的三种事件类型 `Added、Updated、Deleted`，然后将它们传递到 `DeltaFIFO` 中去。
- `DeltaFIFO` 接收到事件和 `watch` 事件对应的对象，然后将他们转换为 `Delta` 对象，这些 `Delta` 对象被附加到队列中去等待处理，对于已经删除的，会检查线程安全的 `store` 中是否已经存在该文件，从而可以避免在不存在某些内容时排队执行删除操作。
- `Cache` 控制器（不要和自定义控制器混淆）调用 `Pop()` 方法从 `DeltaFIFO` 队列中出队列，`Delta` 对象将传递到 `SharedIndexInformer` 的 `HandleDelta()` 方法中以进行进一步处理。
- 根据 `Delta` 对象的操作（事件）类型，首先在 `HandleDeltas` 方法中通过 indexer 的方法将对对象保存到线程安全的 `Store` 中，然后，通过 `SharedIndexInformer` 中的 `sharedProcessor` 的 `distribution()` 方法将这些对象发送到事件 `handlers`，这些事件处理器由自定义控制器通过 `SharedInformer` 的方法比如 `AddEventHandlerWithResyncPeriod()` 进行注册。
- 已注册的事件处理器通过添加或更新时间的 `MetaNamespaceKeyFunc()` 或删除事件的 `DeletionHandingMetaNamespaceKeyFunc()` 将对象转换为格式为 `namespace/name` 或只是 `name` 的 `key`，然后将这个 `key` 添加到自定义控制器的 `workqueue` 中，`workqueues` 的实现可以在 `util/workqueue` 中找到。
- 自定义的控制器通过调用定义的 `handlers` 处理器从 `workqueue` 中 `pop` 一个 `key` 出来进行处理，`handlers` 将调用 `indexer` 的 `GetByKey()` 从线程安全的 `store` 中获取对象，我们的业务逻辑就是在这个 `handlers` 里面实现。

本部分摘录于 Kubernetes 文档中的 [CRD controller 原理实现](https://www.qikqiak.com/k3s/controller/crd/#%E5%8F%82%E8%80%83%E6%96%87%E6%A1%A3)。
