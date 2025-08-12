# GoWeb Frontend

A modern SvelteKit frontend application for the GoWeb backend project.

## Features

- **Authentication**: Email/password login and registration
- **Social Login**: Google and GitHub OAuth integration
- **Posts Management**: Create, read, update, and delete posts
- **Responsive Design**: Modern UI that works on all devices
- **Real-time Updates**: Automatic token refresh and error handling

## Prerequisites

- Node.js 18+ 
- pnpm 8.0+ (recommended package manager)
- The GoWeb backend server running on `http://localhost:8080`

## Installation

1. Install pnpm (if not already installed):
```bash
npm install -g pnpm
```

2. Install dependencies:
```bash
pnpm install
```

3. Start the development server:
```bash
pnpm dev
```

4. Open your browser and navigate to `http://localhost:3000`

## Development

The application uses:
- **SvelteKit** for the framework
- **TypeScript** for type safety
- **Vite** for fast development and building
- **CSS** for styling (no external UI libraries)
- **pnpm** for package management

## Project Structure

```
src/
├── app.css              # Global styles
├── app.html             # HTML template
├── app.d.ts             # TypeScript declarations
├── routes/              # SvelteKit routes
│   ├── +layout.svelte   # Root layout
│   ├── +page.svelte     # Root page (redirects)
│   ├── login/           # Login page
│   ├── register/        # Registration page
│   └── posts/           # Posts management page
└── lib/
    ├── stores/          # Svelte stores
    │   └── auth.ts      # Authentication store
    └── services/        # API services
        └── api.ts       # API client
```

## API Integration

The frontend communicates with the GoWeb backend through:
- **Authentication endpoints**: `/login`, `/register`, `/logout`, `/refresh`
- **Social OAuth**: `/auth/google`, `/auth/github`
- **Posts API**: `/api/post` (CRUD operations)

## Available Scripts

- `pnpm dev` - Start development server
- `pnpm build` - Build for production
- `pnpm preview` - Preview production build

## Building for Production

```bash
pnpm build
```

The built application will be in the `build` directory.

## Configuration

The application is configured to proxy API requests to the backend server running on `http://localhost:8080`. You can modify the proxy settings in `vite.config.ts` if needed.

## pnpm Benefits

This project uses pnpm for the following benefits:
- **Faster installation**: pnpm is significantly faster than npm
- **Disk space efficiency**: Uses hard links and symlinks to save disk space
- **Strict dependency management**: Prevents phantom dependencies
- **Better monorepo support**: Excellent for managing multiple packages
- **Deterministic installs**: Consistent node_modules structure across environments
