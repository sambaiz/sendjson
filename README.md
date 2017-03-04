# SENDJSON

```
$ go install github.com/sambaiz/sendjson
$ sendjson -i 0.5 -u http://localhost:9200/hoge '{"a": {"type": "integer", "min": 10, "max": 100}}'
{"a":27}
{"a":99}
{"a":98}
{"a":80}
...
```

## Example

Input

- not yet supported nested json

```
{
    "aaa": {"type": "string", "value": "foo"},
    "bbb": {"type": "boolean", "or": [false, true] },
    "ccc":  {"type": "integer", "min": 100, "max": 1000},
    "ddd": {"type": "double": "min": 0, "max": 1 } 
}
```

Output(sample)

```
{
    "aaa": "foo",
    "bbb": true,
    "ccc": 133,
    "ddd": 0.6
}
```