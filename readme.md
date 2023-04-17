curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
    "query": "{Pod {spec {nodeName,priority,preemptionPolicy,containers,tolerations}, apiVersion, kind, metadata {name}}}"
}' \
 'http://localhost:8080/graphql'

为解决同一schema下不支持同名字段

- [x] 查询时添加前缀
- [x] 提取数据时特殊处理前缀
- [ ] 结果剔除前缀