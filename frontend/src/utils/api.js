import axios from 'axios'

const api = axios.create({
	baseURL: process.env.API_URL || 'http://localhost:8080/api',
	withCredentials: true, // send cookies (sid) for session auth
})

api.interceptors.request.use((config) => {
	const token = localStorage.getItem('token')
	if (token) {
		config.headers = config.headers || {}
		config.headers.Authorization = `Bearer ${token}`
	}
	return config
})

api.interceptors.response.use(
	response => response,
	error => {
		if (error.response?.status === 401) {
			localStorage.removeItem('token')
			window.location.href = '/login'
		}
		return Promise.reject(error)
	}
)

export const authAPI = {
	register: (data) => api.post('/register', data),
	login: (data) => api.post('/login', data),
	getProfile: () => api.get('/me'),
}

export const listingAPI = {
	getMarket: (params) => api.get('/market', { params }),
	getListing: (id) => api.get(`/listings/${id}`),
	createListing: (data) => api.post('/listings', data),
	buy: (data) => api.post('/buy', data),
	bump: (id) => api.post(`/listings/${id}/bump`),
	uploadImage: (id, data) => api.post(`/listings/${id}/image_from_url`, data),
	getImageMeta: (id) => api.get(`/listings/${id}/image/meta`),
}

export default api