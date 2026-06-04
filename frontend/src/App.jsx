import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import HelloWorld from './components/HelloWorld';
import EventDetail from './pages/EventDetail';
import Register from './pages/Register';
import MyTickets from './pages/MyTickets';

const hasToken = () => Boolean(localStorage.getItem('token'));

const RequireAuth = ({ children }) => (hasToken() ? children : <Navigate to="/login" replace />);

const PublicOnly = ({ children }) => (hasToken() ? <Navigate to="/" replace /> : children);

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<RequireAuth><HelloWorld /></RequireAuth>} />
        <Route path="/eventos/:slug" element={<RequireAuth><EventDetail /></RequireAuth>} />
        <Route path="/mis-entradas" element={<RequireAuth><MyTickets /></RequireAuth>} />
        <Route path="/login" element={<PublicOnly><Login /></PublicOnly>} />
        <Route path="/register" element={<PublicOnly><Register /></PublicOnly>} />
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    </Router>
  );
}

export default App;
