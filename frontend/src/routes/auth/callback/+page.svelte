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

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full">
		<div class="bg-white py-8 px-6 shadow-xl rounded-xl text-center">
			<h2 class="text-2xl font-semibold text-gray-900 mb-2">Social Login Callback</h2>
			<p class="text-gray-600 mb-6">This page should be accessible at /auth/callback</p>
			{#if isLoading}
				<div class="animate-spin w-10 h-10 border-4 border-gray-200 border-t-primary-600 rounded-full mx-auto mb-4"></div>
				<h3 class="text-lg font-medium text-gray-900 mb-2">Authenticating...</h3>
				<p class="text-gray-600">Please wait while we complete your authentication.</p>
			{:else if error}
				<div class="text-5xl mb-4">⚠️</div>
				<h3 class="text-lg font-medium text-gray-900 mb-2">Authentication Failed</h3>
				<p class="text-gray-600 mb-6">{error}</p>
				<a href="/login" class="btn btn-primary">Back to Login</a>
			{/if}
		</div>
	</div>
</div>
