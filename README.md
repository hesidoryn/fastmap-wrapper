# fastmap-wrapper

## how to

```bash
go run main.go fastmap.go osm.go
```

## some benchmarks
```bash
    bbox=27.616,53.853,27.630,53.870:
        fastmap response time:  647.88781ms
        port response time:  4.191980366s
        cgimap response time:  136.19359ms
    
    bbox=27.616,53.853,27.617,53.854:
        fastmap response time:  361.441185ms
   	    port response time:  51.765098ms
   	    cgimap response time:  31.264603ms
```