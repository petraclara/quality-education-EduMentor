import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../api/client';
import { useToast } from '../components/Toast';
import './MyRequests.css';

export default function MyRequests() {
  const [requests, setRequests] = useState([]);
  const [loading, setLoading] = useState(true);
  const { showToast, ToastContainer } = useToast();

  useEffect(() => {
    api.getRequests().then(res => { setRequests(res.data || []); setLoading(false); })
      .catch(err => { showToast(err.message, 'error'); setLoading(false); });
  }, []);

  if (loading) return <div className="page-container"><div className="loading-state"><div className="loading-spinner"></div><p>Loading...</p></div></div>;

  return (
    <div className="page-container">
      <ToastContainer />
      <div className="myr">
        <h1 className="myr-title">Your Requests</h1>
        {requests.length > 0 ? (
          <div className="myr-list">
            {requests.map(r => (
              <div key={r.id} className="myr-card card">
                <div className="myr-header">
                  <strong>{r.mentor_name}</strong>
                  <span className={`myr-status myr-${r.status}`}>
                    {r.status === 'pending' && '⏳ Pending'}
                    {r.status === 'accepted' && '✅ Accepted'}
                    {r.status === 'declined' && '❌ Declined'}
                  </span>
                </div>
                <p className="myr-help">{r.help_with}</p>
                {r.goal && <p className="myr-meta">Goal: {r.goal}</p>}
                {r.status === 'declined' && r.decline_reason && (
                  <div className="myr-reason">
                    <strong>Reason:</strong> {r.decline_reason}
                  </div>
                )}
                {r.status === 'declined' && (
                  <Link to="/explore" className="btn btn-secondary btn-small" style={{ marginTop: 10 }}>
                    Find similar mentors
                  </Link>
                )}
              </div>
            ))}
          </div>
        ) : (
          <div className="empty-state">
            <span className="empty-icon">📭</span>
            <p className="empty-text">You haven't sent any requests yet</p>
            <Link to="/explore" className="btn btn-primary btn-small">Explore Mentors</Link>
          </div>
        )}
      </div>
    </div>
  );
}
