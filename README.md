# Client-TiKV

A client for `TiKV`.

# Get start
It is recommended to use `MacOS` for stable operation
```shell
./tikv-console --pd host1,host2,host3
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

## Info

```sql
>>>> select * from regions order by region_id;
+-----------+----------------------------------------------------------------------+----------------------------------------------------------------------+-----------+-----------------+----------+---------------+------------+------------------+------------------+
| REGION_ID | START_KEY                                                            | END_KEY                                                              | LEADER_ID | LEADER_STORE_ID | PEERS    | WRITTEN_BYTES | READ_BYTES | APPROXIMATE_SIZE | APPROXIMATE_KEYS |
+-----------+----------------------------------------------------------------------+----------------------------------------------------------------------+-----------+-----------------+----------+---------------+------------+------------------+------------------+
|         4 | 696E646578496E666F5F3A5F706530315F3A5F747970655F3A5F363137313330     |                                                                      |         6 |               3 | 5,6,7    |          1632 |         48 |               74 |                0 |
|         8 | 696E646578496E666F5F3A5F706530315F3A5F747970655F3A5F32313230383638   | 696E646578496E666F5F3A5F706530315F3A5F747970655F3A5F32343731333730   |        11 |               2 | 9,10,11  |             0 |          0 |               60 |                0 |
|        12 | 696E646578496E666F5F3A5F706530315F3A5F747970655F3A5F31323037393434   | 696E646578496E666F5F3A5F706530315F3A5F747970655F3A5F31353930383331   |        13 |               1 | 13,14,15 |             0 |          0 |               60 |                0 |
|        16 | 696E646578496E666F5F3A5F706530315F3A5F313339547970655F3A5F3539393035 | 696E646578496E666F5F3A5F706530315F3A5F747970655F3A5F31323037393434   |        17 |               1 | 17,18,19 |             0 |          0 |               54 |                0 |
|        20 | 696E646578496E666F5F3A5F706530315F3A5F747970655F3A5F333033323136     | 696E646578496E666F5F3A5F706530315F3A5F747970655F3A5F363137313330     |        21 |               1 | 21,22,23 |             0 |          0 |               53 |                0 |
|        24 |                                                                      | 696E646578496E666F5F3A5F706530315F3A5F313339547970655F3A5F3539393035 |        25 |               1 | 25,26,27 |          1591 |         10 |              104 |           851503 |
|        28 | 696E646578496E666F5F3A5F706530315F3A5F747970655F3A5F32343731333730   | 696E646578496E666F5F3A5F706530315F3A5F747970655F3A5F333033323136     |        29 |               1 | 29,30,31 |             0 |          0 |               80 |                0 |
|        32 | 696E646578496E666F5F3A5F706530315F3A5F747970655F3A5F31353930383331   | 696E646578496E666F5F3A5F706530315F3A5F747970655F3A5F32313230383638   |        33 |               1 | 33,34,35 |             0 |          0 |               90 |                0 |
+-----------+----------------------------------------------------------------------+----------------------------------------------------------------------+-----------+-----------------+----------+---------------+------------+------------------+------------------+
8 rows in set (0.020785 sec)
```

Continuous iteration...