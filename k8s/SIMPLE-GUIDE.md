# Kubernetes Basics - Simplified Guide

## Core Concepts

### 1. Pod
- **Smallest unit** in Kubernetes
- Runs one or more containers
- Gets its own IP address
- Containers in same pod share network and storage
- Pods are **ephemeral** (temporary) - can be created/destroyed anytime

### 2. Deployment
- **Manages a set of identical pods**
- Ensures desired number of replicas are always running
- If a pod dies, Deployment creates a new one
- Handles rolling updates automatically

**Our deployment.yaml creates 3 pods:**
```yaml
replicas: 3  # Run 3 identical API pods
```

**The selector connects Deployment to its pods:**
```yaml
selector:
  matchLabels:
    app: classconnect-api  # "Manage pods with this label"
```

**Template defines what each pod looks like:**
```yaml
template:
  metadata:
    labels:
      app: classconnect-api  # Label for selector to find
  spec:
    containers:
    - name: api
      image: classconnect-api:latest
      ports:
      - containerPort: 3000  # Pod listens on port 3000
```

### 3. Service
- **Stable endpoint** to access pods
- Pods come and go, but Service IP stays the same
- Load balances traffic across all matching pods

**Two types we use:**

#### LoadBalancer (for API)
- Creates **external IP** accessible from outside cluster
- Routes traffic to pods: `external:80 → pod:3000`
```yaml
type: LoadBalancer
ports:
- port: 80          # External port
  targetPort: 3000  # Pod's port
```

#### ClusterIP (for Database)
- **Internal only** - no external access
- Other pods can reach it by name: `mariadb-service`
```yaml
type: ClusterIP  # Default, internal only
```

**Selector connects Service to pods:**
```yaml
selector:
  app: classconnect-api  # "Send traffic to pods with this label"
```

## How It Works Together

```
Internet
   ↓
LoadBalancer (port 80)
   ↓
Service: classconnect-api-lb
   ↓ (load balances across)
   ├─→ Pod 1 (port 3000)
   ├─→ Pod 2 (port 3000)
   └─→ Pod 3 (port 3000)
        ↓ (connects to)
   Service: mariadb-service (port 3306)
        ↓
   MariaDB Pod (port 3306)
```

## Deploy Step-by-Step

### 1. Build Docker image
```bash
docker build -t classconnect-api:latest .
```

### 2. Load into Kubernetes (minikube)
```bash
# Start minikube if not running
minikube start

# Load image into minikube's Docker
minikube image load classconnect-api:latest
```

### 3. Deploy database first
```bash
kubectl apply -f k8s/mariadb-deployment.yaml
```
**This creates:**
- 1 MariaDB pod
- 1 Service named `mariadb-service` (internal only)

**Check it:**
```bash
kubectl get pods    # See mariadb pod
kubectl get svc     # See mariadb-service
```

### 4. Deploy API
```bash
kubectl apply -f k8s/deployment.yaml
```
**This creates:**
- 3 API pods (replicas: 3)

**Check it:**
```bash
kubectl get pods    # Should see 4 pods total (3 API + 1 DB)
```

### 5. Expose with LoadBalancer
```bash
kubectl apply -f k8s/service.yaml
```
**This creates:**
- LoadBalancer service that routes to 3 API pods

**Check it:**
```bash
kubectl get svc classconnect-api-lb
```

### 6. Access the application
```bash
# Get the URL
minikube service classconnect-api-lb --url

# Example output: http://192.168.49.2:31234
# Use this URL to access your API
```

## Useful Commands

### View Resources
```bash
# List all pods
kubectl get pods

# List all services
kubectl get svc

# List deployments
kubectl get deployments

# Everything at once
kubectl get all
```

### Pod Details
```bash
# Describe a pod (shows events, status)
kubectl describe pod <pod-name>

# View logs
kubectl logs <pod-name>

# Follow logs in real-time
kubectl logs <pod-name> -f

# Shell into a pod
kubectl exec -it <pod-name> -- sh
```

### Testing Load Balancing
```bash
# Watch logs from all API pods simultaneously
kubectl logs -l app=classconnect-api -f

# In another terminal, make requests
curl http://<minikube-service-url>/api/teachers

# You'll see requests distributed across the 3 pods
```

### Scaling
```bash
# Scale to 5 replicas
kubectl scale deployment classconnect-api --replicas=5

# Scale back to 3
kubectl scale deployment classconnect-api --replicas=3

# Watch pods being created/terminated
kubectl get pods -w
```

### Cleanup
```bash
# Delete everything
kubectl delete -f k8s/

# Or delete individual resources
kubectl delete deployment classconnect-api
kubectl delete service classconnect-api-lb
```

## Key Differences: Deployment vs Service

| Deployment | Service |
|------------|---------|
| Manages **pods** | Routes **traffic** |
| Ensures replicas running | Provides stable endpoint |
| Handles updates/rollbacks | Load balances requests |
| Knows about pod lifecycle | Doesn't care which specific pod |

## Why We Need Both

**Without Deployment:** You'd manually create pods, and if one crashes, it's gone.

**Without Service:** Each pod has different IP, changes when restarted, no load balancing.

**Together:** Deployment keeps pods running, Service provides stable access point.

## Common Issues

### Pods not starting
```bash
# Check status
kubectl get pods

# If status is ImagePullBackOff
# → Image not found, reload it:
minikube image load classconnect-api:latest

# If CrashLoopBackOff
# → App is crashing, check logs:
kubectl logs <pod-name>
```

### Can't access LoadBalancer
```bash
# For minikube, run tunnel in separate terminal
minikube tunnel

# Or use the service command
minikube service classconnect-api-lb
```

### Database connection fails
```bash
# Check if database pod is running
kubectl get pods

# Check if service exists
kubectl get svc mariadb-service

# Test DNS from API pod
kubectl exec -it <api-pod-name> -- sh
# Inside pod:
ping mariadb-service
```
