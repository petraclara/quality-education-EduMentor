import './AvailabilityGrid.css';

const DAYS = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
const TIMES = ['Morning', 'Afternoon', 'Evening'];

export default function AvailabilityGrid({ availability = [], onChange }) {
  const toggleSlot = (slot) => {
    if (availability.includes(slot)) {
      onChange(availability.filter(s => s !== slot));
    } else {
      onChange([...availability, slot]);
    }
  };

  return (
    <div className="availability-grid">
      <div className="grid-header">
        <div className="grid-corner"></div>
        {TIMES.map(time => (
          <div key={time} className="grid-time-label">{time}</div>
        ))}
      </div>
      {DAYS.map(day => (
        <div key={day} className="grid-row">
          <div className="grid-day-label">{day.slice(0, 3)}</div>
          {TIMES.map(time => {
            const slot = `${day.toLowerCase()}-${time.toLowerCase()}`;
            const active = availability.includes(slot);
            return (
              <button
                key={slot}
                type="button"
                className={`grid-cell ${active ? 'grid-cell-active' : ''}`}
                onClick={() => toggleSlot(slot)}
                title={`${day} ${time}`}
              />
            );
          })}
        </div>
      ))}
    </div>
  );
}
