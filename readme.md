使用graphql查询k8s资源

```bash
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
    "query": "{Pod {spec {nodeName,priority,preemptionPolicy,containers,tolerations}, apiVersion, kind, metadata {name}}}"
}' \
 'http://localhost:8080/graphql'
```
