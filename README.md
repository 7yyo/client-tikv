# client-tikv

A client for TiKV, Operate TiKV using SQL.

# Get start
It is recommended to use `MacOS` for stable operation
```shell
./tikv-client --pd 172.16.5.133:2379
```

# Build
```shell
 GO111MODULE=on go build .
```

# Usage

```sql
tikv > insert into tikv values ('pingcap','{"name":"tidb","url":"www.tidb.com"}')
Query OK, 1 rows affected (0.062081 sec)
    
tikv > insert into tikv values ('PingCAP','{"name":"tikv","url":"www.tikv.com"}')
Query OK, 1 rows affected (0.009297 sec)

tikv > select kv from tikv where k in ('pingcap','PingCAP');
+---------+--------------------------------------+
| KEY     | VALUE                                |
+---------+--------------------------------------+
| pingcap | {"name":"tidb","url":"www.tidb.com"} |
| PingCAP | {"name":"tikv","url":"www.tikv.com"} |
+---------+--------------------------------------+
Query OK, 2 rows affected (0.021388 sec)

tikv > select name,url from tikv where k in ('pingcap','PingCAP');
+------+--------------+
| NAME | URL          |
+------+--------------+
| tidb | www.tidb.com |
| tikv | www.tikv.com |
+------+--------------+
2 rows in set (0.026630 sec)

```

# Todo
[x] get
[x] batch get
[x] put
[x] batch put
[ ] delete
[ ] batch delete
