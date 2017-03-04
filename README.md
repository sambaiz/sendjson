# sendjson

```
$ go install github.com/sambaiz/sendjson
$ sendjson --help
$ sendjson -i 2s -u http://localhost:9200/hoge '{"a": {"type": "integer", "min": 10, "max": 100}}'
{"a":27}
{"a":99}
{"a":98}
{"a":80}
...
```

## Format

Input

- not yet supported nested json

```
$ sendjson --check '
{
    "aaa": {"type": "string", "value": "foo"},
    "bbb": {"type": "boolean", "or": [false, true]},
    "ccc": {"type": "integer", "min": 100, "max": 1000},
    "ddd": {"type": "double", "min": 0, "max": 1},
    "timestamp": {"type": "time", "time_format": "2006-01-02T15:04:05Z07:00"},
    "unixtime": {"type": "time", "time_format": "unix_epoch"}
}'
```

Output

```
not sended: {"aaa":"foo","bbb":true,"ccc":436,"ddd":0.9866222922441757,"timestamp":"2017-03-04T21:22:39+09:00","unixtime":1488630159}
not sended: {"aaa":"foo","bbb":false,"ccc":896,"ddd":0.8197965426480669,"timestamp":"2017-03-04T21:22:40+09:00","unixtime":1488630160}
not sended: {"aaa":"foo","bbb":false,"ccc":996,"ddd":0.85358913945103,"timestamp":"2017-03-04T21:22:41+09:00","unixtime":1488630161}
not sended: {"aaa":"foo","bbb":true,"ccc":413,"ddd":0.8132848145773628,"timestamp":"2017-03-04T21:22:42+09:00","unixtime":1488630162}
...
```