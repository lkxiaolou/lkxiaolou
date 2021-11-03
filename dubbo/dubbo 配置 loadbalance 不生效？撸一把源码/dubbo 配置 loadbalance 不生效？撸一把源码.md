### 背景

很久之前我给业务方写了一个 dubbo loadbalance 的扩展（为了叙述方便，这个 loadbalance 扩展就叫它 `XLB` 吧），这两天业务方反馈说 `XLB` 不生效了

我心想，不可能啊，都用了大半年了~ 

### 排查

于是我登上不生效的 consumer 机器进行排查，还好我留了一手，当 `XLB` 加载时，会打印一行日志

看了下这个服务，并没有打印日志，说明 `XLB` 并没有加载成功

于是，我就去问对应的开发，有按照我的文档配置 loadbalance 吗？答复：完全按照文档配置

这下我就有点不相信了，但转念一想，配置 loadbalance 如此简单，不应该出错啊，我的文档和他的应用都在 xml 文件中配置了 consumer 的 loadbalance

```xml
<dubbo:consumer loadbalance="xlb"/>
```

抱着试一试的态度，拉取了他们项目的代码，发现配置确实如上，但我发现他们的 application.properties 配置文件也配了一个 consumer 的属性

```xml
dubbo.consumer.check=false
```

以多年和 dubbo 打交道的经验来说，这里有问题，又确认了代码，确实 xml 和 application.properties 都加载了

那这里可能就有问题了，dubbo 从 xml 加载生成了一个 consumer 配置，dubbo-springboot-starter 又从 application.properties 加载配置生成了一个 consumer 配置，这不就冲突了？

别看只配置了 dubbo.consumer.check，它实际上会生成一个完整的 consumer 配置，只不过 loadbalance 为默认值

业务方为什么会这样配置？大概率是因为我的文档里只给出了 xml 形式的配置，没有给 spring-boot 配置，他们原先使用的是 spring-boot 的配置方式，然后看到我的文档是 xml，结果就不会配置了，也写了个 xml，和原先的配置冲突

###  验证

为了验证是这个问题导致，我把他的 application.properties 的 dubbo.consumer.check 配置挪到了 xml 文件中，果然重启后就加载到了 `XLB`

随后我又在本地的测试应用上做了这样一个验证：

```xml
<!-- case 1 -->
<dubbo:consumer />
<dubbo:consumer loadbalance="xlb"/>

<!-- case 2 -->
<dubbo:consumer loadbalance="xlb"/>
<dubbo:consumer />
```

两组配置相同，但顺序不同，测试结果为 case 1 可以加载到 `xlb`，case 2 不行

于是猜测，dubbo consumer 配置以后加载的为准

### 撸源码

显然猜测不符合我的风格，下面开撸源码，不感兴趣可以划过，最下面有总结

首先搞清楚，何时会加载 loadbalance，在 `AbstractClusterInvoker` 的 `invoke` 方法中，加载了 loadbalance

```java
@Override
public Result invoke(final Invocation invocation) throws RpcException {
    ...
    List<Invoker<T>> invokers = list(invocation);
    LoadBalance loadbalance = initLoadBalance(invokers, invocation);
    RpcUtils.attachInvocationIdIfAsync(getUrl(), invocation);
    return doInvoke(invocation, invokers, loadbalance);
}
```

加载代码如下

```java
protected LoadBalance initLoadBalance(List<Invoker<T>> invokers, Invocation invocation) {
    if (CollectionUtils.isNotEmpty(invokers)) {
        return ExtensionLoader.getExtensionLoader(LoadBalance.class).getExtension(invokers.get(0).getUrl()
                .getMethodParameter(RpcUtils.getMethodName(invocation), LOADBALANCE_KEY, DEFAULT_LOADBALANCE));
    } else {
        return ExtensionLoader.getExtensionLoader(LoadBalance.class).getExtension(DEFAULT_LOADBALANCE);
    }
}
```

可以看出

- loadbalance 是发起 dubbo 调用时，且当 `invokers` 非空时（即 providers 非空）会被初始化，后续都从缓存中取
- loadbalance 是根据第一个 invoker 的 loadbalance 参数决定使用哪个 loadbalance 的

于是问题转移到 invoker 的 loadbalance 从哪来？provider 不会配置 loadbalance，所以这个参数一定是从 consumer 的配置上得到的

顺腾摸瓜，在 `RegistryDirectory` 的 `toInvokers` 方法中调用了 `mergeUrl`，它是在注册中心通知时被调用，也就是从注册中心上拿到 provider url 时，还得 merge 一下才能用，merge 了些什么内容？

```java
private URL mergeUrl(URL providerUrl) {
    // 1. merge consumer 参数
    providerUrl = ClusterUtils.mergeUrl(providerUrl, queryMap); 
    // 2. merge configurator 参数
    providerUrl = overrideWithConfigurator(providerUrl);
    ...
    return providerUrl;
}
```

1中 merge 了queryMap 里的参数，这个queryMap 其实就是 consumer 的参数，它来自配置的 reference

再看 reference 配置，当 `ReferenceConfig` 初始化时

```java
// 1
public synchronized void init() {
    ...
    checkAndUpdateSubConfigs();
    ...
    AbstractConfig.appendParameters(map, consumer);
    ...
}

// 2
public void checkAndUpdateSubConfigs() {
    ...
    checkDefault();
    ...
}

// 3
public void checkDefault() throws IllegalStateException {
    if (consumer == null) {
        consumer = ApplicationModel.getConfigManager()
                .getDefaultConsumer()
                .orElse(new ConsumerConfig());
    }
}

// 4
public Optional<ConsumerConfig> getDefaultConsumer() {
    List<ConsumerConfig> consumerConfigs = getDefaultConfigs(getConfigsMap(getTagName(ConsumerConfig.class)));
    if (CollectionUtils.isNotEmpty(consumerConfigs)) {
        return Optional.of(consumerConfigs.get(0));
    }
    return Optional.empty();
}
```

上面调用链从 `1 到 4`，`4` 中获取了第1个 consumer，这就是我们要找的根源

### 总结

- 每配置一个 consumer ，无论是从 xml 文件，或是 spring-boot 配置，或是 api 直接创建，都会生成一个 consumerConfig 对象
- 当消费接口，即配置 reference 时，会将 consumer 的参数 merge 过来，如果存在多个 consumer，会挑第一个，当然我们并不知道谁先加载
- 当 reference 存在 consumer 的配置时，注册中心通知的 provider urls 会和 reference 的参数进行合并，合并后生成可调用的 invoker
-  对于 loadbalance 来说，调用时，如果 invokers 非空，则会尝试通过第一个 invoker 的 loadbalance 参数加载负载均衡算法，第一次调用进行加载，后续调用则使用缓存