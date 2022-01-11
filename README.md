# Client-TiKV

A client for `Ti`.

# Get start
```shell
./ti-client --pd xxx.xx.xxx.x:xxxx
```

# Build
```shell
go build 
```

# Usage

## KV
```sql
>>>> insert into tikv values ('jim','{"id":"1","name":"Jim","url":"www.baidu.com"}')
Query OK, 1 rows affected (0.108496 sec)
>>>> 
>>>> select * from tikv where k = 'jim';
+-----+-----------------------------------------------+
| KEY | VALUE                                         |
+-----+-----------------------------------------------+
| jim | {"id":"1","name":"Jim","url":"www.baidu.com"} |
+-----+-----------------------------------------------+
Query OK, 1 rows affected (0.229609 sec)
>>>> 
>>>> select id,name,url from tikv where k = 'jim';
+----+------+---------------+
| ID | NAME | URL           |
+----+------+---------------+
| 1  | Jim  | www.baidu.com |
+----+------+---------------+
1 rows in set (0.017130 sec)
```