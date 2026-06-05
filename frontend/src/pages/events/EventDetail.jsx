import { useEffect, useState } from 'react';
import { Link, Navigate, useParams, useNavigate } from 'react-router-dom';
import '../../styles/App.css';
import { getEventByID, buyTickets } from '../../services/api/client';

const CalendarIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M7 3v4M17 3v4M4 9h16M5 5h14a1 1 0 0 1 1 1v14H4V6a1 1 0 0 1 1-1Z" />
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

const MinusIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M5 12h14" />
  </svg>
);

const PlusIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M12 5v14M5 12h14" />
  </svg>
);

const formatCurrency = (value) =>
  value.toLocaleString('es-AR', {
    style: 'currency',
    currency: 'ARS',
    maximumFractionDigits: 0,
  });

const slugify = (title, id) => {
  const clean = (title || '')
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)+/g, '');
  return `${clean}-${id}`;
};

const normalizeEvent = (dbEvent) => {
  const isSoldOut = dbEvent.entradas_disponibles === 0;
  return {
    id: dbEvent.id,
    slug: slugify(dbEvent.titulo, dbEvent.id),
    title: dbEvent.titulo,
    subtitle: dbEvent.descripcion || 'Detalles del evento.',
    description: [dbEvent.descripcion || 'Sin descripciÃ³n.'],
    category: dbEvent.categoria || 'MÃºsica',
    date: dbEvent.fecha,
    fullDate: dbEvent.fecha,
    timeRange: `${dbEvent.hora_inicio} - ${dbEvent.hora_fin}`,
    location: dbEvent.ubicacion,
    address: dbEvent.ubicacion,
    capacity: dbEvent.capacidad,
    entradas_disponibles: dbEvent.entradas_disponibles,
    soldOut: isSoldOut,
    badge: isSoldOut ? 'AGOTADO' : '',
    price: '$5.000',
    numericPrice: 5000,
    tickets: [
      {
        name: 'Entrada General',
        description: 'Acceso general al evento.',
      }
    ],
    urlImagen: dbEvent.url_imagen,
  };
};

function EventDetail() {
  const { slug } = useParams();
  const navigate = useNavigate();
  const [quantity, setQuantity] = useState(1);
  const [purchasing, setPurchasing] = useState(false);
  const [modal, setModal] = useState({ isOpen: false, type: '', title: '', message: '' });
  const [event, setEvent] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const isAuthenticated = Boolean(localStorage.getItem('token'));

  useEffect(() => {
    const loadEvent = async () => {
      const eventId = slug.split('-').pop();
      if (!eventId || isNaN(Number(eventId))) {
        setError('ID de evento inválido.');
        setLoading(false);
        return;
      }

      try {
        const response = await getEventByID(eventId);
        if (response.data) {
          setEvent(normalizeEvent(response.data));
        } else {
          setEvent(null);
        }
      } catch (err) {
        console.error('Error fetching event detail:', err);
        setError('No se pudo conectar con el servidor o el evento no existe.');
      } finally {
        setLoading(false);
      }
    };

    void loadEvent();
  }, [slug]);
  useEffect(() => {
    window.scrollTo(0, 0);
  }, [slug]);

  if (loading) {
    return (
      <main className="event-detail-page">
        <header className="event-detail-topbar">
          <Link className="event-detail-brand" to="/">Golden Ticket</Link>
          <nav className="event-detail-nav">
            <Link to="/">Eventos</Link>
          </nav>
        </header>
        <p className="empty-state">Cargando detalles del evento...</p>
      </main>
    );
  }

  if (error || !event) {
    return <Navigate to="/" replace />;
  }

  const total = event.numericPrice * quantity;

  const handlePurchase = () => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }

    setPurchasing(true);

    buyTickets(event.id, { cantidad: quantity })
      .then(() => {
        setModal({
          isOpen: true,
          type: 'success',
          title: 'Â¡Compra Exitosa!',
          message: `Adquiriste ${quantity} entrada(s) para "${event.title}".`
        });
        setEvent((prev) => ({
          ...prev,
          entradas_disponibles: Math.max(0, prev.entradas_disponibles - quantity),
          soldOut: prev.entradas_disponibles - quantity <= 0,
        }));
      })
      .catch((err) => {
        console.error('Error buying tickets:', err);
        const backendError = err.response?.data?.error || 'No se pudo completar la compra.';
        setModal({
          isOpen: true,
          type: 'error',
          title: 'Error en la Compra',
          message: backendError
        });
      })
      .finally(() => {
        setPurchasing(false);
      });
  };

  return (
    <main className="event-detail-page">
      <header className="event-detail-topbar">
        <Link className="event-detail-brand" to="/">Golden Ticket</Link>
        <nav className="event-detail-nav">
          <Link to="/">Eventos</Link>
          {isAuthenticated ? <Link to="/mis-entradas">Mis entradas</Link> : <Link to="/login">Iniciar sesion</Link>}
        </nav>
      </header>

      <>
        <section 
          className="event-detail-hero event-media-concert"
          style={event.urlImagen ? { 
            backgroundImage: `url(${event.urlImagen})`, 
            backgroundSize: 'cover', 
            backgroundPosition: 'center' 
          } : {}}
        >
          <div className="event-detail-hero-overlay" />
          <div className="event-detail-hero-content">
            <span className="event-detail-category">{event.category}</span>
            <h1>{event.title}</h1>
            <p>{event.subtitle}</p>
          </div>
        </section>

        <section className="event-detail-body">
          <div className="event-detail-main">
            <div className="event-detail-info-grid">
              <article className="event-detail-info-card">
                <div className="event-detail-icon-wrap">
                  <CalendarIcon />
                </div>
                <div>
                  <span>Fecha y horario</span>
                  <strong>{event.fullDate}</strong>
                  <p>{event.timeRange}</p>
                </div>
              </article>

              <article className="event-detail-info-card">
                <div className="event-detail-icon-wrap">
                  <PinIcon />
                </div>
                <div>
                  <span>Ubicacion</span>
                  <strong>{event.location}</strong>
                  <p>{event.address}</p>
                </div>
              </article>
            </div>

            <section className="event-detail-copy">
              <h2>Sobre este evento</h2>
              {event.description.map((paragraph) => <p key={paragraph}>{paragraph}</p>)}
              <p>Este evento forma parte del catalogo cliente definido por la consigna y permite completar el flujo de consulta, detalle y compra.</p>
            </section>

            <section className="event-detail-capacity">
              <h2>Capacidad y disponibilidad</h2>
              <div className="event-detail-capacity-row">
                <UsersIcon />
                <span>{event.capacity}</span>
              </div>
              <div className="event-detail-map">
                <div className="event-detail-map-grid" />
                <div className="event-detail-map-park" />
                <div className="event-detail-map-marker marker-one">
                  <PinIcon />
                </div>
                <div className="event-detail-map-marker marker-two">
                  <PinIcon />
                </div>
                <div className="event-detail-map-label">{event.location}</div>
              </div>
            </section>
          </div>

          <aside className="event-detail-sidebar" id="tickets">
            <section className="purchase-panel">
              <div className="purchase-panel-head">
                <div>
                  <h3>{event.tickets[0]?.name || 'Entrada general'}</h3>
                  <p>{event.tickets[0]?.description || 'Compra directa para este evento.'}</p>
                </div>
                <strong>{formatCurrency(event.numericPrice)}</strong>
              </div>

              <div className="purchase-panel-row">
                <span>Cantidad</span>
                <div className="quantity-stepper">
                  <button type="button" onClick={() => setQuantity((current) => Math.max(1, current - 1))}>
                    <MinusIcon />
                  </button>
                  <span>{quantity}</span>
                  <button type="button" onClick={() => setQuantity((current) => current + 1)}>
                    <PlusIcon />
                  </button>
                </div>
              </div>

              <div className="purchase-panel-row total-row">
                <span>Total</span>
                <strong>{formatCurrency(total)}</strong>
              </div>

              <button 
                className="purchase-button" 
                type="button" 
                onClick={handlePurchase}
                disabled={purchasing || event.soldOut}
              >
                {purchasing ? 'Procesando...' : event.soldOut ? 'Agotado' : 'Comprar entrada'}
              </button>

            </section>
          </aside>
        </section>
      </>

      {modal.isOpen && (
        <div className="modal-overlay" onClick={() => setModal({ ...modal, isOpen: false })}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className={`modal-icon-wrap ${modal.type}`}>
              {modal.type === 'success' ? (
                <svg viewBox="0 0 24 24" style={{ width: '36px', height: '36px', stroke: 'currentColor', strokeWidth: '2.5', fill: 'none' }}>
                  <polyline points="20 6 9 17 4 12" />
                </svg>
              ) : (
                <svg viewBox="0 0 24 24" style={{ width: '36px', height: '36px', stroke: 'currentColor', strokeWidth: '2.5', fill: 'none' }}>
                  <line x1="18" y1="6" x2="6" y2="18" />
                  <line x1="6" y1="6" x2="18" y2="18" />
                </svg>
              )}
            </div>
            <h3 className={modal.type}>{modal.title}</h3>
            <p>{modal.message}</p>
            <button 
              type="button" 
              className="modal-button" 
              onClick={() => setModal({ ...modal, isOpen: false })}
            >
              Entendido
            </button>
          </div>
        </div>
      )}
    </main>
  );
}

export default EventDetail;



