# logbot

Bot for watching and parsing logs

## Configuration

| Property            | Description          | Default Value |
|---------------------|----------------------|---------------|
| logbot.node         | name of the node     | 'default'     |
| logbot.[]watch      | list of watch config |               |
| logbot.[]watch.app  | app name             |               |
| logbot.[]watch.file | path of the log file |               |

## API

| Method | Path              | Parameters                                                        | Description               |
|--------|-------------------|-------------------------------------------------------------------|---------------------------|
| POST   | `/log/error/list` | `{ "app" : "app_name" , "page" : { "limit" : 10, "page"  : 1 } }` | List ERROR logs collected |

## Log Pattern

Logs are parsed using following regex:

```
`^([0-9]{4}\-[0-9]{2}\-[0-9]{2} [0-9:\.]+) +(\w+) +\[([\w ]+),([\w ]+)\] ([\w\.]+) +: *((?s).*)`
```

The log looks like this:

```log
2023-06-13 22:16:13.746 ERROR [v2geq7340pbfxcc9,k1gsschfgarpc7no] main.registerWebEndpoints.func2 : Oh on!
continue on a new line :D
```
