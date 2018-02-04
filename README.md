# fastmap-wrapper

## how to

```bash
go run main.go fastmap.go osm.go
```

## some benchmarks
```bash
    # fastmap
    time curl -H "Accept-Encoding: gzip" "http://localhost:3001/api/0.6/map?bbox=27.61649608612061,53.85379229563698,27.671985626220707,53.886459293813054" > /dev/null
    % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                    Dload  Upload   Total   Spent    Left  Speed
    100 10.3M    0 10.3M    0     0  3492k      0 --:--:--  0:00:03 --:--:-- 3493k

    real	0m3.042s
    user	0m0.008s
    sys	0m0.024s
    
    # cgimap
    time curl -H "Accept-Encoding: gzip" "http://localhost:31337/api/0.6/map?bbox=27.61649608612061,53.85379229563698,27.671985626220707,53.886459293813054" > /dev/null
    % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                    Dload  Upload   Total   Spent    Left  Speed
    100 1224k    0 1224k    0     0   602k      0 --:--:--  0:00:02 --:--:--  602k

    real	0m2.039s
    user	0m0.008s
    sys	0m0.000s

    # port
    time curl -H "Accept-Encoding: gzip" "http://localhost:3000/api/0.6/map?bbox=27.61649608612061,53.85379229563698,27.671985626220707,53.886459293813054" > /dev/null
    % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                    Dload  Upload   Total   Spent    Left  Speed
    100 10.2M    0 10.2M    0     0   485k      0 --:--:--  0:00:21 --:--:-- 2430k

    real	0m21.521s
    user	0m0.008s
    sys	0m0.004s
```