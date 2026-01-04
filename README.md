# ClassConnect API

A secure, scalable REST API built with Go for managing educational institutions' executive, teacher, and student data. The system implements enterprise-grade security following the CIA triad (Confidentiality, Integrity, Availability) principles.

## Architecture Overview

### Technology Stack
- **Runtime**: Go 1.24
- **Database**: MariaDB 11.8
- **Containerization**: Docker & Docker Compose
- **Orchestration**: Kubernetes (k8s)
- **Authentication**: JWT (JSON Web Tokens)
- **Encryption**: TLS/HTTPS with self-signed certificates

### System Components

```
┌─────────────────────────────────────────────────────┐
│                 Load Balancer                       │
│            (Kubernetes Service)                     │
└──────────────┬──────────────────┬───────────────────┘
               │                  │
       ┌───────▼──────┐   ┌──────▼───────┐
       │   Pod 1      │   │   Pod 2      │   ... (N Pods)
       │  API Server  │   │  API Server  │
       └───────┬──────┘   └──────┬───────┘
               │                  │
               └──────────┬───────┘
                          │
                   ┌──────▼──────┐
                   │   MariaDB   │
                   │   Database  │
                   └─────────────┘
```

### Directory Structure
```
ClassConnect/
├── cmd/api/              # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/     # HTTP request handlers
│   │   ├── middlewares/  # Security & processing middleware
│   │   └── routers/      # Route definitions
│   ├── models/           # Data models
│   └── repository/       # Database layer
├── pkg/utils/            # Utility functions (JWT, password hashing)
├── k8s/                  # Kubernetes manifests
├── Dockerfile            # Multi-stage container build
└── docker-compose.yaml   # Local development setup
```

## CIA Triad Implementation

### 1. Confidentiality 🔒

**Objective**: Ensure sensitive data is accessible only to authorized users.

#### HTTPS/TLS Encryption
- All API communications encrypted using TLS 1.2+
- Self-signed X.509 certificates (RSA 4096-bit) generated during Docker build
- Prevents man-in-the-middle attacks and eavesdropping
- Certificate files: `/app/server.crt` and `/app/server.key`

#### JWT Authentication
- Token-based authentication for stateless session management
- HMAC-SHA256 signed tokens with configurable expiration
- Tokens stored in secure HTTP-only cookies
- Middleware validates tokens on protected routes

#### Password Security
- Argon2id hashing algorithm (memory-hard, resistant to GPU attacks)
- Per-password random salts (16 bytes)
- Hash parameters: 1 iteration, 64MB memory, 4 parallelism, 32-byte output
- Stored format: `base64(salt).base64(hash)`

#### Security Headers
- `X-Frame-Options: DENY` - Prevents clickjacking
- `X-Content-Type-Options: nosniff` - Blocks MIME sniffing
- `Content-Security-Policy: default-src 'self'` - Mitigates XSS
- `X-XSS-Protection: 1; mode=block` - Browser XSS filter
- `Referrer-Policy: no-referrer` - Prevents referrer leakage

### 2. Integrity ✓

**Objective**: Ensure data remains accurate and unaltered during transmission and storage.

#### JWT Signature Verification
- HMAC-SHA256 signatures prevent token tampering
- Invalid signatures immediately rejected by middleware
- Ensures user identity and permissions cannot be forged

#### TLS Certificate Validation
- Certificates ensure server authenticity
- Cryptographic signatures verify data hasn't been modified in transit
- Handshake process validates endpoint identity

#### Database Constraints
- Foreign key relationships maintain referential integrity
- Transaction support ensures atomic operations
- Schema validation prevents malformed data insertion

#### Request Validation
- Input sanitization through middleware pipeline
- Type-safe JSON unmarshaling with Go structs
- SQL prepared statements prevent injection attacks

#### Compression Middleware
- GZIP compression with integrity checks
- Maintains data consistency during transfer optimization

### 3. Availability ⚡

**Objective**: Ensure the system remains accessible and responsive.

#### Kubernetes Orchestration
- **Replica Sets**: 3 identical pods running simultaneously
- **Self-Healing**: Automatic pod restart on failures
- **Health Checks**: Liveness and readiness probes monitor pod status
- **Rolling Updates**: Zero-downtime deployments

#### Load Balancing
- Kubernetes LoadBalancer service distributes traffic evenly
- External IP exposes API to outside traffic
- Port mapping: External (80/443) → Internal (3000)
- Automatic failover if pods become unhealthy

#### Database High Availability
- MariaDB deployment with persistent volumes
- Health checks ensure database readiness before API starts
- Connection pooling for efficient resource usage
- Configurable connection retries in sqlconnect layer

#### Rate Limiting
- Configurable request throttling (default: 50 requests/minute)
- Prevents DoS attacks and resource exhaustion
- Per-client tracking with token bucket algorithm

#### Resource Optimization
- **Multi-stage Docker builds**: Minimal 15MB final images
- **Distroless base**: Reduced attack surface, faster startup
- **Static binaries**: No runtime dependencies (CGO_ENABLED=0)
- **GZIP compression**: Reduced bandwidth usage

#### Middleware Pipeline
```
Request → CORS → Rate Limiter → Security Headers → 
JWT Validation → Compression → Response Time → Handler
```

## API Endpoints

### Execs (Executives)
- `GET /execs/` - List all executives
- `GET /execs/{id}` - Get executive by ID
- `POST /execs/` - Create new executive(s)
- `POST /execs/login/` - Authenticate executive
- `PUT /execs/{id}` - Update executive
- `DELETE /execs/{id}` - Remove executive
- `POST /execs/forgotPassword/` - Request password reset
- `POST /execs/resetPassword/` - Reset password with code

### Students & Teachers
Similar CRUD operations available for students and teachers.

## Deployment

### Local Development (Docker Compose)
```bash
docker compose up
```
Access API: `https://localhost:3000`

### Production (Kubernetes)
```bash
# Build image
docker build -t classconnect-api:latest .

# Apply Kubernetes manifests
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/mariadb-deployment.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# Get LoadBalancer IP
kubectl get service classconnect-api-lb
```

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `API_PORT` | Server port | `3000` |
| `JWT_SECRET` | Token signing key | `your-secret-key` |
| `JWT_EXPIRES_IN` | Token lifetime | `6000s` |
| `DB_HOST` | Database hostname | `mariadb` |
| `DB_PORT` | Database port | `3307` |
| `DB_USER` | Database user | `admin` |
| `DB_PASSWORD` | Database password | `secure-password` |
| `DB_NAME` | Database name | `ClassConnect` |

## Security Best Practices

1. **Rotate JWT secrets** regularly in production
2. **Use trusted CA certificates** (Let's Encrypt) instead of self-signed
3. **Enable database encryption** at rest
4. **Implement audit logging** for compliance
5. **Set up monitoring** with Prometheus/Grafana
6. **Use secrets management** (Kubernetes Secrets, HashiCorp Vault)
7. **Regular security audits** and dependency updates

## Performance Metrics

- **Startup Time**: < 2 seconds
- **Image Size**: 15MB (distroless)
- **Memory Usage**: ~50MB per pod
- **Response Time**: < 50ms (measured by middleware)
- **Concurrent Connections**: 1000+ (Go's net/http)

## License

MIT License - See LICENSE file for details.
