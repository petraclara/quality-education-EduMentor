import { useState } from 'react';
import './TagInput.css';

export default function TagInput({ tags = [], onChange, placeholder = 'Add tag...' }) {
  const [input, setInput] = useState('');

  const addTag = (e) => {
    e.preventDefault();
    const tag = input.trim().toLowerCase();
    if (tag && !tags.includes(tag)) {
      onChange([...tags, tag]);
      setInput('');
    }
  };

  const removeTag = (index) => {
    onChange(tags.filter((_, i) => i !== index));
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      addTag(e);
    }
    if (e.key === 'Backspace' && input === '' && tags.length > 0) {
      removeTag(tags.length - 1);
    }
  };

  return (
    <div className="tag-input-container">
      <div className="tag-list">
        {tags.map((tag, i) => (
          <span key={i} className="tag">
            {tag}
            <button type="button" className="tag-remove" onClick={() => removeTag(i)}>×</button>
          </span>
        ))}
        <input
          type="text"
          className="tag-input"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={tags.length === 0 ? placeholder : ''}
        />
      </div>
    </div>
  );
}
