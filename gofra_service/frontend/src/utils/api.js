import axios from 'axios'

const API_URL = window.__API_URL__ || (
		(typeof process !== 'undefined' && process && process['env'] && process['env']['API_URL'])
			? process['env']['API_URL']
			: 'http://localhost:8080/api'
)

const api = axios.create({
	baseURL: API_URL,
	withCredentials: true, 
})

export { API_URL }

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
	getMyListings: () => api.get('/my-listings'),
	getMyGofers: () => api.get('/my-gofers'),
	createListing: (data) => api.post('/listings', data),
	buy: (data) => api.post('/buy', data),
	bump: (id) => api.post(`/listings/${id}/bump`),
	uploadImageFromUrl: (id, data) => api.post(`/listings/${id}/image_from_url`, data),
	uploadImageFile: (id, formData) => api.post(`/listings/${id}/image_upload`, formData, {
		headers: { 'Content-Type': 'multipart/form-data' }
	}),
	getImageMeta: (id) => api.get(`/listings/${id}/image/meta`),
}

export const statsAPI = {
	getActiveUsers: () => api.get('/stats/active-users'),
}

export default api