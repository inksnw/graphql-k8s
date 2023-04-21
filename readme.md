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
带参数
```bash
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "query": "query($name: String!,$namespace: String!) {Pod(name: $name,namespace: $namespace) {spec {nodeName,priority,preemptionPolicy,containers,tolerations}, apiVersion, kind, metadata {name}}}",
  "variables": {
    "name": "nginx",
    "namespace":"default"
  }
}' \
 'http://localhost:8080/graphql'
```

label筛选
```bash
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "query": "query($namespace: String!,$label: String!) {Pod(label: $label,namespace: $namespace) {spec {nodeName,priority,preemptionPolicy,containers,tolerations}, apiVersion, kind, metadata {name}}}",
  "variables": {
    "namespace":"default",
    "label":"run=nginx"
  }
}' \
 'http://localhost:8080/graphql'

```