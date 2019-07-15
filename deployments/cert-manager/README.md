#### Installing with regular manifests
```
# Create a namespace to run cert-manager in
kubectl create namespace cert-manager

# Disable resource validation on the cert-manager namespace
kubectl label namespace cert-manager certmanager.k8s.io/disable-validation=true

# Install the CustomResourceDefinitions and cert-manager itself
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v0.8.1/cert-manager.yaml
```
!!! note: 安装过程中请注意文档中的Note,我当前执行的环境为k8s v1.13.6,本文给出的步骤只支持kubernetes v1.12以上。如果遇到以下问题请再执行一遍apply
```
unable to recognize "https://github.com/jetstack/cert-manager/releases/download/v0.8.1/cert-manager.yaml": no matches for kind "Issuer" in version "certmanager.k8s.io/v1alpha1"
unable to recognize "https://github.com/jetstack/cert-manager/releases/download/v0.8.1/cert-manager.yaml": no matches for kind "Certificate" in version "certmanager.k8s.io/v1alpha1"
unable to recognize "https://github.com/jetstack/cert-manager/releases/download/v0.8.1/cert-manager.yaml": no matches for kind "Issuer" in version "certmanager.k8s.io/v1alpha1"
unable to recognize "https://github.com/jetstack/cert-manager/releases/download/v0.8.1/cert-manager.yaml": no matches for kind "Certificate" in version "certmanager.k8s.io/v1alpha1"
```

#### Check pod status
```
kubectl get po -n cert-manager

NAME                                       READY   STATUS    RESTARTS   AGE
cert-manager-68cfd787b6-rbbkt              1/1     Running   1          6d21h
cert-manager-cainjector-5975fd64c5-csd9w   1/1     Running   29         6d21h
cert-manager-webhook-5c7f95fd44-qdplg      1/1     Running   0          6d21h
```

#### Start Our Journey
- Generate a signing key pair
```
openssl genrsa -out ca.key 2048
```
- Create a self signed Certificate, valid for 10yrs with the 'signing' option set
```
openssl req -x509 -new -nodes -key ca.key -subj "/CN=lb-webhook.default.svc" -days 3650 -reqexts v3_req -extensions v3_ca -out ca.crt
```
- Save the signing key pair as a Secret
```
kubectl create secret tls ca-key-pair \
   --cert=ca.crt \
   --key=ca.key \
   --namespace=default
```
- Creating an Issuer referencing the Secret issuer.yaml
```
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  name: ca-issuer
  namespace: default
spec:
  ca:
    secretName: ca-key-pair
```
!!! Note: Certificate需要根据Issuer生成证书,创建的Certificate需要和Issuer在同一个namespace下,如果想所有namespace都可以使用Issuer可以设置Issuer为ClusterIssuer
- Create Certificate certificate.yaml
```
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  name: example-com
  namespace: default
spec:
  secretName: lb-webhook-tls
  issuerRef:
    name: ca-issuer
    # We can reference ClusterIssuers by changing the kind here.
    # The default value is Issuer (i.e. a locally namespaced Issuer)
    kind: Issuer
  commonName: lb-webhook.default.svc
  organization:
  - Wise2c CA
  dnsNames:
  - lb-webhook.default.svc
```

#### Deploy Cert Manager files
apply the files of the below, this will generate a secrets witch name is `lb-webhook-tls` and it include `ca.crt server.crt server.key`
```
> kubectl apply -f issuer.yaml
> kubectl apply -f certificate.yaml
```

### Deploy webhook server
- set `admissionregistration.yaml`'s `caBundle`, the value is `lb-webhook-tls`'s `ca.crt`
- apply 
```
kubectl apply -f admissionregistration.yaml
kubectl apply -f deployment.yaml
```

- test
```
kubectl apply -f demo-namespace.yaml
kubectl apply -f quota.yaml
kubectl apply -f test-fail.yaml #this case will fail, because the deployment not include `io.wise2c.service.type` label
kubectl apply -f test-success.yaml   
