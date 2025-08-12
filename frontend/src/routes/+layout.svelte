<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';

	onMount(() => {
		console.log('=== LAYOUT ONMOUNT ===');
		// Check if user is authenticated on app load
		const token = localStorage.getItem('accessToken');
		console.log('Layout onMount - token found:', !!token);
		if (token) {
			// Parse user data from JWT token
			let userData = null;
			try {
				const payload = JSON.parse(atob(token.split('.')[1]));
				userData = {
					id: payload.userid,
					email: payload.email || '',
					name: payload.username || '',
					avatar: undefined,
					provider: 'local',
					isVerified: true,
					domains: payload.domains || []
				};
			} catch (error) {
				console.error('Error parsing JWT token:', error);
			}
			
			authStore.set({ isAuthenticated: true, token, user: userData, isAuthenticating: false });
			console.log('Layout onMount - auth store set to authenticated');
		}
	});

	// Watch for authentication changes - only on client side
	$: {
		console.log('=== LAYOUT REACTIVE STATEMENT ===');
		console.log('Layout reactive statement - isAuthenticated:', $authStore.isAuthenticated, 'isAuthenticating:', $authStore.isAuthenticating, 'pathname:', $page.url.pathname, 'browser:', browser);
		// Only redirect if we're not currently authenticating and not on auth pages
		if (browser && !$authStore.isAuthenticating && !$authStore.isAuthenticated && 
			$page.url.pathname !== '/login' && $page.url.pathname !== '/register' && $page.url.pathname !== '/' && $page.url.pathname !== '/auth/callback') {
			console.log('Layout redirecting to login - isAuthenticated:', $authStore.isAuthenticated, 'pathname:', $page.url.pathname);
			goto('/login');
		}
	}
</script>

<main class="min-h-screen bg-gradient-to-br from-blue-50 via-blue-25 to-indigo-50">
	<slot />
</main>
