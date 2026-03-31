import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './Navbar.css';

export default function Navbar() {
  const { user, logout } = useAuth();

  return (
    <nav className="navbar">
      <div className="navbar-inner">
        <Link to={user ? '/dashboard' : '/'} className="navbar-brand">
          🎓 <span>MentorConnect</span>
        </Link>

        <div className="navbar-links">
          {user ? (
            <>
              <Link to="/dashboard" className="nav-link">Dashboard</Link>
              <Link to="/explore" className="nav-link">Explore</Link>
              {user.role === 'learner' && (
                <Link to="/my-requests" className="nav-link">My Requests</Link>
              )}
              <Link to="/profile" className="nav-link">Profile</Link>
              <div className="nav-user">
                <span className="nav-user-name">{user.name}</span>
                <span className="nav-role-badge">{user.role}</span>
                <button onClick={logout} className="btn btn-small btn-secondary">Logout</button>
              </div>
            </>
          ) : (
            <>
              <Link to="/login" className="nav-link">Login</Link>
              <Link to="/register" className="btn btn-primary btn-small">Get Started</Link>
            </>
          )}
        </div>
      </div>
    </nav>
  );
}
