<script lang="ts">
	import { goto } from '$app/navigation';
	import { authActions, authStore, type LoginRequest } from '$lib/stores/auth';
	import { onMount } from 'svelte';

	let email = '';
	let password = '';
	let isLoading = false;
	let error = '';

	onMount(() => {
		// Check for error messages in URL parameters
		const urlParams = new URLSearchParams(window.location.search);
		const urlError = urlParams.get('error');
		if (urlError) {
			error = decodeURIComponent(urlError);
			// Clear the error from URL
			window.history.replaceState({}, document.title, '/login');
		}
	});

	// Watch for authentication changes and redirect when authenticated
	$: if ($authStore.isAuthenticated && !isLoading) {
		console.log('Login page - user is authenticated, redirecting to /posts');
		goto('/posts');
	}

	async function handleLogin() {
		if (!email || !password) {
			error = 'Please fill in all fields';
			return;
		}

		isLoading = true;
		error = '';

		console.log('Login page - starting login process');
		const credentials: LoginRequest = { email, password };
		const success = await authActions.login(credentials);
		console.log('Login page - login result:', success);

		if (!success) {
			error = 'Invalid email or password';
		}

		isLoading = false;
	}

	function handleSocialLogin(provider: string) {
		window.location.href = `/auth/${provider}`;
	}
</script>

<svelte:head>
	<title>Login - GoWeb</title>
</svelte:head>

<div class="container">
	<div class="login-container">
		<div class="card">
			<div class="text-center mb-4">
				<h1>Welcome Back</h1>
				<p>Sign in to your account to continue</p>
			</div>

			{#if error}
				<div class="error-message">
					{error}
				</div>
			{/if}

			<form on:submit|preventDefault={handleLogin}>
				<div class="form-group">
					<label for="email" class="form-label">Email</label>
					<input
						type="email"
						id="email"
						class="form-input"
						bind:value={email}
						placeholder="Enter your email"
						required
					/>
				</div>

				<div class="form-group">
					<label for="password" class="form-label">Password</label>
					<input
						type="password"
						id="password"
						class="form-input"
						bind:value={password}
						placeholder="Enter your password"
						required
					/>
				</div>

				<button type="submit" class="btn btn-primary" disabled={isLoading}>
					{isLoading ? 'Signing in...' : 'Sign In'}
				</button>
			</form>

			<div class="divider">
				<span>or</span>
			</div>

			<div class="social-buttons">
				<button
					type="button"
					class="btn btn-secondary social-btn"
					on:click={() => handleSocialLogin('google')}
				>
					<svg width="20" height="20" viewBox="0 0 24 24">
						<path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
						<path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
						<path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
						<path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
					</svg>
					Continue with Google
				</button>

				<button
					type="button"
					class="btn btn-secondary social-btn"
					on:click={() => handleSocialLogin('github')}
				>
					<svg width="20" height="20" viewBox="0 0 24 24">
						<path fill="currentColor" d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
					</svg>
					Continue with GitHub
				</button>
			</div>

			<div class="text-center mt-4">
				<p>
					Don't have an account?
					<a href="/register" class="link">Sign up</a>
				</p>
			</div>
		</div>
	</div>
</div>

<style>
	.login-container {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 2rem 0;
	}

	.card {
		width: 100%;
		max-width: 400px;
	}

	h1 {
		font-size: 2rem;
		font-weight: 700;
		color: #1a202c;
		margin-bottom: 0.5rem;
	}

	p {
		color: #718096;
	}

	.error-message {
		background: #fed7d7;
		color: #c53030;
		padding: 0.75rem;
		border-radius: 8px;
		margin-bottom: 1rem;
		font-size: 0.875rem;
	}

	.divider {
		text-align: center;
		margin: 1.5rem 0;
		position: relative;
	}

	.divider::before {
		content: '';
		position: absolute;
		top: 50%;
		left: 0;
		right: 0;
		height: 1px;
		background: #e2e8f0;
	}

	.divider span {
		background: white;
		padding: 0 1rem;
		color: #718096;
		font-size: 0.875rem;
	}

	.social-buttons {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.social-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		width: 100%;
	}

	.link {
		color: #667eea;
		text-decoration: none;
		font-weight: 500;
	}

	.link:hover {
		text-decoration: underline;
	}

	.btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}
</style>
