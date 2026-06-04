import { useEffect, useMemo, useState } from 'react';
import { Link, Navigate, useParams } from 'react-router-dom';
import '../styles/App.css';
import { getEventBySlug } from '../constants/events';

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

function EventDetail() {
  const { slug } = useParams();
  const [quantity, setQuantity] = useState(1);
  const [feedback, setFeedback] = useState('');
  const event = useMemo(() => getEventBySlug(slug) || null, [slug]);
  const isAuthenticated = Boolean(localStorage.getItem('token'));

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [slug]);

  if (!event) {
    return <Navigate to="/" replace />;
  }

  const total = event.numericPrice * quantity;

  const handlePurchase = () => {
    setFeedback('La pantalla de compra ya esta lista en frontend. La integracion de tickets queda pendiente del backend de eventos.');
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
        <section className={`event-detail-hero event-media-${event.image || 'concert'}`}>
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

              <button className="purchase-button" type="button" onClick={handlePurchase}>
                Comprar entrada
              </button>

              {feedback && <p className="purchase-note">{feedback}</p>}
            </section>
          </aside>
        </section>
      </>
    </main>
  );
}

export default EventDetail;
