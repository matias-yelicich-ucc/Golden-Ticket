import { useEffect, useMemo, useRef, useState } from 'react';
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

const ChevronIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="m7 10 5 5 5-5" />
  </svg>
);

const DashboardIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <rect x="3" y="3" width="8" height="8" rx="2" />
    <rect x="13" y="3" width="8" height="5" rx="2" />
    <rect x="13" y="10" width="8" height="11" rx="2" />
    <rect x="3" y="13" width="8" height="8" rx="2" />
  </svg>
);

const TicketIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M4 8.5A2.5 2.5 0 0 1 6.5 6H20v4a2 2 0 0 0 0 4v4H6.5A2.5 2.5 0 0 1 4 15.5v-7Z" />
    <path d="M9 6v12" />
  </svg>
);

const LogoutIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
    <path d="M16 17l5-5-5-5" />
    <path d="M21 12H9" />
  </svg>
);

const ArrowLeftIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M19 12H5" />
    <path d="m12 19-7-7 7-7" />
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
    description: [dbEvent.descripcion || 'Sin descripción.'],
    category: dbEvent.categoria || 'Musica',
    date: dbEvent.fecha,
    fullDate: dbEvent.fecha,
    timeRange: `${dbEvent.hora_inicio} - ${dbEvent.hora_fin}`,
    location: dbEvent.ubicacion,
    address: dbEvent.ubicacion,
    capacity: dbEvent.capacidad,
    entradas_disponibles: dbEvent.entradas_disponibles,
    soldOut: isSoldOut,
    badge: isSoldOut ? 'AGOTADO' : '',
    price: formatCurrency(dbEvent.precio || 0),
    numericPrice: dbEvent.precio || 0,
    tickets: [
      {
        name: 'Entrada General',
        description: 'Acceso general al evento.',
      },
    ],
    urlImagen: dbEvent.url_imagen,
  };
};

function EventDetail() {
  const { slug } = useParams();
  const navigate = useNavigate();
  const [quantity, setQuantity] = useState(1);
  const [purchasing, setPurchasing] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);
  const [modal, setModal] = useState({ isOpen: false, type: '', title: '', message: '' });
  const [event, setEvent] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const profileMenuRef = useRef(null);
  const isAuthenticated = Boolean(localStorage.getItem('token'));

  const currentUser = useMemo(() => {
    if (!isAuthenticated) return null;
    try {
      return JSON.parse(localStorage.getItem('user') || 'null');
    } catch {
      return null;
    }
  }, [isAuthenticated]);

  const userInitials = useMemo(() => {
    if (!currentUser) return 'GT';
    const rawName = [currentUser.nombre, currentUser.apellido].filter(Boolean).join(' ').trim();
    if (!rawName) return 'GT';
    return rawName
      .split(/\s+/)
      .slice(0, 2)
      .map((part) => part[0]?.toUpperCase() || '')
      .join('');
  }, [currentUser]);

  const isAdmin = ['admin', 'administrador'].includes(currentUser?.rol);

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setMenuOpen(false);
    navigate('/login');
  };

  useEffect(() => {
    if (!menuOpen) return undefined;
    const handleClickOutside = (clickEvent) => {
      if (profileMenuRef.current && !profileMenuRef.current.contains(clickEvent.target)) {
        setMenuOpen(false);
      }
    };
    const handleEscape = (keyEvent) => {
      if (keyEvent.key === 'Escape') {
        setMenuOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    document.addEventListener('keydown', handleEscape);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
      document.removeEventListener('keydown', handleEscape);
    };
  }, [menuOpen]);

  useEffect(() => {
    const loadEvent = async () => {
      const eventId = slug.split('-').pop();
      if (!eventId || Number.isNaN(Number(eventId))) {
        setError('ID de evento invalido.');
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
      } catch (requestError) {
        console.error('Error fetching event detail:', requestError);
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

  useEffect(() => {
    if (!event) return;

    if (event.entradas_disponibles <= 0) {
      setQuantity(0);
      return;
    }

    setQuantity((current) => {
      if (current <= 0) return 1;
      return Math.min(current, event.entradas_disponibles);
    });
  }, [event]);

  if (loading) {
    return (
      <main className="event-detail-page">
        <header className="event-detail-topbar">
          <div className="topbar-title-group">
            <button
              type="button"
              className="topbar-back-button"
              onClick={() => navigate(-1)}
              aria-label="Volver"
            >
              <ArrowLeftIcon />
            </button>
            <h1 className="topbar-page-title">Detalle del evento</h1>
          </div>
        </header>
        <p className="empty-state">Cargando detalles del evento...</p>
      </main>
    );
  }

  if (error || !event) {
    return <Navigate to="/" replace />;
  }

  const reachedAvailabilityLimit = !event.soldOut && quantity >= event.entradas_disponibles;
  const quantityWarning = event.soldOut
    ? 'No quedan entradas disponibles para este evento.'
    : reachedAvailabilityLimit
      ? 'Ya alcanzaste el maximo de entradas disponibles para este evento.'
      : '';
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
          title: 'Compra exitosa',
          message: `Adquiriste ${quantity} entrada(s) para "${event.title}".`,
        });
        setEvent((previousEvent) => ({
          ...previousEvent,
          entradas_disponibles: Math.max(0, previousEvent.entradas_disponibles - quantity),
          soldOut: previousEvent.entradas_disponibles - quantity <= 0,
        }));
      })
      .catch((requestError) => {
        console.error('Error buying tickets:', requestError);
        const backendError = requestError.response?.data?.error || 'No se pudo completar la compra.';
        setModal({
          isOpen: true,
          type: 'error',
          title: 'Error en la compra',
          message: backendError,
        });
      })
      .finally(() => {
        setPurchasing(false);
      });
  };

  return (
    <main className="event-detail-page">
      <header className="event-detail-topbar">
        <div className="topbar-title-group">
          <button
            type="button"
            className="topbar-back-button"
            onClick={() => navigate(-1)}
            aria-label="Volver"
          >
            <ArrowLeftIcon />
          </button>
          <h1 className="topbar-page-title">Detalle del evento</h1>
        </div>
        <div className="event-detail-actions">
          {isAuthenticated ? (
            <div className="topbar-profile" ref={profileMenuRef}>
              <button
                type="button"
                className="topbar-profile-trigger"
                onClick={() => setMenuOpen((current) => !current)}
                aria-expanded={menuOpen}
                aria-haspopup="menu"
              >
                <span className="topbar-profile-avatar">{userInitials}</span>
                <span className="topbar-profile-name">{currentUser?.nombre || 'Mi perfil'}</span>
                <ChevronIcon />
              </button>

              {menuOpen && (
                <div className="topbar-profile-menu" role="menu">
                  {isAdmin && (
                    <Link to="/admin/dashboard" onClick={() => setMenuOpen(false)}>
                      <DashboardIcon />
                      Ir al panel admin
                    </Link>
                  )}
                  {isAdmin && (
                    <Link to="/admin/create-event" onClick={() => setMenuOpen(false)}>
                      <PlusIcon />
                      Crear evento
                    </Link>
                  )}
                  <Link to="/mis-entradas" onClick={() => setMenuOpen(false)}>
                    <TicketIcon />
                    Ver mis entradas
                  </Link>
                  <button type="button" onClick={handleLogout}>
                    <LogoutIcon />
                    Cerrar sesion
                  </button>
                </div>
              )}
            </div>
          ) : (
            <>
              <button type="button" className="event-detail-login" onClick={() => navigate('/login')}>Iniciar sesion</button>
              <button type="button" onClick={() => navigate('/register')}>Crear cuenta</button>
            </>
          )}
        </div>
      </header>

      <section
        className="event-detail-hero event-media-concert"
        style={event.urlImagen ? {
          backgroundImage: `url(${event.urlImagen})`,
          backgroundSize: 'cover',
          backgroundPosition: 'center',
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
                <UsersIcon />
              </div>
              <div className="event-detail-stats-copy">
                <span>Disponibilidad</span>
                <p><span>Capacidad:</span> <strong>{event.capacity} personas</strong></p>
                <p><span>Entradas disponibles:</span> <strong>{event.entradas_disponibles}</strong></p>
              </div>
            </article>
          </div>

          <section className="event-detail-copy">
            <h2>Sobre este evento</h2>
            {event.description.map((paragraph) => <p key={paragraph}>{paragraph}</p>)}
            <p>Este evento forma parte del catalogo cliente definido por la consigna y permite completar el flujo de consulta, detalle y compra.</p>
          </section>

          <section className="event-detail-capacity">
            <h2>Ubicacion</h2>
            <p>{event.location}</p>

            <div
              className="event-detail-map"
              style={{
                marginTop: '1.5rem',
                borderRadius: '12px',
                overflow: 'hidden',
                height: '300px',
                backgroundColor: '#eaeaea',
              }}
            >
              <iframe
                title={`Mapa de ubicacion de ${event.location}`}
                width="100%"
                height="100%"
                style={{ border: 0 }}
                loading="lazy"
                allowFullScreen
                referrerPolicy="no-referrer-when-downgrade"
                src={`https://maps.google.com/maps?q=${encodeURIComponent(event.location)}&t=&z=15&ie=UTF8&iwloc=&output=embed`}
              />
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
                <button
                  type="button"
                  disabled={event.soldOut || quantity <= 1}
                  onClick={() => setQuantity((current) => Math.max(1, current - 1))}
                >
                  <MinusIcon />
                </button>
                <span>{quantity}</span>
                <button
                  type="button"
                  disabled={event.soldOut || quantity >= event.entradas_disponibles}
                  onClick={() => setQuantity((current) => Math.min(event.entradas_disponibles, current + 1))}
                >
                  <PlusIcon />
                </button>
              </div>
            </div>

            {quantityWarning && (
              <p className="purchase-note purchase-note-limit">
                {quantityWarning}
              </p>
            )}

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

      {modal.isOpen && (
        <div className="modal-overlay" onClick={() => setModal({ ...modal, isOpen: false })}>
          <div className="modal-content" onClick={(clickEvent) => clickEvent.stopPropagation()}>
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
