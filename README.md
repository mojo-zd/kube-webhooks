# kube-webhooks

#### Introduce
create webhooks server & generate certificate, this example aim to add default quota values for deployment if the namespace added `resourcequotas`

#### build webhook

- build lb-webhook

make certs depends on https://github.com/tests-always-included/mo, you must install it before make certs 
```
> cd kube-webhooks/cmd/lb-webhook
> make relase=v1.7.0 # v1.7.0 is the image tag you want
```

- generate certs

the scripts will generate certs and `admissionregistration.yaml secret.yaml` files
```
> cd scripts
> make URL=lb-webhook.default.svc  #URL is the webhook dns of kubernetes
```

- deploy

the `deployments` directory  includes deploy files,
`admissionregistration.yaml` is a `MutatingWebhookConfiguration` Object, it's `caBundle` should use the ca.crt file with bas64.
`secret.yaml` include certs information, you should use `ca.crt server.crt server.key`'s base64 value fill `ca.crt server.crt server.key`

```
kubectl apply -f admissionregistration.yaml
kubectl apply -f secret.yaml
kubectl apply -f deployment.yaml
```

- test

```
kubectl apply -f demo-namespace.yaml
kubectl apply -f quota.yaml
kubectl apply -f test-fail.yaml #this case will fail, because the deployment not include `io.wise2c.service.type` label
kubectl apply -f test-success.yaml   
```
> https://mustache.github.io/