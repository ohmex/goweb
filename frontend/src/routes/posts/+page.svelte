<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authActions, authStore } from '$lib/stores/auth';
	import { apiService, type Post, type CreatePostRequest, type UpdatePostRequest } from '$lib/services/api';

	let posts: Post[] = [];
	let isLoading = true;
	let error = '';
	let showCreateModal = false;
	let showEditModal = false;
	let selectedPost: Post | null = null;

	// Form data
	let title = '';
	let content = '';

	onMount(async () => {
		console.log('=== POSTS PAGE ONMOUNT ===');
		console.log('Posts page - onMount called');
		console.log('Posts page - auth store state:', $authStore);
		console.log('Posts page - localStorage accessToken:', localStorage.getItem('accessToken'));
		
		// Check if user has domains
		if ($authStore.user && $authStore.user.domains && $authStore.user.domains.length === 0) {
			console.log('User has no domains, showing error');
			error = 'You do not have access to any domains. Please contact an administrator.';
			isLoading = false;
			return;
		}
		
		console.log('Loading posts...');
		await loadPosts();
		console.log('Posts page - loadPosts completed');
	});

	async function loadPosts() {
		try {
			console.log('Posts page - loadPosts started');
			console.log('Posts page - auth store state:', $authStore);
			console.log('Posts page - access token:', localStorage.getItem('accessToken'));
			isLoading = true;
			posts = await apiService.getPosts();
			console.log('Posts page - posts loaded successfully:', posts.length);
		} catch (err) {
			error = 'Failed to load posts';
			console.error('Error loading posts:', err);
		} finally {
			isLoading = false;
		}
	}

	async function handleLogout() {
		await authActions.logout();
		goto('/login');
	}

	function openCreateModal() {
		title = '';
		content = '';
		showCreateModal = true;
	}

	function openEditModal(post: Post) {
		selectedPost = post;
		title = post.title;
		content = post.content;
		showEditModal = true;
	}

	function closeModals() {
		showCreateModal = false;
		showEditModal = false;
		selectedPost = null;
		title = '';
		content = '';
	}

	async function handleCreatePost() {
		if (!title.trim() || !content.trim()) {
			error = 'Please fill in all fields';
			return;
		}

		try {
			const newPost: CreatePostRequest = { title: title.trim(), content: content.trim() };
			await apiService.createPost(newPost);
			await loadPosts();
			closeModals();
		} catch (err) {
			error = 'Failed to create post';
			console.error('Error creating post:', err);
		}
	}

	async function handleUpdatePost() {
		if (!selectedPost || !title.trim() || !content.trim()) {
			error = 'Please fill in all fields';
			return;
		}

		try {
			const updateData: UpdatePostRequest = { title: title.trim(), content: content.trim() };
			await apiService.updatePost(selectedPost.uuid, updateData);
			await loadPosts();
			closeModals();
		} catch (err) {
			error = 'Failed to update post';
			console.error('Error updating post:', err);
		}
	}

	async function handleDeletePost(post: Post) {
		if (!confirm('Are you sure you want to delete this post?')) {
			return;
		}

		try {
			await apiService.deletePost(post.uuid);
			await loadPosts();
		} catch (err) {
			error = 'Failed to delete post';
			console.error('Error deleting post:', err);
		}
	}

	function formatDate(dateString: string) {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<svelte:head>
	<title>Posts - GoWeb</title>
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
	<header class="bg-white border-b border-gray-200 pb-6 mb-8">
		<div class="flex justify-between items-center">
			<h1 class="text-3xl font-bold text-gray-900">My Posts</h1>
			<div class="flex gap-3">
				<button class="btn btn-primary" on:click={openCreateModal}>
					<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="mr-2">
						<line x1="12" y1="5" x2="12" y2="19"></line>
						<line x1="5" y1="12" x2="19" y2="12"></line>
					</svg>
					New Post
				</button>
				<button class="btn btn-secondary" on:click={handleLogout}>
					Logout
				</button>
			</div>
		</div>
	</header>

	{#if error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6 flex justify-between items-center">
			<span>{error}</span>
			<button class="text-red-400 hover:text-red-600 text-xl font-bold" on:click={() => error = ''}>×</button>
		</div>
	{/if}

	<main class="min-h-96">
		{#if isLoading}
			<div class="flex flex-col items-center justify-center py-16">
				<div class="animate-spin w-10 h-10 border-4 border-gray-200 border-t-primary-600 rounded-full mb-4"></div>
				<p class="text-gray-600">Loading posts...</p>
			</div>
		{:else if posts.length === 0}
			<div class="text-center py-16">
				<svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" class="mx-auto mb-4 text-gray-400">
					<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
					<polyline points="14,2 14,8 20,8"></polyline>
					<line x1="16" y1="13" x2="8" y2="13"></line>
					<line x1="16" y1="17" x2="8" y2="17"></line>
					<polyline points="10,9 9,9 8,9"></polyline>
				</svg>
				<h2 class="text-2xl font-semibold text-gray-900 mb-2">No posts yet</h2>
				<p class="text-gray-600 mb-6">Create your first post to get started!</p>
				<button class="btn btn-primary" on:click={openCreateModal}>Create Post</button>
			</div>
		{:else}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each posts as post (post.uuid)}
					<article class="bg-white rounded-xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow duration-200">
						<div class="flex justify-between items-start mb-4">
							<h3 class="text-xl font-semibold text-gray-900 flex-1 mr-3">{post.title}</h3>
							<div class="flex gap-2">
								<button 
									class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors" 
									on:click={() => openEditModal(post)} 
									title="Edit"
								>
									<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
										<path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
										<path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
									</svg>
								</button>
								<button 
									class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors" 
									on:click={() => handleDeletePost(post)} 
									title="Delete"
								>
									<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
										<polyline points="3,6 5,6 21,6"></polyline>
										<path d="M19,6v14a2,2,0,0,1-2,2H7a2,2,0,0,1-2-2V6m3,0V4a2,2,0,0,1,2-2h4a2,2,0,0,1,2,2V6"></path>
									</svg>
								</button>
							</div>
						</div>
						<div class="mb-4">
							<p class="text-gray-600 leading-relaxed">{post.content}</p>
						</div>
						<div class="flex justify-between items-center pt-4 border-t border-gray-100">
							<small class="text-gray-500 text-sm">Created: {formatDate(post.createdAt)}</small>
							{#if post.updatedAt !== post.createdAt}
								<small class="text-gray-500 text-sm">Updated: {formatDate(post.updatedAt)}</small>
							{/if}
						</div>
					</article>
				{/each}
			</div>
		{/if}
	</main>
</div>

<!-- Create Modal -->
{#if showCreateModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50" on:click={closeModals}>
		<div class="bg-white rounded-xl shadow-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto" on:click|stopPropagation>
			<div class="flex justify-between items-center p-6 border-b border-gray-200">
				<h2 class="text-xl font-semibold text-gray-900">Create New Post</h2>
				<button class="text-gray-400 hover:text-gray-600 text-2xl font-bold" on:click={closeModals}>×</button>
			</div>
			<div class="p-6">
				<div class="form-group">
					<label for="create-title" class="form-label">Title</label>
					<input
						type="text"
						id="create-title"
						class="form-input"
						bind:value={title}
						placeholder="Enter post title"
						required
					/>
				</div>
				<div class="form-group">
					<label for="create-content" class="form-label">Content</label>
					<textarea
						id="create-content"
						class="form-input resize-y min-h-[120px]"
						bind:value={content}
						placeholder="Enter post content"
						required
					></textarea>
				</div>
			</div>
			<div class="flex justify-end gap-3 p-6 border-t border-gray-200">
				<button class="btn btn-secondary" on:click={closeModals}>Cancel</button>
				<button class="btn btn-primary" on:click={handleCreatePost}>Create Post</button>
			</div>
		</div>
	</div>
{/if}

<!-- Edit Modal -->
{#if showEditModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50" on:click={closeModals}>
		<div class="bg-white rounded-xl shadow-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto" on:click|stopPropagation>
			<div class="flex justify-between items-center p-6 border-b border-gray-200">
				<h2 class="text-xl font-semibold text-gray-900">Edit Post</h2>
				<button class="text-gray-400 hover:text-gray-600 text-2xl font-bold" on:click={closeModals}>×</button>
			</div>
			<div class="p-6">
				<div class="form-group">
					<label for="edit-title" class="form-label">Title</label>
					<input
						type="text"
						id="edit-title"
						class="form-input"
						bind:value={title}
						placeholder="Enter post title"
						required
					/>
				</div>
				<div class="form-group">
					<label for="edit-content" class="form-label">Content</label>
					<textarea
						id="edit-content"
						class="form-input resize-y min-h-[120px]"
						bind:value={content}
						placeholder="Enter post content"
						required
					></textarea>
				</div>
			</div>
			<div class="flex justify-end gap-3 p-6 border-t border-gray-200">
				<button class="btn btn-secondary" on:click={closeModals}>Cancel</button>
				<button class="btn btn-primary" on:click={handleUpdatePost}>Update Post</button>
			</div>
		</div>
	</div>
{/if}
