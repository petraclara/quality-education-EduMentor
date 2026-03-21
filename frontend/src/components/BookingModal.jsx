import { useState } from 'react';
import { api } from '../api/client';
import './BookingModal.css';

const DAYS = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
const TIMES = ['Morning', 'Afternoon', 'Evening'];
const TIME_DISPLAY = {
  morning: '9:00 AM - 12:00 PM',
  afternoon: '1:00 PM - 5:00 PM',
  evening: '6:00 PM - 9:00 PM',
};

export default function BookingModal({ mentor, onClose, onSuccess }) {
  const [selectedDate, setSelectedDate] = useState('');
  const [selectedSlot, setSelectedSlot] = useState('');
  const [note, setNote] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  // Generate next 14 days
  const dates = [];
  const today = new Date();
  for (let i = 1; i <= 14; i++) {
    const d = new Date(today);
    d.setDate(today.getDate() + i);
    dates.push(d);
  }

  // Get available slots for selected date
  const getAvailableSlots = () => {
    if (!selectedDate) return [];
    const date = new Date(selectedDate);
    const dayName = DAYS[date.getDay() === 0 ? 6 : date.getDay() - 1].toLowerCase();
    return TIMES.filter(time => {
      const slot = `${dayName}-${time.toLowerCase()}`;
      return mentor.availability?.includes(slot);
    }).map(time => ({
      key: `${dayName}-${time.toLowerCase()}`,
      label: time,
      timeRange: TIME_DISPLAY[time.toLowerCase()],
    }));
  };

  const availableSlots = getAvailableSlots();

  const handleSubmit = async () => {
    if (!selectedDate || !selectedSlot) {
      setError('Please select both a date and time slot');
      return;
    }

    setSubmitting(true);
    setError('');
    try {
      await api.createBooking({
        mentor_id: mentor.id,
        date: selectedDate,
        time_slot: selectedSlot,
        note: note.trim(),
      });
      onSuccess();
    } catch (err) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  };

  const formatDate = (d) => {
    const days = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    return `${days[d.getDay()]}, ${months[d.getMonth()]} ${d.getDate()}`;
  };

  const toISODate = (d) => {
    return d.toISOString().split('T')[0];
  };

  // Check if the date has any available slots
  const dateHasSlots = (d) => {
    const dayName = DAYS[d.getDay() === 0 ? 6 : d.getDay() - 1].toLowerCase();
    return TIMES.some(time => {
      const slot = `${dayName}-${time.toLowerCase()}`;
      return mentor.availability?.includes(slot);
    });
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-card" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <div>
            <h2 className="modal-title">Book a Session</h2>
            <p className="modal-subtitle">with {mentor.name}</p>
          </div>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>

        {error && <div className="modal-error">{error}</div>}

        {/* Step 1: Select Date */}
        <div className="modal-section">
          <h3 className="modal-step-label">
            <span className="step-num">1</span> Select a Date
          </h3>
          <div className="date-grid">
            {dates.map(d => {
              const iso = toISODate(d);
              const hasSlots = dateHasSlots(d);
              return (
                <button
                  key={iso}
                  className={`date-btn ${selectedDate === iso ? 'date-btn-selected' : ''} ${!hasSlots ? 'date-btn-disabled' : ''}`}
                  onClick={() => { if (hasSlots) { setSelectedDate(iso); setSelectedSlot(''); } }}
                  disabled={!hasSlots}
                >
                  <span className="date-day">{['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'][d.getDay()]}</span>
                  <span className="date-num">{d.getDate()}</span>
                </button>
              );
            })}
          </div>
        </div>

        {/* Step 2: Select Time Slot */}
        {selectedDate && (
          <div className="modal-section">
            <h3 className="modal-step-label">
              <span className="step-num">2</span> Select a Time
            </h3>
            {availableSlots.length > 0 ? (
              <div className="slot-list">
                {availableSlots.map(slot => (
                  <button
                    key={slot.key}
                    className={`slot-btn ${selectedSlot === slot.key ? 'slot-btn-selected' : ''}`}
                    onClick={() => setSelectedSlot(slot.key)}
                  >
                    <span className="slot-label">{slot.label}</span>
                    <span className="slot-time">{slot.timeRange}</span>
                  </button>
                ))}
              </div>
            ) : (
              <p className="no-slots-msg">No available slots on this date</p>
            )}
          </div>
        )}

        {/* Step 3: Add Note */}
        {selectedSlot && (
          <div className="modal-section">
            <h3 className="modal-step-label">
              <span className="step-num">3</span> Add a Note <span className="optional-label">(optional)</span>
            </h3>
            <textarea
              className="modal-textarea"
              value={note}
              onChange={(e) => setNote(e.target.value)}
              placeholder="What would you like to discuss in this session?"
              rows={3}
            />
          </div>
        )}

        {/* Confirm */}
        {selectedSlot && (
          <div className="modal-footer">
            <div className="booking-summary">
              <span>📅 {formatDate(new Date(selectedDate))}</span>
              <span>🕐 {TIME_DISPLAY[selectedSlot.split('-')[1]]}</span>
            </div>
            <button
              id="confirm-booking"
              className="btn btn-primary btn-large"
              onClick={handleSubmit}
              disabled={submitting}
              style={{ width: '100%' }}
            >
              {submitting ? 'Booking...' : '✓ Confirm Booking'}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
