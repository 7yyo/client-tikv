# client-tikv

```shell
 ./tikv-client --pd 127.0.0.1:2379,127.0.0.2:2379,127.0.0.3:2379
```

# usage
You can query the value directly according to the key.
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

If the value is in the standard json format, you can query according to the label, like this
```sql
tikv> select kv from tikv where k = 'Green';
+-------+------------------------------------------+
| KEY   | VALUE                                    |
+-------+------------------------------------------+
| Green | {"id":"000810","password":"pingcap@123"} |
+-------+------------------------------------------+
Query OK, 1 rows affected (0.014478 sec)

tikv> select id from tikv where k = 'Green';
+--------+
| ID     |
+--------+
| 000810 |
+--------+
1 rows in set (0.008572 sec)

tikv> select id,password from tikv where k = 'Green';
+--------+-------------+
| ID     | PASSWORD    |
+--------+-------------+
| 000810 | pingcap@123 |
+--------+-------------+
1 rows in set (0.009683 sec)

tikv> select id,password from tikv where k in ('Green','Jim');
+--------+-------------+
| ID     | PASSWORD    |
+--------+-------------+
| 000810 | pingcap@123 |
| 000820 | pingcap@456 |
+--------+-------------+
2 rows in set (0.013796 sec)
```
