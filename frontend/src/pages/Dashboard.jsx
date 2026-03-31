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
  const [showAccept, setShowAccept] = useState(false);
  const [meetingType, setMeetingType] = useState('Google Meet');
  const [meetingLink, setMeetingLink] = useState('');
  const [slots, setSlots] = useState([{ date: '', time: '' }]);

  const handleAddSlot = () => setSlots([...slots, { date: '', time: '' }]);
  const handleSlotChange = (index, field, value) => {
    const newSlots = [...slots];
    newSlots[index][field] = value;
    setSlots(newSlots);
  };
  const handleRemoveSlot = (index) => setSlots(slots.filter((_, i) => i !== index));

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
          {req.status === 'scheduled' && '📅 Scheduled'}
          {req.status === 'declined' && '❌ Declined'}
        </span>
      </div>
      <p className="rq-help"><strong>Needs help with:</strong> {req.help_with}</p>
      {req.goal && <p className="rq-goal"><strong>Goal:</strong> {req.goal}</p>}
      {req.message && <p className="rq-msg">"{req.message}"</p>}

      {req.status === 'scheduled' && req.proposed_slots?.length > 0 && (
        <div className="myr-scheduled" style={{ marginTop: '15px', padding: '15px', background: '#ecfdf5', border: '1px solid #a7f3d0', borderRadius: '8px' }}>
          <h4 style={{ margin: '0 0 10px', color: '#059669' }}>Meeting Confirmed 🎉</h4>
          <div style={{ display: 'grid', gridTemplateColumns: '80px 1fr', gap: '5px', fontSize: '14px' }}>
            <strong>Date:</strong> <span>{req.proposed_slots[0].date}</span>
            <strong>Time:</strong> <span>{req.proposed_slots[0].time}</span>
            <strong>Type:</strong> <span>{req.meeting_type}</span>
            <strong>Link:</strong> <a href={req.meeting_link} target="_blank" rel="noreferrer" style={{ color: '#2563eb' }}>{req.meeting_link}</a>
          </div>
        </div>
      )}

      {req.status === 'pending' && (
        <div className="rq-actions" style={{ flexDirection: 'column', alignItems: 'flex-start' }}>
          {!showDecline && !showAccept && (
            <div style={{ display: 'flex', gap: '10px' }}>
              <button className="btn btn-primary btn-small" onClick={() => setShowAccept(true)}>Accept</button>
              <button className="btn btn-secondary btn-small" onClick={() => setShowDecline(true)}>Decline</button>
            </div>
          )}
          
          {showAccept && (
            <div className="accept-flow" style={{ width: '100%', marginTop: '10px', background: '#f8fafc', padding: '15px', borderRadius: '8px' }}>
              <h4 style={{ margin: '0 0 10px' }}>Set Available Slots</h4>
              {slots.map((slot, idx) => (
                <div key={idx} style={{ display: 'flex', gap: '10px', marginBottom: '10px' }}>
                  <input className="form-input" type="date" value={slot.date} onChange={e => handleSlotChange(idx, 'date', e.target.value)} required style={{ flex: 1, margin: 0 }} />
                  <input className="form-input" type="text" placeholder="e.g. 2:00 PM - 3:00 PM" value={slot.time} onChange={e => handleSlotChange(idx, 'time', e.target.value)} required style={{ flex: 1, margin: 0 }} />
                  {slots.length > 1 && <button type="button" onClick={() => handleRemoveSlot(idx)} style={{ background: 'none', border: 'none', cursor: 'pointer', fontSize: '18px' }}>❌</button>}
                </div>
              ))}
              <button type="button" className="btn btn-secondary btn-small" onClick={handleAddSlot}>+ Add Slot</button>
              
              <div style={{ marginTop: '15px' }}>
                <label style={{ fontWeight: '500', marginRight: '10px' }}>Meeting Type:</label>
                <select className="form-input" value={meetingType} onChange={e => setMeetingType(e.target.value)} style={{ width: 'auto', display: 'inline-block', padding: '5px' }}>
                  <option>Google Meet</option>
                  <option>Zoom</option>
                  <option>In-person</option>
                </select>
              </div>
              
              <div style={{ marginTop: '10px' }}>
                <label style={{ fontWeight: '500' }}>Meeting Link / Location:</label>
                <input className="form-input" type="text" placeholder="https://meet.google.com/xyz" value={meetingLink} onChange={e => setMeetingLink(e.target.value)} required style={{ width: '100%', marginTop: '5px' }} />
              </div>

              <div style={{ marginTop: '15px', display: 'flex', gap: '10px' }}>
                <button className="btn btn-primary btn-small" 
                  disabled={!meetingLink || slots.some(s => !s.date || !s.time)}
                  onClick={() => onAccept(req.id, { meeting_type: meetingType, meeting_link: meetingLink, proposed_slots: slots })}>
                  Send Invitation
                </button>
                <button className="btn btn-secondary btn-small" onClick={() => setShowAccept(false)}>Cancel</button>
              </div>
            </div>
          )}

          {showDecline && (
            <div className="decline-flow" style={{ marginTop: '10px' }}>
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
              <button className="btn btn-secondary btn-small" style={{ marginLeft: '10px' }} onClick={() => setShowDecline(false)}>Cancel</button>
            </div>
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

  const handleAccept = async (id, payload) => {
    try {
      await api.acceptRequest(id, payload);
      showToast('Invitation sent! 🎉', 'success');
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
            <span className="stat-val">{requests.filter(r => r.status === 'accepted' || r.status === 'scheduled').length}</span>
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
