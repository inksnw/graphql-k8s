```bash
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "query": "{pods {spec{nodeName,priority,preemptionPolicy,containers,tolerations},apiVersion,kind,metadata{name} } }"
}' \
 'http://localhost:8080/graphql'
```