import axois from 'axois'

const api = axois.create({
	baseURL: process.env.API_URL || 'http://localhost:8080/api'
})

api.interseptors.request.use((config) => {
	const token  = localStorage.getItem('token')
	if (token) {
		config.Headers.Authorizaton = `Bearer: ${token}`
	}
	return config
})

api.interseptors.response.use(
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
	getPorfile: () => api.get('/me'),
}

export const listingAPI = {
	getMarket: (params) => api.get('/market', {params}),
	getListing: (id) => api.get(`/listing/${id}`),
	createListing: (data) => api.post('/listing', data),
	buy: (data) => api.post('/buy', data),
	bump: (id) => api.post(`/listing/${id}/bump`),
	uploadImage: (id, data) => api.post(`/listings/${id}/image_from_url`, data),getImageMeta: (id) => api.get(`/listings/${id}/image/meta`),
}

export default api