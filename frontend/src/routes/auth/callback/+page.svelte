<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authStore, authActions } from '$lib/stores/auth';
	import { page } from '$app/stores';

	let isLoading = true;
	let error = '';

	onMount(async () => {
		console.log('=== SOCIAL LOGIN CALLBACK STARTED ===');
		console.log('Current URL:', window.location.href);
		console.log('Current auth store state:', $authStore);
		
		try {
			// Get tokens from URL parameters
			const urlParams = new URLSearchParams(window.location.search);
			const accessToken = urlParams.get('access_token');
			const refreshToken = urlParams.get('refresh_token');
			const error = urlParams.get('error');

			console.log('URL Parameters:', {
				accessToken: accessToken ? 'PRESENT' : 'MISSING',
				refreshToken: refreshToken ? 'PRESENT' : 'MISSING',
				error: error || 'NONE'
			});

			if (error) {
				console.error('Social auth error:', error);
				await goto('/login?error=' + encodeURIComponent(error));
				return;
			}

			if (!accessToken || !refreshToken) {
				console.error('Missing tokens in callback');
				await goto('/login?error=' + encodeURIComponent('Authentication failed - missing tokens'));
				return;
			}

			console.log('Tokens received, calling handleSocialAuth...');
			
			// Handle social authentication using auth store method
			await authActions.handleSocialAuth(accessToken, refreshToken);

			console.log('handleSocialAuth completed, auth store state:', $authStore);
			console.log('isAuthenticated:', $authStore.isAuthenticated);
			console.log('token:', $authStore.token ? 'PRESENT' : 'MISSING');
			console.log('user:', $authStore.user);
			
			// Test: Show the parsed data for debugging
			if ($authStore.user) {
				console.log('User ID:', $authStore.user.id);
				console.log('User Name:', $authStore.user.name);
				console.log('User Email:', $authStore.user.email);
				console.log('User Domains:', $authStore.user.domains);
				console.log('User Domains Length:', $authStore.user.domains ? $authStore.user.domains.length : 'NULL');
			}
			
			// Clear URL parameters for security
			window.history.replaceState({}, document.title, '/auth/callback');
			
			console.log('Redirecting to /posts...');
			
			// Add a small delay to ensure the store is updated
			await new Promise(resolve => setTimeout(resolve, 100));
			
			// Redirect to posts page
			await goto('/posts');
			
			console.log('Redirect to /posts completed');
		} catch (err) {
			console.error('Callback error:', err);
			error = 'Authentication failed';
			await goto('/login?error=' + encodeURIComponent('Authentication failed'));
		} finally {
			isLoading = false;
		}
	});
</script>

<svelte:head>
	<title>Authenticating... - GoWeb</title>
</svelte:head>

<div class="container">
	<div class="loading-container">
		<div class="card">
			<div class="text-center">
				<h2>Social Login Callback</h2>
				<p>This page should be accessible at /auth/callback</p>
				{#if isLoading}
					<div class="spinner"></div>
					<h3>Authenticating...</h3>
					<p>Please wait while we complete your authentication.</p>
				{:else if error}
					<div class="error-icon">⚠️</div>
					<h3>Authentication Failed</h3>
					<p>{error}</p>
					<a href="/login" class="btn btn-primary">Back to Login</a>
				{/if}
			</div>
		</div>
	</div>
</div>

<style>
	.loading-container {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 2rem 0;
	}

	.card {
		width: 100%;
		max-width: 400px;
		text-align: center;
	}

	h2 {
		font-size: 1.5rem;
		font-weight: 600;
		color: #1a202c;
		margin-bottom: 0.5rem;
	}

	p {
		color: #718096;
		margin-bottom: 1.5rem;
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 4px solid #e2e8f0;
		border-top: 4px solid #667eea;
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin: 0 auto 1rem;
	}

	.error-icon {
		font-size: 3rem;
		margin-bottom: 1rem;
	}

	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}

	.btn {
		display: inline-block;
		padding: 0.75rem 1.5rem;
		border-radius: 8px;
		text-decoration: none;
		font-weight: 500;
		transition: all 0.2s;
	}

	.btn-primary {
		background: #667eea;
		color: white;
	}

	.btn-primary:hover {
		background: #5a67d8;
	}
</style>
