import { useState, useEffect, useRef } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import '../../styles/App.css';
import { createEvent, updateEvent, getEventByID } from '../../services/api/client';

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

const ChevronIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="m7 10 5 5 5-5" />
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
  price: '',
  location: '',
  coordinates: '',
  imageUrl: '',
});

const MAX_CAPACITY_DIGITS = 7;
const ADMIN_CATEGORIES = ['Musica', 'Deportes', 'Teatro', 'Charlas', 'Cursos', 'E-Sports', 'Comedia'];

function AdminCreateEvent() {
  const { id } = useParams();
  const [form, setForm] = useState(createInitialForm);
  const [errors, setErrors] = useState({});
  const [feedback, setFeedback] = useState('');
  const [isCategoryOpen, setIsCategoryOpen] = useState(false);
  const categoryMenuRef = useRef(null);
  const navigate = useNavigate();

  useEffect(() => {
    const { body, documentElement } = document;
    const scrollY = window.scrollY;
    const previousBodyStyle = {
      overflow: body.style.overflow,
      position: body.style.position,
      top: body.style.top,
      width: body.style.width,
      left: body.style.left,
      right: body.style.right,
    };
    const previousHtmlOverflow = documentElement.style.overflow;

    documentElement.style.overflow = 'hidden';
    body.style.overflow = 'hidden';
    body.style.position = 'fixed';
    body.style.top = `-${scrollY}px`;
    body.style.left = '0';
    body.style.right = '0';
    body.style.width = '100%';

    return () => {
      documentElement.style.overflow = previousHtmlOverflow;
      body.style.overflow = previousBodyStyle.overflow;
      body.style.position = previousBodyStyle.position;
      body.style.top = previousBodyStyle.top;
      body.style.width = previousBodyStyle.width;
      body.style.left = previousBodyStyle.left;
      body.style.right = previousBodyStyle.right;
      window.scrollTo(0, scrollY);
    };
  }, []);

  useEffect(() => {
    if (id) {
      getEventByID(id)
        .then((response) => {
          const ev = response.data;
          setForm({
            title: ev.titulo || '',
            description: ev.descripcion || '',
            category: ev.categoria || '',
            eventDate: ev.fecha || '',
            startTime: ev.hora_inicio || '',
            endTime: ev.hora_fin || '',
            capacity: ev.capacidad ? ev.capacidad.toString() : '',
            price: ev.precio !== undefined ? ev.precio.toString() : '',
            location: ev.ubicacion || '',
            coordinates: ev.coordenadas || '',
            imageUrl: ev.url_imagen || '',
          });
        })
        .catch((err) => {
          console.error('Error fetching event details:', err);
          setFeedback('Error al cargar los datos del evento.');
        });
    }
  }, [id]);

  useEffect(() => {
    if (!isCategoryOpen) return undefined;

    const handlePointerDown = (event) => {
      if (!categoryMenuRef.current?.contains(event.target)) {
        setIsCategoryOpen(false);
      }
    };

    const handleEscape = (event) => {
      if (event.key === 'Escape') {
        setIsCategoryOpen(false);
      }
    };

    document.addEventListener('mousedown', handlePointerDown);
    document.addEventListener('keydown', handleEscape);

    return () => {
      document.removeEventListener('mousedown', handlePointerDown);
      document.removeEventListener('keydown', handleEscape);
    };
  }, [isCategoryOpen]);

  const handleChange = (field, value) => {
    if (field === 'capacity') {
      const sanitized = value.replace(/\D/g, '').slice(0, MAX_CAPACITY_DIGITS);
      setForm((current) => ({ ...current, [field]: sanitized }));
      setErrors((current) => ({ ...current, [field]: '' }));
      setFeedback('');
      return;
    }

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
    else if (!/^\d+$/.test(form.capacity)) nextErrors.capacity = 'La capacidad debe contener solo numeros.';
    else if (form.capacity.length > MAX_CAPACITY_DIGITS) nextErrors.capacity = `La capacidad no puede superar los ${MAX_CAPACITY_DIGITS} digitos.`;
    if (!form.price.trim()) nextErrors.price = 'El precio es obligatorio.';
    else if (isNaN(form.price) || parseFloat(form.price) < 0) nextErrors.price = 'El precio debe ser un número mayor o igual a 0.';
    setErrors(nextErrors);
    return Object.keys(nextErrors).length === 0;
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    if (!validate()) return;

    try {
      const payload = {
        titulo: form.title,
        descripcion: form.description,
        categoria: form.category,
        fecha: form.eventDate,
        hora_inicio: form.startTime,
        hora_fin: form.endTime,
        ubicacion: form.location,
        coordenadas: form.coordinates,
        url_imagen: form.imageUrl,
        capacidad: parseInt(form.capacity, 10),
        precio: parseFloat(form.price),
      };

      if (id) {
        await updateEvent(id, payload);
        setFeedback('¡Evento actualizado con éxito en el servidor!');
      } else {
        await createEvent(payload);
        setFeedback('¡Evento creado con éxito en el servidor!');
      }
      
      setTimeout(() => {
        navigate('/admin');
      }, 1500);
    } catch (error) {
      const errorMsg = error.response?.data?.error || 'Error al conectar con el servidor';
      setFeedback(`Error: ${errorMsg}`);
    }
  };

  return (
    <div className="admin-dialog-page">
      <div className="admin-dialog-shell">
        <Link className="admin-dialog-backdrop" to="/admin/dashboard" aria-label="Cerrar dialogo" />

        <form className="admin-form-card admin-dialog-card" onSubmit={handleSubmit}>
          <div className="admin-dialog-header">
            <h1>{id ? 'Editar evento' : 'Crear evento'}</h1>

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

            <div className="admin-field admin-field-full" ref={categoryMenuRef}>
              <label htmlFor="event-category">Categoria</label>
              <button
                id="event-category"
                type="button"
                className={`admin-custom-select ${errors.category ? 'has-error' : ''} ${isCategoryOpen ? 'is-open' : ''}`}
                onClick={() => setIsCategoryOpen((current) => !current)}
                aria-haspopup="listbox"
                aria-expanded={isCategoryOpen}
              >
                <span className={form.category ? 'has-value' : ''}>
                  {form.category || 'Selecciona una...'}
                </span>
                <ChevronIcon />
              </button>
              {isCategoryOpen && (
                <div className="admin-custom-select-menu" role="listbox" aria-labelledby="event-category">
                  {ADMIN_CATEGORIES.map((category) => (
                    <button
                      key={category}
                      type="button"
                      className={`admin-custom-select-option ${form.category === category ? 'is-selected' : ''}`}
                      onClick={() => {
                        handleChange('category', category);
                        setIsCategoryOpen(false);
                      }}
                    >
                      {category}
                    </button>
                  ))}
                </div>
              )}
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
                  type="text"
                  inputMode="numeric"
                  maxLength={MAX_CAPACITY_DIGITS}
                  className={errors.capacity ? 'has-error' : ''}
                  placeholder="Ej. 5000"
                  value={form.capacity}
                  onChange={(event) => handleChange('capacity', event.target.value)}
                />
              </div>
              {errors.capacity && <span className="admin-field-error">{errors.capacity}</span>}
            </div>

            <div className="admin-field">
              <label htmlFor="event-price">Precio *</label>
              <div className="admin-input-icon">
                <span>$</span>
                <input
                  id="event-price"
                  type="number"
                  min="0"
                  step="0.01"
                  className={errors.price ? 'has-error' : ''}
                  placeholder="Ej. 1500.00"
                  value={form.price}
                  onChange={(event) => handleChange('price', event.target.value)}
                />
              </div>
              {errors.price && <span className="admin-field-error">{errors.price}</span>}
            </div>

          </div>

          <div className="admin-form-footer">
            <button type="button" className="admin-ghost-button" onClick={() => navigate('/admin/dashboard')}>Cancelar</button>
            <button type="submit" className="admin-submit-button">{id ? 'Guardar cambios' : 'Publicar evento'}</button>
          </div>

          {feedback && <p className="purchase-note success-state admin-inline-feedback">{feedback}</p>}
        </form>
      </div>
    </div>
  );
}

export default AdminCreateEvent;
