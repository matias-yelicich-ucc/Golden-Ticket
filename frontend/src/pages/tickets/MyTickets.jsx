import { useEffect, useState, useMemo, useRef } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import '../../styles/App.css';
import { getMyTickets, transferTicket } from '../../services/api/client';

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

const ChevronIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true" style={{ width: '18px', height: '18px', fill: 'none', stroke: 'currentColor', strokeWidth: '2' }}>
    <path d="m7 10 5 5 5-5" />
  </svg>
);

const DashboardIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true" style={{ width: '18px', height: '18px', fill: 'none', stroke: 'currentColor', strokeWidth: '2' }}>
    <rect x="3" y="3" width="8" height="8" rx="2" />
    <rect x="13" y="3" width="8" height="5" rx="2" />
    <rect x="13" y="10" width="8" height="11" rx="2" />
    <rect x="3" y="13" width="8" height="8" rx="2" />
  </svg>
);

const TicketIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true" style={{ width: '18px', height: '18px', fill: 'none', stroke: 'currentColor', strokeWidth: '2' }}>
    <path d="M4 8.5A2.5 2.5 0 0 1 6.5 6H20v4a2 2 0 0 0 0 4v4H6.5A2.5 2.5 0 0 1 4 15.5v-7Z" />
    <path d="M9 6v12" />
  </svg>
);

const LogoutIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true" style={{ width: '18px', height: '18px', fill: 'none', stroke: 'currentColor', strokeWidth: '2' }}>
    <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
    <path d="M16 17l5-5-5-5" />
    <path d="M21 12H9" />
  </svg>
);

const PlusIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true" style={{ width: '18px', height: '18px', fill: 'none', stroke: 'currentColor', strokeWidth: '2' }}>
    <path d="M12 5v14" />
    <path d="M5 12h14" />
  </svg>
);

function MyTickets() {
  const navigate = useNavigate();
  const [tickets, setTickets] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [menuOpen, setMenuOpen] = useState(false);
  const [selectedTicket, setSelectedTicket] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [transferModal, setTransferModal] = useState({
    isOpen: false,
    ticket: null,
    dni: '',
    loading: false,
    error: '',
    success: '',
  });
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
    const handleClickOutside = (event) => {
      if (profileMenuRef.current && !profileMenuRef.current.contains(event.target)) {
        setMenuOpen(false);
      }
    };
    const handleEscape = (event) => {
      if (event.key === 'Escape') {
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
          <Link to="/mis-entradas" className="is-active">Mis entradas</Link>
        </nav>
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
                  <button 
                    type="button" 
                    style={{ 
                      background: 'linear-gradient(135deg, #d4af37 0%, #b38728 100%)', 
                      color: '#111', 
                      cursor: 'pointer' 
                    }}
                    onClick={() => {
                      setSelectedTicket(ticket);
                      setIsModalOpen(true);
                    }}
                  >
                    Ver Entrada
                  </button>
                  <button 
                    type="button" 
                    disabled={ticket.estado !== 'activo'} 
                    style={{ 
                      opacity: ticket.estado === 'activo' ? 1 : 0.6, 
                      cursor: ticket.estado === 'activo' ? 'pointer' : 'not-allowed',
                      background: 'rgba(255, 255, 255, 0.05)',
                      color: ticket.estado === 'activo' ? '#f6df9c' : '#7f7f7f',
                      border: '1px solid rgba(212, 175, 55, 0.2)'
                    }}
                    onClick={() => {
                      if (ticket.estado === 'activo') {
                        setTransferModal({
                          isOpen: true,
                          ticket,
                          dni: '',
                          loading: false,
                          error: '',
                          success: ''
                        });
                      }
                    }}
                  >
                    Transferir entrada
                  </button>
                  <button type="button" disabled style={{ opacity: 0.6, cursor: 'not-allowed', background: 'rgba(255, 255, 255, 0.05)', color: '#7f7f7f' }}>Cancelar compra</button>
                </div>
              </article>
            ))}
          </div>
        )}
      </section>

      {/* Ticket Pass Modal */}
      {isModalOpen && selectedTicket && (
        <div className="modal-overlay" onClick={() => { setIsModalOpen(false); setSelectedTicket(null); }}>
          <div 
            className="modal-content" 
            onClick={(e) => e.stopPropagation()} 
            style={{ 
              background: 'linear-gradient(135deg, #1f1d1a 0%, #0c0b0a 100%)',
              border: '2px solid #d4af37',
              padding: '28px 24px',
              borderRadius: '20px',
              width: 'min(440px, 95%)',
              position: 'relative',
              boxShadow: '0 24px 70px rgba(0, 0, 0, 0.9)'
            }}
          >
            {/* Punch notches */}
            <div style={{ position: 'absolute', top: '55%', left: '-12px', width: '24px', height: '24px', borderRadius: '50%', background: '#0d0d0d', borderRight: '2px solid #d4af37', zIndex: 10 }} />
            <div style={{ position: 'absolute', top: '55%', right: '-12px', width: '24px', height: '24px', borderRadius: '50%', background: '#0d0d0d', borderLeft: '2px solid #d4af37', zIndex: 10 }} />

            {/* Ticket Header */}
            <div style={{ textTransform: 'uppercase', letterSpacing: '0.2em', fontSize: '0.8rem', color: '#d4af37', fontWeight: '800', marginBottom: '16px' }}>
              ★ Golden Ticket Pass ★
            </div>
            
            {/* Ticket Event Content */}
            <div style={{ textAlign: 'left', marginBottom: '24px' }}>
              <span style={{ 
                background: 'rgba(212, 175, 55, 0.15)', 
                color: '#f6df9c', 
                padding: '4px 10px', 
                borderRadius: '12px', 
                fontSize: '0.75rem', 
                fontWeight: 'bold',
                textTransform: 'uppercase',
                display: 'inline-block',
                marginBottom: '8px'
              }}>
                {selectedTicket.event?.categoria || 'Espectáculo'}
              </span>
              <h2 style={{ fontSize: '1.6rem', color: '#fff', margin: '0 0 12px 0', lineHeight: '1.2', fontWeight: '800' }}>
                {selectedTicket.event?.titulo || 'Evento'}
              </h2>
              
              <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', fontSize: '0.9rem', color: '#c4bcae' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                  <CalendarIcon />
                  <span>{selectedTicket.event?.fecha || 'Fecha pendiente'} — {selectedTicket.event?.hora_inicio || ''}</span>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                  <PinIcon />
                  <span>{selectedTicket.event?.ubicacion || 'Ubicación pendiente'}</span>
                </div>
              </div>
            </div>

            {/* Ticket Stub Separator */}
            <div style={{ 
              borderTop: '2px dashed rgba(212, 175, 55, 0.4)', 
              margin: '20px 0 24px 0',
              position: 'relative'
            }} />

            {/* Stub Content / QR & Code */}
            <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '16px' }}>
              {/* Fake QR code */}
              <div style={{ background: '#fff', padding: '12px', borderRadius: '12px', display: 'inline-block', boxShadow: '0 8px 20px rgba(0,0,0,0.5)' }}>
                <svg viewBox="0 0 100 100" style={{ width: '130px', height: '130px', display: 'block' }}>
                  {/* Finder pattern Top-Left */}
                  <rect x="0" y="0" width="28" height="28" fill="black" />
                  <rect x="4" y="4" width="20" height="20" fill="white" />
                  <rect x="8" y="8" width="12" height="12" fill="black" />
                  
                  {/* Finder pattern Top-Right */}
                  <rect x="72" y="0" width="28" height="28" fill="black" />
                  <rect x="76" y="4" width="20" height="20" fill="white" />
                  <rect x="80" y="8" width="12" height="12" fill="black" />
                  
                  {/* Finder pattern Bottom-Left */}
                  <rect x="0" y="72" width="28" height="28" fill="black" />
                  <rect x="4" y="76" width="20" height="20" fill="white" />
                  <rect x="8" y="80" width="12" height="12" fill="black" />
                  
                  {/* Alignment pattern Bottom-Right */}
                  <rect x="76" y="76" width="8" height="8" fill="black" />
                  <rect x="78" y="78" width="4" height="4" fill="white" />
                  <rect x="79" y="79" width="2" height="2" fill="black" />
                  
                  {/* Random pixels */}
                  <rect x="36" y="4" width="4" height="4" fill="black" />
                  <rect x="48" y="0" width="4" height="4" fill="black" />
                  <rect x="56" y="8" width="4" height="4" fill="black" />
                  <rect x="40" y="16" width="4" height="4" fill="black" />
                  <rect x="64" y="20" width="4" height="4" fill="black" />
                  <rect x="8" y="36" width="4" height="4" fill="black" />
                  <rect x="20" y="44" width="4" height="4" fill="black" />
                  <rect x="0" y="52" width="4" height="4" fill="black" />
                  <rect x="36" y="36" width="8" height="4" fill="black" />
                  <rect x="48" y="40" width="4" height="8" fill="black" />
                  <rect x="60" y="36" width="12" height="4" fill="black" />
                  <rect x="84" y="36" width="4" height="8" fill="black" />
                  <rect x="36" y="52" width="4" height="12" fill="black" />
                  <rect x="48" y="56" width="16" height="4" fill="black" />
                  <rect x="72" y="48" width="8" height="4" fill="black" />
                  <rect x="88" y="56" width="4" height="4" fill="black" />
                  <rect x="36" y="76" width="4" height="4" fill="black" />
                  <rect x="44" y="84" width="8" height="4" fill="black" />
                  <rect x="60" y="72" width="4" height="12" fill="black" />
                  <rect x="56" y="88" width="12" height="4" fill="black" />
                  <rect x="88" y="72" width="4" height="4" fill="black" />
                  <rect x="84" y="88" width="8" height="4" fill="black" />
                </svg>
              </div>

              {/* Barcode SVG */}
              <svg viewBox="0 0 100 20" style={{ width: '180px', height: '32px', display: 'block', opacity: 0.85 }}>
                <rect x="0" y="0" width="2" height="20" fill="white" />
                <rect x="4" y="0" width="1" height="20" fill="white" />
                <rect x="7" y="0" width="3" height="20" fill="white" />
                <rect x="12" y="0" width="1" height="20" fill="white" />
                <rect x="15" y="0" width="2" height="20" fill="white" />
                <rect x="19" y="0" width="4" height="20" fill="white" />
                <rect x="25" y="0" width="1" height="20" fill="white" />
                <rect x="28" y="0" width="2" height="20" fill="white" />
                <rect x="32" y="0" width="3" height="20" fill="white" />
                <rect x="37" y="0" width="1" height="20" fill="white" />
                <rect x="40" y="0" width="2" height="20" fill="white" />
                <rect x="44" y="0" width="4" height="20" fill="white" />
                <rect x="50" y="0" width="1" height="20" fill="white" />
                <rect x="53" y="0" width="3" height="20" fill="white" />
                <rect x="58" y="0" width="2" height="20" fill="white" />
                <rect x="62" y="0" width="1" height="20" fill="white" />
                <rect x="65" y="0" width="4" height="20" fill="white" />
                <rect x="71" y="0" width="2" height="20" fill="white" />
                <rect x="75" y="0" width="1" height="20" fill="white" />
                <rect x="78" y="0" width="3" height="20" fill="white" />
                <rect x="83" y="0" width="2" height="20" fill="white" />
                <rect x="87" y="0" width="4" height="20" fill="white" />
                <rect x="93" y="0" width="1" height="20" fill="white" />
                <rect x="96" y="0" width="3" height="20" fill="white" />
              </svg>

              {/* Ticket metadata */}
              <div style={{ textAlign: 'center', width: '100%' }}>
                <div style={{ fontSize: '0.8rem', color: '#8f8778', textTransform: 'uppercase', letterSpacing: '0.05em', marginBottom: '4px' }}>
                  Código de entrada
                </div>
                <div style={{ fontSize: '1.2rem', color: '#f6df9c', fontWeight: 'bold', marginBottom: '12px' }}>
                  #{selectedTicket.id}
                </div>
                
                <div style={{ fontSize: '0.8rem', color: '#8f8778', textTransform: 'uppercase', letterSpacing: '0.05em', marginBottom: '4px' }}>
                  Titular
                </div>
                <div style={{ fontSize: '1.1rem', color: '#fff', fontWeight: '600', marginBottom: '12px' }}>
                  {[currentUser?.nombre, currentUser?.apellido].filter(Boolean).join(' ') || 'Cliente Golden Ticket'}
                </div>

                <div style={{ 
                  display: 'inline-block', 
                  padding: '4px 12px', 
                  borderRadius: '20px', 
                  fontSize: '0.8rem', 
                  fontWeight: 'bold',
                  background: selectedTicket.estado === 'activo' ? 'rgba(34, 197, 94, 0.15)' : 'rgba(239, 68, 68, 0.15)',
                  color: selectedTicket.estado === 'activo' ? '#22c55e' : '#ef4444',
                  border: selectedTicket.estado === 'activo' ? '1px solid rgba(34, 197, 94, 0.3)' : '1px solid rgba(239, 68, 68, 0.3)'
                }}>
                  {selectedTicket.estado === 'activo' ? 'ACTIVO' : 'CANCELADO'}
                </div>
              </div>
            </div>

            <button 
              type="button" 
              className="modal-button" 
              onClick={() => { setIsModalOpen(false); setSelectedTicket(null); }}
              style={{ marginTop: '24px', background: 'rgba(255, 255, 255, 0.05)', color: '#fff', border: '1px solid rgba(212, 175, 55, 0.3)', cursor: 'pointer' }}
            >
              Cerrar Entrada
            </button>
          </div>
        </div>
      )}

      {/* Transfer Modal */}
      {transferModal.isOpen && transferModal.ticket && (
        <div className="modal-overlay" onClick={() => setTransferModal((prev) => ({ ...prev, isOpen: false, ticket: null }))}>
          <div 
            className="modal-content" 
            onClick={(e) => e.stopPropagation()} 
            style={{ 
              background: '#121212',
              border: '1px solid rgba(212, 175, 55, 0.3)',
              borderRadius: '16px',
              padding: '28px',
              width: 'min(400px, 90%)',
              textAlign: 'center',
              boxShadow: '0 24px 60px rgba(0, 0, 0, 0.8)'
            }}
          >
            <h3 style={{ color: '#fff', fontSize: '1.5rem', marginBottom: '8px' }}>Transferir Entrada</h3>
            <p style={{ color: '#b5ae9d', fontSize: '0.9rem', marginBottom: '20px' }}>
              Vas a transferir la entrada para <strong>{transferModal.ticket.event?.titulo}</strong> (Código #{transferModal.ticket.id}).
            </p>

            {transferModal.error && (
              <div style={{ color: '#ef4444', background: 'rgba(239, 68, 68, 0.1)', border: '1px solid rgba(239, 68, 68, 0.3)', padding: '10px', borderRadius: '8px', fontSize: '0.85rem', marginBottom: '16px' }}>
                {transferModal.error}
              </div>
            )}

            {transferModal.success && (
              <div style={{ color: '#22c55e', background: 'rgba(34, 197, 94, 0.1)', border: '1px solid rgba(34, 197, 94, 0.3)', padding: '10px', borderRadius: '8px', fontSize: '0.85rem', marginBottom: '16px' }}>
                {transferModal.success}
              </div>
            )}

            <form onSubmit={async (e) => {
              e.preventDefault();
              if (!transferModal.dni.trim()) return;

              setTransferModal(prev => ({ ...prev, loading: true, error: '', success: '' }));
              try {
                await transferTicket(transferModal.ticket.id, { dni: transferModal.dni.trim() });
                setTransferModal(prev => ({ ...prev, success: '¡Entrada transferida con éxito!' }));
                
                // Refresh tickets after 1.5 seconds
                setTimeout(() => {
                  setTransferModal(prev => ({ ...prev, isOpen: false, ticket: null }));
                  // Fetch tickets list again
                  setLoading(true);
                  getMyTickets()
                    .then((response) => {
                      setTickets(response.data || []);
                    })
                    .catch((err) => {
                      console.error('Error fetching tickets:', err);
                    })
                    .finally(() => {
                      setLoading(false);
                    });
                }, 1500);
              } catch (err) {
                const errMsg = err.response?.data?.error || 'No se pudo realizar la transferencia.';
                setTransferModal(prev => ({ ...prev, error: errMsg, loading: false }));
              }
            }}>
              <div style={{ display: 'grid', gap: '8px', textAlign: 'left', marginBottom: '20px' }}>
                <label htmlFor="transfer-dni" style={{ color: '#f6df9c', fontSize: '0.85rem', fontWeight: 'bold' }}>DNI del destinatario</label>
                <input 
                  id="transfer-dni" 
                  type="text" 
                  style={{ 
                    background: '#1a1a1a', 
                    border: '1px solid #333', 
                    borderRadius: '8px', 
                    color: '#fff', 
                    padding: '12px',
                    fontSize: '1rem',
                    outline: 'none'
                  }} 
                  value={transferModal.dni} 
                  onChange={(e) => setTransferModal(prev => ({ ...prev, dni: e.target.value }))}
                  placeholder="Ej: 12345678"
                  required
                  disabled={transferModal.loading || !!transferModal.success}
                />
                <span style={{ color: '#8f8778', fontSize: '0.75rem' }}>
                  Asegurate de que el destinatario tenga una cuenta creada con este DNI.
                </span>
              </div>

              <div style={{ display: 'flex', gap: '12px' }}>
                <button 
                  type="button" 
                  className="modal-button" 
                  style={{ background: 'rgba(255, 255, 255, 0.05)', color: '#fff', border: '1px solid #333', cursor: 'pointer' }}
                  onClick={() => setTransferModal((prev) => ({ ...prev, isOpen: false, ticket: null }))}
                  disabled={transferModal.loading}
                >
                  Cancelar
                </button>
                <button 
                  type="submit" 
                  className="modal-button" 
                  disabled={transferModal.loading || !transferModal.dni.trim() || !!transferModal.success}
                  style={{ cursor: 'pointer' }}
                >
                  {transferModal.loading ? 'Transfiriendo...' : 'Transferir'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </main>
  );
}

export default MyTickets;
