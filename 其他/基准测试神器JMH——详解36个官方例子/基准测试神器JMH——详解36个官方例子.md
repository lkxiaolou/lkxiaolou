> 本文已收录 https://github.com/lkxiaolou/lkxiaolou 欢迎star。

## 简介
基准测试是指通过设计科学的测试方法、测试工具和测试系统，实现对一类测试对象的某项性能指标进行定量的和可对比的测试。
而JMH是一个用来构建，运行，分析Java或其他运行在JVM之上的语言的 纳秒/微秒/毫秒/宏观 级别基准测试的工具。

>JMH is a Java harness for building, running, and analysing nano/micro/milli/macro benchmarks written in Java and other languages targetting the JVM.

## 为什么需要
有人可能会说，我可以在代码的前后打点计算代码运行时间，为什么还需要JMH？下面着重介绍一下JMH的常用功能，都来源于官方提供的sample，相信看完这些例子，就能回答这个问题了，[官方sample地址点击可查看](http://hg.openjdk.java.net/code-tools/jmh/file/tip/jmh-samples/src/main/java/org/openjdk/jmh/samples/)。

## 官方sample解读
###### （1）JMHSample_01_HelloWorld
第一个例子教我们如何使用，在开始前，只需要引入依赖。类似单元测试，常放在test目录下运行。
```
<dependencies>
    <dependency>
        <groupId>org.openjdk.jmh</groupId>
        <artifactId>jmh-core</artifactId>
        <version>1.23</version>
        <scope>test</scope>
    </dependency>
    <dependency>
        <groupId>org.openjdk.jmh</groupId>
        <artifactId>jmh-generator-annprocess</artifactId>
        <version>1.23</version>
        <scope>test</scope>
    </dependency>
</dependencies>
```

这里精简一下simple的代码，使用 @Benchmark 来标记需要基准测试的方法，然后需要写一个main方法来启动基准测试。
```
public class Sample1 {

    @Benchmark
    public void wellHelloThere() {
        // this method was intentionally left blank.
    }

    public static void main(String[] args) throws RunnerException {
        Options opt = new OptionsBuilder()
                .include(Sample1.class.getSimpleName())
                .forks(1)
                .build();

        new Runner(opt).run();
    }
}
```
可以在IDE中直接运行main方法；如果是在服务器上可以打成jar包运行。运行后控制台输出如下格式的报告：
报告的第一部分是此次运行的环境和配置，包括JDK、JMH版本，基准测试的配置（后面会详细介绍）等。第二部分则是每次运行的报告输出，第三部分是汇总的报告，包括最小值、平均值、最大值。最后一部分则是本次基准测试的最终报告。
```
/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/bin/java -Dfile.encoding=UTF-8 -classpath /Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/charsets.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/deploy.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/ext/cldrdata.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/ext/dnsns.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/ext/jaccess.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/ext/jfxrt.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/ext/localedata.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/ext/nashorn.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/ext/sunec.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/ext/sunjce_provider.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/ext/sunpkcs11.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/ext/zipfs.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/javaws.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/jce.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/jfr.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/jfxswt.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/jsse.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/management-agent.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/plugin.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/resources.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/lib/rt.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/lib/ant-javafx.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/lib/dt.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/lib/javafx-mx.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/lib/jconsole.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/lib/packager.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/lib/sa-jdi.jar:/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/lib/tools.jar:/Users/lkxx/projects/leetcode/target/test-classes:/Users/lkxx/projects/leetcode/target/classes:/Users/lkxx/.m2/repository/org/openjdk/jmh/jmh-core/1.23/jmh-core-1.23.jar:/Users/lkxx/.m2/repository/net/sf/jopt-simple/jopt-simple/4.6/jopt-simple-4.6.jar:/Users/lkxx/.m2/repository/org/apache/commons/commons-math3/3.2/commons-math3-3.2.jar:/Users/lkxx/.m2/repository/org/openjdk/jmh/jmh-generator-annprocess/1.23/jmh-generator-annprocess-1.23.jar jmh.simple1
# JMH version: 1.23
# VM version: JDK 1.8.0_111, Java HotSpot(TM) 64-Bit Server VM, 25.111-b14
# VM invoker: /Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre/bin/java
# VM options: -Dfile.encoding=UTF-8
# Warmup: 5 iterations, 10 s each
# Measurement: 5 iterations, 10 s each
# Timeout: 10 min per iteration
# Threads: 1 thread, will synchronize iterations
# Benchmark mode: Throughput, ops/time
# Benchmark: jmh.simple1.wellHelloThere

# Run progress: 0.00% complete, ETA 00:01:40
# Fork: 1 of 1
# Warmup Iteration   1: 3658920710.925 ops/s
# Warmup Iteration   2: 3461583245.233 ops/s
# Warmup Iteration   3: 3789034272.009 ops/s
# Warmup Iteration   4: 3861706248.346 ops/s
# Warmup Iteration   5: 3828576338.676 ops/s
Iteration   1: 3824989172.956 ops/s
Iteration   2: 3854633115.447 ops/s
Iteration   3: 3825194474.967 ops/s
Iteration   4: 3835680543.019 ops/s
Iteration   5: 3842258462.856 ops/s

Result "jmh.simple1.wellHelloThere":
  3836551153.849 ±(99.9%) 48053746.458 ops/s [Average]
  (min, avg, max) = (3824989172.956, 3836551153.849, 3854633115.447), stdev = 12479405.354
  CI (99.9%): [3788497407.391, 3884604900.306] (assumes normal distribution)

# Run complete. Total time: 00:01:45

REMEMBER: The numbers below are just data. To gain reusable insights, you need to follow up on
why the numbers are the way they are. Use profilers (see -prof, -lprof), design factorial
experiments, perform baseline and negative tests that provide experimental control, make sure
the benchmarking environment is safe on JVM/OS/HW level, ask for reviews from the domain experts.
Do not assume the numbers tell you what you want them to tell.

Benchmark        Mode  Cnt           Score          Error  Units
wellHelloThere  thrpt    5  3836551153.849 ± 48053746.458  ops/s

Process finished with exit code 0
```
###### （2）JMHSample_02_BenchmarkModes
本例介绍了注解 @OutputTimeUnit 和 @BenchmarkMode。
```
public class Simple2 {

    @Benchmark
    @BenchmarkMode(Mode.Throughput)
    @OutputTimeUnit(TimeUnit.SECONDS)
    public void measureThroughput() throws InterruptedException {
        TimeUnit.MILLISECONDS.sleep(100);
    }

    @Benchmark
    @BenchmarkMode(Mode.AverageTime)
    @OutputTimeUnit(TimeUnit.MICROSECONDS)
    public void measureAvgTime() throws InterruptedException {
        TimeUnit.MILLISECONDS.sleep(100);
    }

    @Benchmark
    @BenchmarkMode(Mode.SampleTime)
    @OutputTimeUnit(TimeUnit.MICROSECONDS)
    public void measureSamples() throws InterruptedException {
        TimeUnit.MILLISECONDS.sleep(100);
    }

    @Benchmark
    @BenchmarkMode(Mode.SingleShotTime)
    @OutputTimeUnit(TimeUnit.MICROSECONDS)
    public void measureSingleShot() throws InterruptedException {
        TimeUnit.MILLISECONDS.sleep(100);
    }

    @Benchmark
    @BenchmarkMode({Mode.Throughput, Mode.AverageTime, Mode.SampleTime, Mode.SingleShotTime})
    @OutputTimeUnit(TimeUnit.MICROSECONDS)
    public void measureMultiple() throws InterruptedException {
        TimeUnit.MILLISECONDS.sleep(100);
    }

    @Benchmark
    @BenchmarkMode(Mode.All)
    @OutputTimeUnit(TimeUnit.MICROSECONDS)
    public void measureAll() throws InterruptedException {
        TimeUnit.MILLISECONDS.sleep(100);
    }
}
```

- @OutputTimeUnit 可以指定输出的时间单位，可以传入 java.util.concurrent.TimeUnit 中的时间单位，最小可以到纳秒级别；
- @BenchmarkMode 指明了基准测试的模式
- Mode.Throughput ：吞吐量，单位时间内执行的次数
- Mode.AverageTime：平均时间，一次执行需要的单位时间，其实是吞吐量的倒数
- Mode.SampleTime：是基于采样的执行时间，采样频率由JMH自动控制，同时结果中也会统计出p90、p95的时间
- Mode.SingleShotTime：单次执行时间，只执行一次，可用于冷启动的测试

这些模式可以自由组合，甚至可以使用全部。
```
Benchmark                                          Mode  Cnt       Score      Error   Units
Simple2.measureAll                                thrpt    5      ≈ 10⁻⁵             ops/us
Simple2.measureMultiple                           thrpt    5      ≈ 10⁻⁵             ops/us
Simple2.measureThroughput                         thrpt    5       9.795 ±    0.142   ops/s
Simple2.measureAll                                 avgt    5  103306.387 ± 2606.185   us/op
Simple2.measureAvgTime                             avgt    5  106399.561 ± 1447.102   us/op
Simple2.measureMultiple                            avgt    5  106154.350 ±  525.031   us/op
Simple2.measureAll                               sample  476  105773.452 ±  417.365   us/op
Simple2.measureAll:measureAll·p0.00              sample       100007.936              us/op
Simple2.measureAll:measureAll·p0.50              sample       106430.464              us/op
Simple2.measureAll:measureAll·p0.90              sample       108920.832              us/op
Simple2.measureAll:measureAll·p0.95              sample       109707.264              us/op
Simple2.measureAll:measureAll·p0.99              sample       110100.480              us/op
Simple2.measureAll:measureAll·p0.999             sample       110624.768              us/op
Simple2.measureAll:measureAll·p0.9999            sample       110624.768              us/op
Simple2.measureAll:measureAll·p1.00              sample       110624.768              us/op
Simple2.measureMultiple                          sample  475  106084.986 ±  390.356   us/op
Simple2.measureMultiple:measureMultiple·p0.00    sample       100007.936              us/op
Simple2.measureMultiple:measureMultiple·p0.50    sample       106561.536              us/op
Simple2.measureMultiple:measureMultiple·p0.90    sample       108920.832              us/op
Simple2.measureMultiple:measureMultiple·p0.95    sample       109445.120              us/op
Simple2.measureMultiple:measureMultiple·p0.99    sample       110100.480              us/op
Simple2.measureMultiple:measureMultiple·p0.999   sample       110100.480              us/op
Simple2.measureMultiple:measureMultiple·p0.9999  sample       110100.480              us/op
Simple2.measureMultiple:measureMultiple·p1.00    sample       110100.480              us/op
Simple2.measureSamples                           sample  490  102720.324 ±  228.966   us/op
Simple2.measureSamples:measureSamples·p0.00      sample       100007.936              us/op
Simple2.measureSamples:measureSamples·p0.50      sample       102891.520              us/op
Simple2.measureSamples:measureSamples·p0.90      sample       104857.600              us/op
Simple2.measureSamples:measureSamples·p0.95      sample       105119.744              us/op
Simple2.measureSamples:measureSamples·p0.99      sample       105119.744              us/op
Simple2.measureSamples:measureSamples·p0.999     sample       105119.744              us/op
Simple2.measureSamples:measureSamples·p0.9999    sample       105119.744              us/op
Simple2.measureSamples:measureSamples·p1.00      sample       105119.744              us/op
Simple2.measureAll                                   ss       100173.622              us/op
Simple2.measureMultiple                              ss       103086.645              us/op
Simple2.measureSingleShot                            ss       105189.463              us/op
```
######（3）JMHSample_03_States
```
public class Sample3 {

    @State(Scope.Benchmark)
    public static class BenchmarkState {
        volatile double x = Math.PI;
    }

    @State(Scope.Thread)
    public static class ThreadState {
        volatile double x = Math.PI;
    }

    @Benchmark
    public void measureUnshared(ThreadState state) {
        state.x++;
    }

    @Benchmark
    public void measureShared(BenchmarkState state) {
        state.x++;
    }
}
```
本例介绍了注解 @State 的用法，用于多线程的测试
- @State(Scope.Thread)：作用域为线程，可以理解为一个ThreadLocal变量
- @State(Scope.Benchmark)：作用域为本次JMH测试，线程共享
- @State(Scope.Group)：作用域为group，将在后文看到

而且JMH可以像spring一样自动注入这些变量。
###### （4）JMHSample_04_DefaultState
本例介绍了 @State 注解可以直接写在 Benchmark 的测试类上，表明类的所有属性的作用域。
###### （5）JMHSample_05_StateFixtures
```
@State(Scope.Thread)
public class Sample5 {

    double x;

    @Setup
    public void prepare() {
        x = Math.PI;
    }

    @TearDown
    public void check() {
        assert x > Math.PI : "Nothing changed?";
    }

    @Benchmark
    public void measureRight() {
        x++;
    }

    @Benchmark
    public void measureWrong() {
        double x = 0;
        x++;
    }
}
```
本例中介绍了两个注解 @Setup 和 @TearDown。 @Setup 用于基准测试前的初始化动作， @TearDown 用于基准测试后的动作
###### （6）JMHSample_06_FixtureLevel
@Setup 和 @TearDown两个注解都可以传入 Level 参数，Level参数表明粒度，粒度从粗到细分别是
- Level.Trial：Benchmark级别
- Level.Iteration：执行迭代级别
- Level.Invocation：每次方法调用级别
###### （7）JMHSample_07_FixtureLevelInvocation
本例中主要介绍了使用Level.Invocation达到每次方法执行完成后sleep一段时间，模拟在需要唤醒线程的情况下耗时更多。
###### （8）JMHSample_08_DeadCode
```
public class Sample8 {

    private double x = Math.PI;

    @Benchmark
    public void baseline() {
        // do nothing, this is a baseline
    }

    @Benchmark
    public void measureWrong() {
        // This is wrong: result is not used and the entire computation is optimized away.
        Math.log(x);
    }

    @Benchmark
    public double measureRight() {
        // This is correct: the result is being used.
        return Math.log(x);
    }
}
```
本例主要介绍了一个知识点：Dead-Code Elimination (DCE) ，即死码消除，文档上说编译器非常聪明，有的代码没啥用，就在编译器被消除了，但这给我做基准测试带了一些麻烦，比如上面的代码中，baseline 和 measureWrong 有着相同的性能，因为编译器觉得 measureWrong这段代码执行后没有任何影响，为了效率，就直接消除掉这段代码，但是如果加上return语句，就不会在编译期被去掉，这是我们在写基准测试时需要注意的点。
###### （9）JMHSample_09_Blackholes
本例是为了解决（8）中死码消除问题，JMH提供了一个 Blackholes （黑洞），这样写就不会被编译器消除了。
```
@Benchmark
public void measureRight1(Blackhole blackhole) {
    blackhole.consume(Math.log(x));
}
```
###### （10）JMHSample_10_ConstantFold
```
public class Sample10 {
    
    private double x = Math.PI;
    
    private final double wrongX = Math.PI;

    @Benchmark
    public double baseline() {
        // simply return the value, this is a baseline
        return Math.PI;
    }

    @Benchmark
    public double measureWrong_1() {
        // This is wrong: the source is predictable, and computation is foldable.
        return Math.log(Math.PI);
    }

    @Benchmark
    public double measureWrong_2() {
        // This is wrong: the source is predictable, and computation is foldable.
        return Math.log(wrongX);
    }

    @Benchmark
    public double measureRight() {
        // This is correct: the source is not predictable.
        return Math.log(x);
    }
}
```
本例介绍了 constant-folding，即常量折叠，上述代码的 measureWrong_1 和  measureWrong_2 中的运算都是可以预测的值，所以也会在编译期直接替换为计算结果，从而导致基准测试失败，注意 final 修饰的变量也会被折叠。
###### （11）JMHSample_11_Loops
本例直接给出一个结论，不要在基准测试的时候使用循环，使用循环就会导致测试结果不准确，原因很复杂，甚至可以单独写一篇文章来介绍。简单能理解的一点是如果使用循环，预热可能就会存在问题。
```    
/*
     * 永远不要在基准测试的时候使用循环
     */
    private int reps(int reps) {
        int s = 0;
        for (int i = 0; i < reps; i++) {
            s += (x + y);
        }
        return s;
    }
```
###### （12）JMHSample_12_Forking
本例介绍了 @Fork 注解，@Fork 可以指定代码运行时是否需要 fork 出一个JVM进程，如果在同一个JVM中测试则会相互影响，一般fork进程设置为1。
###### （13）JMHSample_13_RunToRun
由于JVM的复杂性，每次测试结果都有差异，可以使用 @Fork 注解启动多个 JVM 经过多次测试来消除这种差异。
###### （15）JMHSample_15_Asymmetric
原来没有14，直接跳到了15
```
@State(Scope.Group)
@BenchmarkMode(Mode.AverageTime)
@OutputTimeUnit(TimeUnit.NANOSECONDS)
public class Sample15 {

    private AtomicInteger counter;

    @Setup
    public void up() {
        counter = new AtomicInteger();
    }

    @Benchmark
    @Group("g")
    @GroupThreads(3)
    public int inc() {
        return counter.incrementAndGet();
    }

    @Benchmark
    @Group("g")
    @GroupThreads(1)
    public int get() {
        return counter.get();
    }

    public static void main(String[] args) throws RunnerException {
        Options opt = new OptionsBuilder()
                .include(Sample15.class.getSimpleName())
                .forks(1)
                .build();

        new Runner(opt).run();
    }
}
```
本例是对 @Group 和 @GroupThreads 使用的介绍，@Group 定义了一个线程组， @GroupThreads 可以分配线程给测试用例，可以测试线程执行不均衡的情况，比如三个线程写，一个线程读，这里用 @State(Scope.Group) 定义了counter 作用域是这个线程组。
```
Benchmark       Mode  Cnt   Score   Error  Units
Sample15.g      avgt    5  49.339 ± 3.703  ns/op
Sample15.g:get  avgt    5  28.546 ± 2.414  ns/op
Sample15.g:inc  avgt    5  56.269 ± 4.320  ns/op
```
执行完的数据包含get、inc、和整个组的统计，这样数据更直观，更全面。
###### （16）JMHSample_16_CompilerControl
本例提到了JVM的方法内联，简单来说比较短但是执行频率又很高的方法，在执行多次后，JVM将该方法的调用替换为本身，以减少出栈入栈，从而减少性能的消耗。但是Java方法内联是无法人为控制的。
```
@State(Scope.Thread)
@BenchmarkMode(Mode.AverageTime)
@OutputTimeUnit(TimeUnit.NANOSECONDS)
public class Sample16 {

    public void target_blank() {
        // this method was intentionally left blank
    }

    @CompilerControl(CompilerControl.Mode.DONT_INLINE)
    public void target_dontInline() {
        // this method was intentionally left blank
    }

    @CompilerControl(CompilerControl.Mode.INLINE)
    public void target_inline() {
        // this method was intentionally left blank
    }

    @CompilerControl(CompilerControl.Mode.EXCLUDE)
    public void target_exclude() {
        // this method was intentionally left blank
    }
    
    @Benchmark
    public void baseline() {
        // this method was intentionally left blank
    }

    @Benchmark
    public void blank() {
        target_blank();
    }

    @Benchmark
    public void dontinline() {
        target_dontInline();
    }

    @Benchmark
    public void inline() {
        target_inline();
    }

    @Benchmark
    public void exclude() {
        target_exclude();
    }
}
```
JMH提供了可以控制是否使用内联的注解 @CompilerControl ，它的参数有如下可选：
- CompilerControl.Mode.DONT_INLINE：不使用内联
- CompilerControl.Mode.INLINE：强制使用内联
- CompilerControl.Mode.EXCLUDE：不编译
```
Benchmark            Mode  Cnt   Score   Error  Units
Sample16.baseline    avgt    3   0.261 ± 0.004  ns/op
Sample16.blank       avgt    3   0.262 ± 0.005  ns/op
Sample16.dontinline  avgt    3   1.711 ± 4.063  ns/op
Sample16.exclude     avgt    3  52.937 ± 5.268  ns/op
Sample16.inline      avgt    3   0.262 ± 0.027  ns/op
```
从执行结果可以看到内联方法和空方法执行速度一样，不编译执行最慢。
###### （17）JMHSample_17_SyncIterations
本例阐述了在多线程条件下，线程池的启动与销毁都会影响基准测试的准确性，如果自己来实现需要让线程同时开始启动工作，但这又比较难做到，如果在启动和关闭线程池时，无法做到同时，那么测量必定不准确，因为无法确定开始和结束时间；JMH提供了多线程基准测试的方法，先让线程池预热，都预热完成后让所有线程同时进行基准测试，测试完等待所有线程都结束再关闭线程池。
```
@State(Scope.Thread)
@OutputTimeUnit(TimeUnit.MILLISECONDS)
public class Sample17 {

    private double src;

    @Benchmark
    public double test() {
        double s = src;
        for (int i = 0; i < 1000; i++) {
            s = Math.sin(s);
        }
        return s;
    }

    public static void main(String[] args) throws RunnerException {
        Options opt = new OptionsBuilder()
                .include(Sample17.class.getSimpleName())
                .warmupTime(TimeValue.seconds(1))
                .measurementTime(TimeValue.seconds(1))
                .threads(Runtime.getRuntime().availableProcessors()*16)
                .forks(1)
                .syncIterations(true) // try to switch to "false"
                .build();

        new Runner(opt).run();
    }

}
```
这里warmupTime是预热时间，measurementTime是测量时间，threads是线程数，forks之前说过，是fork出一个子进程进行测试，syncIterations是是否需要同步预热，前面几个参数好理解，看了下代码才知道syncIterations如果设置为true代表等所有线程预热完成，然后所有线程一起进入测量阶段，等所有线程执行完测试后，再一起进入关闭；看一下设置为false时跑出的结果：
```
Benchmark       Mode  Cnt    Score    Error   Units
Sample17.test  thrpt    5  189.855 ± 23.415  ops/ms
```
再看一下为true的结果：
```
Benchmark       Mode  Cnt    Score    Error   Units
Sample17.test  thrpt    5  173.435 ± 41.659  ops/ms
```
当syncIterations设置为true时更准确地反应了多线程下被测试方法的性能，这个参数默认为true，无需手动设置。
###### （18）JMHSample_18_Control
本例介绍了使用 Control 的场景
```
@State(Scope.Group)
public class Sample18 {
    public final AtomicBoolean flag = new AtomicBoolean();

    @Benchmark
    @Group("pingpong")
    public void ping(Control cnt) {
        while (!cnt.stopMeasurement && !flag.compareAndSet(false, true)) {
            // this body is intentionally left blank
        }
    }

    @Benchmark
    @Group("pingpong")
    public void pong(Control cnt) {
        while (!cnt.stopMeasurement && !flag.compareAndSet(true, false)) {
            // this body is intentionally left blank
        }
    }
}
```
如果测试一个线程组对一个AtomicBoolean分别进行set true 和 set false操作，我们知道只有一个线程set true成功，另一个线程才能对其set false，否则另一个线程就陷入死锁，但我们的测试用例两个方法的执行不是均匀成对的，所以极大概率测试会陷入死锁，这时需要JMH提供的Control进行控制，当测量结束，双方都退出循环。
###### （20）JMHSample_20_Annotations
19也不存在。本例介绍了所有在main方法中通过Options提供的参数都可以通过注解写在需要测试的方法上，这在编写大量需要不同运行环境的基准测试时显得非常方便，比如这样
```
@Benchmark
@Warmup(iterations = 5, time = 100, timeUnit = TimeUnit.MILLISECONDS)
@Measurement(iterations = 5, time = 100, timeUnit = TimeUnit.MILLISECONDS)
public double measure() {
    return Math.log(x1);
}
```

###### （21）JMHSample_21_ConsumeCPU
本例介绍了Blackhole的另一种用法，前面提到的Blackhole可以用来干掉“死码消除”，同时Blackhole也可以“吞噬”cpu时间片（怪不得起这个名字）
```
@Benchmark
public void consume_0001() {
    Blackhole.consumeCPU(1);
}
```
Blackhole.consumeCPU的参数是时间片的tokens，和时间片成线性关系。
###### （22）JMHSample_22_FalseSharing
从名字可知本例是说false-sharing，即伪共享，关于伪共享这里就不再展开，后续会写一篇文章专门介绍伪共享。尽管JMH中提供的 @State 会自动填充缓存行，但是对于对象中特定的单个变量还是无法填充，所以本例介绍了4种方式来消除伪共享，需要我们自己注意。
###### （23）JMHSample_23_AuxCounters
本例介绍了注解 @AuxCounters ，它是一个辅助计数器，可以统计 @State 修饰的对象中的 public 属性被执行的情况，它的参数有两个
- AuxCounters.Type.EVENTS： 统计发生的次数
- AuxCounters.Type.OPERATIONS：按我们指定的格式统计，如按吞吐量统计
``` 
@OutputTimeUnit(TimeUnit.SECONDS)
@Warmup(iterations = 5, time = 1, timeUnit = TimeUnit.SECONDS)
@Measurement(iterations = 5, time = 1, timeUnit = TimeUnit.SECONDS)
@Fork(1)
public class Sample23 {
    
    @State(Scope.Thread)
    @AuxCounters(AuxCounters.Type.OPERATIONS)
    public static class OpCounters {
        // These fields would be counted as metrics
        public int case1;
        public int case2;

        // This accessor will also produce a metric
        public int total() {
            return case1 + case2;
        }
    }

    @State(Scope.Thread)
    @AuxCounters(AuxCounters.Type.EVENTS)
    public static class EventCounters {
        // This field would be counted as metric
        public int wows;
    }

    @Benchmark
    public void splitBranch(OpCounters counters) {
        if (Math.random() < 0.1) {
            counters.case1++;
        } else {
            counters.case2++;
        }
    }

    @Benchmark
    public void runSETI(EventCounters counters) {
        float random = (float) Math.random();
        float wowSignal = (float) Math.PI / 4;
        if (random == wowSignal) {
            // WOW, that's unusual.
            counters.wows++;
        }
    }
}
```
运行结果如下：
```
Benchmark                    Mode  Cnt         Score        Error  Units
Sample23.runSETI            thrpt    5  55524497.719 ± 815647.423  ops/s
Sample23.runSETI:wows       thrpt    5        19.000                   #
Sample23.splitBranch        thrpt    5  54429984.316 ± 355356.288  ops/s
Sample23.splitBranch:case1  thrpt    5   5443290.189 ±  29639.142  ops/s
Sample23.splitBranch:case2  thrpt    5  48986694.126 ± 326180.908  ops/s
Sample23.splitBranch:total  thrpt    5  54429984.316 ± 355356.288  ops/s
```
###### （24）JMHSample_24_Inheritance
本例介绍了JMH在父类中定义了Benchmark，子类也会继承并合成自己的Benchmark，这是JMH在编译期完成的，如果不是由JMH来编译，就无法享受这种继承。
```
public class Sample24 {

    @BenchmarkMode(Mode.AverageTime)
    @Fork(1)
    @State(Scope.Thread)
    @OutputTimeUnit(TimeUnit.NANOSECONDS)
    public static abstract class AbstractBenchmark {
        int x;

        @Setup
        public void setup() {
            x = 42;
        }

        @Benchmark
        @Warmup(iterations = 5, time = 100, timeUnit = TimeUnit.MILLISECONDS)
        @Measurement(iterations = 5, time = 100, timeUnit = TimeUnit.MILLISECONDS)
        public double bench() {
            return doWork() * doWork();
        }

        protected abstract double doWork();
    }

    public static class BenchmarkLog extends AbstractBenchmark {
        @Override
        protected double doWork() {
            return Math.log(x);
        }
    }

    public static class BenchmarkSin extends AbstractBenchmark {
        @Override
        protected double doWork() {
            return Math.sin(x);
        }
    }

    public static class BenchmarkCos extends AbstractBenchmark {
        @Override
        protected double doWork() {
            return Math.cos(x);
        }
    }
}
```
运行结果：
```
Benchmark                    Mode  Cnt    Score    Error  Units
Sample24.BenchmarkCos.bench  avgt    5  104.208 ±  5.578  ns/op
Sample24.BenchmarkLog.bench  avgt    5   23.998 ±  0.464  ns/op
Sample24.BenchmarkSin.bench  avgt    5  137.769 ± 22.275  ns/op
```
###### （25）JMHSample_25_API_GA
本例介绍了如何用API的方式来写基准测试代码，比较复杂，个人觉得还是注解的方式简单。
###### （26）JMHSample_26_BatchSize
本例介绍了 batchSize，指定一次迭代方法需要执行 batchSize 次
```
@Benchmark
@Warmup(iterations = 5, batchSize = 5000)
@Measurement(iterations = 5, batchSize = 5000)
@BenchmarkMode(Mode.SingleShotTime)
public List<String> measureRight() {
    list.add(list.size() / 2, "something");
    return list;
}
```
###### （27）JMHSample_27_Params
本例介绍 @Param ，@Param 允许使用一份基准测试代码跑多组数据，特别适合测量方法性能和参数取值的关系
```
@BenchmarkMode(Mode.AverageTime)
@OutputTimeUnit(TimeUnit.NANOSECONDS)
@Warmup(iterations = 5, time = 1, timeUnit = TimeUnit.SECONDS)
@Measurement(iterations = 5, time = 1, timeUnit = TimeUnit.SECONDS)
@Fork(1)
@State(Scope.Benchmark)
public class Sample27 {

    @Param({"1", "31", "65", "101", "103"})
    public int arg;

    @Param({"0", "1", "2", "4", "8", "16", "32"})
    public int certainty;

    @Benchmark
    public boolean bench() {
        return BigInteger.valueOf(arg).isProbablePrime(certainty);
    }
}
```
运行结果：
```
Benchmark       (arg)  (certainty)  Mode  Cnt     Score     Error  Units
Sample27.bench      1            0  avgt    5     4.135 ±   0.177  ns/op
Sample27.bench      1            1  avgt    5     6.327 ±   0.366  ns/op
Sample27.bench      1            2  avgt    5     6.212 ±   0.154  ns/op
Sample27.bench      1            4  avgt    5     6.572 ±   0.215  ns/op
Sample27.bench      1            8  avgt    5     6.842 ±   1.713  ns/op
Sample27.bench      1           16  avgt    5     6.238 ±   0.297  ns/op
Sample27.bench      1           32  avgt    5     6.558 ±   1.604  ns/op
Sample27.bench     31            0  avgt    5     4.258 ±   0.084  ns/op
Sample27.bench     31            1  avgt    5   466.786 ±  14.588  ns/op
Sample27.bench     31            2  avgt    5   470.098 ±  10.044  ns/op
Sample27.bench     31            4  avgt    5   902.594 ±  49.249  ns/op
Sample27.bench     31            8  avgt    5  1768.349 ± 262.570  ns/op
Sample27.bench     31           16  avgt    5  3352.993 ±  96.689  ns/op
Sample27.bench     31           32  avgt    5  6649.535 ± 438.881  ns/op
Sample27.bench     65            0  avgt    5     4.268 ±   0.115  ns/op
Sample27.bench     65            1  avgt    5  1092.193 ±  18.749  ns/op
Sample27.bench     65            2  avgt    5  1083.709 ±  47.573  ns/op
Sample27.bench     65            4  avgt    5  1177.795 ±  38.753  ns/op
Sample27.bench     65            8  avgt    5  1174.277 ±  10.592  ns/op
Sample27.bench     65           16  avgt    5  1173.642 ±  11.522  ns/op
Sample27.bench     65           32  avgt    5  1181.138 ±  18.386  ns/op
Sample27.bench    101            0  avgt    5     4.303 ±   0.345  ns/op
Sample27.bench    101            1  avgt    5   619.921 ±  16.203  ns/op
Sample27.bench    101            2  avgt    5   619.270 ±  29.136  ns/op
Sample27.bench    101            4  avgt    5  1204.741 ±  23.560  ns/op
Sample27.bench    101            8  avgt    5  2376.725 ±  54.460  ns/op
Sample27.bench    101           16  avgt    5  4640.182 ± 104.423  ns/op
Sample27.bench    101           32  avgt    5  9006.358 ± 176.451  ns/op
Sample27.bench    103            0  avgt    5     4.247 ±   0.120  ns/op
Sample27.bench    103            1  avgt    5   562.567 ±  14.095  ns/op
Sample27.bench    103            2  avgt    5   571.643 ±  44.721  ns/op
Sample27.bench    103            4  avgt    5  1091.065 ±  36.354  ns/op
Sample27.bench    103            8  avgt    5  2100.759 ±  31.438  ns/op
Sample27.bench    103           16  avgt    5  4220.375 ±  42.815  ns/op
Sample27.bench    103           32  avgt    5  8336.839 ±  58.304  ns/op
```
###### （28）JMHSample_28_BlackholeHelpers
本例介绍了 Blackhole 不仅可以用在 Benchmark 修饰的方法上，也可以用在其他JMH提供的方法上，比如 @Setup 和 @TearDown 等等
```
@Setup
public void setup(final Blackhole bh) {
    workerBaseline = new Worker() {
        double x;

        @Override
        public void work() {
            // do nothing
        }
    };

    workerWrong = new Worker() {
        double x;

        @Override
        public void work() {
            Math.log(x);
        }
    };

    workerRight = new Worker() {
        double x;

        @Override
        public void work() {
            bh.consume(Math.log(x));
        }
    };

}
```
###### （29）JMHSample_29_StatesDAG
本例展示了 @State 中嵌套 @State 的情况，不过想不出为啥需要这样做，例子中也说这是个实验性质的Feature。
###### （30）JMHSample_30_Interrupts
本例类似（18），也是在测试时遇到死锁的情况，JMH可对能被Interrupt的超时动作发起Interrupt动作，让测试正常结束。如对一个BlockQueue进行put和take操作，如果put和get不对称，就会死锁，此时JMH会自行打断。
```
@BenchmarkMode(Mode.AverageTime)
@OutputTimeUnit(TimeUnit.NANOSECONDS)
@State(Scope.Group)
public class Sample30 {
    private BlockingQueue<Integer> q;

    @Setup
    public void setup() {
        q = new ArrayBlockingQueue<Integer>(1);
    }

    @Group("Q")
    @Benchmark
    public Integer take() throws InterruptedException {
        return q.take();
    }

    @Group("Q")
    @Benchmark
    public void put() throws InterruptedException {
        q.put(42);
    }
}
```
运行时会提示被Interrupt：
```
Iteration   2: (*interrupt*) 8604.342 ns/op
                 put:  8604.338 ns/op
                 take: 8604.345 ns/op
```

###### （31）JMHSample_31_InfraParams
本例介绍了在方法中可覆盖的三种参数，这给在测试时获取配置以及动态修改配置提供了可能
- BenchmarkParams：基准测试级别
- IterationParams：迭代级别
- ThreadParams：线程级别
```
@BenchmarkMode(Mode.Throughput)
@OutputTimeUnit(TimeUnit.SECONDS)
@State(Scope.Benchmark)
public class Sample31 {
    
    static final int THREAD_SLICE = 1000;

    private ConcurrentHashMap<String, String> mapSingle;
    private ConcurrentHashMap<String, String> mapFollowThreads;

    @Setup
    public void setup(BenchmarkParams params) {
        int capacity = 16 * THREAD_SLICE * params.getThreads();
        mapSingle        = new ConcurrentHashMap<String, String>(capacity, 0.75f, 1);
        mapFollowThreads = new ConcurrentHashMap<String, String>(capacity, 0.75f, params.getThreads());
    }

    @State(Scope.Thread)
    public static class Ids {
        private List<String> ids;

        @Setup
        public void setup(ThreadParams threads) {
            ids = new ArrayList<String>();
            for (int c = 0; c < THREAD_SLICE; c++) {
                ids.add("ID" + (THREAD_SLICE * threads.getThreadIndex() + c));
            }
        }
    }

    @Benchmark
    public void measureDefault(Ids ids) {
        for (String s : ids.ids) {
            mapSingle.remove(s);
            mapSingle.put(s, s);
        }
    }

    @Benchmark
    public void measureFollowThreads(Ids ids) {
        for (String s : ids.ids) {
            mapFollowThreads.remove(s);
            mapFollowThreads.put(s, s);
        }
    }
}
```
###### （32）JMHSample_32_BulkWarmup
本例介绍了三种预热方式：
- WarmupMode.INDI：每个Benchmark单独预热
- WarmupMode.BULK：在每个Benchmark执行前都预热所有的Benchmark
- WarmupMode.BULK_INDI：在每个Benchmark执行前都预热所有的Benchmark，且需要再预热本次执行的Benchmark
###### （33）JMHSample_33_SecurityManager
关于 SecurityManager ，可以通过注入的方式来修改 SecurityManager，用处不大
###### （34）JMHSample_34_SafeLooping
这节专门讲如何构造安全的循环，和前面的内容稍有重复
```
@State(Scope.Thread)
@Warmup(iterations = 5, time = 1, timeUnit = TimeUnit.SECONDS)
@Measurement(iterations = 5, time = 1, timeUnit = TimeUnit.SECONDS)
@Fork(3)
@BenchmarkMode(Mode.AverageTime)
@OutputTimeUnit(TimeUnit.NANOSECONDS)
public class Sample34 {

    static final int BASE = 42;

    static int work(int x) {
        return BASE + x;
    }

    @Param({"1", "10", "100", "1000"})
    int size;

    int[] xs;

    @Setup
    public void setup() {
        xs = new int[size];
        for (int c = 0; c < size; c++) {
            xs[c] = c;
        }
    }

    @Benchmark
    public int measureWrong_1() {
        int acc = 0;
        for (int x : xs) {
            acc = work(x);
        }
        return acc;
    }

    @Benchmark
    public int measureWrong_2() {
        int acc = 0;
        for (int x : xs) {
            acc += work(x);
        }
        return acc;
    }

    @Benchmark
    public void measureRight_1(Blackhole bh) {
        for (int x : xs) {
            bh.consume(work(x));
        }
    }

    @Benchmark
    public void measureRight_2() {
        for (int x : xs) {
            sink(work(x));
        }
    }

    @CompilerControl(CompilerControl.Mode.DONT_INLINE)
    public static void sink(int v) {
        // IT IS VERY IMPORTANT TO MATCH THE SIGNATURE TO AVOID AUTOBOXING.
        // The method intentionally does nothing.
    }
}
```
measureWrong_1和measureWrong_2都可能在编译时被展开，后面两种正确的方式在前面都提过。

###### （35）JMHSample_35_Profilers
本例讲解了如何使用JMH内置的性能剖析工具查看基准测试消耗在什么地方，具体的剖析方式内置的有如下几种：
- ClassloaderProfiler：类加载剖析
- CompilerProfiler：JIT编译剖析
- GCProfiler：GC剖析
- StackProfiler：栈剖析
- PausesProfiler：停顿剖析
- HotspotThreadProfiler：Hotspot线程剖析
- HotspotRuntimeProfiler：Hotspot运行时剖析
- HotspotMemoryProfiler：Hotspot内存剖析
- HotspotCompilationProfiler：Hotspot编译剖析
- HotspotClassloadingProfiler：Hotspot 类加载剖析

用法
```
public static void main(String[] args) throws RunnerException {
    Options opt = new OptionsBuilder()
        .include(JMHSample_35_Profilers.Maps.class.getSimpleName())
        .addProfiler(StackProfiler.class)
//                    .addProfiler(GCProfiler.class)
        .build();

    new Runner(opt).run();
}
```
###### （36）JMHSample_36_BranchPrediction
本例提醒我们要注意“分支预测”，简单来说，分支预测是CPU在处理有规律的数据比没有规律的数据要快，CPU可以“预测”这种规律。我们在基准测试时需要注意样本数据的规律性对结果也会产生影响。
###### （37）JMHSample_37_CacheAccess
本例提醒我们对内存的顺序访问与非顺序访问会对测试结果产生影响，这点也是因为CPU存在缓存行的缘故，与之前提到的伪共享类似。
###### （38）JMHSample_38_PerInvokeSetup
本例也是之前提到过的在每次调用前执行Setup，使用它可以测量排序性能，如果Setup不在每次执行排序时执行，那么只有第一次排序是执行了排序，后面每次排序的都是相同顺序的数据。
##总结
通过这么多例子我们概括出为什么需要JMH
- 方便：使用方便，配置一些注解即可测试，且测量纬度、内置的工具丰富；
- 专业：JMH自动地帮我们避免了一些测试上的“坑”，避免不了的也在例子中给出了说明；
- 准确：预热、fork隔离、避免方法内联、避免常量折叠等很多方法可以提高测量的准确性。