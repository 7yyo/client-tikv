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
tikv > insert into tikv values ('pingcap','{ "name":"google" , "url":"www.google.com" }');
Query OK, 1 rows affected (0.009565 sec)

tikv > select kv from tikv where k = 'pingcap';
+---------+----------------------------------------------+
| KEY     | VALUE                                        |
+---------+----------------------------------------------+
| pingcap | { "name":"google" , "url":"www.google.com" } |
+---------+----------------------------------------------+
Query OK, 1 rows affected (0.218763 sec)

tikv > select name, url from tikv where k = 'pingcap';
+--------+----------------+
| NAME   | URL            |
+--------+----------------+
| google | www.google.com |
+--------+----------------+
1 rows in set (0.007734 sec)
```
