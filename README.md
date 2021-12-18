# Client-TiKV

A client for TiKV, Operate TiKV using SQL.

# Get start
It is recommended to use `MacOS` for stable operation
```shell
./tikv-client --pd host1,host2,host3
```

# Build
```shell
 GO111MODULE=on go build .
```

# Usage

## Put
![image](https://github.com/7yyo/client-tikv/blob/master/gif/batch-put.gif)

## Get
![image](https://github.com/7yyo/client-tikv/blob/master/gif/batch-get.gif)

## Get fields
> Just support `JSON` format value

![image](https://github.com/7yyo/client-tikv/blob/master/gif/batch-get-field.gif)

# Todo
- [x] get
- [x] batch get
- [x] get for fields
- [x] put
- [x] batch put
- [ ] delete
- [ ] batch delete
