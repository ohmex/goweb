# Social Login Setup Guide

This guide explains how to configure Google and GitHub OAuth for social login in your Go web MVC application.

## Environment Variables

Add the following environment variables to your `.env` file:

```bash
# Google OAuth Configuration
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback

# GitHub OAuth Configuration
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=http://localhost:8080/auth/github/callback
```

## Google OAuth Setup

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google+ API
4. Go to "Credentials" → "Create Credentials" → "OAuth 2.0 Client IDs"
5. Configure the OAuth consent screen:
   - Add your domain to authorized domains
   - Add scopes: `email`, `profile`
6. Create OAuth 2.0 Client ID:
   - Application type: Web application
   - Authorized redirect URIs: `http://localhost:8080/auth/google/callback` (for development)
   - Copy the Client ID and Client Secret to your `.env` file

## GitHub OAuth Setup

1. Go to [GitHub Developer Settings](https://github.com/settings/developers)
2. Click "New OAuth App"
3. Fill in the application details:
   - Application name: Your app name
   - Homepage URL: `http://localhost:8080` (for development)
   - Authorization callback URL: `http://localhost:8080/auth/github/callback`
4. Click "Register application"
5. Copy the Client ID and Client Secret to your `.env` file

## API Endpoints

### Google OAuth
- **Initiate Login**: `GET /auth/google`
- **Callback**: `GET /auth/google/callback`

### GitHub OAuth
- **Initiate Login**: `GET /auth/github`
- **Callback**: `GET /auth/github/callback`

## Usage

1. Users can visit `/auth/google` or `/auth/github` to initiate social login
2. They will be redirected to the respective OAuth provider
3. After authorization, they'll be redirected back to the callback URL
4. The system will create or link their account and return JWT tokens

## Database Changes

The User model has been updated with the following new fields:
- `avatar`: User's profile picture URL
- `provider`: OAuth provider (google, github, local)
- `provider_id`: Unique ID from the OAuth provider
- `is_verified`: Email verification status

## Security Notes

- The system automatically links existing accounts by email address
- Social login users don't have passwords (empty password field)
- Email verification is automatically set based on provider verification status
- Consider implementing state parameter validation for additional security

## Testing

1. Start your application
2. Visit `http://localhost:8080/auth/google` or `http://localhost:8080/auth/github`
3. Complete the OAuth flow
4. You should receive JWT tokens upon successful authentication
