import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import '../../styles/Login.css';
import { loginUser } from '../../services/api/client';

const ErrorIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <circle cx="12" cy="12" r="10"></circle>
    <line x1="12" y1="8" x2="12" y2="12"></line>
    <line x1="12" y1="16" x2="12.01" y2="16"></line>
  </svg>
);

const ArrowIcon = () => (
  <svg className="arrow-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <line x1="5" y1="12" x2="19" y2="12"></line>
    <polyline points="12 5 19 12 12 19"></polyline>
  </svg>
);

const EyeIcon = ({ open }) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M2 12s3.6-6 10-6 10 6 10 6-3.6 6-10 6-10-6-10-6Z"></path>
    <circle cx="12" cy="12" r="3"></circle>
    {!open && <path d="M4 4 20 20"></path>}
  </svg>
);

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(false);
  const [generalError, setGeneralError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    if (params.get('expired') === 'true') {
      setGeneralError('Su sesión ha expirado. Debe iniciar nuevamente');
    }
  }, []);

  const validateForm = () => {
    const nextErrors = {};
    if (!email) {
      nextErrors.email = 'El correo electronico es obligatorio';
    } else if (!/\S+@\S+\.\S+/.test(email)) {
      nextErrors.email = 'El formato del correo es invalido';
    }
    if (!password) {
      nextErrors.password = 'La contraseña es obligatoria';
    }
    setErrors(nextErrors);
    return Object.keys(nextErrors).length === 0;
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setGeneralError('');

    if (!validateForm()) {
      return;
    }

    setLoading(true);
    try {
      const response = await loginUser({
        email,
        password,
      });

      localStorage.setItem('token', response.data.token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
      navigate('/');
    } catch (error) {
      if (error.response?.status === 401) {
        setGeneralError('Email o contraseña incorrectos');
      } else {
        setGeneralError(error.response?.data?.error || 'No se pudo conectar con el servidor');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-left">
        <div className="login-form-container">
          <div className="login-logo">Golden Ticket</div>

          <h1 className="login-title">Bienvenido</h1>
          <p className="login-subtitle">
            Ingresa tus datos para acceder a tu panel de gestion de eventos.
          </p>

          {generalError && <div className="general-error">{generalError}</div>}

          <form className="login-form" onSubmit={handleSubmit}>
            <div className="input-group">
              <div className="input-header">
                <label className="input-label" htmlFor="email">Correo Electronico</label>
              </div>
              <div className="input-wrapper">
                <input
                  id="email"
                  type="email"
                  className={`login-input ${errors.email ? 'error' : ''}`}
                  placeholder="usuario@ejemplo.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                />
                {errors.email && (
                  <div className="error-icon">
                    <ErrorIcon />
                  </div>
                )}
              </div>
              {errors.email && <span className="error-message">{errors.email}</span>}
            </div>

            <div className="input-group">
              <div className="input-header">
                <label className="input-label" htmlFor="password">Contraseña</label>
              </div>
              <div className="input-wrapper">
                <input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  className={`login-input ${errors.password ? 'error' : ''}`}
                  placeholder="********"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
                <button
                  type="button"
                  className={`password-toggle ${errors.password ? 'has-error' : ''}`}
                  onClick={() => setShowPassword((current) => !current)}
                  aria-label={showPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'}
                  aria-pressed={showPassword}
                >
                  <EyeIcon open={showPassword} />
                </button>
                {errors.password && (
                  <div className="error-icon">
                    <ErrorIcon />
                  </div>
                )}
              </div>
              {errors.password && <span className="error-message">{errors.password}</span>}
            </div>

            <button type="submit" className="login-button" disabled={loading}>
              {loading ? 'Ingresando...' : (
                <>
                  Iniciar sesion
                  <ArrowIcon />
                </>
              )}
            </button>
          </form>

          <hr className="login-separator" />

          <div className="register-text">
            No tienes una cuenta?
            <span className="register-link" onClick={() => navigate('/register')}>
              Crear cuenta
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

export default Login;
