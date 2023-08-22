package logbot

import (
	"testing"

	"github.com/curtisnewbie/gocommon/common"
)

func TestParseLine(t *testing.T) {
	line := `2023-06-13 12:58:35.509 INFO  [                ,                ] consul.DeregisterService      : Deregistering current instance on Consul, service_id: 'goauth-8081'`
	logLine, err := parseLogLine(common.EmptyRail(), line, "go")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", logLine)

	line = `2023-06-13 22:16:13.746 ERROR [v2geq7340pbfxcc9,k1gsschfgarpc7no] main.registerWebEndpoints.func2 : Oh on!
continue on a new line :D`
	logLine, err = parseLogLine(common.EmptyRail(), line, "go")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", logLine)

	line = `2023-06-14 09:50:30.500 DEBUG [ptqnta70npjjxfz8,114lkur90ui6ywqt] redis.TimedRLockRun.func1     : Released lock for key 'rcache:POST:/goauth/open/api/path/update'
`
	logLine, err = parseLogLine(common.EmptyRail(), line, "go")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", logLine)

	line = `2023-06-16 16:29:48.811 INFO  [3pdac13hagg9v8cs,pskalvoqmaets17f] server.BootstrapServer        :



---------------------------------------------- goauth started (took: 59ms) --------------------------------------------

`

	logLine, err = parseLogLine(common.EmptyRail(), line, "go")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", logLine)

	line = `2023-06-17 17:34:48.762  INFO [auth-service,,,] 78446 --- [           main] .c.m.r.c.YamlBasedRedissonClientProvider : Loading RedissonClient from yaml config file, reading environment property: redisson-config`

	logLine, err = parseLogLine(common.EmptyRail(), line, "java")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", logLine)

	line = `2023-06-17 18:01:11.489 ERROR [auth-service,,,] 84063 --- [onPool-worker-1] c.c.goauth.client.RestPathReporter       : Failed to report path to goauth, req: AddPathReq(type=PROTECTED, url=/auth-service/error, method=TRACE, group=auth-service, desc=, resCode=)

java.lang.RuntimeException: com.netflix.client.ClientException: Load balancer does not have available server for client: goauth
	at org.springframework.cloud.openfeign.ribbon.LoadBalancerFeignClient.execute(LoadBalancerFeignClient.java:90)
	at org.springframework.cloud.sleuth.instrument.web.client.feign.TraceLoadBalancerFeignClient.execute(TraceLoadBalancerFeignClient.java:78)
	at feign.SynchronousMethodHandler.executeAndDecode(SynchronousMethodHandler.java:119)
	at feign.SynchronousMethodHandler.invoke(SynchronousMethodHandler.java:89)
	at feign.ReflectiveFeign$FeignInvocationHandler.invoke(ReflectiveFeign.java:100)
	at jdk.proxy2/jdk.proxy2.$Proxy178.addPath(Unknown Source)
	at com.curtisnewbie.goauth.client.RestPathReporter.reportPath(RestPathReporter.java:92)
	at com.curtisnewbie.goauth.client.RestPathReporter.lambda$reportPaths$11(RestPathReporter.java:87)
	at java.base/java.util.stream.ForEachOps$ForEachOp$OfRef.accept(ForEachOps.java:183)
	at java.base/java.util.stream.ReferencePipeline$3$1.accept(ReferencePipeline.java:197)
	at java.base/java.util.ArrayList$ArrayListSpliterator.forEachRemaining(ArrayList.java:1625)
	at java.base/java.util.stream.AbstractPipeline.copyInto(AbstractPipeline.java:509)
	at java.base/java.util.stream.AbstractPipeline.wrapAndCopyInto(AbstractPipeline.java:499)
	at java.base/java.util.stream.ForEachOps$ForEachOp.evaluateSequential(ForEachOps.java:150)
	at java.base/java.util.stream.ForEachOps$ForEachOp$OfRef.evaluateSequential(ForEachOps.java:173)
	at java.base/java.util.stream.AbstractPipeline.evaluate(AbstractPipeline.java:234)
	at java.base/java.util.stream.ReferencePipeline.forEach(ReferencePipeline.java:596)
	at com.curtisnewbie.goauth.client.RestPathReporter.reportPaths(RestPathReporter.java:87)
	at com.curtisnewbie.goauth.client.RestPathReporter.lambda$afterPropertiesSet$1(RestPathReporter.java:51)
	at java.base/java.util.concurrent.CompletableFuture$AsyncRun.run(CompletableFuture.java:1804)
	at java.base/java.util.concurrent.CompletableFuture$AsyncRun.exec(CompletableFuture.java:1796)
	at java.base/java.util.concurrent.ForkJoinTask.doExec(ForkJoinTask.java:387)
	at java.base/java.util.concurrent.ForkJoinPool$WorkQueue.topLevelExec(ForkJoinPool.java:1311)
	at java.base/java.util.concurrent.ForkJoinPool.scan(ForkJoinPool.java:1840)
	at java.base/java.util.concurrent.ForkJoinPool.runWorker(ForkJoinPool.java:1806)
	at java.base/java.util.concurrent.ForkJoinWorkerThread.run(ForkJoinWorkerThread.java:177)
Caused by: com.netflix.client.ClientException: Load balancer does not have available server for client: goauth
	at com.netflix.loadbalancer.LoadBalancerContext.getServerFromLoadBalancer(LoadBalancerContext.java:483)
	at com.netflix.loadbalancer.reactive.LoadBalancerCommand$1.call(LoadBalancerCommand.java:184)
	at com.netflix.loadbalancer.reactive.LoadBalancerCommand$1.call(LoadBalancerCommand.java:180)
	at rx.Observable.unsafeSubscribe(Observable.java:10327)
	at rx.internal.operators.OnSubscribeConcatMap.call(OnSubscribeConcatMap.java:94)
	at rx.internal.operators.OnSubscribeConcatMap.call(OnSubscribeConcatMap.java:42)
	at rx.Observable.unsafeSubscribe(Observable.java:10327)
	at rx.internal.operators.OperatorRetryWithPredicate$SourceSubscriber$1.call(OperatorRetryWithPredicate.java:127)
	at rx.internal.schedulers.TrampolineScheduler$InnerCurrentThreadScheduler.enqueue(TrampolineScheduler.java:73)
	at rx.internal.schedulers.TrampolineScheduler$InnerCurrentThreadScheduler.schedule(TrampolineScheduler.java:52)
	at rx.internal.operators.OperatorRetryWithPredicate$SourceSubscriber.onNext(OperatorRetryWithPredicate.java:79)
	at rx.internal.operators.OperatorRetryWithPredicate$SourceSubscriber.onNext(OperatorRetryWithPredicate.java:45)
	at rx.internal.util.ScalarSynchronousObservable$WeakSingleProducer.request(ScalarSynchronousObservable.java:276)
	at rx.Subscriber.setProducer(Subscriber.java:209)
	at rx.internal.util.ScalarSynchronousObservable$JustOnSubscribe.call(ScalarSynchronousObservable.java:138)
	at rx.internal.util.ScalarSynchronousObservable$JustOnSubscribe.call(ScalarSynchronousObservable.java:129)
	at rx.internal.operators.OnSubscribeLift.call(OnSubscribeLift.java:48)
	at rx.internal.operators.OnSubscribeLift.call(OnSubscribeLift.java:30)
	at rx.internal.operators.OnSubscribeLift.call(OnSubscribeLift.java:48)
	at rx.internal.operators.OnSubscribeLift.call(OnSubscribeLift.java:30)
	at rx.internal.operators.OnSubscribeLift.call(OnSubscribeLift.java:48)
	at rx.internal.operators.OnSubscribeLift.call(OnSubscribeLift.java:30)
	at rx.Observable.subscribe(Observable.java:10423)
	at rx.Observable.subscribe(Observable.java:10390)
	at rx.observables.BlockingObservable.blockForSingle(BlockingObservable.java:443)
	at rx.observables.BlockingObservable.single(BlockingObservable.java:340)
	at com.netflix.client.AbstractLoadBalancerAwareClient.executeWithLoadBalancer(AbstractLoadBalancerAwareClient.java:112)
	at org.springframework.cloud.openfeign.ribbon.LoadBalancerFeignClient.execute(LoadBalancerFeignClient.java:83)
	... 25 common frames omitted`

	logLine, err = parseLogLine(common.EmptyRail(), line, "java")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", logLine)
}
