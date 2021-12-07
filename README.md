# client-tikv

A TiKV client, which allows you to operate TiKV like TiDB.

# Get start
```shell
./tikv-client --pd 172.16.5.133:2379
```

# Build
```shell
 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a
```

# Usage

## Get
You can query the value directly according to the `key`.
```sql
tikv> select kv from tikv where k = 'client';
+--------+-------+
| KEY    | VALUE |
+--------+-------+
| client | tikv  |
+--------+-------+
Query OK, 1 rows affected (0.008222 sec)

tikv> select kv from tikv where k in ('client','pingcap');
+---------+-------+
| KEY     | VALUE |
+---------+-------+
| client  | tikv  |
| pingcap | tikv  |
+---------+-------+
Query OK, 2 rows affected (0.012975 sec)
```
If the `value` is in the standard `JSON` format, you can query according to the `label`, like this
```sql
tikv> select kv from tikv where k = 'Green';
+-------+------------------------------------------+
| KEY   | VALUE                                    |
+-------+------------------------------------------+
| Green | {"id":"000810","password":"pingcap@123"} |
+-------+------------------------------------------+
Query OK, 1 rows affected (0.014478 sec)

tikv> select id, password from tikv where k = 'Green';
+--------+-------------+
| ID     | PASSWORD    |
+--------+-------------+
| 000810 | pingcap@123 |
+--------+-------------+
1 rows in set (0.009683 sec)

tikv> select id, password from tikv where k in ('Green','Jim');
+--------+-------------+
| ID     | PASSWORD    |
+--------+-------------+
| 000810 | pingcap@123 |
| 000820 | pingcap@456 |
+--------+-------------+
2 rows in set (0.013796 sec)
```
