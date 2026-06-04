import { useMemo, useState } from 'react';
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
  const navigate = useNavigate();

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
    navigate('/login');
  };

  return (
    <main className="home-screen">
      <header className="event-detail-topbar">
        <Link className="event-detail-brand" to="/">Golden Ticket</Link>
        <nav className="event-detail-nav">
          <Link to="/">Eventos</Link>
          <Link to="/mis-entradas">Mis entradas</Link>
        </nav>
        <div className="event-detail-actions">
          <button type="button" onClick={handleLogout}>Cerrar sesion</button>
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
