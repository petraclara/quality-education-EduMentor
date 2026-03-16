const API_BASE = '/api';

async function request(endpoint, options = {}) {
  const config = {
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    ...options,
  };

  const response = await fetch(`${API_BASE}${endpoint}`, config);
  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || 'Something went wrong');
  }

  return data;
}

export const api = {
  // Auth
  register: (body) => request('/register', { method: 'POST', body: JSON.stringify(body) }),
  login: (body) => request('/login', { method: 'POST', body: JSON.stringify(body) }),
  logout: () => request('/logout', { method: 'POST' }),
  me: () => request('/me'),

  // Profile
  getProfile: () => request('/profile'),
  updateProfile: (body) => request('/profile', { method: 'PUT', body: JSON.stringify(body) }),

  // Matches
  findMatches: () => request('/matches/find'),
  getMyMatches: () => request('/matches'),
  acceptMatch: (id) => request(`/matches/${id}/accept`, { method: 'POST' }),
  rejectMatch: (id) => request(`/matches/${id}/reject`, { method: 'POST' }),

  // Dashboard
  getDashboard: () => request('/dashboard'),
};
