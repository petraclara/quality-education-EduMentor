import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api/client';
import { useAuth } from '../context/AuthContext';
import { useToast } from '../components/Toast';
import './Preferences.css';

const INTEREST_OPTIONS = [
  'Python', 'JavaScript', 'Go', 'React', 'Data Science', 'Machine Learning',
  'Web Development', 'Mobile Apps', 'UI Design', 'DevOps', 'Cybersecurity',
  'Algorithms', 'Cloud Computing', 'Databases', 'API Design',
];

export default function Preferences() {
  const { user } = useAuth();
  const [interests, setInterests] = useState([]);
  const [level, setLevel] = useState('');
  const [goal, setGoal] = useState('');
  const [saving, setSaving] = useState(false);
  const navigate = useNavigate();
  const { showToast, ToastContainer } = useToast();

  const toggleInterest = (item) => {
    setInterests(prev =>
      prev.includes(item) ? prev.filter(i => i !== item) : [...prev, item]
    );
  };

  const handleSave = async (e) => {
    e.preventDefault();
    if (interests.length === 0) {
      showToast('Please select at least one interest', 'error');
      return;
    }
    setSaving(true);
    try {
      const profileData = {
        bio: '',
        skills: user?.role === 'mentor' ? interests : [],
        interests: user?.role === 'learner' ? interests : [],
        level,
        goal,
        availability: [],
        max_mentees: 5,
      };
      await api.updateProfile(profileData);
      showToast('Preferences saved!', 'success');
      navigate('/dashboard');
    } catch (err) {
      showToast(err.message, 'error');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="pref-page">
      <ToastContainer />
      <div className="pref-card">
        <div className="pref-header">
          <span className="pref-emoji">👋</span>
          <h1 className="pref-title">Tell us about you</h1>
          <p className="pref-sub">
            {user?.role === 'learner'
              ? "Let's find the best mentor for you"
              : "Help learners find you"}
          </p>
        </div>

        <form onSubmit={handleSave} className="pref-form">
          <div className="pref-section">
            <label className="pref-label">
              {user?.role === 'learner' ? 'What do you want to learn?' : 'What can you teach?'}
            </label>
            <div className="interest-grid">
              {INTEREST_OPTIONS.map(item => (
                <button
                  key={item}
                  type="button"
                  className={`interest-btn ${interests.includes(item) ? 'interest-btn-active' : ''}`}
                  onClick={() => toggleInterest(item)}
                >
                  {item}
                </button>
              ))}
            </div>
          </div>

          <div className="pref-section">
            <label className="pref-label">Your level</label>
            <div className="level-options">
              {['beginner', 'intermediate', 'advanced'].map(l => (
                <button
                  key={l}
                  type="button"
                  className={`level-btn ${level === l ? 'level-btn-active' : ''}`}
                  onClick={() => setLevel(l)}
                >
                  {l === 'beginner' && '🌱 '}
                  {l === 'intermediate' && '🌿 '}
                  {l === 'advanced' && '🌳 '}
                  {l.charAt(0).toUpperCase() + l.slice(1)}
                </button>
              ))}
            </div>
          </div>

          <div className="pref-section">
            <label className="pref-label">Your goal <span className="optional">(optional)</span></label>
            <input
              className="form-input"
              type="text"
              placeholder='e.g. "Pass exam", "Build projects", "Learn basics"'
              value={goal}
              onChange={(e) => setGoal(e.target.value)}
            />
          </div>

          <button type="submit" className="btn btn-primary btn-large btn-full" disabled={saving}>
            {saving ? 'Saving...' : 'Save & Continue'}
          </button>
        </form>
      </div>
    </div>
  );
}
