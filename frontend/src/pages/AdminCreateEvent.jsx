import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import '../styles/App.css';
import { adminCategories } from '../constants/admin';

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

const initialForm = {
  title: '',
  description: '',
  imageUrl: '',
  date: '',
  time: '',
  duration: '',
  location: '',
  category: '',
  capacity: '',
  price: '',
};

function AdminCreateEvent() {
  const [form, setForm] = useState(initialForm);
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
    if (!form.date.trim()) nextErrors.date = 'La fecha es obligatoria para publicar el evento.';
    setErrors(nextErrors);
    return Object.keys(nextErrors).length === 0;
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    if (!validate()) return;
    setFeedback('Frontend listo: este formulario quedo preparado para conectar luego con el backend de eventos.');
  };

  return (
    <main className="admin-page">
      <header className="event-detail-topbar">
        <Link className="event-detail-brand" to="/">Golden Ticket</Link>
        <nav className="event-detail-nav">
          <Link to="/">Eventos</Link>
          <Link to="/admin">Panel admin</Link>
          <Link to="/admin/eventos/nuevo" className="is-active">Nuevo evento</Link>
        </nav>
        <div className="event-detail-actions">
          <Link className="event-detail-login" to="/admin">Volver al panel</Link>
        </div>
      </header>

      <section className="admin-shell admin-form-shell">
        <div className="admin-heading admin-heading-compact">
          <div>
            <span className="admin-kicker">Crear nuevo evento</span>
            <h1>Configura una experiencia lista para publicar en la plataforma.</h1>
            <p>Esta vista es solo frontend y replica el flujo visual de alta de eventos con validaciones basicas y placeholders realistas.</p>
          </div>
        </div>

        <form className="admin-form-card" onSubmit={handleSubmit}>
          <div className="admin-form-hero">
            <div className="admin-form-hero-overlay" />
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
                placeholder="Escribe los detalles principales del evento aqui..."
                value={form.description}
                onChange={(event) => handleChange('description', event.target.value)}
              />
            </div>

            <div className="admin-field admin-field-full">
              <label htmlFor="event-image">URL de la imagen destacada</label>
              <div className="admin-input-icon">
                <span><ImageIcon /></span>
                <input
                  id="event-image"
                  placeholder="https://ejemplo.com/imagen.jpg"
                  value={form.imageUrl}
                  onChange={(event) => handleChange('imageUrl', event.target.value)}
                />
              </div>
            </div>

            <div className="admin-field">
              <label htmlFor="event-date">Fecha *</label>
              <div className="admin-input-icon">
                <span><CalendarIcon /></span>
                <input
                  id="event-date"
                  type="date"
                  className={errors.date ? 'has-error' : ''}
                  value={form.date}
                  onChange={(event) => handleChange('date', event.target.value)}
                />
              </div>
              {errors.date && <span className="admin-field-error">{errors.date}</span>}
            </div>

            <div className="admin-field">
              <label htmlFor="event-time">Hora de inicio</label>
              <div className="admin-input-icon">
                <span><ClockIcon /></span>
                <input
                  id="event-time"
                  type="time"
                  value={form.time}
                  onChange={(event) => handleChange('time', event.target.value)}
                />
              </div>
            </div>

            <div className="admin-field">
              <label htmlFor="event-duration">Duracion (hs)</label>
              <input
                id="event-duration"
                type="number"
                min="0"
                placeholder="Ej. 2"
                value={form.duration}
                onChange={(event) => handleChange('duration', event.target.value)}
              />
            </div>

            <div className="admin-field admin-field-full">
              <label htmlFor="event-location">Ubicacion / lugar</label>
              <div className="admin-input-icon">
                <span><PinIcon /></span>
                <input
                  id="event-location"
                  placeholder="Nombre del recinto o direccion completa"
                  value={form.location}
                  onChange={(event) => handleChange('location', event.target.value)}
                />
              </div>
            </div>

            <div className="admin-field">
              <label htmlFor="event-category">Categoria</label>
              <select
                id="event-category"
                value={form.category}
                onChange={(event) => handleChange('category', event.target.value)}
              >
                <option value="">Selecciona una...</option>
                {adminCategories.map((category) => (
                  <option key={category} value={category}>{category}</option>
                ))}
              </select>
            </div>

            <div className="admin-field">
              <label htmlFor="event-capacity">Cupo maximo</label>
              <div className="admin-input-icon">
                <span><UsersIcon /></span>
                <input
                  id="event-capacity"
                  type="number"
                  min="0"
                  placeholder="Ilimitado si se deja vacio"
                  value={form.capacity}
                  onChange={(event) => handleChange('capacity', event.target.value)}
                />
              </div>
            </div>

            <div className="admin-field">
              <label htmlFor="event-price">Precio ticket base</label>
              <div className="admin-input-prefix">
                <span>$</span>
                <input
                  id="event-price"
                  type="number"
                  min="0"
                  step="0.01"
                  placeholder="0.00"
                  value={form.price}
                  onChange={(event) => handleChange('price', event.target.value)}
                />
              </div>
            </div>
          </div>

          <div className="admin-form-footer">
            <button type="button" className="admin-ghost-button" onClick={() => navigate('/admin')}>Cancelar</button>
            <button type="submit" className="admin-submit-button">Guardar evento</button>
          </div>

          {feedback && <p className="purchase-note success-state admin-inline-feedback">{feedback}</p>}
        </form>
      </section>
    </main>
  );
}

export default AdminCreateEvent;
