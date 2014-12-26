MOCK TCP SERVER
===============

A tcp server mock, if receive some matched bytes, it will response the specific file data. 

#### Config it
`server.conf`
```json
{
    "host": "127.0.0.1",
    "port": 8080,
    "dump_request": true,
    "matchs": [{
        "type": "string",
        "match_data": "efg",
        "response_file": "test_1.txt"
    }, {
        "type": "byte",
        "match_data": "616263",
        "response_file": "test_2.txt"
    }]
}
```

- type
    string
    byte
- match_data
    string: just input match string
    byte:ascii
- dump_request
    if configured the dump_request = true, will dump the request data to file `./dump/{timestamp}/{ID}.dat`

#### Start it
```bash
$ go run main.go
```

#### Test it
Open browser, and test `http://127.0.0.1:8080/abc` and `http://127.0.0.1:8080/efg`, you will get different data.


