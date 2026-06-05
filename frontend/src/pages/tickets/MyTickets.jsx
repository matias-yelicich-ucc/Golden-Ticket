import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import '../../styles/App.css';
import { getMyTickets } from '../../services/api/client';

const CalendarIcon = () => (
  <svg viewBox="0 0 24 24" style={{ width: '18px', height: '18px', fill: 'none', stroke: 'currentColor', strokeWidth: '2' }}>
    <path d="M7 3v4M17 3v4M4 9h16M5 5h14a1 1 0 0 1 1 1v14H4V6a1 1 0 0 1 1-1Z" />
  </svg>
);

const PinIcon = () => (
  <svg viewBox="0 0 24 24" style={{ width: '18px', height: '18px', fill: 'none', stroke: 'currentColor', strokeWidth: '2' }}>
    <path d="M12 21s7-6.1 7-12a7 7 0 0 0-14 0c0 5.9 7 12 7 12Z" />
    <circle cx="12" cy="9" r="2.3" />
  </svg>
);

function MyTickets() {
  const [tickets, setTickets] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    getMyTickets()
      .then((response) => {
        setTickets(response.data || []);
      })
      .catch((err) => {
        console.error('Error fetching tickets:', err);
        setError('No se pudieron cargar tus entradas. Intenta de nuevo más tarde.');
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  return (
    <main className="home-screen">
      <header className="event-detail-topbar">
        <Link className="event-detail-brand" to="/">Golden Ticket</Link>
        <nav className="event-detail-nav">
          <Link to="/">Eventos</Link>
          <Link to="/mis-entradas">Mis entradas</Link>
        </nav>
      </header>

      <section className="events-section" style={{ padding: '40px 24px', maxWidth: '1100px', margin: '0 auto' }}>
        <h2 style={{ fontSize: '2.2rem', marginBottom: '24px', color: '#fff' }}>Mis entradas</h2>

        {loading ? (
          <p className="empty-state">Cargando tus entradas...</p>
        ) : error ? (
          <p className="empty-state" style={{ color: '#ef4444' }}>{error}</p>
        ) : tickets.length === 0 ? (
          <div style={{ textAlign: 'center', padding: '60px 20px' }}>
            <p className="empty-state" style={{ marginBottom: '20px' }}>No tenés entradas adquiridas todavía.</p>
            <Link to="/" className="modal-button" style={{ display: 'inline-block', width: 'auto', textDecoration: 'none' }}>
              Explorar eventos
            </Link>
          </div>
        ) : (
          <div className="tickets-list">
            {tickets.map((ticket) => (
              <article key={ticket.id} className="ticket-history-card">
                {/* Column 1: Event Info */}
                <div>
                  <span className="category-pill" style={{ display: 'inline-block', marginBottom: '8px' }}>
                    {ticket.event?.categoria || 'Espectáculo'}
                  </span>
                  <h3>{ticket.event?.titulo || 'Evento Desconocido'}</h3>
                  <div className="event-date" style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '6px' }}>
                    <CalendarIcon />
                    <span>{ticket.event?.fecha || 'Fecha pendiente'} — {ticket.event?.hora_inicio || ''}</span>
                  </div>
                  <div className="event-location" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                    <PinIcon />
                    <span>{ticket.event?.ubicacion || 'Ubicación pendiente'}</span>
                  </div>
                </div>

                {/* Column 2: Ticket details */}
                <div className="ticket-history-meta">
                  <span>Código de entrada</span>
                  <strong>#{ticket.id}</strong>
                  <span style={{ 
                    display: 'inline-block', 
                    padding: '4px 10px', 
                    borderRadius: '20px', 
                    fontSize: '0.85rem', 
                    fontWeight: 'bold', 
                    textAlign: 'center',
                    background: ticket.estado === 'activo' ? 'rgba(34, 197, 94, 0.15)' : 'rgba(239, 68, 68, 0.15)',
                    color: ticket.estado === 'activo' ? '#22c55e' : '#ef4444',
                    border: ticket.estado === 'activo' ? '1px solid rgba(34, 197, 94, 0.3)' : '1px solid rgba(239, 68, 68, 0.3)'
                  }}>
                    {ticket.estado === 'activo' ? 'Activo' : 'Cancelado'}
                  </span>
                </div>

                {/* Column 3: Actions */}
                <div className="ticket-history-actions">
                  <button type="button" disabled style={{ opacity: 0.6, cursor: 'not-allowed' }}>Transferir entrada</button>
                  <button type="button" disabled style={{ opacity: 0.6, cursor: 'not-allowed', background: 'rgba(255, 255, 255, 0.05)', color: '#7f7f7f' }}>Cancelar compra</button>
                </div>
              </article>
            ))}
          </div>
        )}
      </section>
    </main>
  );
}

export default MyTickets;
