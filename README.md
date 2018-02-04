# fastmap-wrapper

## how to

```bash
go run main.go fastmap.go osm.go
```

## some benchmarks
```bash
    bbox=27.616,53.853,27.630,53.870:
        fastmap response time: 948.199091ms
        port response time: 4.558749065s
        cgimap response time: 230.034071ms
    
    bbox=27.616,53.853,27.617,53.854:
        fastmap response time: 362.441185ms
   	    port response time: 49.829756ms
   	    cgimap response time:  38.804292ms
```