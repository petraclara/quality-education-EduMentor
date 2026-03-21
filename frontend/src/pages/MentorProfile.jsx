import { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { api } from '../api/client';
import { useAuth } from '../context/AuthContext';
import { useToast } from '../components/Toast';
import './MentorProfile.css';

const AVATAR_COLORS = ['#14b8a6','#8b5cf6','#f59e0b','#ec4899','#06b6d4','#10b981','#f43f5e','#6366f1'];
const DAYS = ['Monday','Tuesday','Wednesday','Thursday','Friday','Saturday','Sunday'];
const TIMES = ['Morning','Afternoon','Evening'];

export default function MentorProfile() {
  const { id } = useParams();
  const { user } = useAuth();
  const navigate = useNavigate();
  const [mentor, setMentor] = useState(null);
  const [loading, setLoading] = useState(true);
  const { showToast, ToastContainer } = useToast();

  useEffect(() => { fetchMentor(); }, [id]);

  const fetchMentor = async () => {
    try {
      const res = await api.getMentor(id);
      setMentor(res.data);
    } catch (err) { showToast(err.message, 'error'); }
    finally { setLoading(false); }
  };

  if (loading) return <div className="page-container"><div className="loading-state"><div className="loading-spinner"></div><p>Loading...</p></div></div>;
  if (!mentor) return <div className="page-container"><div className="empty-state"><span className="empty-icon">😕</span><p className="empty-text">Mentor not found</p><Link to="/explore" className="btn btn-primary btn-small">Back</Link></div></div>;

  const bg = AVATAR_COLORS[mentor.id % AVATAR_COLORS.length];
  const initials = mentor.name.split(' ').map(n => n[0]).join('').slice(0, 2);
  const availGrid = {};
  DAYS.forEach(d => TIMES.forEach(t => {
    const key = `${d.toLowerCase()}-${t.toLowerCase()}`;
    availGrid[key] = mentor.availability?.includes(key);
  }));

  return (
    <div className="page-container">
      <ToastContainer />
      <div className="mp">
        <Link to="/explore" className="mp-back">← Back to Mentors</Link>

        <div className="mp-hero card">
          <div className="mp-avatar" style={{ background: bg }}>{initials}</div>
          <div className="mp-info">
            <h1 className="mp-name">{mentor.name}</h1>
            {mentor.level && <span className="mp-level">{mentor.level} level</span>}
            {mentor.rating > 0 && (
              <div className="mp-rating">
                <span className="mp-stars">{'★'.repeat(Math.round(mentor.rating))}</span>
                <span>{mentor.rating.toFixed(1)} ({mentor.rating_count} reviews)</span>
              </div>
            )}
          </div>
          {user && user.role === 'learner' && (
            <button className="btn btn-primary btn-large" onClick={() => navigate(`/request/${mentor.id}`)}>
              Request Mentorship
            </button>
          )}
        </div>

        <div className="mp-body">
          <div className="mp-section card">
            <h2>About</h2>
            <p className="mp-bio">{mentor.bio || 'No bio provided.'}</p>
          </div>

          {mentor.skills?.length > 0 && (
            <div className="mp-section card">
              <h2>Skills</h2>
              <div className="mp-tags">
                {mentor.skills.map((s, i) => <span key={i} className="mp-tag">{s}</span>)}
              </div>
            </div>
          )}

          {mentor.interests?.length > 0 && (
            <div className="mp-section card">
              <h2>Interests</h2>
              <div className="mp-tags">
                {mentor.interests.map((s, i) => <span key={i} className="mp-tag mp-tag-purple">{s}</span>)}
              </div>
            </div>
          )}

          {mentor.availability?.length > 0 && (
            <div className="mp-section card">
              <h2>Availability</h2>
              <div className="avail-grid">
                <div className="avail-header"><div className="avail-corner"></div>
                  {TIMES.map(t => <div key={t} className="avail-th">{t}</div>)}
                </div>
                {DAYS.map(d => (
                  <div key={d} className="avail-row">
                    <div className="avail-day">{d.slice(0,3)}</div>
                    {TIMES.map(t => {
                      const key = `${d.toLowerCase()}-${t.toLowerCase()}`;
                      return <div key={key} className={`avail-cell ${availGrid[key] ? 'avail-yes' : 'avail-no'}`}>
                        {availGrid[key] ? '✓' : '—'}
                      </div>;
                    })}
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
