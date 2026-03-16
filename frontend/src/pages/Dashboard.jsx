import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../api/client';
import { useAuth } from '../context/AuthContext';
import MatchCard from '../components/MatchCard';
import { useToast } from '../components/Toast';
import './Dashboard.css';

export default function Dashboard() {
  const { user } = useAuth();
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const { showToast, ToastContainer } = useToast();

  useEffect(() => {
    fetchDashboard();
  }, []);

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
      await api.acceptMatch(id);
      showToast('Match accepted!', 'success');
      fetchDashboard();
    } catch (err) {
      showToast(err.message, 'error');
    }
  };

  const handleReject = async (id) => {
    try {
      await api.rejectMatch(id);
      showToast('Match declined', 'info');
      fetchDashboard();
    } catch (err) {
      showToast(err.message, 'error');
    }
  };

  if (loading) {
    return (
      <div className="page-container">
        <div className="loading-state">
          <div className="loading-spinner"></div>
          <p>Loading dashboard...</p>
        </div>
      </div>
    );
  }

  const stats = data?.stats || {};
  const profile = data?.profile || {};
  const recentMatches = data?.recent_matches || [];
  const hasProfile = profile.skills?.length > 0 || profile.interests?.length > 0;

  return (
    <div className="page-container">
      <ToastContainer />

      <div className="dashboard">
        {/* Welcome Section */}
        <div className="dashboard-welcome">
          <div className="welcome-text">
            <h1 className="welcome-title">
              Welcome back, <span className="gradient-text">{user?.name}</span>
            </h1>
            <p className="welcome-subtitle">
              {hasProfile
                ? "Here's your matching overview"
                : "Complete your profile to find great matches"}
            </p>
          </div>
          <div className="welcome-actions">
            {!hasProfile && (
              <Link to="/profile" className="btn btn-primary">
                Complete Profile →
              </Link>
            )}
            <Link to="/matches" className="btn btn-secondary">
              Find Matches
            </Link>
          </div>
        </div>

        {/* Stats Grid */}
        <div className="stats-grid">
          <div className="stat-card">
            <div className="stat-value">{stats.total_matches || 0}</div>
            <div className="stat-label">Total Matches</div>
          </div>
          <div className="stat-card stat-card-accepted">
            <div className="stat-value">{stats.accepted_matches || 0}</div>
            <div className="stat-label">Accepted</div>
          </div>
          <div className="stat-card stat-card-pending">
            <div className="stat-value">{stats.pending_matches || 0}</div>
            <div className="stat-label">Pending</div>
          </div>
          <div className="stat-card stat-card-rejected">
            <div className="stat-value">{stats.rejected_matches || 0}</div>
            <div className="stat-label">Declined</div>
          </div>
        </div>

        {/* Profile Summary */}
        {hasProfile && (
          <div className="dashboard-section">
            <h2 className="section-heading">Your Profile</h2>
            <div className="profile-summary-card">
              <div className="profile-summary-tags">
                {profile.skills?.length > 0 && (
                  <div>
                    <span className="tags-label">Skills:</span>
                    <div className="tags-list">
                      {profile.skills.map((s, i) => (
                        <span key={i} className="tag-pill skill-pill">{s}</span>
                      ))}
                    </div>
                  </div>
                )}
                {profile.interests?.length > 0 && (
                  <div>
                    <span className="tags-label">Interests:</span>
                    <div className="tags-list">
                      {profile.interests.map((s, i) => (
                        <span key={i} className="tag-pill interest-pill">{s}</span>
                      ))}
                    </div>
                  </div>
                )}
              </div>
              <Link to="/profile" className="btn btn-small btn-secondary">Edit Profile</Link>
            </div>
          </div>
        )}

        {/* Recent Matches */}
        <div className="dashboard-section">
          <div className="section-header">
            <h2 className="section-heading">Recent Matches</h2>
            <Link to="/matches" className="section-link">View All →</Link>
          </div>
          {recentMatches.length > 0 ? (
            <div className="matches-list">
              {recentMatches.map(match => (
                <MatchCard
                  key={match.id}
                  match={match}
                  onAccept={handleAccept}
                  onReject={handleReject}
                  currentUserId={user?.id}
                />
              ))}
            </div>
          ) : (
            <div className="empty-state">
              <span className="empty-icon">🔍</span>
              <p className="empty-text">No matches yet.</p>
              <Link to="/matches" className="btn btn-primary btn-small">Find Matches</Link>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
