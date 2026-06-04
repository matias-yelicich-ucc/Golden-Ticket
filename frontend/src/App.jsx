import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Login from './pages/auth/Login';
import Register from './pages/auth/Register';
import HomePage from './pages/home/HomePage';
import AdminCreateEvent from './pages/admin/AdminCreateEvent';
import AdminDashboard from './pages/admin/AdminDashboard';
import EventDetail from './pages/events/EventDetail';
import MyTickets from './pages/tickets/MyTickets';

const hasToken = () => Boolean(localStorage.getItem('token'));

const RequireAuth = ({ children }) => (hasToken() ? children : <Navigate to="/login" replace />);

const PublicOnly = ({ children }) => (hasToken() ? <Navigate to="/" replace /> : children);

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/eventos/:slug" element={<EventDetail />} />
        <Route path="/mis-entradas" element={<RequireAuth><MyTickets /></RequireAuth>} />
        <Route path="/admin" element={<RequireAuth><AdminDashboard /></RequireAuth>} />
        <Route path="/admin/dashboard" element={<RequireAuth><AdminDashboard /></RequireAuth>} />
        <Route
          path="/admin/eventos/nuevo"
          element={<RequireAuth><><AdminDashboard /><AdminCreateEvent /></></RequireAuth>}
        />
        <Route
          path="/admin/create-event"
          element={<RequireAuth><><AdminDashboard /><AdminCreateEvent /></></RequireAuth>}
        />
        <Route path="/login" element={<PublicOnly><Login /></PublicOnly>} />
        <Route path="/register" element={<PublicOnly><Register /></PublicOnly>} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Router>
  );
}

export default App;
