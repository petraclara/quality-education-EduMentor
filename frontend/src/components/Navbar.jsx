import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './Navbar.css';

export default function Navbar() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/');
  };

  return (
    <nav className="navbar">
      <div className="navbar-container">
        <Link to={user ? '/dashboard' : '/'} className="navbar-brand">
          <span className="navbar-logo">🎓</span>
          <span className="navbar-title">EduMentor</span>
        </Link>

        <div className="navbar-links">
          {user ? (
            <>
              <Link to="/dashboard" className="nav-link">Dashboard</Link>
              <Link to="/matches" className="nav-link">Matches</Link>
              <Link to="/profile" className="nav-link">Profile</Link>
              <div className="nav-user">
                <span className="nav-user-name">{user.name}</span>
                <span className={`nav-role-badge role-${user.role}`}>{user.role}</span>
                <button onClick={handleLogout} className="nav-logout-btn">Logout</button>
              </div>
            </>
          ) : (
            <>
              <Link to="/login" className="nav-link">Login</Link>
              <Link to="/register" className="nav-link nav-link-primary">Get Started</Link>
            </>
          )}
        </div>
      </div>
    </nav>
  );
}
