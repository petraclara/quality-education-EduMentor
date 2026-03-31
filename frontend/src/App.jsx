import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import Navbar from './components/Navbar';
import Landing from './pages/Landing';
import Login from './pages/Login';
import Register from './pages/Register';
import Preferences from './pages/Preferences';
import Dashboard from './pages/Dashboard';
import Profile from './pages/Profile';
import MentorProfile from './pages/MentorProfile';
import RequestMentorship from './pages/RequestMentorship';
import MyRequests from './pages/MyRequests';
import ExploreMentors from './pages/ExploreMentors';
import './App.css';

function ProtectedRoute({ children, requirePreferences = true }) {
  const { user, preferencesSet, loading } = useAuth();
  if (loading) return <div className="page-container"><div className="loading-state"><div className="loading-spinner"></div></div></div>;
  if (!user) return <Navigate to="/login" />;
  // If preferences are required and not set, redirect to preferences page
  if (requirePreferences && !preferencesSet) return <Navigate to="/preferences" />;
  return children;
}

function GuestRoute({ children }) {
  const { user, loading } = useAuth();
  if (loading) return null;
  return user ? <Navigate to="/dashboard" /> : children;
}

function AppRoutes() {
  return (
    <>
      <Navbar />
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/login" element={<GuestRoute><Login /></GuestRoute>} />
        <Route path="/register" element={<GuestRoute><Register /></GuestRoute>} />
        {/* Preferences page: protected but does NOT require preferences (avoids redirect loop) */}
        <Route path="/preferences" element={<ProtectedRoute requirePreferences={false}><Preferences /></ProtectedRoute>} />
        <Route path="/dashboard" element={<ProtectedRoute><Dashboard /></ProtectedRoute>} />
        <Route path="/profile" element={<ProtectedRoute><Profile /></ProtectedRoute>} />
        <Route path="/explore" element={<ProtectedRoute><ExploreMentors /></ProtectedRoute>} />
        <Route path="/mentor/:id" element={<ProtectedRoute><MentorProfile /></ProtectedRoute>} />
        <Route path="/request/:mentorId" element={<ProtectedRoute><RequestMentorship /></ProtectedRoute>} />
        <Route path="/my-requests" element={<ProtectedRoute><MyRequests /></ProtectedRoute>} />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </BrowserRouter>
  );
}
