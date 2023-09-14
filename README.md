# admisson-webhook
It is a k8s [admission webhook controller](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/) demo. 

this demo covers 'ValidatingWebhook' and 'MutatingWebhook'
the creating 'pod' must have 'Time' env key for 'Validating'; when pod created, you will find 'mutate-timestamp' annotation key and 'mutated-app' label key inject.

# Step
1. build image.(this demo's image is arm arch, so you'd better build image by yourself)
```shell
make build-image
```
2. load image to kind cluster (you could skip this step if not use 'kind' k8s cluster)
```shell
 kind load docker-image xxx:xxx  --name clustername
```
3. apply
```shell
make apply
```

# Test
- without 'Time' env key
```shell
$ kubectl run nginx --image nginx --env='FOO=BAR' -n webhook
Error from server (container nginx validate failed.env vars doesn't have 'Time' key): admission webhook "admission-demo.xiaoxlm.dev" denied the request: pod validating invalid
```

- with 'Time' env key
```shell
$ kubectl run nginx --image nginx --env='Time=BAR' -n webhook
pod/nginx created
```

- get mutated info
```shell
$ kubectl get pods nginx -o jsonpath={..labels}
{"mutated-app":"true","run":"nginx"}

$ kubectl get pods nginx -o jsonpath={..annotations}
{"mutate-timestamp":"1694660459"}
```


