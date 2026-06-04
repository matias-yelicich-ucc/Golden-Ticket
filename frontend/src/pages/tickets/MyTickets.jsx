import { Link } from 'react-router-dom';
import '../../styles/App.css';

function MyTickets() {
  return (
    <main className="event-detail-page">
      <header className="event-detail-topbar">
        <Link className="event-detail-brand" to="/">Golden Ticket</Link>
        <nav className="event-detail-nav">
          <Link to="/">Eventos</Link>
          <Link to="/mis-entradas">Mis entradas</Link>
        </nav>
      </header>

      <section className="events-section">
        <h2>Mis entradas</h2>
        <p className="home-intro-copy">Esta pantalla queda reservada para la siguiente integracion del backend de tickets.</p>
        <p className="empty-state">La rama unificada conserva el backend de autenticacion y tus pantallas principales de cliente: login, home y detalle de evento.</p>
      </section>
    </main>
  );
}

export default MyTickets;
