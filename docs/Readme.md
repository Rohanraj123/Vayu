## config.yaml

```yaml
server:
  port: 8087
  read_timeout: 5s
  write_timeout: 10s 
  idle_timeout: 120s

logging:
  format: text
  level: info

auth:
  mode: api_key
  enabled: true

routes:
  - path: /users
    upstream: http://localhost:9001
    auth_required: true
    methods: ["GET"]
    service: "users"
  - path: /orders
    upstream: http://localhost:9002
    auth_required: false
    methods: ["GET"]
    service: "orders"
```

## role.yaml

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: vayu-secret-role
  namespace: vayu-system
rules:
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "list", "create", "update", "patch"]

```

## role-binding.yaml
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: vayu-secret-binding
  namespace: vayu-system
subjects:
  - kind: ServiceAccount
    name: vayu-sa
    namespace: vayu-system
roleRef:
  kind: Role
  name: vayu-secret-role
  apiGroup: rbac.authorization.k8s.io

```

## service-account.yaml
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: vayu-sa
  namespace: vayu-system

```

## pod.yaml
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: vayu
  namespace: vayu-system
  labels:
    app: vayu
spec:
  serviceAccountName: vayu-sa
  volumes:
    - name: logs
      configMap:
        name: config
  containers:
    - name: vayu
      image: rohanraj123/vayu:latest
      imagePullPolicy: Always
      ports:
        - containerPort: 8080
      volumeMounts:
      - name: logs
        mountPath: /app/config.yaml
        subPath: config.yaml
```

We need to configure api-gateway this way, providing required permission to create secrets in a specific namespace