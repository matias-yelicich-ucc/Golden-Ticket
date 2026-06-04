import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/Login.css';
import { registerUser } from '../services/api/client';

function Register() {
  const [form, setForm] = useState({
    nombre: '',
    apellido: '',
    email: '',
    password: '',
  });
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleChange = (key, value) => {
    setForm((current) => ({ ...current, [key]: value }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError('');
    setSuccess('');
    setLoading(true);
    try {
      await registerUser(form);
      setSuccess('Cuenta creada con exito. Te llevamos a la home.');
      setTimeout(() => navigate('/'), 1200);
    } catch (requestError) {
      setError(requestError.response?.data?.error || 'No se pudo registrar el usuario');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-left">
        <div className="login-form-container">
          <div className="login-logo">Golden Ticket</div>
          <h1 className="login-title">Crear cuenta</h1>
          <p className="login-subtitle">Registrate para comprar entradas, administrarlas y seguir tus eventos.</p>

          {error && <div className="general-error">{error}</div>}
          {success && <div className="general-error" style={{ borderColor: 'rgba(201, 168, 95, 0.5)', color: '#d8bf7a' }}>{success}</div>}

          <form className="login-form" onSubmit={handleSubmit}>
            <div className="input-group">
              <label className="input-label" htmlFor="name">Nombre</label>
              <input id="name" className="login-input" value={form.nombre} onChange={(e) => handleChange('nombre', e.target.value)} />
            </div>

            <div className="input-group">
              <label className="input-label" htmlFor="lastname">Apellido</label>
              <input id="lastname" className="login-input" value={form.apellido} onChange={(e) => handleChange('apellido', e.target.value)} />
            </div>

            <div className="input-group">
              <label className="input-label" htmlFor="register-email">Correo Electronico</label>
              <input id="register-email" type="email" className="login-input" value={form.email} onChange={(e) => handleChange('email', e.target.value)} />
            </div>

            <div className="input-group">
              <label className="input-label" htmlFor="register-password">Contrasena</label>
              <input id="register-password" type="password" className="login-input" value={form.password} onChange={(e) => handleChange('password', e.target.value)} />
            </div>

            <button type="submit" className="login-button" disabled={loading}>
              {loading ? 'Creando...' : 'Crear cuenta'}
            </button>
          </form>

          <hr className="login-separator" />
          <div className="register-text">
            Ya tenes cuenta?
            <span className="register-link" onClick={() => navigate('/login')}>
              Iniciar sesion
            </span>
          </div>
        </div>
      </div>

      <div className="login-right">
        <div className="login-right-overlay"></div>
      </div>
    </div>
  );
}

export default Register;
