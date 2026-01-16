# Deployment Strategy

**Choice:** AWS Lightsail (Double Instance)
**Cost:** ~$8.50/month

## Why Lightsail?

| Requirement | Lightsail |
|-------------|-----------|
| WebSocket support | Full (no timeout limits) |
| Low starting cost | $3.50-5/month per instance |
| Docker support | Yes |
| Easy setup | Yes |
| AWS skills transfer | Yes |

## Options Considered

| Option | Cost | WebSockets | Notes |
|--------|------|------------|-------|
| **Lightsail VM** | $5/mo | Full | Chosen - best balance |
| ECS Fargate | $30+/mo | Full | Overkill for starting out |
| App Runner | $10/mo | 30-min limit | WebSocket timeout is a dealbreaker |
| Fly.io | $0-10/mo | Full | Not AWS |

## Why Double Instance?

| Benefit | Description |
|---------|-------------|
| Independent deploys | Update frontend without disrupting WebSocket connections |
| Resource isolation | SSR and WebSocket don't compete |
| Security | Frontend has no public IP |
| Scaling | Upgrade each instance independently |

## Scaling Path

```
Phase 1 (Now)          Phase 2 (Growth)         Phase 3 (Scale)
     │                      │                        │
     ▼                      ▼                        ▼
2x Lightsail  ────►  Add Load Balancer  ────►  ECS Fargate
($8.50/mo)            + more instances          (auto-scaling)
```

## When to Scale

| Trigger | Action |
|---------|--------|
| CPU > 80% consistently | Upgrade instance size |
| Need redundancy | Add Lightsail load balancer ($18/mo) |
| Need auto-scaling | Migrate to ECS Fargate |
