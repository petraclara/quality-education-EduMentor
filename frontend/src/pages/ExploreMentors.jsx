import { useState, useEffect, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../api/client';
import { useToast } from '../components/Toast';
import './ExploreMentors.css';

const AVATAR_COLORS = ['#14b8a6','#8b5cf6','#f59e0b','#ec4899','#06b6d4','#10b981','#f43f5e','#6366f1'];

export default function ExploreMentors() {
  const [mentors, setMentors] = useState([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [skillFilter, setSkillFilter] = useState('');
  const [levelFilter, setLevelFilter] = useState('');
  const { showToast, ToastContainer } = useToast();

  useEffect(() => {
    api.getMentors().then(res => { setMentors(res.data || []); setLoading(false); })
      .catch(err => { showToast(err.message, 'error'); setLoading(false); });
  }, []);

  const allSkills = useMemo(() => {
    const s = new Set();
    mentors.forEach(m => m.skills?.forEach(sk => s.add(sk)));
    return Array.from(s).sort();
  }, [mentors]);

  const filtered = useMemo(() => {
    let r = [...mentors];
    if (search) {
      const q = search.toLowerCase();
      r = r.filter(m => m.name.toLowerCase().includes(q) || m.skills?.some(s => s.toLowerCase().includes(q)));
    }
    if (skillFilter) r = r.filter(m => m.skills?.includes(skillFilter));
    if (levelFilter) r = r.filter(m => m.level === levelFilter);
    return r;
  }, [mentors, search, skillFilter, levelFilter]);

  if (loading) return <div className="page-container"><div className="loading-state"><div className="loading-spinner"></div><p>Loading...</p></div></div>;

  return (
    <div className="page-container">
      <ToastContainer />
      <div className="explore">
        <h1 className="explore-title">Explore Mentors</h1>

        <div className="explore-filters">
          <input className="form-input explore-search" type="text" placeholder="Search by name or skill..."
            value={search} onChange={(e) => setSearch(e.target.value)} />
          <select className="explore-select" value={skillFilter} onChange={(e) => setSkillFilter(e.target.value)}>
            <option value="">All Skills</option>
            {allSkills.map(s => <option key={s} value={s}>{s}</option>)}
          </select>
          <select className="explore-select" value={levelFilter} onChange={(e) => setLevelFilter(e.target.value)}>
            <option value="">All Levels</option>
            <option value="beginner">Beginner</option>
            <option value="intermediate">Intermediate</option>
            <option value="advanced">Advanced</option>
          </select>
        </div>

        {filtered.length > 0 ? (
          <div className="explore-grid">
            {filtered.map(m => {
              const bg = AVATAR_COLORS[m.id % AVATAR_COLORS.length];
              const initials = m.name.split(' ').map(n => n[0]).join('').slice(0, 2);
              return (
                <Link key={m.id} to={`/mentor/${m.id}`} className="explore-card card">
                  <div className="ec-top">
                    <div className="ec-avatar" style={{ background: bg }}>{initials}</div>
                    <div className="ec-info">
                      <h3 className="ec-name">{m.name}</h3>
                      {m.level && <span className="ec-level">{m.level}</span>}
                    </div>
                  </div>
                  <p className="ec-bio">{m.bio?.slice(0, 100)}{m.bio?.length > 100 ? '...' : ''}</p>
                  <div className="ec-skills">
                    {m.skills?.slice(0, 3).map((s, i) => <span key={i} className="ec-skill">{s}</span>)}
                    {m.skills?.length > 3 && <span className="ec-skill ec-more">+{m.skills.length - 3}</span>}
                  </div>
                  <span className="ec-cta">View Profile →</span>
                </Link>
              );
            })}
          </div>
        ) : (
          <div className="empty-state">
            <span className="empty-icon">🔍</span>
            <p className="empty-text">No mentors match your filters</p>
            <button className="btn btn-secondary btn-small" onClick={() => { setSearch(''); setSkillFilter(''); setLevelFilter(''); }}>Clear Filters</button>
          </div>
        )}
      </div>
    </div>
  );
}
