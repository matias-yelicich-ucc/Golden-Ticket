import { useEffect, useMemo, useRef, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import '../styles/App.css';
import { events } from '../constants/events';

const SearchIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <circle cx="11" cy="11" r="7" />
    <path d="m20 20-4.2-4.2" />
  </svg>
);

const PinIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M12 21s7-6.1 7-12a7 7 0 0 0-14 0c0 5.9 7 12 7 12Z" />
    <circle cx="12" cy="9" r="2.3" />
  </svg>
);

const CalendarIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M7 3v4M17 3v4M4 9h16M5 5h14a1 1 0 0 1 1 1v14H4V6a1 1 0 0 1 1-1Z" />
  </svg>
);

const ChevronIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="m7 10 5 5 5-5" />
  </svg>
);

const UserIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M18 21v-2a6 6 0 0 0-12 0v2" />
    <circle cx="12" cy="8" r="4" />
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

const PlusIcon = () => (
  <svg viewBox="0 0 24 24" aria-hidden="true">
    <path d="M12 5v14" />
    <path d="M5 12h14" />
  </svg>
);

const EventCard = ({ event }) => (
  <Link className="event-card-link" to={`/eventos/${event.slug}`}>
    <article className="event-card">
      <div className={`event-media event-media-${event.image || 'concert'}`}>
        <span className="category-pill">{event.category}</span>
        {(event.badge || event.soldOut) && <span className="status-pill">{event.badge || 'AGOTADO'}</span>}
      </div>
      <div className="event-content">
        <div className="event-date">
          <CalendarIcon />
          {event.date}
        </div>
        <h3>{event.title}</h3>
        <div className="event-location">
          <PinIcon />
          {event.location}
        </div>
        <div className="event-divider" />
        <div className="event-footer">
          <div>
            <span>DESDE</span>
            <strong>{event.price}</strong>
          </div>
          <span className="event-detail-cta">Ver detalle</span>
        </div>
      </div>
    </article>
  </Link>
);

function HelloWorld() {
  const [search, setSearch] = useState('');
  const [category, setCategory] = useState('Todos');
  const [menuOpen, setMenuOpen] = useState(false);
  const navigate = useNavigate();
  const isAuthenticated = Boolean(localStorage.getItem('token'));
  const profileMenuRef = useRef(null);

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

  const categories = useMemo(() => {
    const dynamicCategories = Array.from(new Set(events.map((event) => event.category)));
    return ['Todos', ...dynamicCategories];
  }, []);

  const filteredEvents = useMemo(() => {
    return events.filter((event) => {
      const matchesCategory = category === 'Todos' || event.category === category;
      const matchesSearch = !search || `${event.title} ${event.location}`.toLowerCase().includes(search.toLowerCase());
      return matchesCategory && matchesSearch;
    });
  }, [category, search]);

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

  return (
    <main className="home-screen">
      <header className="event-detail-topbar">
        <Link className="event-detail-brand" to="/">Golden Ticket</Link>
        <nav className="event-detail-nav">
          <Link to="/">Eventos</Link>
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

      <section className="home-hero">
        <div className="hero-content">
          <h1>Encontra los mejores eventos cerca tuyo</h1>
          <p>Descubri conciertos, charlas y espectaculos. Busca, filtra y compra tus entradas desde un solo lugar.</p>
          <form className="hero-search" onSubmit={(event) => event.preventDefault()}>
            <label className="search-field">
              <SearchIcon />
              <input type="search" placeholder="Que queres ver?" value={search} onChange={(e) => setSearch(e.target.value)} />
            </label>
            <button className="location-field" type="button">
              <PinIcon />
              <span>Ciudad de Cordoba</span>
              <ChevronIcon />
            </button>
            <button className="search-button" type="submit">Buscar</button>
          </form>
        </div>
      </section>

      <section className="filters-bar" aria-label="Filtros de eventos">
        <div className="category-filters">
          {categories.map((item) => (
            <button className={item === category ? 'active' : ''} key={item} type="button" onClick={() => setCategory(item)}>
              {item}
            </button>
          ))}
        </div>
      </section>

      <section className="events-section">
        <h2>Catalogo de eventos</h2>
        {filteredEvents.length === 0 && <p className="empty-state">No hay eventos para ese filtro.</p>}
        <div className="events-grid">
          {filteredEvents.map((event) => (
            <EventCard event={event} key={event.id} />
          ))}
        </div>
      </section>
    </main>
  );
}

export default HelloWorld;
