# 强悍的LongAdder

LongAdder是jdk8引入的适用于统计场景的线程安全的计数器。

在此之前，实现一款线程安全的计数器要么加锁，要么使用AtomicLong，加锁性能必然很差，AtomicLong性能要好很多，但是在高并发、多线程下，也显得吃力。于是就有了LongAdder，LongAdder有两个重要的方法：add和sum，add是线程安全的加，sum是返回结果，之所以叫sum是因为LongAdder通过分段的思想维护了一组变量，多线程并发更新时被散列到不同的变量上执行，减少冲突，所以最后获取返回值是将这些变量求和。通过这点也能看出sum获取的结果是不准确的，所以它只适用于统计场景，如果要获取精确的返回值，还是得用AtomicLong，性能和准确不可兼得。

通过JMH测试LongAdder、AtomicLong以及加锁的计数器的性能，感受一下LongAdder的强大。（如无特殊说明，本文后续JMH测试均以此为标准：fork1进程，4线程，预热2次，正式测量2次，测试机器4核，完整代码已上传github，文末有地址）

```java
private final AtomicLong atomicLong = new AtomicLong();
    private final LongAdder longAdder = new LongAdder();
    private long counter = 0;

    public static void main(String[] args) throws RunnerException {
        Options opt = new OptionsBuilder()
                .include(LongAdderTest.class.getSimpleName())
                .forks(1)
                .threads(4)
                .warmupIterations(2)
                .measurementIterations(2)
                .mode(Mode.Throughput)
                .syncIterations(false)
                .build();

        new Runner(opt).run();
    }

    @Benchmark
    public void testAtomic() {
        atomicLong.incrementAndGet();
    }

    @Benchmark
    public void testLongAdder() {
        longAdder.increment();
    }

    @Benchmark
    public synchronized void testLockAdder() {
        counter++;
    }
```
运行后
```
Benchmark                          Mode  Cnt          Score   Error  Units
LongAdderTest.testAtomic          thrpt    2   73520672.658          ops/s
LongAdderTest.testLockAdder       thrpt    2   23456856.867          ops/s
LongAdderTest.testLongAdder       thrpt    2  300013067.345          ops/s
```
可以看到LongAdder和另外两种实现完全不在一个量级上，性能及其恐怖。既然知道LongAdder的大致原理，那我们能不能实现一个MyLongAdder，保证写入线程安全的同时，性能比肩甚至超越LongAdder呢？

# AtomicLong分段(V0)
性能优化中很多都是依靠LongAdder这种分段的方式，如ConcurrentHashMap就是采用分段锁，于是我们也实现一个V0版本的MyLongAdder
```java
public class MyLongAdderV0 {

    private final int coreSize;
    private final AtomicLong[] counts;

    public MyLongAdderV0(int coreSize) {
        this.coreSize = coreSize;
        this.counts = new AtomicLong[coreSize];
        for (int i = 0; i < coreSize; i++) {
            this.counts[i] = new AtomicLong();
        }
    }

    public void increment() {
        int index = (int) (Thread.currentThread().getId() % coreSize);
        counts[index].incrementAndGet();
    }
}
```
使用一个AtomicLong数组，线程执行时，按线程id散列开，coreSize这里期望是cpu核数，和LongAdder、AtomicLong对比一下看看（测试代码省略，后同）

```
Benchmark                         Mode  Cnt          Score   Error  Units
LongAdderTest.testAtomic         thrpt    2   73391661.579          ops/s
LongAdderTest.testLongAdder      thrpt    2  309539056.885          ops/s
LongAdderTest.testMyLongAdderV0  thrpt    2   83737867.380          ops/s
```
emmm，V0性能仅仅比AtomicLong好一点点，跟LongAdder还是不在一个量级上，难道是数组不够大？将coreSize作为参数，测试一下 4, 8, 16, 32的情况，我测试了好几次，每次结果都不一样但又差不多在一个量级（偶尔会上亿），无法总结结果与coreSize的关系，这里给出其中一组

```
Benchmark                        (coreSize)   Mode  Cnt          Score   Error  Units
LongAdderTest.testMyLongAdderV0           4  thrpt    2   62328997.667          ops/s
LongAdderTest.testMyLongAdderV0           8  thrpt    2  124725716.902          ops/s
LongAdderTest.testMyLongAdderV0          16  thrpt    2   84718415.566          ops/s
LongAdderTest.testMyLongAdderV0          32  thrpt    2   85321816.442          ops/s
```
猜想是因为依赖了线程的id，分散的不够均匀导致，而且还有一个有意思的情况，有时候V0居然比AtomicLong的性能还低。

# 取模优化(V1)
注意到V0里面有一个取模的操作，这个操作可能比较耗时，可能会导致V0的性能甚至不如单个AtomicLong，可以通过移位操作来代替，但代替的前提是coreSize必须为2的n次方，如2，4，8，16（我们假定后续coreSize只取2的n次方），V1版本的代码如下：

```java
public class MyLongAdderV1 {

    private final int coreSize;
    private final AtomicLong[] counts;

    public MyLongAdderV1(int coreSize) {
        this.coreSize = coreSize;
        this.counts = new AtomicLong[coreSize];
        for (int i = 0; i < coreSize; i++) {
            this.counts[i] = new AtomicLong();
        }
    }

    public void increment() {
        int index = (int) (Thread.currentThread().getId() & (coreSize - 1));
        counts[index].incrementAndGet();
    }

}
```
测试一下性能
```
Benchmark                         Mode  Cnt          Score   Error  Units
LongAdderTest.testLongAdder      thrpt    2  312683635.190          ops/s
LongAdderTest.testMyLongAdderV0  thrpt    2   60641758.648          ops/s
LongAdderTest.testMyLongAdderV1  thrpt    2  100887869.829          ops/s
```
性能稍微好了一点，但是跟LongAdder比还是差了一大截

# 消除伪共享(V2)

在cpu面前内存太慢了，所以cpu有三级缓存 L3，L2，L1。L1最接近cpu，速度也最快，cpu查找的顺序是先L1，再L2，再L3，最后取不到会去内存取。通常来说每个缓存由很多缓存行组成，缓存行通常是64个字节，java的long是8字节，因此一个缓存行可以缓存8个long变量。如果多个核的线程在操作同一个缓存行中的不同变量数据，那么就会出现频繁的缓存失效，即使在代码层面看这两个线程操作的数据之间完全没有关系。这种不合理的资源竞争情况学名伪共享（False Sharing），会严重影响机器的并发执行效率。

在V1中，AtomicLong中有一个value，每次incrementAndGet会改变这个value，同时AtomicLong是一个数组，数组的内存地址也是连续的，这样就会导致相邻的AtomicLong的value缓存失效，其他线程读取这个value就会变得很慢。优化的方法就是填充AtomicLong，让每个AtomicLong的value相互隔离，不要相互影响。

通常填充缓存行有如下几种方式：

- （1）java8可以在类属性上使用 @sun.misc.Contended，jvm参数需要指定-XX:-RestrictContended
- （2）使用继承的方式在类的属性前后插入变量实现，这里举一个通过继承来实现的，如果不用继承，这些填充的无用变量会被编译器优化掉，当然也可以通过数组来构造填充，这里就不多说。

```java
abstract class RingBufferPad {
    protected long p1, p2, p3, p4, p5, p6, p7;
}

abstract class RingBufferFields<E> extends RingBufferPad {
    protected long value;
}

public final class RingBuffer<E> extends RingBufferFields<E> {
    protected long p1, p2, p3, p4, p5, p6, p7;
}
```
我们直接用java8的`@sun.misc.Contended`来对V1进行优化

```java
public class MyLongAdderV2 {

    private static class AtomicLongWrap {
        @Contended
        private final AtomicLong value = new AtomicLong();
    }

    private final int coreSize;
    private final AtomicLongWrap[] counts;

    public MyLongAdderV2(int coreSize) {
        this.coreSize = coreSize;
        this.counts = new AtomicLongWrap[coreSize];
        for (int i = 0; i < coreSize; i++) {
            this.counts[i] = new AtomicLongWrap();
        }
    }

    public void increment() {
        int index = (int) (Thread.currentThread().getId() & (coreSize - 1));
        counts[index].value.incrementAndGet();
    }

}
```
执行后神奇的情况出现了

```
Benchmark                         Mode  Cnt          Score   Error  Units
LongAdderTest.testLongAdder      thrpt    2  272733686.330          ops/s
LongAdderTest.testMyLongAdderV2  thrpt    2  307754425.667          ops/s
```
居然V2版本比LongAdder还快！但这是真的吗？为此，我多测试了几组，分别在线程数为8时

```
Benchmark                         Mode  Cnt          Score   Error  Units
LongAdderTest.testLongAdder      thrpt    2  260909722.754          ops/s
LongAdderTest.testMyLongAdderV2  thrpt    2  215785206.276          ops/s
```
线程数为16时：

```
Benchmark                         Mode  Cnt          Score   Error  Units
LongAdderTest.testLongAdder      thrpt    2  307269737.067          ops/s
LongAdderTest.testMyLongAdderV2  thrpt    2  185774540.302          ops/s
```
发现随着线程数的增加，V2的性能越来越低，但LongAdder纹丝不动，不得不佩服写jdk的大佬。

# 改进hash算法

V0到V2版本均使用了线程id作为hash值来散列到不同的槽点，线程id生成后不会改变，这样就会导致每次执行的测试可能结果都不太一样，如果比较聚焦，性能必然会很差，当线程数增多后必然会造成更多的冲突，有没有更好的hash算法？

- 尝试hashCode
java的每个对象都有一个hashCode，我们使用线程对象的hashCode来散列试试，版本V3关键改动如下

```java
public void increment() {
    int index = Thread.currentThread().hashCode() & (coreSize - 1);
    counts[index].incrementAndGet();
}
```
结果如下
```
Benchmark                         Mode  Cnt          Score   Error  Units
LongAdderTest.testLongAdder      thrpt    2  277084413.669          ops/s
LongAdderTest.testMyLongAdderV3  thrpt    2  103351246.650          ops/s
```
性能似乎不尽如人意。

- 尝试随机数

当然使用Random当然不行，用性能更好的ThreadLocalRandom，V4版本关键改动如下

```java
public void increment() {
      counts[ThreadLocalRandom.current().nextInt(coreSize)].value.incrementAndGet();
  }
```
结果如下
```
Benchmark                         Mode  Cnt          Score   Error  Units
LongAdderTest.testLongAdder      thrpt    2  292807355.101          ops/s
LongAdderTest.testMyLongAdderV4  thrpt    2   95200307.226          ops/s
```
性能也上不去，猜想是因为生成随机数比较耗时。

- 冲突时重新计算hash

为了优化V4版本，参考了LongAdder，算是一个黑科技，生成一个随机数存在Thread对象中，可以看一下Thread类，刚好有这个变量

```java
/** Probe hash value; nonzero if threadLocalRandomSeed initialized */
@sun.misc.Contended("tlr")
int threadLocalRandomProbe;
```

但是这个变量是不对外开放，只能通过反射（性能太差）或者UNSAFE来取，它在 ThreadLocalRandomSeed 中被初始化，发生冲突时重新生成并修改它（生成的方法可以参考ThreadLocalRandomSeed），也是通过UNSAFE可以搞定。既然要在冲突时重新hash，那必须能检测出冲突，AtomicLong就不能用incrementAndGet了，使用AtomicLong的compareAndSet方法，返回false时代表有冲突，冲突时重新hash，并用incrementAndGet兜底，保证一定能成功。如此一来，既可以均匀地散列开，也能保证随机数生成的效率。V5版本代码如下

```java
public class MyLongAdderV5 {

    private static sun.misc.Unsafe UNSAFE = null;
    private static final long PROBE;
    static {

        try {
            // 反射获取unsafe
            Field f = Unsafe.class.getDeclaredField("theUnsafe");
            f.setAccessible(true);
            UNSAFE = (Unsafe) f.get(null);
        } catch (Exception e) {

        }

        try {
            Class<?> tk = Thread.class;
            PROBE = UNSAFE.objectFieldOffset
                    (tk.getDeclaredField("threadLocalRandomProbe"));
        } catch (Exception e) {
            throw new Error(e);
        }
    }

    static final int getProbe() {
        // 获取thread的threadLocalRandomProbe属性值
        return UNSAFE.getInt(Thread.currentThread(), PROBE);
    }

    static final int advanceProbe(int probe) {
        // 重新生成随机数并写入thread对象
        probe ^= probe << 13;   // xorshift
        probe ^= probe >>> 17;
        probe ^= probe << 5;
        UNSAFE.putInt(Thread.currentThread(), PROBE, probe);
        return probe;
    }

    private static class AtomicLongWrap {
        @Contended
        private final AtomicLong value = new AtomicLong();
    }

    private final int coreSize;
    private final AtomicLongWrap[] counts;

    public MyLongAdderV5(int coreSize) {
        this.coreSize = coreSize;
        this.counts = new AtomicLongWrap[coreSize];
        for (int i = 0; i < coreSize; i++) {
            this.counts[i] = new AtomicLongWrap();
        }
    }

    public void increment() {

        int h = getProbe();

        int index = getProbe() & (coreSize - 1);
        long r;
        if (!counts[index].value.compareAndSet(r = counts[index].value.get(), r + 1)) {
            if (h == 0) {
                // 初始化随机数
                ThreadLocalRandom.current();
                h = getProbe();
            }
            // 冲突后重新生成随机数
            advanceProbe(h);
            // 用getAndIncrement来兜底
            counts[index].value.getAndIncrement();
        }
    }

}
```
结果如下

```
Benchmark                         Mode  Cnt          Score   Error  Units
LongAdderTest.testLongAdder      thrpt    2  274131797.300          ops/s
LongAdderTest.testMyLongAdderV5  thrpt    2  298402832.456          ops/s
```

效果还可以，试试8线程：

```
Benchmark                         Mode  Cnt          Score   Error  Units
LongAdderTest.testLongAdder      thrpt    2  324982482.774          ops/s
LongAdderTest.testMyLongAdderV5  thrpt    2  290476796.289          ops/s
```

16线程：

```
Benchmark                         Mode  Cnt          Score   Error  Units
LongAdderTest.testLongAdder      thrpt    2  291180444.998          ops/s
LongAdderTest.testMyLongAdderV5  thrpt    2  282745610.470          ops/s
```

32线程：

```
Benchmark                         Mode  Cnt          Score   Error  Units
LongAdderTest.testLongAdder      thrpt    2  294237473.396          ops/s
LongAdderTest.testMyLongAdderV5  thrpt    2  301187346.873          ops/s
```

果然这个方法很牛皮，无论在多少个线程下都能稳如🐶。

# 总结

实现一款超越LongAdder性能的多线程计数器非常难，折腾了两天也只是达到和LongAdder相当的性能，其中对性能影响最大的几个改动点是

- 分段：基础优化，一般人都能想到
- 取模优化：也比较基础
- 消除伪共享：这个优化提升很大
- hash算法：这条保证了稳定性，无论多少线程都是最高吞吐量

其中前三条比较常规，第四条可以算得上是`黑科技`

---
所有的测试代码已上传 https://github.com/lkxiaolou/all-in-one/tree/master/src/main/java/org/newboo/longadder


---

欢迎关注我的公众号

![捉虫大师](../../qrcode_small.jpg)

- 原文链接: https://mp.weixin.qq.com/s/N1scBcRr3Zz6kzanMNyN0Q
- 发布时间: 2020.05.19















