import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../api/client';
import { useAuth } from '../context/AuthContext';
import { useToast } from '../components/Toast';
import './Dashboard.css';

const AVATAR_COLORS = ['#14b8a6','#8b5cf6','#f59e0b','#ec4899','#06b6d4','#10b981','#f43f5e','#6366f1'];

function MentorCardSmall({ mentor }) {
  const bg = AVATAR_COLORS[mentor.id % AVATAR_COLORS.length];
  const initials = mentor.name.split(' ').map(n => n[0]).join('').slice(0, 2);
  return (
    <Link to={`/mentor/${mentor.id}`} className="mentor-card-sm">
      <div className="mc-avatar" style={{ background: bg }}>{initials}</div>
      <div className="mc-info">
        <h3 className="mc-name">{mentor.name}</h3>
        <div className="mc-meta">
          {mentor.skills?.slice(0, 3).map((s, i) => (
            <span key={i} className="mc-skill">{s}</span>
          ))}
        </div>
        <p className="mc-bio">{mentor.bio?.slice(0, 80)}{mentor.bio?.length > 80 ? '...' : ''}</p>
      </div>
      <span className="mc-arrow">→</span>
    </Link>
  );
}

function RequestCard({ req, onAccept, onDecline }) {
  const [reason, setReason] = useState('');
  const [showDecline, setShowDecline] = useState(false);

  return (
    <div className={`request-card rq-${req.status}`}>
      <div className="rq-header">
        <div>
          <strong className="rq-name">{req.learner_name}</strong>
          {req.learner_level && <span className="rq-level">{req.learner_level}</span>}
        </div>
        <span className={`rq-status rq-status-${req.status}`}>
          {req.status === 'pending' && '⏳ Pending'}
          {req.status === 'accepted' && '✅ Accepted'}
          {req.status === 'declined' && '❌ Declined'}
        </span>
      </div>
      <p className="rq-help"><strong>Needs help with:</strong> {req.help_with}</p>
      {req.goal && <p className="rq-goal"><strong>Goal:</strong> {req.goal}</p>}
      {req.message && <p className="rq-msg">"{req.message}"</p>}

      {req.status === 'pending' && (
        <div className="rq-actions">
          <button className="btn btn-primary btn-small" onClick={() => onAccept(req.id)}>Accept</button>
          {showDecline ? (
            <div className="decline-flow">
              <select className="decline-select" value={reason} onChange={(e) => setReason(e.target.value)}>
                <option value="">Select reason...</option>
                <option value="Busy">Busy</option>
                <option value="Not my area">Not my area</option>
                <option value="Not a good fit">Not a good fit</option>
              </select>
              <button className="btn btn-danger btn-small" disabled={!reason}
                onClick={() => { onDecline(req.id, reason); setShowDecline(false); }}>
                Confirm Decline
              </button>
            </div>
          ) : (
            <button className="btn btn-secondary btn-small" onClick={() => setShowDecline(true)}>Decline</button>
          )}
        </div>
      )}
    </div>
  );
}

export default function Dashboard() {
  const { user } = useAuth();
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const { showToast, ToastContainer } = useToast();

  useEffect(() => { fetchDashboard(); }, []);

  const fetchDashboard = async () => {
    try {
      const res = await api.getDashboard();
      setData(res.data);
    } catch (err) {
      showToast(err.message, 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleAccept = async (id) => {
    try {
      await api.acceptRequest(id);
      showToast('Request accepted! 🎉', 'success');
      fetchDashboard();
    } catch (err) { showToast(err.message, 'error'); }
  };

  const handleDecline = async (id, reason) => {
    try {
      await api.declineRequest(id, reason);
      showToast('Request declined', 'info');
      fetchDashboard();
    } catch (err) { showToast(err.message, 'error'); }
  };

  if (loading) {
    return <div className="page-container"><div className="loading-state"><div className="loading-spinner"></div><p>Loading...</p></div></div>;
  }

  // LEARNER DASHBOARD
  if (data?.role === 'learner') {
    const mentors = data.mentors || [];
    return (
      <div className="page-container">
        <ToastContainer />
        <div className="dash">
          <div className="dash-welcome">
            <h1>Hi {data.user?.name?.split(' ')[0]} 👋</h1>
            <p className="dash-sub">Mentors matched for you</p>
          </div>

          {mentors.length > 0 ? (
            <div className="mentor-list">
              {mentors.map(m => <MentorCardSmall key={m.id} mentor={m} />)}
            </div>
          ) : (
            <div className="empty-state">
              <span className="empty-icon">😕</span>
              <p className="empty-text">No mentors matched your preferences yet</p>
              <Link to="/explore" className="btn btn-primary btn-small">Explore All Mentors</Link>
            </div>
          )}

          <div className="dash-explore-link">
            <Link to="/explore" className="btn btn-secondary">Explore All Mentors</Link>
          </div>
        </div>
      </div>
    );
  }

  // MENTOR DASHBOARD
  const requests = data?.requests || [];
  const pending = requests.filter(r => r.status === 'pending');
  return (
    <div className="page-container">
      <ToastContainer />
      <div className="dash">
        <div className="dash-welcome">
          <h1>Welcome, {data?.user?.name?.split(' ')[0]} 🎓</h1>
          <p className="dash-sub">Your mentor dashboard</p>
        </div>

        <div className="dash-stats">
          <div className="stat-card">
            <span className="stat-val">{pending.length}</span>
            <span className="stat-label">New Requests</span>
          </div>
          <div className="stat-card">
            <span className="stat-val">{requests.filter(r => r.status === 'accepted').length}</span>
            <span className="stat-label">Accepted</span>
          </div>
          <div className="stat-card">
            <span className="stat-val">{requests.length}</span>
            <span className="stat-label">Total</span>
          </div>
        </div>

        {pending.length > 0 && <h2 className="section-title">New Requests ({pending.length})</h2>}

        <div className="requests-list">
          {requests.length > 0 ? (
            requests.map(r => (
              <RequestCard key={r.id} req={r} onAccept={handleAccept} onDecline={handleDecline} />
            ))
          ) : (
            <div className="empty-state">
              <span className="empty-icon">📭</span>
              <p className="empty-text">No requests yet. Learners will find you based on your profile.</p>
              <Link to="/profile" className="btn btn-primary btn-small">Update Profile</Link>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
