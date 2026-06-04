import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import HelloWorld from './components/HelloWorld';
import EventDetail from './pages/EventDetail';
import Register from './pages/Register';
import MyTickets from './pages/MyTickets';
import AdminDashboard from './pages/AdminDashboard';
import AdminCreateEvent from './pages/AdminCreateEvent';

const hasToken = () => Boolean(localStorage.getItem('token'));

const RequireAuth = ({ children }) => (hasToken() ? children : <Navigate to="/login" replace />);

const PublicOnly = ({ children }) => (hasToken() ? <Navigate to="/" replace /> : children);

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HelloWorld />} />
        <Route path="/eventos/:slug" element={<EventDetail />} />
        <Route path="/mis-entradas" element={<RequireAuth><MyTickets /></RequireAuth>} />
        <Route path="/admin" element={<RequireAuth><AdminDashboard /></RequireAuth>} />
        <Route path="/admin/eventos/nuevo" element={<RequireAuth><AdminCreateEvent /></RequireAuth>} />
        <Route path="/login" element={<PublicOnly><Login /></PublicOnly>} />
        <Route path="/register" element={<PublicOnly><Register /></PublicOnly>} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Router>
  );
}

export default App;
