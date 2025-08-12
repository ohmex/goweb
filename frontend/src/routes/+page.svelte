<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';

	onMount(() => {
		// Redirect based on authentication status
		if ($authStore.isAuthenticated) {
			goto('/posts');
		} else {
			goto('/login');
		}
	});
</script>

<div class="min-h-screen flex flex-col items-center justify-center bg-gray-50">
	<div class="animate-spin w-10 h-10 border-4 border-gray-200 border-t-primary-600 rounded-full mb-4"></div>
	<p class="text-gray-600 text-lg">Loading...</p>
</div>
