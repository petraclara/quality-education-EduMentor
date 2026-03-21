import { useState, useEffect } from 'react';
import { api } from '../api/client';
import { useAuth } from '../context/AuthContext';
import TagInput from '../components/TagInput';
import AvailabilityGrid from '../components/AvailabilityGrid';
import { useToast } from '../components/Toast';
import './Profile.css';

export default function Profile() {
  const { user } = useAuth();
  const [bio, setBio] = useState('');
  const [skills, setSkills] = useState([]);
  const [interests, setInterests] = useState([]);
  const [level, setLevel] = useState('');
  const [goal, setGoal] = useState('');
  const [availability, setAvailability] = useState([]);
  const [maxMentees, setMaxMentees] = useState(5);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const { showToast, ToastContainer } = useToast();

  useEffect(() => { fetchProfile(); }, []);

  const fetchProfile = async () => {
    try {
      const res = await api.getProfile();
      const p = res.data.profile;
      setBio(p.bio || '');
      setSkills(p.skills || []);
      setInterests(p.interests || []);
      setLevel(p.level || '');
      setGoal(p.goal || '');
      setAvailability(p.availability || []);
      setMaxMentees(p.max_mentees || 5);
    } catch (err) { showToast(err.message, 'error'); }
    finally { setLoading(false); }
  };

  const handleSave = async (e) => {
    e.preventDefault();
    setSaving(true);
    try {
      await api.updateProfile({ bio, skills, interests, level, goal, availability, max_mentees: maxMentees });
      showToast('Profile updated!', 'success');
    } catch (err) { showToast(err.message, 'error'); }
    finally { setSaving(false); }
  };

  if (loading) return <div className="page-container"><div className="loading-state"><div className="loading-spinner"></div><p>Loading...</p></div></div>;

  return (
    <div className="page-container">
      <ToastContainer />
      <div className="profile-page">
        <h1 className="prof-title">Edit Profile</h1>
        <p className="prof-sub">Update your profile to improve your matches</p>

        <form onSubmit={handleSave} className="prof-form">
          <div className="prof-section card">
            <h2>About You</h2>
            <textarea className="form-input form-textarea" value={bio} onChange={(e) => setBio(e.target.value)}
              placeholder="Tell others about yourself..." rows={4} />
          </div>

          <div className="prof-section card">
            <h2>Level</h2>
            <div className="level-options">
              {['beginner','intermediate','advanced'].map(l => (
                <button key={l} type="button" className={`level-btn ${level === l ? 'level-btn-active' : ''}`}
                  onClick={() => setLevel(l)}>
                  {l.charAt(0).toUpperCase() + l.slice(1)}
                </button>
              ))}
            </div>
          </div>

          <div className="prof-section card">
            <h2>Goal</h2>
            <input className="form-input" value={goal} onChange={(e) => setGoal(e.target.value)}
              placeholder='e.g. "Build real projects", "Pass exams"' />
          </div>

          <div className="prof-section card">
            <h2>{user?.role === 'mentor' ? 'Skills You Teach' : 'Skills'}</h2>
            <TagInput tags={skills} onChange={setSkills} placeholder="Type a skill and press Enter..." />
          </div>

          <div className="prof-section card">
            <h2>Interests</h2>
            <TagInput tags={interests} onChange={setInterests} placeholder="Type an interest and press Enter..." />
          </div>

          {user?.role === 'mentor' && (
            <div className="prof-section card">
              <h2>Availability</h2>
              <p className="prof-hint">Click slots when you're available</p>
              <AvailabilityGrid availability={availability} onChange={setAvailability} />
            </div>
          )}

          {user?.role === 'mentor' && (
            <div className="prof-section card">
              <h2>Max Mentees</h2>
              <p className="prof-hint">How many learners can you mentor at once?</p>
              <input className="form-input" type="number" min="1" max="20" value={maxMentees}
                onChange={(e) => setMaxMentees(parseInt(e.target.value) || 1)} style={{maxWidth: 120}} />
            </div>
          )}

          <button type="submit" className="btn btn-primary btn-large btn-full" disabled={saving}>
            {saving ? 'Saving...' : 'Save Profile'}
          </button>
        </form>
      </div>
    </div>
  );
}
