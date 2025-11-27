# GitHub Secrets Configuration Guide

This file documents all secrets needed for the CI/CD pipeline.
DO NOT commit actual secret values to version control!

## Required Secrets

### GITHUB_TOKEN (Automatic - No Setup Needed!)
Description: Used for pushing to GitHub Container Registry (GHCR)
Required: Yes (automatically provided by GitHub Actions)
Where to get: Automatically available in all workflows
Note: No configuration needed - GitHub provides this automatically!

### SNYK_TOKEN (Optional - for security scanning)
Description: Snyk API token for vulnerability scanning
Example: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
Required: No (optional for security scans)
Where to get:
  1. Sign up at https://snyk.io (free tier available)
  2. Go to Account Settings
  3. Copy your API token
  4. Add as SNYK_TOKEN secret in GitHub

## Optional Secrets for AWS ECR Deployment

### AWS_REGION
Description: AWS region where ECR repositories will be created
Example: us-east-1
Required: No (only if using AWS ECR)
Common values: us-east-1, us-west-2, eu-west-1, ap-south-1

### AWS_ACCESS_KEY_ID
Description: AWS IAM access key ID
Example: AKIAIOSFODNN7EXAMPLE
Required: No (only if using AWS ECR)
Where to get:
  1. AWS Console → IAM
  2. Create user with ECR permissions
  3. Generate access key
  4. Copy Access Key ID

### AWS_SECRET_ACCESS_KEY
Description: AWS IAM secret access key
Example: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
Required: No (only if using AWS ECR)
Where to get:
  1. Generated at same time as Access Key ID
  2. Only shown once - save it immediately!

## How to Add Secrets to GitHub (Optional Secrets Only)

### For Snyk Token (Optional):
1. Go to your GitHub repository
2. Click Settings
3. In left sidebar, click "Secrets and variables" → "Actions"
4. Click "New repository secret"
5. Name: SNYK_TOKEN
6. Value: Your Snyk API token
7. Click "Add secret"

### For AWS ECR (Optional):
Follow the same steps for AWS_REGION, AWS_ACCESS_KEY_ID, and AWS_SECRET_ACCESS_KEY

## Testing Configuration

Push to main branch - Docker images will automatically push to GHCR:

```bash
git add .
git commit -m "Update CI/CD workflows"
git push origin main

# Check GitHub Container Registry:
# https://github.com/YOUR_USERNAME?tab=packages
```

## View Your Docker Images

After successful workflow run:
- Go to: https://github.com/YOUR_USERNAME/HMS/pkgs/container
- All images will be public by default
- Image names: ghcr.io/YOUR_USERNAME/hms/hms-{service}:latest

## Security Features

### GitHub Container Registry (GHCR)
- ✅ Automatic authentication via GITHUB_TOKEN
- ✅ No manual secrets to manage
- ✅ Built-in access control
- ✅ Free for public repositories
- ✅ Integrated with GitHub UI

### Snyk Security Scanning
- ✅ Scans for vulnerabilities in Docker images
- ✅ Checks dependencies for known CVEs
- ✅ Provides remediation advice
- ✅ Free tier available
- ✅ Continues workflow even if issues found

## Security Best Practices

1. ✅ Images automatically pushed to GHCR (no manual tokens!)
2. ✅ Snyk scans every image for vulnerabilities
3. ✅ Multi-stage Docker builds for smaller attack surface
4. ✅ Non-root users in containers
5. ✅ Automatic security updates via dependabot
6. ✅ Monitor vulnerabilities in GitHub Security tab

## IAM Policy for AWS ECR (if using)

Minimum required permissions for AWS user:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:PutImage",
        "ecr:InitiateLayerUpload",
        "ecr:UploadLayerPart",
        "ecr:CompleteLayerUpload",
        "ecr:CreateRepository",
        "ecr:DescribeRepositories",
        "ecr:ListImages"
      ],
      "Resource": "*"
    }
  ]
}
```

## Troubleshooting

### "permission denied" for GHCR
- Check repository Settings → Actions → General
- Ensure "Read and write permissions" is enabled
- Workflow permissions should allow package writes

### Snyk scan failing
- Verify SNYK_TOKEN is correct
- Check Snyk account is active
- Review Snyk dashboard for issues
- Note: Scans continue even with vulnerabilities found

### AWS ECR authentication fails
- Check AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
- Verify IAM user has ECR permissions
- Ensure AWS_REGION is valid
- Check if access key is active in AWS Console

### Secrets not working in workflows
- Secret names must match exactly (case-sensitive)
- Secrets are only available after push to main
- Check workflow logs for "Secret not found" errors
- Verify secrets are set at repository level, not environment level

## Verification Checklist

Before pushing to production:

- [x] GITHUB_TOKEN automatically configured (no setup needed!)
- [ ] (Optional) SNYK_TOKEN added for security scanning
- [ ] Workflow permissions set to "Read and write"
- [ ] (Optional) AWS secrets configured if using ECR
- [x] GitHub Container Registry enabled (automatic)
- [x] Images will be public by default
- [x] No manual secrets needed for basic CI/CD!
- [ ] No secrets committed to git

## Contact

If you need help configuring secrets:
1. Check GitHub Actions documentation
2. Review workflow logs for specific errors
3. Verify credentials work locally with Docker CLI
4. Open an issue with error details (never post actual secrets!)
