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

<div class="container">
	<header class="header">
		<div class="header-content">
			<h1>My Posts</h1>
			<div class="header-actions">
				<button class="btn btn-primary" on:click={openCreateModal}>
					<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
		<div class="error-message">
			{error}
			<button class="error-close" on:click={() => error = ''}>×</button>
		</div>
	{/if}

	<main class="main-content">
		{#if isLoading}
			<div class="loading">
				<div class="spinner"></div>
				<p>Loading posts...</p>
			</div>
		{:else if posts.length === 0}
			<div class="empty-state">
				<svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
					<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
					<polyline points="14,2 14,8 20,8"></polyline>
					<line x1="16" y1="13" x2="8" y2="13"></line>
					<line x1="16" y1="17" x2="8" y2="17"></line>
					<polyline points="10,9 9,9 8,9"></polyline>
				</svg>
				<h2>No posts yet</h2>
				<p>Create your first post to get started!</p>
				<button class="btn btn-primary" on:click={openCreateModal}>Create Post</button>
			</div>
		{:else}
			<div class="posts-grid">
				{#each posts as post (post.uuid)}
					<article class="post-card">
						<div class="post-header">
							<h3 class="post-title">{post.title}</h3>
							<div class="post-actions">
								<button class="btn-icon" on:click={() => openEditModal(post)} title="Edit">
									<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
										<path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
										<path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
									</svg>
								</button>
								<button class="btn-icon btn-icon-danger" on:click={() => handleDeletePost(post)} title="Delete">
									<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
										<polyline points="3,6 5,6 21,6"></polyline>
										<path d="M19,6v14a2,2,0,0,1-2,2H7a2,2,0,0,1-2-2V6m3,0V4a2,2,0,0,1,2-2h4a2,2,0,0,1,2,2V6"></path>
									</svg>
								</button>
							</div>
						</div>
						<div class="post-content">
							<p>{post.content}</p>
						</div>
						<div class="post-footer">
							<small>Created: {formatDate(post.createdAt)}</small>
							{#if post.updatedAt !== post.createdAt}
								<small>Updated: {formatDate(post.updatedAt)}</small>
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
	<div class="modal-overlay" on:click={closeModals}>
		<div class="modal" on:click|stopPropagation>
			<div class="modal-header">
				<h2>Create New Post</h2>
				<button class="modal-close" on:click={closeModals}>×</button>
			</div>
			<div class="modal-body">
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
						class="form-input form-textarea"
						bind:value={content}
						placeholder="Enter post content"
						required
					></textarea>
				</div>
			</div>
			<div class="modal-footer">
				<button class="btn btn-secondary" on:click={closeModals}>Cancel</button>
				<button class="btn btn-primary" on:click={handleCreatePost}>Create Post</button>
			</div>
		</div>
	</div>
{/if}

<!-- Edit Modal -->
{#if showEditModal}
	<div class="modal-overlay" on:click={closeModals}>
		<div class="modal" on:click|stopPropagation>
			<div class="modal-header">
				<h2>Edit Post</h2>
				<button class="modal-close" on:click={closeModals}>×</button>
			</div>
			<div class="modal-body">
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
						class="form-input form-textarea"
						bind:value={content}
						placeholder="Enter post content"
						required
					></textarea>
				</div>
			</div>
			<div class="modal-footer">
				<button class="btn btn-secondary" on:click={closeModals}>Cancel</button>
				<button class="btn btn-primary" on:click={handleUpdatePost}>Update Post</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.header {
		background: white;
		border-bottom: 1px solid #e2e8f0;
		padding: 1rem 0;
		margin-bottom: 2rem;
	}

	.header-content {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.header h1 {
		font-size: 2rem;
		font-weight: 700;
		color: #1a202c;
		margin: 0;
	}

	.header-actions {
		display: flex;
		gap: 1rem;
	}

	.main-content {
		min-height: 60vh;
	}

	.error-message {
		background: #fed7d7;
		color: #c53030;
		padding: 1rem;
		border-radius: 8px;
		margin-bottom: 1rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.error-close {
		background: none;
		border: none;
		font-size: 1.5rem;
		cursor: pointer;
		color: #c53030;
	}

	.loading {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 4rem 0;
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 4px solid #e2e8f0;
		border-top: 4px solid #667eea;
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin-bottom: 1rem;
	}

	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}

	.empty-state {
		text-align: center;
		padding: 4rem 0;
		color: #718096;
	}

	.empty-state svg {
		margin-bottom: 1rem;
		color: #cbd5e0;
	}

	.empty-state h2 {
		font-size: 1.5rem;
		font-weight: 600;
		margin-bottom: 0.5rem;
		color: #4a5568;
	}

	.posts-grid {
		display: grid;
		gap: 1.5rem;
	}

	.post-card {
		background: white;
		border-radius: 12px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
		padding: 1.5rem;
		transition: transform 0.2s ease, box-shadow 0.2s ease;
	}

	.post-card:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
	}

	.post-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 1rem;
	}

	.post-title {
		font-size: 1.25rem;
		font-weight: 600;
		color: #1a202c;
		margin: 0;
		flex: 1;
	}

	.post-actions {
		display: flex;
		gap: 0.5rem;
	}

	.btn-icon {
		background: none;
		border: none;
		padding: 0.5rem;
		border-radius: 6px;
		cursor: pointer;
		color: #718096;
		transition: all 0.2s ease;
	}

	.btn-icon:hover {
		background: #f7fafc;
		color: #4a5568;
	}

	.btn-icon-danger:hover {
		background: #fed7d7;
		color: #c53030;
	}

	.post-content {
		margin-bottom: 1rem;
	}

	.post-content p {
		color: #4a5568;
		line-height: 1.6;
		margin: 0;
	}

	.post-footer {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding-top: 1rem;
		border-top: 1px solid #e2e8f0;
	}

	.post-footer small {
		color: #718096;
		font-size: 0.875rem;
	}

	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		padding: 1rem;
	}

	.modal {
		background: white;
		border-radius: 12px;
		width: 100%;
		max-width: 500px;
		max-height: 90vh;
		overflow-y: auto;
	}

	.modal-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1.5rem;
		border-bottom: 1px solid #e2e8f0;
	}

	.modal-header h2 {
		font-size: 1.25rem;
		font-weight: 600;
		margin: 0;
	}

	.modal-close {
		background: none;
		border: none;
		font-size: 1.5rem;
		cursor: pointer;
		color: #718096;
		padding: 0;
		width: 24px;
		height: 24px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.modal-body {
		padding: 1.5rem;
	}

	.modal-footer {
		display: flex;
		justify-content: flex-end;
		gap: 1rem;
		padding: 1.5rem;
		border-top: 1px solid #e2e8f0;
	}

	@media (min-width: 768px) {
		.posts-grid {
			grid-template-columns: repeat(2, 1fr);
		}
	}

	@media (min-width: 1024px) {
		.posts-grid {
			grid-template-columns: repeat(3, 1fr);
		}
	}
</style>
