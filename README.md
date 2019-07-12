# kube-webhooks

#### Introduce
create webhooks & generate certificate for wisecloud

#### build wehbook

- build lb-webhook 
```
> cd kube-webhooks/cmd/lb-webhook
> make relase=v1.7.0 # v1.7.0 is the image tag you want
```

- generate certs
```
> cd scripts
> make URL=lb-webhook.wisecloud-agent.svc  #URL is the webhook dns of kubernetes
```

- deploy

the `deployments` directory  includes deploy files,
`admissionregistration.yaml` is a `MutatingWebhookConfiguration` Object, it's `caBundle` should use the ca.crt file with bas64
`secret.yaml` include certs information, you should use `ca.crt server.crt server.key`'s base64 value fill `ca.crt tls.crt tls.key`

!!! note: ca.crt generate with `scripts/Makefile`