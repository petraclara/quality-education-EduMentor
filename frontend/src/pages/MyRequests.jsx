import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../api/client';
import { useToast } from '../components/Toast';
import './MyRequests.css';

export default function MyRequests() {
  const [requests, setRequests] = useState([]);
  const [loading, setLoading] = useState(true);
  const { showToast, ToastContainer } = useToast();

  const fetchRequests = () => {
    api.getRequests().then(res => { setRequests(res.data || []); setLoading(false); })
      .catch(err => { showToast(err.message, 'error'); setLoading(false); });
  };

  useEffect(() => {
    fetchRequests();
  }, []);

  const handleConfirm = async (id, slot) => {
    try {
      await api.confirmSlot(id, { date: slot.date, time: slot.time });
      showToast('Meeting confirmed! 🎉', 'success');
      fetchRequests();
    } catch (err) {
      showToast(err.message, 'error');
    }
  };

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
                    {r.status === 'accepted' && '⏳ Action Required'}
                    {r.status === 'scheduled' && '📅 Scheduled'}
                    {r.status === 'declined' && '❌ Declined'}
                  </span>
                </div>
                <p className="myr-help">{r.help_with}</p>
                {r.goal && <p className="myr-meta">Goal: {r.goal}</p>}
                
                {r.status === 'accepted' && r.proposed_slots?.length > 0 && (
                  <div className="myr-slots" style={{ marginTop: '15px', padding: '15px', background: '#f8fafc', borderRadius: '8px' }}>
                    <h4 style={{ margin: '0 0 10px', color: '#3b82f6' }}>Mentor Proposed Slots</h4>
                    <p style={{ fontSize: '14px', marginBottom: '10px' }}>Please choose a time that works for you:</p>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                      {r.proposed_slots.map((slot, i) => (
                        <div key={i} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', background: 'white', padding: '10px', borderRadius: '6px', border: '1px solid #e2e8f0' }}>
                          <span><strong>{slot.date}</strong> at {slot.time}</span>
                          <button className="btn btn-primary btn-small" onClick={() => handleConfirm(r.id, slot)}>Select</button>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {r.status === 'scheduled' && r.proposed_slots?.length > 0 && (
                  <div className="myr-scheduled" style={{ marginTop: '15px', padding: '15px', background: '#ecfdf5', border: '1px solid #a7f3d0', borderRadius: '8px' }}>
                    <h4 style={{ margin: '0 0 10px', color: '#059669' }}>Meeting Confirmed 🎉</h4>
                    <div style={{ display: 'grid', gridTemplateColumns: '80px 1fr', gap: '5px', fontSize: '14px' }}>
                      <strong>Date:</strong> <span>{r.proposed_slots[0].date}</span>
                      <strong>Time:</strong> <span>{r.proposed_slots[0].time}</span>
                      <strong>Type:</strong> <span>{r.meeting_type}</span>
                      <strong>Link:</strong> <a href={r.meeting_link} target="_blank" rel="noreferrer" style={{ color: '#2563eb' }}>{r.meeting_link}</a>
                    </div>
                  </div>
                )}

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
