# TiKV-client

This is a client for tikv (rawkv), you can use some supported SQL statements to query data in tikv. 

## Example

### Put
```mysql
>>>> insert into t values ('japan','tokyo')
Query OK, 1 rows affected (0.02 sec)
>>>> insert into t values ('japan','tokyo'),('england','london')
Query OK, 2 rows affected (0.03 sec)
```

### Get
```mysql
>>>> select * from t where k = 'china'
+-------+---------+
| key   | value   |
+-------+---------+
| china | beijing |
+-------+---------+
Query OK, 1 rows affected (0.02 sec)
>>>> select * from t where k in ('china','japan')
+-------+---------+
| key   | value   |
+-------+---------+
| china | beijing |
| japan | tokyo   |
+-------+---------+
Query OK, 2 rows affected (0.02 sec)
>>>> select * from t limit 10
+---------+---------+
| key     | value   |
+---------+---------+
| china   | beijing |
| england | london  |
| japan   | tokyo   |
+---------+---------+
Query OK, 3 rows affected (0.11 sec)

```