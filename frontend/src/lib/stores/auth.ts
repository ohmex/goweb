import { writable } from 'svelte/store';

export interface AuthState {
	isAuthenticated: boolean;
	token: string | null;
	user: User | null;
	isAuthenticating: boolean;
}

export interface User {
	id: number;
	email: string;
	name: string;
	avatar?: string;
	provider: string;
	isVerified: boolean;
	domains?: Array<{
		UUID: string;
		Name: string;
	}>;
}

export interface LoginRequest {
	email: string;
	password: string;
}

export interface RegisterRequest {
	email: string;
	password: string;
	name: string;
}

export interface LoginResponse {
	accessToken: string;
	refreshToken: string;
	exp: number;
}

// Create the auth store
export const authStore = writable<AuthState>({
	isAuthenticated: false,
	token: null,
	user: null,
	isAuthenticating: false
});

// Auth actions
export const authActions = {
	login: async (credentials: LoginRequest): Promise<boolean> => {
		try {
			console.log('Attempting login with:', credentials);
			
			// Set authenticating flag
			authStore.update(state => ({ ...state, isAuthenticating: true }));
			
			const response = await fetch('/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(credentials)
			});

			console.log('Login response status:', response.status);
			if (!response.ok) {
				const errorData = await response.json();
				console.error('Login failed with error:', errorData);
				throw new Error('Login failed');
			}

			const data: LoginResponse = await response.json();
			console.log('Login successful, received data:', data);
			
			// Store tokens
			localStorage.setItem('accessToken', data.accessToken);
			localStorage.setItem('refreshToken', data.refreshToken);
			
			// Parse user data from JWT token
			let userData = null;
			try {
				const payload = JSON.parse(atob(data.accessToken.split('.')[1]));
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

			// Update store immediately
			authStore.set({
				isAuthenticated: true,
				token: data.accessToken,
				user: userData,
				isAuthenticating: false
			});

			console.log('Auth store updated, isAuthenticated should be true');
			return true;
		} catch (error) {
			console.error('Login error:', error);
			// Reset authenticating flag on error
			authStore.update(state => ({ ...state, isAuthenticating: false }));
			return false;
		}
	},

	register: async (userData: RegisterRequest): Promise<boolean> => {
		try {
			const response = await fetch('/register', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(userData)
			});

			if (!response.ok) {
				throw new Error('Registration failed');
			}

			const data: LoginResponse = await response.json();
			
			// Store tokens
			localStorage.setItem('accessToken', data.accessToken);
			localStorage.setItem('refreshToken', data.refreshToken);
			
			// Parse user data from JWT token
			let parsedUserData = null;
			try {
				const payload = JSON.parse(atob(data.accessToken.split('.')[1]));
				parsedUserData = {
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

			// Update store
			authStore.set({
				isAuthenticated: true,
				token: data.accessToken,
				user: parsedUserData,
				isAuthenticating: false
			});

			return true;
		} catch (error) {
			console.error('Registration error:', error);
			return false;
		}
	},

	logout: async (): Promise<void> => {
		try {
			const token = localStorage.getItem('accessToken');
			if (token) {
				await fetch('/logout', {
					method: 'POST',
					headers: {
						'Authorization': `Bearer ${token}`
					}
				});
			}
		} catch (error) {
			console.error('Logout error:', error);
		} finally {
			// Clear local storage and store
			localStorage.removeItem('accessToken');
			localStorage.removeItem('refreshToken');
			authStore.set({
				isAuthenticated: false,
				token: null,
				user: null,
				isAuthenticating: false
			});
		}
	},

	refreshToken: async (): Promise<boolean> => {
		try {
			const refreshToken = localStorage.getItem('refreshToken');
			if (!refreshToken) {
				return false;
			}

			const response = await fetch('/refresh', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ token: refreshToken })
			});

			if (!response.ok) {
				throw new Error('Token refresh failed');
			}

			const data: LoginResponse = await response.json();
			
			// Update tokens
			localStorage.setItem('accessToken', data.accessToken);
			localStorage.setItem('refreshToken', data.refreshToken);
			
			// Update store
			authStore.update(state => ({
				...state,
				token: data.accessToken
			}));

			return true;
		} catch (error) {
			console.error('Token refresh error:', error);
			// If refresh fails, logout
			authActions.logout();
			return false;
		}
	},

	getAuthHeaders: (): Record<string, string> => {
		const token = localStorage.getItem('accessToken');
		return token ? { 'Authorization': `Bearer ${token}` } : {};
	},

	// Handle social authentication tokens
	handleSocialAuth: async (accessToken: string, refreshToken: string): Promise<void> => {
		console.log('=== HANDLE SOCIAL AUTH STARTED ===');
		console.log('Storing tokens in localStorage...');
		
		// Store tokens
		localStorage.setItem('accessToken', accessToken);
		localStorage.setItem('refreshToken', refreshToken);
		
		console.log('Parsing JWT token...');
		console.log('Access token length:', accessToken.length);
		console.log('Access token first 50 chars:', accessToken.substring(0, 50));
		
		// Parse user data from JWT token
		let userData: any = null;
		try {
			// Split the token and get the payload part
			const parts = accessToken.split('.');
			if (parts.length !== 3) {
				throw new Error(`Invalid JWT format: expected 3 parts, got ${parts.length}`);
			}
			
			const payload = parts[1];
			console.log('JWT payload part:', payload);
			
			// Decode base64
			const decodedPayload = atob(payload);
			console.log('Decoded payload:', decodedPayload);
			
			// Parse JSON
			const userPayload = JSON.parse(decodedPayload);
			console.log('Parsed user payload:', userPayload);
			
			userData = {
				id: userPayload.userid || userPayload.sub || userPayload.id,
				email: userPayload.email || '',
				name: userPayload.username || userPayload.name || '',
				avatar: undefined,
				provider: 'social',
				isVerified: true,
				domains: userPayload.domains || []
			};
			console.log('User data parsed successfully:', userData);
		} catch (error) {
			console.error('Error parsing JWT token:', error);
			// Set a default user data if parsing fails
			userData = {
				id: 0,
				email: 'unknown@example.com',
				name: 'Unknown User',
				avatar: undefined,
				provider: 'social',
				isVerified: false,
				domains: []
			};
		}

		console.log('Updating auth store...');
		// Update store
		authStore.set({
			isAuthenticated: true,
			token: accessToken,
			user: userData,
			isAuthenticating: false
		});
		
		console.log('Auth store updated, new state:', {
			isAuthenticated: true,
			token: accessToken ? 'PRESENT' : 'MISSING',
			user: userData ? 'PRESENT' : 'MISSING',
			isAuthenticating: false
		});
		
		// Wait a bit to ensure the store update is processed
		await new Promise(resolve => setTimeout(resolve, 50));
		
		console.log('=== HANDLE SOCIAL AUTH COMPLETED ===');
	}
};
