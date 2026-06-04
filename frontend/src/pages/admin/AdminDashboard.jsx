import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import '../../styles/App.css';
import { adminMetrics } from '../../constants/admin';
import { getEvents } from '../../services/api/client';

const iconMap = {
  calendar: (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M7 3v4M17 3v4M4 9h16M5 5h14a1 1 0 0 1 1 1v14H4V6a1 1 0 0 1 1-1Z" />
    </svg>
  ),
  ticket: (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M4 8.5A2.5 2.5 0 0 1 6.5 6H20v4a2 2 0 0 0 0 4v4H6.5A2.5 2.5 0 0 1 4 15.5v-7Z" />
      <path d="M9 6v12" />
    </svg>
  ),
  users: (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M16 21v-2a4 4 0 0 0-4-4H7a4 4 0 0 0-4 4v2" />
      <circle cx="9.5" cy="7" r="3" />
      <path d="M22 21v-2a4 4 0 0 0-3-3.87" />
      <path d="M16 4.13a4 4 0 0 1 0 7.75" />
    </svg>
  ),
  wallet: (
    <svg viewBox="0 0 24 24" aria-hidden="true">
      <path d="M4 7.5A2.5 2.5 0 0 1 6.5 5H19a1 1 0 0 1 1 1v3H15a2 2 0 0 0 0 4h5v5a1 1 0 0 1-1 1H6.5A2.5 2.5 0 0 1 4 16.5v-9Z" />
      <path d="M15 13h6" />
      <circle cx="16.5" cy="13" r=".5" fill="currentColor" stroke="none" />
    </svg>
  ),
};

const BarsIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M5 21V10" />
    <path d="M12 21V4" />
    <path d="M19 21v-7" />
  </svg>
);

const PencilIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="m12 20 8-8" />
    <path d="M18 6a2.8 2.8 0 1 1 4 4L10 22l-5 1 1-5L18 6Z" />
  </svg>
);

const TrashIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M3 6h18" />
    <path d="M8 6V4h8v2" />
    <path d="M19 6v14H5V6" />
    <path d="M10 11v6" />
    <path d="M14 11v6" />
  </svg>
);

function AdminDashboard() {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    getEvents()
      .then((response) => {
        setEvents(response.data || []);
      })
      .catch((err) => {
        console.error('Error fetching events:', err);
        setError('No se pudieron cargar los eventos del servidor.');
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  return (
    <main className="admin-page">
      <header className="event-detail-topbar">
        <Link className="event-detail-brand" to="/">Golden Ticket</Link>
        <nav className="event-detail-nav">
          <Link to="/">Eventos</Link>
          <Link to="/admin" className="is-active">Panel admin</Link>
          <Link to="/mis-entradas">Mis entradas</Link>
        </nav>
        <div className="event-detail-actions">
          <Link className="event-detail-login" to="/admin/eventos/nuevo">Crear evento</Link>
        </div>
      </header>

      <section className="admin-shell">
        <div className="admin-heading">
          <div>
            <span className="admin-kicker">Panel de administracion</span>
            <h1>Controla tu calendario de eventos y sus metricas clave.</h1>
            <p>Gestiona los proximos lanzamientos, capacidad y ventas en tiempo real con datos de tu base de datos.</p>
          </div>
          <Link className="admin-primary-link" to="/admin/eventos/nuevo">+ Crear evento</Link>
        </div>

        <section className="admin-metrics-grid">
          {adminMetrics.map((metric) => (
            <article className="admin-metric-card" key={metric.id}>
              <div className="admin-metric-top">
                <span className="admin-metric-icon">{iconMap[metric.icon]}</span>
                <span className={`admin-delta admin-delta-${metric.tone}`}>{metric.delta}</span>
              </div>
              <span className="admin-metric-label">{metric.label}</span>
              <strong>{metric.value}</strong>
            </article>
          ))}
        </section>

        <section className="admin-table-card">
          <div className="admin-table-header">
            <div>
              <h2>Eventos proximos</h2>
              <p>Listado en tiempo real desde la base de datos.</p>
            </div>
            <Link to="/admin/eventos/nuevo">Ver todos</Link>
          </div>

          <div className="admin-table">
            <div className="admin-table-head">
              <span>Titulo</span>
              <span>Fecha</span>
              <span>Cupo</span>
              <span>Vendidas</span>
              <span>Estado</span>
              <span>Acciones</span>
            </div>

            {loading && <p className="empty-state">Cargando eventos...</p>}
            {error && <p className="empty-state" style={{ color: '#ff6b6b' }}>{error}</p>}
            {!loading && !error && events.length === 0 && (
              <p className="empty-state">No hay eventos creados. ¡Crea el primero para comenzar!</p>
            )}

            {!loading && !error && events.map((event) => {
              const sold = event.capacidad - event.entradas_disponibles;
              const fill = event.capacidad > 0 ? (sold / event.capacidad) * 100 : 0;
              const isSoldOut = event.entradas_disponibles === 0;
              const status = isSoldOut ? 'Agotado' : 'Activo';
              const statusTone = isSoldOut ? 'danger' : 'success';

              return (
                <article className="admin-table-row" key={event.id}>
                  <div className="admin-event-main">
                    <div 
                      className="admin-event-thumb" 
                      style={event.url_imagen ? { 
                        backgroundImage: `url(${event.url_imagen})`, 
                        backgroundSize: 'cover', 
                        backgroundPosition: 'center' 
                      } : {}}
                    />
                    <div>
                      <strong>{event.titulo}</strong>
                      <p>{event.ubicacion}</p>
                    </div>
                  </div>
                  <span className="admin-table-cell">{event.fecha} {event.hora_inicio}</span>
                  <span className="admin-table-cell">{event.capacidad}</span>
                  <div className="admin-sales-cell">
                    <span>{sold}</span>
                    <div className="admin-progress">
                      <span style={{ width: `${fill}%` }} />
                    </div>
                  </div>
                  <span className={`admin-status admin-status-${statusTone}`}>{status}</span>
                  <div className="admin-actions-cell">
                    <button type="button" aria-label={`Ver metricas de ${event.titulo}`}>
                      <BarsIcon />
                    </button>
                    <button type="button" aria-label={`Editar ${event.titulo}`}>
                      <PencilIcon />
                    </button>
                    <button type="button" aria-label={`Eliminar ${event.titulo}`}>
                      <TrashIcon />
                    </button>
                  </div>
                </article>
              );
            })}
          </div>
        </section>
      </section>
    </main>
  );
}

export default AdminDashboard;
