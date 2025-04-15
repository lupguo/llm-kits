## 测试

```shell

# 上传文件
curl -u admin:123456 -T local.config.txt http://localhost:8080/dev.config

# 下载文件
curl -u admin:123456 -o local_file.txt http://localhost:8080/dev.config
```

## 锁

### 加锁

> 加锁（LOCK）操作必须遵循 RFC 4918，依赖 XML 定义锁参数

```shell
 curl -u admin:123456 -X LOCK \
  -H "Content-Type: application/xml" \
  -H "Depth: infinity" \
  --data-binary @./lock.xml \
  http://localhost:8080/dev.config
<?xml version="1.0" encoding="utf-8"?>
<D:prop xmlns:D="DAV:"><D:lockdiscovery><D:activelock>
        <D:locktype><D:write/></D:locktype>
        <D:lockscope><D:exclusive/></D:lockscope>
        <D:depth>infinity</D:depth>
        <D:owner></D:owner>
        <D:timeout>Second-0</D:timeout>
        <D:locktoken><D:href>1743331814</D:href></D:locktoken>
        <D:lockroot><D:href>/dev.config</D:href></D:lockroot>
</D:activelock></D:lockdiscovery></D:prop> 
```

### 解锁

> 解锁（UNLOCK）仅需通过请求头传递锁令牌。

```shell
curl -u admin:123456 -X UNLOCK \
  -H "Lock-Token: <urn:uuid:abcd1234-5678-90ef-ghijklmnopqr>" \
  http://localhost:8080/dev.config
```

### 查看锁状态
```shell
$ curl -u admin:123456 -X PROPFIND http://localhost:8080/dev.config
<?xml version="1.0" encoding="UTF-8"?><D:multistatus xmlns:D="DAV:"><D:response><D:href>/dev.config</D:href><D:propstat><D:prop><D:getcontenttype>text/plain; charset=utf-8</D:getcontenttype><D:getetag>"18318f7657acfce52c"</D:getetag><D:resourcetype></D:resourcetype><D:displayname>dev.config</D:displayname><D:getcontentlength>44</D:getcontentlength><D:getlastmodified>Sun, 30 Mar 2025 10:57:49 GMT</D:getlastmodified><D:supportedlock><D:lockentry xmlns:D="DAV:"><D:lockscope><D:exclusive/></D:lockscope><D:locktype><D:write/></D:locktype></D:lockentry></D:supportedlock></D:prop><D:status>HTTP/1.1 200 OK</D:status></D:propstat></D:response></D:multistatus>
```