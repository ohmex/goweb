import { authActions } from '$lib/stores/auth';

export interface Post {
	id: number;
	uuid: string;
	title: string;
	content: string;
	userID: number;
	user: {
		id: number;
		name: string;
		email: string;
	};
	created_at: string;
	updated_at: string;
}

export interface CreatePostRequest {
	title: string;
	content: string;
}

export interface UpdatePostRequest {
	title: string;
	content: string;
}

class ApiService {
	private baseUrl = '/api';

	private async makeRequest<T>(
		endpoint: string,
		options: RequestInit = {}
	): Promise<T> {
		const url = `${this.baseUrl}${endpoint}`;
		const headers = {
			'Content-Type': 'application/json',
			...authActions.getAuthHeaders(),
			...options.headers
		};

		console.log('API Service - Making request to:', url);
		console.log('API Service - Headers:', headers);

		// Add domain header for API requests (excluding auth endpoints)
		if (!endpoint.startsWith('/auth') && !endpoint.includes('/login') && !endpoint.includes('/register') && !endpoint.includes('/refresh') && !endpoint.includes('/logout')) {
			// For now, use the first domain from the user's domains
			// In a real app, you might want to let users select a domain
			const token = localStorage.getItem('accessToken');
			if (token) {
				try {
					const payload = JSON.parse(atob(token.split('.')[1]));
					if (payload.domains && payload.domains.length > 0) {
						headers['domain'] = payload.domains[0].UUID;
					} else {
						console.warn('User has no domains assigned');
					}
				} catch (error) {
					console.error('Error parsing JWT token:', error);
				}
			}
		}

		try {
			const response = await fetch(url, {
				...options,
				headers
			});

			console.log('API Service - Response status:', response.status);
			console.log('API Service - Response headers:', response.headers);

			if (response.status === 401) {
				// Try to refresh token
				const refreshed = await authActions.refreshToken();
				if (refreshed) {
					// Retry the request with new token
					const retryHeaders = {
						'Content-Type': 'application/json',
						...authActions.getAuthHeaders(),
						...options.headers
					};
					const retryResponse = await fetch(url, {
						...options,
						headers: retryHeaders
					});
					
					if (!retryResponse.ok) {
						throw new Error(`HTTP error! status: ${retryResponse.status}`);
					}
					
					return await retryResponse.json();
				} else {
					throw new Error('Authentication failed');
				}
			}

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			return await response.json();
		} catch (error) {
			console.error('API request failed:', error);
			throw error;
		}
	}

	// Posts API
	async getPosts(): Promise<Post[]> {
		return this.makeRequest<Post[]>('/post');
	}

	async getPost(uuid: string): Promise<Post> {
		return this.makeRequest<Post>(`/post/${uuid}`);
	}

	async createPost(data: CreatePostRequest): Promise<Post> {
		return this.makeRequest<Post>('/post', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async updatePost(uuid: string, data: UpdatePostRequest): Promise<Post> {
		return this.makeRequest<Post>(`/post/${uuid}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	}

	async deletePost(uuid: string): Promise<void> {
		return this.makeRequest<void>(`/post/${uuid}`, {
			method: 'DELETE'
		});
	}

	// User API
	async getCurrentUser(): Promise<any> {
		// This endpoint doesn't exist in the backend yet
		// For now, return the user data from the JWT token
		const token = localStorage.getItem('accessToken');
		if (token) {
			try {
				const payload = JSON.parse(atob(token.split('.')[1]));
				return {
					id: payload.userid,
					email: payload.email || '',
					name: payload.username || '',
					domains: payload.domains || []
				};
			} catch (error) {
				console.error('Error parsing JWT token:', error);
				throw new Error('Failed to get current user');
			}
		}
		throw new Error('No access token found');
	}
}

export const apiService = new ApiService();
