import { Link } from 'react-router-dom';
import './Landing.css';

export default function Landing() {
  return (
    <div className="landing">
      <div className="landing-hero">
        <h1 className="hero-title">Find the Right <span className="highlight">Mentor</span> for You</h1>
        <p className="hero-sub">
          Get personalized mentorship based on your goals and interests.
          Connect with experienced mentors who can guide your learning journey.
        </p>
        <div className="hero-actions">
          <Link to="/register" className="btn btn-primary btn-large">Get Started</Link>
          <Link to="/login" className="btn btn-secondary btn-large">I have an account</Link>
        </div>
      </div>

      <div className="landing-features">
        <div className="feature-card">
          <span className="feature-icon">🎯</span>
          <h3>Personalized Matches</h3>
          <p>We match you with mentors based on your interests, level, and goals. No searching needed.</p>
        </div>
        <div className="feature-card">
          <span className="feature-icon">📩</span>
          <h3>Simple Requests</h3>
          <p>Tell mentors what you need help with. They'll review your request and respond quickly.</p>
        </div>
        <div className="feature-card">
          <span className="feature-icon">🤝</span>
          <h3>Meaningful Connections</h3>
          <p>Build real mentorship relationships. Learn from people who've been where you want to go.</p>
        </div>
      </div>

    </div>
  );
}
