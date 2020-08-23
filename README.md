# istio-authz-example
Example using Istio and Envoy's authz (external authorization) filter.

Inspired by https://github.com/salrashid123/istio_external_authorization_server.

The use case is as follows:
You've got your kubernetes (k8s) cluster.
You want to route traffic into the cluster.
But before traffic gets routed to upstream (deeply internal) services,
it should get "checked" by a service to see if the bearer token in the 
Authorization header checks out.

## Getting started
### Aliases
This README uses aliases from the [oh-my-zsh kubectl plugin](https://github.com/ohmyzsh/ohmyzsh/tree/master/plugins/kubectl).
### Build Docker Image
Let's build our app.

```
DOCKER_BUILDKIT=1 docker build \
  --force-rm \
  --no-cache \
  -t kevinmichaelchen/istio-authz-example:0.0.1 \
  .
```

### Create a custom namespace
```
k create namespace istio-authz-example-ns

istioctl install --set profile=demo
```

### Enabling VirtualService delegation
We need to tell Istio to allow VirtualServices to use delegation with the [`PILOT_ENABLE_VIRTUAL_SERVICE_DELEGATE`](https://istio.io/latest/docs/reference/commands/pilot-agent/#envvars) env var.
You can [edit the istiod deployment directly](https://discuss.istio.io/t/try-to-create-a-delegate-virtual-service-but-got-error-configuration-invalid-virtual-service-must-have-at-least-one-host/7133/4).
```
KUBE_EDITOR=vim k edit deployment istiod -n istio-system
```
Paste in
```yaml
        - name: PILOT_ENABLE_VIRTUAL_SERVICE_DELEGATE
          value: "true"
```
View all env vars with
```
kgd istiod -n istio-system -o json | jq '.spec.template.spec.containers[0].env'
```

### Enabling Istio Sidecar Injection
Sidecars are necessary for the mesh to work.
The easiest way to accomplish sidecar injection is to automatically enable it.
```
k label namespace istio-authz-example-ns istio-injection=enabled
```

We should now be able to see all namespaces with or without that label:
```
kgns -L istio-injection
```

By default, we should see the sidecar injector applying to all namespaces with the `istio-injection=enabled` label.
```
k get mutatingwebhookconfiguration istio-sidecar-injector -o yaml | grep "namespaceSelector:" -A 3
```

### Installing our backend
For microk8s clusters:
```
helm install api -n istio-authz-example-ns --set image.repository=localhost:32000/kevinmichaelchen/istio-authz-example ./helm/api
```

Otherwise:
```
helm install api -n istio-authz-example-ns ./helm/api
```

### Installing Istio resources
```
kaf ./helm/istio/gateway.yaml
kaf ./helm/istio/virtualservice.yaml
kaf ./helm/istio/envoyfilter.yaml

kdelf ./helm/istio/gateway.yaml
kdelf ./helm/istio/virtualservice.yaml
kdelf ./helm/istio/envoyfilter.yaml

# List our Istio resources
k get gateway,virtualservice -A

export INGRESS_HOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].port}')
curl -s -I -HHost:api.example.com "http://$INGRESS_HOST:$INGRESS_PORT/insecure"
```

### Kiali
Visualize traffic with Kiali.
```
istioctl dashboard kiali
```
The default credentials are `admin:admin`.

### Checking logs
Check logs with
```
kl -n istio-authz-example-ns -l app.kubernetes.io/name=api -c api
kl -n istio-authz-example-ns -l app.kubernetes.io/name=api --all-containers=true
```

### Helm cheatsheet
#### List Helm releases
```
helm list -n istio-authz-example-ns
```

#### Uninstall Helm charts
```
helm uninstall <release_name> -n istio-authz-example-ns
```