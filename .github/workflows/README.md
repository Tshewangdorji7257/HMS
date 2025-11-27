# CI/CD Pipeline Documentation

## Overview

This repository includes comprehensive GitHub Actions workflows for automated testing, building, and deployment of the HMS (Hostel Management System) microservices.

## Workflows

### 1. Backend CI Pipeline (`backend-ci.yml`)
**Trigger**: Push or PR to `main` branch with changes in `Backend/` directory

**Jobs**:
- **Test & Build**: Tests and builds all 4 microservices (auth, booking, building, api-gateway) in parallel
- **Verify gRPC**: Validates Protocol Buffer definitions
- **Docker Build & Push**: Builds and pushes Docker images to Docker Hub (on push to main)
- **Integration Tests**: Runs docker-compose to test service integration
- **CI Summary**: Provides overall pipeline status

### 2. Frontend CI Pipeline (`frontend-ci.yml`)
**Trigger**: Push or PR to `main` branch with changes in `Frontend/` directory

**Jobs**:
- **Build & Test**: Installs dependencies, runs TypeScript checks, linting, and builds Next.js app
- **Docker Build**: Builds and pushes frontend Docker image to Docker Hub
- **Analyze Build**: Analyzes bundle size and build artifacts
- **Integration Test**: Tests frontend with backend services
- **CI Summary**: Provides overall pipeline status

### 3. Full Stack CI Pipeline (`fullstack-ci.yml`)
**Trigger**: Push or PR to `main` branch

**Jobs**:
- **Detect Changes**: Determines which parts of the stack changed
- **Backend CI**: Runs tests for all backend services in parallel
- **Frontend CI**: Runs frontend build and tests
- **Integration Tests**: Tests full stack integration with docker-compose
- **CI Status**: Overall pipeline status

### 4. Docker Build and Push (`docker-push.yml`)
**Trigger**: Push to `main`, version tags, or manual dispatch

**Jobs**:
- **Backend Docker Hub**: Builds and pushes all 4 backend services to Docker Hub
- **Frontend Docker Hub**: Builds and pushes frontend to Docker Hub
- **Backend AWS ECR**: Builds and pushes all 4 backend services to AWS ECR
- **Frontend AWS ECR**: Builds and pushes frontend to AWS ECR
- **Docker Summary**: Overall deployment status

## Required GitHub Secrets

### For Docker Hub
```
DOCKER_USERNAME   # Your Docker Hub username
DOCKER_PASSWORD   # Your Docker Hub password or access token
```

### For AWS ECR (Optional)
```
AWS_REGION              # AWS region (e.g., us-east-1)
AWS_ACCESS_KEY_ID       # AWS access key ID
AWS_SECRET_ACCESS_KEY   # AWS secret access key
```

## Setting Up Secrets

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add the required secrets listed above

### Docker Hub Setup
1. Create a Docker Hub account at https://hub.docker.com
2. Generate an access token:
   - Go to Account Settings → Security → New Access Token
   - Give it a name (e.g., "github-actions")
   - Copy the token
3. Add `DOCKER_USERNAME` and `DOCKER_PASSWORD` (use the token) to GitHub secrets

### AWS ECR Setup (Optional)
1. Create an AWS account if you don't have one
2. Create an IAM user with ECR permissions:
   - `AmazonEC2ContainerRegistryFullAccess` policy
3. Create access keys for the IAM user
4. Add `AWS_REGION`, `AWS_ACCESS_KEY_ID`, and `AWS_SECRET_ACCESS_KEY` to GitHub secrets

## Docker Images

### Backend Services
- `{username}/hms-auth-service:latest`
- `{username}/hms-booking-service:latest`
- `{username}/hms-building-service:latest`
- `{username}/hms-api-gateway:latest`

### Frontend
- `{username}/hms-frontend:latest`

## Local Testing

### Backend Services
```bash
cd Backend/auth-service
go test -v ./...
go build -o main .
```

### Frontend
```bash
cd Frontend
pnpm install
pnpm build
pnpm start
```

### Docker Compose (Full Stack)
```bash
cd Backend
docker-compose up --build
```

## Workflow Features

### Matrix Strategy
- Backend workflows use matrix strategy to run all 4 services in parallel
- Reduces CI time significantly
- Fail-fast disabled to see all service results

### Caching
- **Go**: Module cache for faster dependency installation
- **pnpm**: Store cache for faster npm package installation
- **Docker**: GitHub Actions cache for layer caching

### Multi-Platform Builds
- Docker images built for both `linux/amd64` and `linux/arm64`
- Ensures compatibility across different architectures

### Smart Triggers
- Path-based triggers only run relevant workflows
- Change detection prevents unnecessary builds
- Manual workflow dispatch available for `docker-push.yml`

## CI/CD Flow

```
┌─────────────────┐
│   Push to main  │
└────────┬────────┘
         │
    ┌────┴────┐
    │ Detect  │
    │ Changes │
    └────┬────┘
         │
    ┌────┴──────────────────┐
    │                       │
┌───▼────┐           ┌──────▼───┐
│Backend │           │ Frontend │
│  CI    │           │    CI    │
└───┬────┘           └──────┬───┘
    │                       │
    └───────┬───────────────┘
            │
      ┌─────▼─────┐
      │Integration│
      │   Tests   │
      └─────┬─────┘
            │
      ┌─────▼─────┐
      │  Docker   │
      │Build&Push │
      └───────────┘
```

## Troubleshooting

### Tests Failing
- Check the test logs in the GitHub Actions tab
- Run tests locally to reproduce the issue
- Ensure all dependencies are properly installed

### Docker Build Failing
- Verify Dockerfile syntax
- Check if base images are accessible
- Ensure all required files are in the build context

### Secrets Not Working
- Verify secret names match exactly (case-sensitive)
- Check if secrets are set at repository level, not environment level
- Regenerate tokens if they've expired

### Integration Tests Timing Out
- Increase sleep time in workflow for services to start
- Check docker-compose.yml health checks
- Review service logs for startup issues

## Best Practices

1. **Always test locally before pushing**
   ```bash
   # Backend
   go test ./...
   
   # Frontend
   pnpm build
   ```

2. **Keep secrets secure**
   - Never commit secrets to code
   - Use GitHub Secrets for sensitive data
   - Rotate credentials regularly

3. **Monitor CI performance**
   - Review workflow run times
   - Optimize slow steps
   - Use caching effectively

4. **Version your images**
   - Use semantic versioning for releases
   - Tag images with git SHA for traceability
   - Keep `latest` tag updated

## Next Steps

1. Add unit test coverage reporting
2. Implement E2E tests with Playwright/Cypress
3. Add performance testing
4. Set up staging/production deployments
5. Implement blue-green deployment strategy
6. Add security scanning (Trivy, Snyk)
7. Set up monitoring and alerting

## Support

For issues or questions:
1. Check workflow logs in GitHub Actions tab
2. Review this documentation
3. Open an issue in the repository
4. Contact the development team
