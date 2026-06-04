import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Login from './pages/auth/Login';
import Register from './pages/auth/Register';
import HomePage from './pages/home/HomePage';
import AdminCreateEvent from './pages/admin/AdminCreateEvent';
import AdminDashboard from './pages/admin/AdminDashboard';
import EventDetail from './pages/events/EventDetail';
import MyTickets from './pages/tickets/MyTickets';

const getStoredUser = () => {
  try {
    return JSON.parse(localStorage.getItem('user') || 'null');
  } catch {
    return null;
  }
};

const hasToken = () => Boolean(localStorage.getItem('token'));
const getUserRole = () => getStoredUser()?.rol;
const isAdmin = () => getUserRole() === 'admin';

const RequireAuth = ({ children }) => (hasToken() ? children : <Navigate to="/login" replace />);

const RequireAdmin = ({ children }) => {
  if (!hasToken()) return <Navigate to="/login" replace />;
  return isAdmin() ? children : <Navigate to="/" replace />;
};

const ClientOnly = ({ children }) => {
  if (isAdmin()) return <Navigate to="/admin/dashboard" replace />;
  return children;
};

const PublicOnly = ({ children }) => {
  if (!hasToken()) return children;
  return <Navigate to={isAdmin() ? '/admin/dashboard' : '/'} replace />;
};

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<ClientOnly><HomePage /></ClientOnly>} />
        <Route path="/eventos/:slug" element={<ClientOnly><EventDetail /></ClientOnly>} />
        <Route path="/mis-entradas" element={<RequireAuth><ClientOnly><MyTickets /></ClientOnly></RequireAuth>} />
        <Route path="/admin" element={<RequireAdmin><AdminDashboard /></RequireAdmin>} />
        <Route path="/admin/dashboard" element={<RequireAdmin><AdminDashboard /></RequireAdmin>} />
        <Route
          path="/admin/eventos/nuevo"
          element={<RequireAdmin><><AdminDashboard /><AdminCreateEvent /></></RequireAdmin>}
        />
        <Route
          path="/admin/create-event"
          element={<RequireAdmin><><AdminDashboard /><AdminCreateEvent /></></RequireAdmin>}
        />
        <Route path="/login" element={<PublicOnly><Login /></PublicOnly>} />
        <Route path="/register" element={<PublicOnly><Register /></PublicOnly>} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Router>
  );
}

export default App;
