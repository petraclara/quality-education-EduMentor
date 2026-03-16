import { useState, useEffect } from 'react';
import { api } from '../api/client';
import TagInput from '../components/TagInput';
import AvailabilityGrid from '../components/AvailabilityGrid';
import { useToast } from '../components/Toast';
import './Profile.css';

export default function Profile() {
  const [bio, setBio] = useState('');
  const [skills, setSkills] = useState([]);
  const [interests, setInterests] = useState([]);
  const [availability, setAvailability] = useState([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const { showToast, ToastContainer } = useToast();

  useEffect(() => {
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    try {
      const res = await api.getProfile();
      const profile = res.data.profile;
      setBio(profile.bio || '');
      setSkills(profile.skills || []);
      setInterests(profile.interests || []);
      setAvailability(profile.availability || []);
    } catch (err) {
      showToast(err.message, 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async (e) => {
    e.preventDefault();
    setSaving(true);
    try {
      await api.updateProfile({ bio, skills, interests, availability });
      showToast('Profile updated successfully!', 'success');
    } catch (err) {
      showToast(err.message, 'error');
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="page-container">
        <div className="loading-state">
          <div className="loading-spinner"></div>
          <p>Loading profile...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="page-container">
      <ToastContainer />

      <div className="profile-page">
        <div className="profile-header">
          <h1 className="page-title">Edit Profile</h1>
          <p className="page-subtitle">
            Add your skills, interests, and availability to improve your match quality
          </p>
        </div>

        <form onSubmit={handleSave} className="profile-form">
          {/* Bio */}
          <div className="profile-section">
            <h2 className="section-label">About You</h2>
            <textarea
              id="profile-bio"
              className="form-input form-textarea"
              value={bio}
              onChange={(e) => setBio(e.target.value)}
              placeholder="Tell potential mentors/mentees about yourself, your background, and what you're looking to achieve..."
              rows={4}
            />
          </div>

          {/* Skills */}
          <div className="profile-section">
            <h2 className="section-label">Skills</h2>
            <p className="section-hint">Add your technical skills and areas of expertise (press Enter to add)</p>
            <TagInput
              tags={skills}
              onChange={setSkills}
              placeholder="e.g. JavaScript, Python, Machine Learning..."
            />
          </div>

          {/* Interests */}
          <div className="profile-section">
            <h2 className="section-label">Interests</h2>
            <p className="section-hint">What topics are you interested in learning or teaching?</p>
            <TagInput
              tags={interests}
              onChange={setInterests}
              placeholder="e.g. Web Development, Data Science, UI Design..."
            />
          </div>

          {/* Availability */}
          <div className="profile-section">
            <h2 className="section-label">Availability</h2>
            <p className="section-hint">Click time slots when you're available for mentoring sessions</p>
            <AvailabilityGrid
              availability={availability}
              onChange={setAvailability}
            />
          </div>

          {/* Save Button */}
          <button
            id="profile-save"
            type="submit"
            className="btn btn-primary btn-large btn-full"
            disabled={saving}
          >
            {saving ? 'Saving...' : '💾 Save Profile'}
          </button>
        </form>
      </div>
    </div>
  );
}
