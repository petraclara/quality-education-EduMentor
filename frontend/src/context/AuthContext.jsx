import { createContext, useContext, useState, useEffect } from 'react';
import { api } from '../api/client';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [preferencesSet, setPreferencesSet] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => { checkAuth(); }, []);

  const checkAuth = async () => {
    try {
      const res = await api.me();
      // Backend now returns { user: {...}, preferences_set: bool }
      setUser(res.data.user);
      setPreferencesSet(res.data.preferences_set);
    } catch {
      setUser(null);
      setPreferencesSet(false);
    } finally {
      setLoading(false);
    }
  };

  // login handles both login and register
  const login = async (credentials, isRegister = false) => {
    const res = isRegister
      ? await api.register(credentials)
      : await api.login(credentials);
    // After login/register, fetch full user state including preferences flag
    await checkAuth();
    return res;
  };

  const logout = async () => {
    await api.logout();
    setUser(null);
    setPreferencesSet(false);
  };

  // Called after preferences are saved so the gate updates immediately
  const refreshPreferences = async () => {
    try {
      const res = await api.me();
      setPreferencesSet(res.data.preferences_set);
    } catch {
      // ignore
    }
  };

  return (
    <AuthContext.Provider value={{ user, preferencesSet, loading, login, logout, checkAuth, refreshPreferences }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within an AuthProvider');
  return context;
}
