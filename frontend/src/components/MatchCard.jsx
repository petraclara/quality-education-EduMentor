import './MatchCard.css';

export default function MatchCard({ match, onAccept, onReject, currentUserId }) {
  const scorePercent = Math.round(match.score * 100);

  const getScoreColor = (score) => {
    if (score >= 70) return '#10b981';
    if (score >= 40) return '#f59e0b';
    return '#ef4444';
  };

  return (
    <div className={`match-card match-status-${match.status}`}>
      <div className="match-card-header">
        <div className="match-avatar">
          {match.user_name?.charAt(0)?.toUpperCase() || '?'}
        </div>
        <div className="match-info">
          <h3 className="match-name">{match.user_name}</h3>
          <span className={`match-role role-${match.user_role}`}>{match.user_role}</span>
        </div>
        <div className="match-score" style={{ borderColor: getScoreColor(scorePercent) }}>
          <span className="match-score-value" style={{ color: getScoreColor(scorePercent) }}>
            {scorePercent}%
          </span>
          <span className="match-score-label">match</span>
        </div>
      </div>

      {match.bio && (
        <p className="match-bio">{match.bio}</p>
      )}

      <div className="match-tags-section">
        {match.skills && match.skills.length > 0 && (
          <div className="match-tags-group">
            <span className="match-tags-label">Skills</span>
            <div className="match-tags">
              {match.skills.map((skill, i) => (
                <span key={i} className="match-tag skill-tag">{skill}</span>
              ))}
            </div>
          </div>
        )}
        {match.interests && match.interests.length > 0 && (
          <div className="match-tags-group">
            <span className="match-tags-label">Interests</span>
            <div className="match-tags">
              {match.interests.map((interest, i) => (
                <span key={i} className="match-tag interest-tag">{interest}</span>
              ))}
            </div>
          </div>
        )}
      </div>

      {match.status === 'pending' && (
        <div className="match-actions">
          <button className="match-btn match-btn-accept" onClick={() => onAccept(match.id)}>
            ✓ Accept
          </button>
          <button className="match-btn match-btn-reject" onClick={() => onReject(match.id)}>
            ✗ Decline
          </button>
        </div>
      )}

      {match.status !== 'pending' && (
        <div className={`match-status-badge status-${match.status}`}>
          {match.status === 'accepted' ? '✓ Connected' : '✗ Declined'}
        </div>
      )}
    </div>
  );
}
