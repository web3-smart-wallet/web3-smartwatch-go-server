# Web3 Smartwatch Go Server - Kubernetes Deployment

This directory contains Kubernetes configuration files for deploying the Web3 Smartwatch Go Server.

## Prerequisites

- Kubernetes cluster
- kubectl configured to communicate with your cluster
- cert-manager installed on your cluster
- Docker registry where your container images are stored

## Configuration Files

- `namespace.yaml`: Defines the Kubernetes namespace for the application
- `deployment.yaml`: Defines the application deployment
- `service.yaml`: Exposes the application through a service
- `ingress.yaml`: Sets up ingress with HTTPS support
- `secrets.yaml`: Contains sensitive environment variables
- `configmap.yaml`: Contains non-sensitive configuration
- `cluster-issuer.yaml`: Sets up Let's Encrypt certificate issuer
- `kustomization.yaml`: Manages all resources together

## Deployment Steps

1. Update placeholder values:

   - In `secrets.yaml`, replace placeholder API keys and URLs with actual values
   - In `cluster-issuer.yaml`, replace `your-email@example.com` with your actual email
   - In `deployment.yaml`, replace `${DOCKER_REGISTRY}` with your Docker registry URL
   - In `ingress.yaml`, replace `api.web3smartwatch.com` with your actual domain

2. Apply the namespace first:

   ```bash
   kubectl apply -f namespace.yaml
   ```

3. Apply the ClusterIssuer:

   ```bash
   kubectl apply -f cluster-issuer.yaml
   ```

4. Apply secrets (consider using sealed-secrets or another secure method):

   ```bash
   kubectl apply -f secrets.yaml
   ```

5. Apply all other resources using kustomize:
   ```bash
   kubectl apply -k ./
   ```

## Verification

Check the status of the deployment:

```bash
kubectl get all -n web3-smartwatch
```

Check the ingress status:

```bash
kubectl get ingress -n web3-smartwatch
```

Check certificate status:

```bash
kubectl get certificates -n web3-smartwatch
```

## Scaling

To scale the deployment:

```bash
kubectl scale deployment web3-smartwatch-server -n web3-smartwatch --replicas=3
```
