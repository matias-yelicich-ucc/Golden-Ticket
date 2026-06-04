import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import '../../styles/App.css';
import { adminCategories } from '../../constants/admin';

const ImageIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <rect x="3" y="5" width="18" height="14" rx="2" />
    <circle cx="9" cy="10" r="1.5" />
    <path d="m21 16-5-5-7 7" />
  </svg>
);

const CalendarIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M7 3v4M17 3v4M4 9h16M5 5h14a1 1 0 0 1 1 1v14H4V6a1 1 0 0 1 1-1Z" />
  </svg>
);

const ClockIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <circle cx="12" cy="12" r="9" />
    <path d="M12 7v5l3 2" />
  </svg>
);

const PinIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M12 21s7-6.1 7-12a7 7 0 0 0-14 0c0 5.9 7 12 7 12Z" />
    <circle cx="12" cy="9" r="2.3" />
  </svg>
);

const UsersIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M16 21v-2a4 4 0 0 0-4-4H7a4 4 0 0 0-4 4v2" />
    <circle cx="9.5" cy="7" r="3" />
    <path d="M22 21v-2a4 4 0 0 0-3-3.87" />
    <path d="M16 4.13a4 4 0 0 1 0 7.75" />
  </svg>
);

const CrosshairIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <circle cx="12" cy="12" r="3" />
    <path d="M12 2v3" />
    <path d="M12 19v3" />
    <path d="M2 12h3" />
    <path d="M19 12h3" />
  </svg>
);

const CloseIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M18 6 6 18" />
    <path d="m6 6 12 12" />
  </svg>
);

const createInitialForm = () => ({
  title: '',
  description: '',
  category: '',
  eventDate: '',
  startTime: '',
  endTime: '',
  capacity: '',
  location: '',
  coordinates: '',
  imageUrl: '',
});

function AdminCreateEvent() {
  const [form, setForm] = useState(createInitialForm);
  const [errors, setErrors] = useState({});
  const [feedback, setFeedback] = useState('');
  const navigate = useNavigate();

  const handleChange = (field, value) => {
    setForm((current) => ({ ...current, [field]: value }));
    setErrors((current) => ({ ...current, [field]: '' }));
    setFeedback('');
  };

  const validate = () => {
    const nextErrors = {};
    if (!form.title.trim()) nextErrors.title = 'El titulo es obligatorio para guardar el evento.';
    if (!form.description.trim()) nextErrors.description = 'La descripcion es obligatoria.';
    if (!form.category.trim()) nextErrors.category = 'La categoria es obligatoria.';
    if (!form.imageUrl.trim()) nextErrors.imageUrl = 'La imagen destacada es obligatoria.';
    if (!form.eventDate.trim()) nextErrors.eventDate = 'La fecha es obligatoria para publicar el evento.';
    if (!form.startTime.trim()) nextErrors.startTime = 'La hora de inicio es obligatoria.';
    if (!form.endTime.trim()) nextErrors.endTime = 'La hora de fin es obligatoria.';
    if (!form.location.trim()) nextErrors.location = 'La ubicacion es obligatoria.';
    if (!form.coordinates.trim()) nextErrors.coordinates = 'Las coordenadas son obligatorias.';
    if (!form.capacity.trim()) nextErrors.capacity = 'La capacidad es obligatoria.';
    setErrors(nextErrors);
    return Object.keys(nextErrors).length === 0;
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    if (!validate()) return;
    setFeedback('Frontend listo: este formulario quedo preparado para conectar luego con el backend de eventos.');
  };

  return (
    <div className="admin-dialog-page">
      <div className="admin-dialog-shell">
        <Link className="admin-dialog-backdrop" to="/admin/dashboard" aria-label="Cerrar dialogo" />

        <form className="admin-form-card admin-dialog-card" onSubmit={handleSubmit}>
          <div className="admin-dialog-header">
            <div>
              <span className="admin-kicker">Crear nuevo evento</span>
              <h1>Completa la ficha del evento</h1>
              <p>Dialogo frontend para modelar la entidad de eventos y dejar lista la integracion futura del backend.</p>
            </div>

            <button type="button" className="admin-dialog-close" onClick={() => navigate('/admin/dashboard')} aria-label="Cerrar dialogo">
              <CloseIcon />
            </button>
          </div>

          <div className="admin-form-grid">
            <div className="admin-field admin-field-full">
              <label htmlFor="event-title">Titulo del evento *</label>
              <input
                id="event-title"
                className={errors.title ? 'has-error' : ''}
                placeholder="Ej. Conferencia de Tecnologia 2024"
                value={form.title}
                onChange={(event) => handleChange('title', event.target.value)}
              />
              {errors.title && <span className="admin-field-error">{errors.title}</span>}
            </div>

            <div className="admin-field admin-field-full">
              <label htmlFor="event-description">Descripcion</label>
              <textarea
                id="event-description"
                className={errors.description ? 'has-error' : ''}
                placeholder="Escribe los detalles principales del evento aqui..."
                value={form.description}
                onChange={(event) => handleChange('description', event.target.value)}
              />
              {errors.description && <span className="admin-field-error">{errors.description}</span>}
            </div>

            <div className="admin-field admin-field-full">
              <label htmlFor="event-category">Categoria</label>
              <select
                id="event-category"
                className={errors.category ? 'has-error' : ''}
                value={form.category}
                onChange={(event) => handleChange('category', event.target.value)}
              >
                <option value="">Selecciona una...</option>
                {adminCategories.map((category) => (
                  <option key={category} value={category}>{category}</option>
                ))}
              </select>
              {errors.category && <span className="admin-field-error">{errors.category}</span>}
            </div>

            <div className="admin-field admin-field-full">
              <label htmlFor="event-image">URL de la imagen destacada</label>
              <div className="admin-input-icon">
                <span><ImageIcon /></span>
                <input
                  id="event-image"
                  className={errors.imageUrl ? 'has-error' : ''}
                  placeholder="https://ejemplo.com/imagen.jpg"
                  value={form.imageUrl}
                  onChange={(event) => handleChange('imageUrl', event.target.value)}
                />
              </div>
              {errors.imageUrl && <span className="admin-field-error">{errors.imageUrl}</span>}
            </div>

            <div className="admin-field">
              <label htmlFor="event-date">Fecha del evento *</label>
              <div className="admin-input-icon">
                <span><CalendarIcon /></span>
                <input
                  id="event-date"
                  type="date"
                  className={errors.eventDate ? 'has-error' : ''}
                  value={form.eventDate}
                  onChange={(event) => handleChange('eventDate', event.target.value)}
                />
              </div>
              <span className="admin-field-error">{errors.eventDate || ' '}</span>
            </div>

            <div className="admin-field">
              <label htmlFor="event-time">Hora de inicio</label>
              <div className="admin-input-icon">
                <span><ClockIcon /></span>
                <input
                  id="event-time"
                  type="time"
                  className={errors.startTime ? 'has-error' : ''}
                  value={form.startTime}
                  onChange={(event) => handleChange('startTime', event.target.value)}
                />
              </div>
              <span className="admin-field-error">{errors.startTime || ' '}</span>
            </div>

            <div className="admin-field">
              <label htmlFor="event-end-time">Hora de fin</label>
              <div className="admin-input-icon">
                <span><ClockIcon /></span>
                <input
                  id="event-end-time"
                  type="time"
                  className={errors.endTime ? 'has-error' : ''}
                  value={form.endTime}
                  onChange={(event) => handleChange('endTime', event.target.value)}
                />
              </div>
              <span className="admin-field-error">{errors.endTime || ' '}</span>
            </div>

            <div className="admin-field admin-field-full">
              <label htmlFor="event-location">Ubicacion</label>
              <div className="admin-input-icon">
                <span><PinIcon /></span>
                <input
                  id="event-location"
                  className={errors.location ? 'has-error' : ''}
                  placeholder="Nombre del recinto o direccion completa"
                  value={form.location}
                  onChange={(event) => handleChange('location', event.target.value)}
                />
              </div>
              {errors.location && <span className="admin-field-error">{errors.location}</span>}
            </div>

            <div className="admin-field">
              <label htmlFor="event-coordinates">Coordenadas</label>
              <div className="admin-input-icon">
                <span><CrosshairIcon /></span>
                <input
                  id="event-coordinates"
                  className={errors.coordinates ? 'has-error' : ''}
                  placeholder="-31.4201, -64.1888"
                  value={form.coordinates}
                  onChange={(event) => handleChange('coordinates', event.target.value)}
                />
              </div>
              {errors.coordinates && <span className="admin-field-error">{errors.coordinates}</span>}
            </div>

            <div className="admin-field">
              <label htmlFor="event-capacity">Capacidad</label>
              <div className="admin-input-icon">
                <span><UsersIcon /></span>
                <input
                  id="event-capacity"
                  type="number"
                  min="0"
                  className={errors.capacity ? 'has-error' : ''}
                  placeholder="Ej. 5000"
                  value={form.capacity}
                  onChange={(event) => handleChange('capacity', event.target.value)}
                />
              </div>
              {errors.capacity && <span className="admin-field-error">{errors.capacity}</span>}
            </div>

          </div>

          <div className="admin-form-footer">
            <button type="button" className="admin-ghost-button" onClick={() => navigate('/admin/dashboard')}>Cancelar</button>
            <button type="submit" className="admin-submit-button">Guardar evento</button>
          </div>

          {feedback && <p className="purchase-note success-state admin-inline-feedback">{feedback}</p>}
        </form>
      </div>
    </div>
  );
}

export default AdminCreateEvent;
