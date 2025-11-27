# GitHub Secrets Configuration Guide

This file documents all secrets needed for the CI/CD pipeline.
DO NOT commit actual secret values to version control!

## Required Secrets for Docker Hub Deployment

### DOCKER_USERNAME
Description: Your Docker Hub username
Example: johndoe
Required: Yes (for Docker Hub push)
Where to get: Your Docker Hub account name

### DOCKER_PASSWORD
Description: Docker Hub access token (NOT your account password)
Example: dckr_pat_xxxxxxxxxxxxxxxxxxxxx
Required: Yes (for Docker Hub push)
Where to get:
  1. Login to hub.docker.com
  2. Account Settings → Security
  3. New Access Token
  4. Give it Read, Write, Delete permissions
  5. Copy the generated token

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

## How to Add Secrets to GitHub

1. Go to your GitHub repository
2. Click Settings
3. In left sidebar, click "Secrets and variables" → "Actions"
4. Click "New repository secret"
5. Enter the Name (exactly as shown above, case-sensitive)
6. Enter the Value (your actual secret)
7. Click "Add secret"

## Testing Secret Configuration

After adding secrets, push to main branch and check:

```bash
# Check if workflows can access secrets
git add .
git commit -m "Test CI/CD with secrets"
git push origin main

# Then go to GitHub → Actions tab
# Look for successful Docker push jobs
```

## Security Best Practices

1. ✅ Never commit secrets to git
2. ✅ Use access tokens, not passwords
3. ✅ Rotate secrets regularly (every 90 days)
4. ✅ Use separate tokens for different purposes
5. ✅ Set minimum required permissions
6. ✅ Monitor secret usage in GitHub Actions logs
7. ✅ Revoke tokens when no longer needed

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

### "Invalid credentials" error
- Verify DOCKER_USERNAME is correct (case-sensitive)
- Regenerate Docker Hub access token
- Make sure you're using the token, not your password
- Re-add the secret to GitHub

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

- [ ] DOCKER_USERNAME added to GitHub Secrets
- [ ] DOCKER_PASSWORD (token) added to GitHub Secrets
- [ ] Docker Hub token has write permissions
- [ ] Tested workflow runs successfully
- [ ] Docker images appear in Docker Hub
- [ ] (Optional) AWS secrets configured if using ECR
- [ ] (Optional) ECR repositories created
- [ ] All secret names match exactly
- [ ] No secrets committed to git

## Contact

If you need help configuring secrets:
1. Check GitHub Actions documentation
2. Review workflow logs for specific errors
3. Verify credentials work locally with Docker CLI
4. Open an issue with error details (never post actual secrets!)
