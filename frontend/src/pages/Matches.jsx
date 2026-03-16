import { useState, useEffect } from 'react';
import { api } from '../api/client';
import { useAuth } from '../context/AuthContext';
import MatchCard from '../components/MatchCard';
import { useToast } from '../components/Toast';
import './Matches.css';

export default function Matches() {
  const { user } = useAuth();
  const [matches, setMatches] = useState([]);
  const [loading, setLoading] = useState(true);
  const [finding, setFinding] = useState(false);
  const [filter, setFilter] = useState('all');
  const { showToast, ToastContainer } = useToast();

  useEffect(() => {
    fetchMatches();
  }, []);

  const fetchMatches = async () => {
    try {
      const res = await api.getMyMatches();
      setMatches(res.data || []);
    } catch (err) {
      showToast(err.message, 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleFindMatches = async () => {
    setFinding(true);
    try {
      const res = await api.findMatches();
      const found = res.data || [];
      showToast(`Found ${found.length} potential matches!`, 'success');
      fetchMatches();
    } catch (err) {
      showToast(err.message, 'error');
    } finally {
      setFinding(false);
    }
  };

  const handleAccept = async (id) => {
    try {
      await api.acceptMatch(id);
      showToast('Match accepted! 🎉', 'success');
      setMatches(prev =>
        prev.map(m => m.id === id ? { ...m, status: 'accepted' } : m)
      );
    } catch (err) {
      showToast(err.message, 'error');
    }
  };

  const handleReject = async (id) => {
    try {
      await api.rejectMatch(id);
      showToast('Match declined', 'info');
      setMatches(prev =>
        prev.map(m => m.id === id ? { ...m, status: 'rejected' } : m)
      );
    } catch (err) {
      showToast(err.message, 'error');
    }
  };

  const filteredMatches = matches.filter(m => {
    if (filter === 'all') return true;
    return m.status === filter;
  });

  if (loading) {
    return (
      <div className="page-container">
        <div className="loading-state">
          <div className="loading-spinner"></div>
          <p>Loading matches...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="page-container">
      <ToastContainer />

      <div className="matches-page">
        <div className="matches-header">
          <div>
            <h1 className="page-title">Your Matches</h1>
            <p className="page-subtitle">
              Find and connect with compatible {user?.role === 'mentor' ? 'mentees' : 'mentors'}
            </p>
          </div>
          <button
            id="find-matches-btn"
            className="btn btn-primary"
            onClick={handleFindMatches}
            disabled={finding}
          >
            {finding ? '🔍 Searching...' : '🔍 Find New Matches'}
          </button>
        </div>

        {/* Filter Tabs */}
        <div className="filter-tabs">
          {['all', 'pending', 'accepted', 'rejected'].map(f => (
            <button
              key={f}
              className={`filter-tab ${filter === f ? 'filter-tab-active' : ''}`}
              onClick={() => setFilter(f)}
            >
              {f === 'all' ? 'All' : f.charAt(0).toUpperCase() + f.slice(1)}
              <span className="filter-count">
                {f === 'all' ? matches.length : matches.filter(m => m.status === f).length}
              </span>
            </button>
          ))}
        </div>

        {/* Match Results */}
        {filteredMatches.length > 0 ? (
          <div className="matches-grid">
            {filteredMatches.map(match => (
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
            <span className="empty-icon">
              {filter === 'all' ? '🔍' : '📭'}
            </span>
            <p className="empty-text">
              {filter === 'all'
                ? 'No matches yet. Click "Find New Matches" to discover compatible partners!'
                : `No ${filter} matches.`}
            </p>
            {filter === 'all' && (
              <button
                className="btn btn-primary btn-small"
                onClick={handleFindMatches}
                disabled={finding}
              >
                Find Matches
              </button>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
