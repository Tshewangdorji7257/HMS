# 🚀 CI/CD Quick Start Guide

This guide will help you set up the complete CI/CD pipeline for your HMS project in under 10 minutes.

## ✅ Prerequisites

- GitHub account
- Docker Hub account (for Docker image hosting)
- AWS account (optional, for ECR)

## 📋 Step-by-Step Setup

### Step 1: Verify Workflow Files

Ensure these workflow files exist in `.github/workflows/`:
- ✅ `backend-ci.yml` - Backend testing and building
- ✅ `frontend-ci.yml` - Frontend testing and building  
- ✅ `fullstack-ci.yml` - Combined full-stack pipeline
- ✅ `docker-push.yml` - Docker image building and pushing

### Step 2: Set Up Docker Hub

1. **Create Docker Hub Account**
   - Go to https://hub.docker.com
   - Sign up for a free account
   - Verify your email

2. **Generate Access Token**
   ```
   - Login to Docker Hub
   - Click on your username → Account Settings
   - Security → New Access Token
   - Name: "github-actions"
   - Permissions: Read, Write, Delete
   - Copy the token (you won't see it again!)
   ```

### Step 3: Add GitHub Secrets

1. **Navigate to Repository Settings**
   ```
   GitHub Repository → Settings → Secrets and variables → Actions
   ```

2. **Add Required Secrets**
   
   Click "New repository secret" for each:
   
   **For Docker Hub (Required)**:
   ```
   Name: DOCKER_USERNAME
   Value: your-docker-hub-username
   
   Name: DOCKER_PASSWORD
   Value: your-docker-hub-access-token
   ```
   
   **For AWS ECR (Optional)**:
   ```
   Name: AWS_REGION
   Value: us-east-1 (or your preferred region)
   
   Name: AWS_ACCESS_KEY_ID
   Value: your-aws-access-key
   
   Name: AWS_SECRET_ACCESS_KEY
   Value: your-aws-secret-key
   ```

### Step 4: Test the Pipeline

1. **Push to GitHub**
   ```bash
   git add .
   git commit -m "Add CI/CD workflows"
   git push origin main
   ```

2. **Check Workflow Status**
   ```
   - Go to your GitHub repository
   - Click on "Actions" tab
   - You should see workflows running
   - Click on a workflow to see detailed logs
   ```

3. **Expected Results**
   - ✅ All workflows should appear
   - ✅ `fullstack-ci.yml` runs automatically
   - ✅ Tests compile and run
   - ✅ Docker images build successfully
   - ✅ Images pushed to Docker Hub (if secrets configured)

## 🔍 Verify Everything Works

### Check Workflow Status
```bash
# All these should show green checkmarks in GitHub Actions:
- ✅ Backend CI Pipeline
- ✅ Frontend CI Pipeline  
- ✅ Full Stack CI Pipeline
- ✅ Docker Build and Push (if on main branch)
```

### Verify Docker Images
After successful push to main branch, check Docker Hub:

1. Login to https://hub.docker.com
2. Go to "Repositories"
3. You should see:
   - `hms-auth-service`
   - `hms-booking-service`
   - `hms-building-service`
   - `hms-api-gateway`
   - `hms-frontend`

### Test Local Pull
```bash
docker pull yourusername/hms-auth-service:latest
docker pull yourusername/hms-frontend:latest
```

## 🎯 What Each Workflow Does

### Backend CI (`backend-ci.yml`)
- ✅ Tests all 4 Go microservices
- ✅ Builds binaries
- ✅ Verifies gRPC proto files
- ✅ Runs integration tests with docker-compose
- ✅ Builds and pushes Docker images (on main)

### Frontend CI (`frontend-ci.yml`)
- ✅ Installs dependencies with pnpm
- ✅ Runs TypeScript type checking
- ✅ Runs ESLint
- ✅ Builds Next.js application
- ✅ Analyzes bundle size
- ✅ Builds and pushes Docker image (on main)

### Full Stack CI (`fullstack-ci.yml`)
- ✅ Detects which parts changed (backend/frontend)
- ✅ Runs appropriate tests in parallel
- ✅ Tests full integration with docker-compose
- ✅ Reports overall status

### Docker Push (`docker-push.yml`)
- ✅ Builds multi-platform images (amd64/arm64)
- ✅ Pushes to Docker Hub
- ✅ Pushes to AWS ECR (if configured)
- ✅ Tags with version/branch/SHA

## 🏃 Running Locally

### Test Backend Services
```bash
cd Backend/auth-service
go test -v ./...
go build -o main .
./main
```

### Test Frontend
```bash
cd Frontend
pnpm install
pnpm build
pnpm start
```

### Test Full Stack
```bash
cd Backend
docker-compose up --build
```

### Test Docker Builds
```bash
# Backend service
docker build -t hms-auth:test Backend/auth-service

# Frontend
docker build -t hms-frontend:test Frontend
```

## 🐛 Troubleshooting

### Workflows Not Running
**Problem**: No workflows appear in Actions tab
**Solution**: 
- Ensure workflow files are in `.github/workflows/`
- Check YAML syntax with online validator
- Push to `main` branch

### Docker Push Fails
**Problem**: "denied: requested access to the resource is denied"
**Solution**:
- Verify DOCKER_USERNAME is exact (case-sensitive)
- Regenerate Docker Hub access token
- Ensure token has Write permissions
- Re-add secrets to GitHub

### Tests Failing
**Problem**: Go tests fail
**Solution**:
```bash
# Run locally to see actual error
cd Backend/auth-service
go test -v ./...

# Fix any import or syntax issues
go mod tidy
```

### Build Timeout
**Problem**: Workflow times out
**Solution**:
- Check if services are starting properly
- Increase timeout in workflow file
- Review docker-compose health checks

## 📊 Monitoring CI/CD

### View Workflow History
```
GitHub Repository → Actions → All workflows
```

### Check Individual Run
```
Click on workflow run → Click on job → View logs
```

### Download Artifacts
```
Workflow run → Artifacts section → Download
```

## 🎓 Next Steps

Now that CI/CD is set up:

1. **Add More Tests**
   - Unit tests for handlers
   - Integration tests for APIs
   - E2E tests for critical paths

2. **Improve Coverage**
   - Add code coverage reporting
   - Set coverage thresholds
   - Generate coverage badges

3. **Add Security**
   - Scan Docker images with Trivy
   - Check dependencies with Snyk
   - Add SAST with CodeQL

4. **Setup Deployments**
   - Deploy to staging on PR merge
   - Deploy to production on tag
   - Add environment-specific configs

5. **Monitoring & Alerts**
   - Set up Slack/Discord notifications
   - Add performance monitoring
   - Track deployment metrics

## 📚 Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Hub Documentation](https://docs.docker.com/docker-hub/)
- [AWS ECR Documentation](https://docs.aws.amazon.com/ecr/)
- [Go Testing Guide](https://go.dev/doc/tutorial/add-a-test)
- [Next.js Deployment](https://nextjs.org/docs/deployment)

## 🆘 Need Help?

1. Check workflow logs in GitHub Actions
2. Review the main README in `.github/workflows/`
3. Test locally to reproduce issues
4. Open an issue in the repository

## ✨ Success Indicators

You'll know everything is working when:

- ✅ All workflow badges are green
- ✅ Docker images appear in Docker Hub
- ✅ Tests pass consistently
- ✅ No authentication errors
- ✅ Integration tests complete successfully

**Congratulations! Your CI/CD pipeline is now live! 🎉**
