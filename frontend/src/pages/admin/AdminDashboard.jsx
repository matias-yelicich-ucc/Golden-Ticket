import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import '../../styles/App.css';
import { getAdminDashboardStats, getEvents, deleteEvent } from '../../services/api/client';

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

const buildMetrics = (stats) => [
  {
    id: 'events',
    label: 'Total eventos',
    value: String(stats.total_eventos ?? 0),
    icon: 'calendar',
  },
  {
    id: 'tickets',
    label: 'Entradas vendidas',
    value: new Intl.NumberFormat('es-AR').format(stats.entradas_vendidas ?? 0),
    icon: 'ticket',
  },
  {
    id: 'occupancy',
    label: 'Ocupacion media',
    value: `${Math.round(stats.ocupacion_media ?? 0)}%`,
    icon: 'users',
  },
  {
    id: 'revenue',
    label: 'Recaudacion',
    value: new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: 'ARS',
      maximumFractionDigits: 0,
    }).format(stats.recaudacion_total ?? 0),
    icon: 'wallet',
  },
];

function AdminDashboard() {
  const [events, setEvents] = useState([]);
  const [metrics, setMetrics] = useState(buildMetrics({}));
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const navigate = useNavigate();
  const [deleteModal, setDeleteModal] = useState({
    isOpen: false,
    event: null,
    loading: false,
    error: '',
  });

  useEffect(() => {
    getAdminDashboardStats()
      .then((response) => {
        setMetrics(buildMetrics(response.data || {}));
      })
      .catch((err) => {
        console.error('Error fetching admin dashboard stats:', err);
        setMetrics(buildMetrics({}));
      });

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

  useEffect(() => {
    if (!deleteModal.isOpen) return undefined;

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
  }, [deleteModal.isOpen]);

  const handleDeleteClick = (event) => {
    setDeleteModal({ isOpen: true, event, loading: false, error: '' });
  };

  const handleDeleteConfirm = async () => {
    if (!deleteModal.event) return;
    setDeleteModal((prev) => ({ ...prev, loading: true, error: '' }));
    try {
      await deleteEvent(deleteModal.event.id);
      setEvents((prev) => prev.filter((eventItem) => eventItem.id !== deleteModal.event.id));
      setDeleteModal({ isOpen: false, event: null, loading: false, error: '' });
    } catch (err) {
      const msg = err.response?.data?.error || 'Error al eliminar el evento.';
      setDeleteModal((prev) => ({ ...prev, loading: false, error: msg }));
    }
  };

  const handleDeleteClose = () => {
    if (!deleteModal.loading) {
      setDeleteModal({ isOpen: false, event: null, loading: false, error: '' });
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    navigate('/login');
  };

  return (
    <>
      <main className="admin-page">
        <header className="event-detail-topbar">
          <Link className="event-detail-brand" to="/">Golden Ticket</Link>
          <div className="event-detail-actions">
            <Link className="admin-topbar-cta" to="/admin/eventos/nuevo">+ Crear evento</Link>
            <button type="button" onClick={handleLogout}>Cerrar sesion</button>
          </div>
        </header>

        <section className="admin-shell">
          <div className="admin-heading">
            <div>
              <span className="admin-kicker">Panel de administracion</span>
            </div>
          </div>

          <section className="admin-metrics-grid">
            {metrics.map((metric) => (
              <article className="admin-metric-card" key={metric.id}>
                <div className="admin-metric-top">
                  <span className="admin-metric-icon">{iconMap[metric.icon]}</span>
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
              </div>
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
                <div className="admin-empty-state">
                  <p>No hay eventos creados. Crea el primero para comenzar.</p>
                </div>
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
                          backgroundPosition: 'center',
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
                      <Link className="admin-action-button" to={`/admin/eventos/${event.id}/editar`} aria-label={`Editar ${event.titulo}`}>
                        <PencilIcon />
                      </Link>
                      <button className="admin-action-button" type="button" aria-label={`Eliminar ${event.titulo}`} onClick={() => handleDeleteClick(event)}>
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

      {deleteModal.isOpen && deleteModal.event && (
        <div className="modal-overlay" onClick={handleDeleteClose}>
          <div
            className="modal-content"
            onClick={(event) => event.stopPropagation()}
            style={{
              background: 'linear-gradient(135deg, #1f1d1a 0%, #0c0b0a 100%)',
              border: '2px solid #d4af37',
              padding: '32px 28px',
              borderRadius: '20px',
              width: 'min(460px, 95%)',
              boxShadow: '0 24px 70px rgba(0, 0, 0, 0.9)',
              textAlign: 'center',
            }}
          >
            <div style={{ marginBottom: '20px' }}>
              <div style={{
                width: '56px',
                height: '56px',
                borderRadius: '50%',
                background: 'rgba(239, 68, 68, 0.15)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                margin: '0 auto 16px',
                border: '1px solid rgba(239, 68, 68, 0.3)',
              }}
              >
                <svg viewBox="0 0 24 24" style={{ width: '28px', height: '28px', fill: 'none', stroke: '#ef4444', strokeWidth: '2' }}>
                  <path d="M3 6h18" />
                  <path d="M8 6V4h8v2" />
                  <path d="M19 6v14H5V6" />
                  <path d="M10 11v6" />
                  <path d="M14 11v6" />
                </svg>
              </div>
              <h3 style={{ color: '#fff', fontSize: '1.3rem', marginBottom: '8px' }}>Cancelar evento</h3>
              <p style={{ color: 'rgba(255,255,255,0.6)', fontSize: '0.95rem', lineHeight: '1.5' }}>
                Estas por cancelar <strong style={{ color: '#d4af37' }}>{deleteModal.event.titulo}</strong>.
                Todas las entradas activas seran canceladas y el dinero sera reintegrado a los compradores.
                Esta accion no se puede deshacer.
              </p>
            </div>

            {deleteModal.error && (
              <p style={{ color: '#ef4444', fontSize: '0.9rem', marginBottom: '12px' }}>{deleteModal.error}</p>
            )}

            <div style={{ display: 'flex', gap: '12px', justifyContent: 'center' }}>
              <button
                type="button"
                onClick={handleDeleteClose}
                disabled={deleteModal.loading}
                style={{
                  padding: '10px 24px',
                  borderRadius: '10px',
                  border: '1px solid rgba(255,255,255,0.1)',
                  background: 'rgba(255,255,255,0.05)',
                  color: '#fff',
                  cursor: deleteModal.loading ? 'not-allowed' : 'pointer',
                  fontSize: '0.95rem',
                }}
              >
                Volver
              </button>
              <button
                type="button"
                onClick={handleDeleteConfirm}
                disabled={deleteModal.loading}
                style={{
                  padding: '10px 24px',
                  borderRadius: '10px',
                  border: 'none',
                  background: deleteModal.loading ? '#666' : 'linear-gradient(135deg, #ef4444 0%, #dc2626 100%)',
                  color: '#fff',
                  cursor: deleteModal.loading ? 'not-allowed' : 'pointer',
                  fontWeight: 'bold',
                  fontSize: '0.95rem',
                }}
              >
                {deleteModal.loading ? 'Cancelando...' : 'Si, cancelar evento'}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}

export default AdminDashboard;
