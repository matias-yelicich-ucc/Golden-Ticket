import { Link } from 'react-router-dom';
import { CalendarIcon, PinIcon } from '../common/Icons';

const EventCard = ({ event }) => (
  <Link className="event-card-link" to={`/eventos/${event.slug}`}>
    <article className="event-card">
      <div 
        className="event-media event-media-concert"
        style={event.urlImagen ? { 
          backgroundImage: `url(${event.urlImagen})`, 
          backgroundSize: 'cover', 
          backgroundPosition: 'center' 
        } : {}}
      >
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

export default EventCard;
