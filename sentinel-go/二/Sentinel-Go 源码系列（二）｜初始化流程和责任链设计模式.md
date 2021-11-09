上节中我们知道了 Sentinel-Go 大概能做什么事情，最简单的例子如何跑起来

其实我早就写好了本系列的第二篇，但迟迟没有发布，感觉光初始化流程显得有些单一，于是又补充了责任链模式，二合一，内容显得丰富一些。

# 初始化流程

## 初始化做了什么

Sentinel-Go 初始化时主要做了以下2件事情：

- 通过各种方式（文件、环境变量等）载入全局配置
- 启动异步的定时任务或服务，如机器 cpu、内存信息收集、metric log 写入等等

## 初始化流程详解

### 提供的 API

上节例子中，我们使用了最简单的初始化方式

```go
func InitDefault() error
```

除此之外，它还提供了另外几种初始化方式

```go
// 使用给定的 parser 方法解析配置的方式来初始化
func InitWithParser(configBytes []byte, parser func([]byte) (*config.Entity, error)) (err error)

// 使用已解析好的配置对象初始化
func InitWithConfig(confEntity *config.Entity) (err error)

// 从 yaml 文件加载配置初始化 
func InitWithConfigFile(configPath string) error
```

从命名能看出它们只是配置的获取方式不一样，其中`InitWithParser` 有点意思，传入的 `parser` 是个函数指针，对于 Java 写惯了的我来说还是有点陌生，比如通过 `json` 解析可以写出如下 `parser`

```go
parser := func(configBytes []byte) (*config.Entity, error) {
	conf := &config.Entity{}
	err := json.Unmarshal(configBytes, conf)
	return conf, err
}
conf := "{\"Version\":\"v1\",\"Sentinel\":{\"App\":{\"Name\":\"roshi-app\",\"Type\":0}}}"
err := api.InitWithParser([]byte(conf), parser)
```

### 配置项

简单看一下 Sentinel-Go 的配置项，首先配置被包装在一个 `Entity` 中，包含了一个 `Version` 和 真正的配置信息 `SentinelConfig`

```go
type Entity struct {
	Version string
	Sentinel SentinelConfig
}
```

接着， `SentinelConfig` 是这样：

```go
type SentinelConfig struct {
	App struct {
		// 应用名
		Name string
		// 应用类型：普通应用，网关
		Type int32
	}
	// Exporter 配置
	Exporter ExporterConfig
	// 日志配置
	Log LogConfig
	// 统计配置
	Stat StatConfig
	// 是否缓存时间戳
	UseCacheTime bool `yaml:"useCacheTime"`
}
```

- App 应用信息
  - 应用名
  - 应用类型：如普通应用、网关应用等
- ExporterConfig：prometheus exporter 暴露服务的端口和 path

```go
type ExporterConfig struct {
	Metric MetricExporterConfig
}

type MetricExporterConfig struct {
	// http 服务地址，如 ":8080"
	HttpAddr string `yaml:"http_addr"`
	//  http 服务 path，如"/metrics".
	HttpPath string `yaml:"http_path"`
}
```

- LogConfig：包括使用什么logger，日志目录，文件是否使用 pid（防止一台机器部署两个应用日志混合），以及 metric log 的单个文件大小、最多保留文件个数、刷新时间

```go
type LogConfig struct {
	// logger，可自定义
	Logger logging.Logger
	// 日志目录
	Dir string
	// 是否在日志文件后加 PID
	UsePid bool `yaml:"usePid"`
	// metric 日志配置
	Metric MetricLogConfig
}

type MetricLogConfig struct {
  // 单个文件最大占用空间
	SingleFileMaxSize uint64 `yaml:"singleFileMaxSize"`
	// 最多文件个数
	MaxFileCount      uint32 `yaml:"maxFileCount"`
	// 刷新间隔
	FlushIntervalSec  uint32 `yaml:"flushIntervalSec"`
}
```

- StatConfig：统计配置包括资源采集窗口配置，metric 统计的窗口、系统信息收集间隔

```go
type StatConfig struct {
	// 全局统计资源的窗口（后续文章再解释）
	GlobalStatisticSampleCountTotal uint32 `yaml:"globalStatisticSampleCountTotal"`
	GlobalStatisticIntervalMsTotal  uint32 `yaml:"globalStatisticIntervalMsTotal"`
	// metric 统计的窗口（后续文章再解释）
	MetricStatisticSampleCount uint32 `yaml:"metricStatisticSampleCount"`
	MetricStatisticIntervalMs  uint32 `yaml:"metricStatisticIntervalMs"`
	// 系统采集配置
	System SystemStatConfig `yaml:"system"`
}

type SystemStatConfig struct {
	// 采集默认间隔
	CollectIntervalMs uint32 `yaml:"collectIntervalMs"`
	// 采集 cpu load 间隔
	CollectLoadIntervalMs uint32 `yaml:"collectLoadIntervalMs"`
	// 采集 cpu 使用间隔
	CollectCpuIntervalMs uint32 `yaml:"collectCpuIntervalMs"`
	// 采集内存间隔使用
	CollectMemoryIntervalMs uint32 `yaml:"collectMemoryIntervalMs"`
}
```

### 配置覆盖

从上文知道，参数可以通过自定义 `parser` /  `文件` / `默认` 的方式来传入配置，但后面这个配置还可以用系统的`环境变量`覆盖，覆盖项目前只包括应用名、应用类型、日志文件使用使用` PID` 结尾、日志目录

```go
func OverrideConfigFromEnvAndInitLog() error {
	// 系统环境变量可覆盖传入的配置
	err := overrideItemsFromSystemEnv()
	if err != nil {
		return err
	}
	...
	return nil
}
```

### 启动后台服务

- 启动 聚合 metric 定时任务，聚合后发送到 chan，聚合后的格式如下：

```go
_, err := fmt.Fprintf(&b, "%d|%s|%s|%d|%d|%d|%d|%d|%d|%d|%d",
		m.Timestamp, timeStr, finalName, m.PassQps,
		m.BlockQps, m.CompleteQps, m.ErrorQps, m.AvgRt,
		m.OccupiedPassQps, m.Concurrency, m.Classification)
```

 `时间戳|时间字符串|名称|通过QPS|阻断QPS|完成QPS|出错QPS|平均RT|已经通过QPS|并发|类别`

- 启动 metric 写入日志定时任务，可配置间隔时间（秒级），接受上个任务写入 chan 的数据

- 启动单独 goroutine 收集 cpu 使用率 / load、内存使用，收集间隔可配置，收集到的信息存放在 `system_metric` 下的私有变量

```go
var (
	currentLoad        atomic.Value
	currentCpuUsage    atomic.Value
	currentMemoryUsage atomic.Value
)
```

- 若开启，则启动单独 goroutine 缓存时间戳，间隔是 1ms，这个主要是为了高并发下提高获取时间戳的性能

```go
func (t *RealClock) CurrentTimeMillis() uint64 {
  // 从缓存获取时间戳
	tickerNow := CurrentTimeMillsWithTicker()
	if tickerNow > uint64(0) {
		return tickerNow
	}
	return uint64(time.Now().UnixNano()) / UnixTimeUnitOffset
}
```

获取时，如果拿到 0 则说明未开启缓存时间戳，取当前，如果拿到值说明已开启，可直接使用

- 若配置了 metric exporter，则启动服务，监听端口，暴露 prometheus 的 exporter

# 责任链模式

## 什么是责任链模式

可以用这样一张图形象地解释什么是责任链：

<img src="/Users/didi/Documents/lkxiaolou/zaizai/dubbo/dubbo的前世今生/img19.jpg" alt="img19" style="zoom:50%;" /> 

责任链模式为每次请求创建了一个`链`，链上有 N 多个处理者，处理者可在不同阶段处理不同的事情，就像这幅图上的小人，拿到一桶水（请求）后都可以完成各自的事情，比如往头上浇，然后再传递给下一个。

为什么叫责任？因为每个处理者只关心自己的`责任`，跟自己没关系就递交给链上的下一个处理者。

责任链在哪里有用到？很多开源产品都是用了责任链模式，如 `Dubbo`、`Spring MVC `等等

这么设计有什么好处？

- 简化编码难度，抽象出处理模型，只需关注关心的点即可
- 扩展性好，如果需要自定义责任链中的一环或者插拔某一环，非常容易实现

> 关于扩展性除了大家理解的软件设计中的扩展性外，这里还想提两点，阿里开源的软件其实都有高扩展性这个特性，一是因为是开源，别人使用场景未必和自己一致，留出扩展接口，不符合要求的，用户可以自行实现，二是如果要追溯，阿里开源扩展性 Dubbo 可能算是祖师爷（未考证），Dubbo 作者（梁飞）的博客中说过为什么 Dubbo 要设计这么强的扩展性，他对代码有一定的追求，在他维护时期，代码能保证高质量，但如果项目交给别人，如何才能保持现在的水准呢？于是他设计出一套很强的扩展，后面开发基于这个扩展去做，代码就不会差到哪里去

- 可动态，可针对每个请求构造不同的责任链

## Sentinel-Go 责任链设计

先看责任链的数据结构定义，Sentinel-Go 把处理者叫 `Slot`（插槽），将 Slot 分为了前置统计、规则校验、统计三组，且每组是有有序的

```go
type SlotChain struct {
	// 前置准备（有序）
	statPres []StatPrepareSlot
	// 规则校验（有序）
	ruleChecks []RuleCheckSlot
	// 统计（有序）
	stats []StatSlot
	// 上线文对象池（复用对象）
	ctxPool *sync.Pool
}
```

在调用 `Entry` 开始进入 Sentinel 逻辑时，如果没有手动构造 SlotChain，则使用默认。

为什么这里要设计成三个 Slot组呢？因为每组 Slot 的行为稍有不同，比如前置准备的 Slot 不需要返回值，规则校验组需要返回值，如果校验当前流量不通过，还需要返回原因、类型等信息，统计 Slot 还会有一些入参，比如请求是否失败等等

```go
type BaseSlot interface {
	Order() uint32
}

type StatPrepareSlot interface {
	BaseSlot
	Prepare(ctx *EntryContext)
}

type RuleCheckSlot interface {
	BaseSlot
	Check(ctx *EntryContext) *TokenResult
}

type StatSlot interface {
	BaseSlot
	OnEntryPassed(ctx *EntryContext)
	OnEntryBlocked(ctx *EntryContext, blockError *BlockError)
	OnCompleted(ctx *EntryContext)
}
```

# 总结

本文从源码角度分析了 Sentinel-Go 的初始化流程和责任链的设计，总体上来说还是比较简单，接下来的系列文章将会分析 Sentinel-Go 的限流熔断等的核心设计与实现。
