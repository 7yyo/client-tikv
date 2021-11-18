# client-tikv

```shell
 ./tikv-client --pd 127.0.0.1:2379,127.0.0.2:2379,127.0.0.3:2379
```

```sql
tikv> select kv from tikv where k in ('123','1')
┌─────┬───────┐
│ Key │ Value │
├─────┼───────┤
│ 123 │ 456   │
│ 1   │ 2     │
└─────┴───────┘
2 rows in set (0.047940 sec)
```
