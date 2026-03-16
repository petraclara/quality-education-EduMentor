import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './Landing.css';

export default function Landing() {
  const { user } = useAuth();

  return (
    <div className="landing">
      {/* Hero Section */}
      <section className="hero">
        <div className="hero-bg-effects">
          <div className="hero-orb hero-orb-1"></div>
          <div className="hero-orb hero-orb-2"></div>
          <div className="hero-orb hero-orb-3"></div>
        </div>
        <div className="hero-content">
          <div className="hero-badge">🎓 Smart Mentor Matching</div>
          <h1 className="hero-title">
            Find Your Perfect
            <span className="hero-gradient-text"> Learning Partner</span>
          </h1>
          <p className="hero-description">
            EduMentor uses intelligent matching to connect mentors and mentees based on 
            shared skills, interests, and availability. Start your learning journey today.
          </p>
          <div className="hero-cta">
            {user ? (
              <Link to="/dashboard" className="btn btn-primary btn-large">
                Go to Dashboard →
              </Link>
            ) : (
              <>
                <Link to="/register" className="btn btn-primary btn-large">
                  Get Started Free
                </Link>
                <Link to="/login" className="btn btn-secondary btn-large">
                  Sign In
                </Link>
              </>
            )}
          </div>
          <div className="hero-stats">
            <div className="hero-stat">
              <span className="hero-stat-value">Smart</span>
              <span className="hero-stat-label">AI Matching</span>
            </div>
            <div className="hero-stat-divider"></div>
            <div className="hero-stat">
              <span className="hero-stat-value">3-Way</span>
              <span className="hero-stat-label">Compatibility</span>
            </div>
            <div className="hero-stat-divider"></div>
            <div className="hero-stat">
              <span className="hero-stat-value">Free</span>
              <span className="hero-stat-label">Forever</span>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="features">
        <div className="features-container">
          <h2 className="section-title">How It Works</h2>
          <p className="section-subtitle">Three simple steps to find your perfect match</p>
          
          <div className="features-grid">
            <div className="feature-card">
              <div className="feature-icon">📝</div>
              <h3 className="feature-title">Create Your Profile</h3>
              <p className="feature-desc">
                Add your skills, interests, and availability. The more detail you add, 
                the better your matches will be.
              </p>
            </div>
            
            <div className="feature-card">
              <div className="feature-icon">🔍</div>
              <h3 className="feature-title">Find Matches</h3>
              <p className="feature-desc">
                Our algorithm analyzes skill overlap, shared interests, and schedule 
                compatibility to rank your best matches.
              </p>
            </div>
            
            <div className="feature-card">
              <div className="feature-icon">🤝</div>
              <h3 className="feature-title">Connect & Learn</h3>
              <p className="feature-desc">
                Accept matches to connect with mentors or mentees. 
                Build meaningful learning relationships.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Matching Preview */}
      <section className="matching-preview">
        <div className="matching-container">
          <h2 className="section-title">Intelligent Matching</h2>
          <p className="section-subtitle">Our weighted algorithm considers three key factors</p>
          
          <div className="matching-factors">
            <div className="factor-card">
              <div className="factor-bar" style={{ width: '100%' }}>
                <div className="factor-fill" style={{ width: '40%', background: 'linear-gradient(90deg, #14b8a6, #0d9488)' }}></div>
              </div>
              <div className="factor-info">
                <span className="factor-weight">40%</span>
                <span className="factor-name">Skills Match</span>
              </div>
              <p className="factor-desc">Technical and domain expertise overlap</p>
            </div>
            
            <div className="factor-card">
              <div className="factor-bar" style={{ width: '100%' }}>
                <div className="factor-fill" style={{ width: '30%', background: 'linear-gradient(90deg, #8b5cf6, #7c3aed)' }}></div>
              </div>
              <div className="factor-info">
                <span className="factor-weight">30%</span>
                <span className="factor-name">Shared Interests</span>
              </div>
              <p className="factor-desc">Common learning goals and topics</p>
            </div>
            
            <div className="factor-card">
              <div className="factor-bar" style={{ width: '100%' }}>
                <div className="factor-fill" style={{ width: '30%', background: 'linear-gradient(90deg, #f59e0b, #d97706)' }}></div>
              </div>
              <div className="factor-info">
                <span className="factor-weight">30%</span>
                <span className="factor-name">Availability</span>
              </div>
              <p className="factor-desc">Compatible schedules and time zones</p>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="footer">
        <div className="footer-content">
          <span className="footer-brand">🎓 EduMentor</span>
          <span className="footer-text">Quality Education Through Mentorship</span>
        </div>
      </footer>
    </div>
  );
}
