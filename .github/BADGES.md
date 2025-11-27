# CI/CD Status Badges

Add these badges to your main README.md to show workflow status:

## Backend CI
```markdown
[![Backend CI](https://github.com/YOUR_USERNAME/HMS/actions/workflows/backend-ci.yml/badge.svg)](https://github.com/YOUR_USERNAME/HMS/actions/workflows/backend-ci.yml)
```

## Frontend CI
```markdown
[![Frontend CI](https://github.com/YOUR_USERNAME/HMS/actions/workflows/frontend-ci.yml/badge.svg)](https://github.com/YOUR_USERNAME/HMS/actions/workflows/frontend-ci.yml)
```

## Full Stack CI
```markdown
[![Full Stack CI](https://github.com/YOUR_USERNAME/HMS/actions/workflows/fullstack-ci.yml/badge.svg)](https://github.com/YOUR_USERNAME/HMS/actions/workflows/fullstack-ci.yml)
```

## Docker Build
```markdown
[![Docker Build](https://github.com/YOUR_USERNAME/HMS/actions/workflows/docker-push.yml/badge.svg)](https://github.com/YOUR_USERNAME/HMS/actions/workflows/docker-push.yml)
```

## All Badges Together

```markdown
[![Backend CI](https://github.com/YOUR_USERNAME/HMS/actions/workflows/backend-ci.yml/badge.svg)](https://github.com/YOUR_USERNAME/HMS/actions/workflows/backend-ci.yml)
[![Frontend CI](https://github.com/YOUR_USERNAME/HMS/actions/workflows/frontend-ci.yml/badge.svg)](https://github.com/YOUR_USERNAME/HMS/actions/workflows/frontend-ci.yml)
[![Full Stack CI](https://github.com/YOUR_USERNAME/HMS/actions/workflows/fullstack-ci.yml/badge.svg)](https://github.com/YOUR_USERNAME/HMS/actions/workflows/fullstack-ci.yml)
[![Docker Build](https://github.com/YOUR_USERNAME/HMS/actions/workflows/docker-push.yml/badge.svg)](https://github.com/YOUR_USERNAME/HMS/actions/workflows/docker-push.yml)
```

**Note**: Replace `YOUR_USERNAME` with your actual GitHub username (e.g., `Tshewangdorji7257`)

## Alternative: Dynamic Badge with Branch

For main branch specifically:
```markdown
![Backend CI](https://github.com/YOUR_USERNAME/HMS/actions/workflows/backend-ci.yml/badge.svg?branch=main)
```

## Docker Hub Badges

Add these to show Docker image status:

```markdown
![Docker Pulls](https://img.shields.io/docker/pulls/YOUR_DOCKERHUB_USERNAME/hms-auth-service)
![Docker Image Size](https://img.shields.io/docker/image-size/YOUR_DOCKERHUB_USERNAME/hms-auth-service/latest)
```
