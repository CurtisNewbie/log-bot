# logbot

Bot for watching and parsing logs

## Requirements

- MySQL
- Redis
- Consul
- RabbitMQ
- [Goauth](https://github.com/CurtisNewbie/goauth)

## Configuration

| Property              | Description                                          | Default Value |
|-----------------------|------------------------------------------------------|---------------|
| logbot.node           | name of the node                                     | 'default'     |
| logbot.[]watch        | list of watch config                                 |               |
| logbot.[]watch.app    | app name                                             |               |
| logbot.[]watch.file   | path of the log file                                 |               |
| logbot.[]watch.type   | type of log pattern [ 'go', 'java' ]                 | 'go'          |
| task.remove-error-log | enable task to remove error logs reported 7 days ago | false         |

## API

| Method | Path              | Parameters                                                        | Description               |
|--------|-------------------|-------------------------------------------------------------------|---------------------------|
| POST   | `/log/error/list` | `{ "app" : "app_name" , "page" : { "limit" : 10, "page"  : 1 } }` | List ERROR logs collected |

## Log Pattern

There are two log patterns available:

### Pattern for go

Logs are parsed using following regex:

```
`^([0-9]{4}\-[0-9]{2}\-[0-9]{2} [0-9:\.]+) +(\w+) +\[([\w ]+),([\w ]+)\] ([\w\.]+) +: *((?s).*)`
```

The log looks like this:

```log
2023-06-13 22:16:13.746 ERROR [v2geq7340pbfxcc9,k1gsschfgarpc7no] main.registerWebEndpoints.func2 : Oh on!
continue on a new line :D
```

### Pattern for java

Logs are parsed using following regex:

```
`^([0-9]{4}\-[0-9]{2}\-[0-9]{2} [0-9:\.]+) +(\w+) +\[[\w \-]+,([\w ]*),([\w ]*),[\w ]*\] [\w\.]+ \-\-\- \[[\w\- ]+\] ([\w\-\.]+) +: *((?s).*)`
```

The log looks like this:

```log
2023-06-17 17:34:48.762  INFO [auth-service,,,] 78446 --- [           main] .c.m.r.c.YamlBasedRedissonClientProvider : Loading RedissonClient from yaml config file, reading environment property: redisson-config
```
