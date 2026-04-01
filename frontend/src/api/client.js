const API_BASE = import.meta.env.VITE_API_URL || '/api';

async function request(endpoint, options = {}) {
  const config = {
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    ...options,
  };
  const response = await fetch(`${API_BASE}${endpoint}`, config);
  
  let data = {};
  const text = await response.text();
  if (text) {
    try {
      data = JSON.parse(text);
    } catch (e) {
      console.warn('API returned non-JSON response:', text);
    }
  }

  if (!response.ok) throw new Error(data.message || 'Something went wrong');
  return data;
}

export const api = {
  register: (body) => request('/register', { method: 'POST', body: JSON.stringify(body) }),
  login: (body) => request('/login', { method: 'POST', body: JSON.stringify(body) }),
  logout: () => request('/logout', { method: 'POST' }),
  me: () => request('/me'),

  getProfile: () => request('/profile'),
  updateProfile: (body) => request('/profile', { method: 'PUT', body: JSON.stringify(body) }),

  getMentors: () => request('/mentors'),
  getMentor: (id) => request(`/mentors/${id}`),

  getDashboard: () => request('/dashboard'),

  sendRequest: (body) => request('/requests', { method: 'POST', body: JSON.stringify(body) }),
  getRequests: () => request('/requests'),
  acceptRequest: (id, payload) => request(`/requests/${id}/accept`, { method: 'POST', body: JSON.stringify(payload) }),
  declineRequest: (id, reason) => request(`/requests/${id}/decline`, { method: 'POST', body: JSON.stringify({ decline_reason: reason }) }),
  confirmSlot: (id, payload) => request(`/requests/${id}/confirm`, { method: 'POST', body: JSON.stringify(payload) }),

  createBooking: (body) => request('/bookings', { method: 'POST', body: JSON.stringify(body) }),
  getBookings: () => request('/bookings'),
};
