import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useToast } from '../components/Toast';
import './Auth.css';

export default function Register() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState('learner');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();
  const { showToast, ToastContainer } = useToast();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      await login({ name, email, password, role }, true);
      navigate('/preferences');
    } catch (err) {
      showToast(err.message, 'error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      <ToastContainer />
      <div className="auth-card">
        <h1 className="auth-title">Create Account</h1>
        <p className="auth-sub">Join MentorConnect and start your journey</p>

        <form onSubmit={handleSubmit} className="auth-form">
          <input id="reg-name" className="form-input" type="text" placeholder="Full Name"
            value={name} onChange={(e) => setName(e.target.value)} required />
          <input id="reg-email" className="form-input" type="email" placeholder="Email"
            value={email} onChange={(e) => setEmail(e.target.value)} required />
          <input id="reg-password" className="form-input" type="password" placeholder="Password (min 6 chars)"
            value={password} onChange={(e) => setPassword(e.target.value)} required minLength={6} />

          <div className="role-selection">
            <label className="role-label">I want to be a:</label>
            <div className="role-options">
              <button type="button" className={`role-btn ${role === 'learner' ? 'role-btn-active' : ''}`}
                onClick={() => setRole('learner')}>
                📚 Learner
              </button>
              <button type="button" className={`role-btn ${role === 'mentor' ? 'role-btn-active' : ''}`}
                onClick={() => setRole('mentor')}>
                🎓 Mentor
              </button>
            </div>
          </div>

          <button id="reg-submit" type="submit" className="btn btn-primary btn-large btn-full" disabled={loading}>
            {loading ? 'Creating account...' : 'Continue'}
          </button>
        </form>

        <p className="auth-footer">
          Already have an account? <Link to="/login">Login</Link>
        </p>
      </div>
    </div>
  );
}
