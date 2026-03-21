import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { api } from '../api/client';
import { useToast } from '../components/Toast';
import './RequestMentorship.css';

export default function RequestMentorship() {
  const { mentorId } = useParams();
  const navigate = useNavigate();
  const [mentor, setMentor] = useState(null);
  const [helpWith, setHelpWith] = useState('');
  const [goal, setGoal] = useState('');
  const [message, setMessage] = useState('');
  const [sending, setSending] = useState(false);
  const { showToast, ToastContainer } = useToast();

  useEffect(() => {
    api.getMentor(mentorId).then(res => setMentor(res.data)).catch(() => {});
  }, [mentorId]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!helpWith.trim()) { showToast('Please tell us what you need help with', 'error'); return; }
    setSending(true);
    try {
      await api.sendRequest({ mentor_id: parseInt(mentorId), help_with: helpWith, goal, message });
      showToast('Request sent! 🎉', 'success');
      setTimeout(() => navigate('/my-requests'), 1200);
    } catch (err) { showToast(err.message, 'error'); }
    finally { setSending(false); }
  };

  return (
    <div className="page-container">
      <ToastContainer />
      <div className="rqm-page">
        <Link to={`/mentor/${mentorId}`} className="mp-back">← Back to profile</Link>
        <div className="rqm-card card">
          <h1 className="rqm-title">Request Mentorship</h1>
          {mentor && <p className="rqm-sub">with {mentor.name}</p>}

          <form onSubmit={handleSubmit} className="rqm-form">
            <div className="rqm-field">
              <label>What do you need help with? *</label>
              <input className="form-input" value={helpWith} onChange={(e) => setHelpWith(e.target.value)}
                placeholder='e.g. "Understanding loops in Python"' required />
            </div>
            <div className="rqm-field">
              <label>Your goal <span className="optional">(optional)</span></label>
              <input className="form-input" value={goal} onChange={(e) => setGoal(e.target.value)}
                placeholder='e.g. "Build a small project"' />
            </div>
            <div className="rqm-field">
              <label>Message <span className="optional">(optional)</span></label>
              <textarea className="form-input form-textarea" value={message} onChange={(e) => setMessage(e.target.value)}
                placeholder="Tell the mentor a bit about yourself and what you're hoping to learn..." rows={4} />
            </div>
            <button type="submit" className="btn btn-primary btn-large btn-full" disabled={sending}>
              {sending ? 'Sending...' : 'Send Request'}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
