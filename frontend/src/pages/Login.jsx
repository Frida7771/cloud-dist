import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { useApp } from '../contexts/AppContext'
import './Auth.css'

function Login() {
  const [name, setName] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { login } = useAuth()
  const { t } = useApp()
  const navigate = useNavigate()

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const result = await login(name, password)
      console.log('Login result:', result)
      
      if (result.success) {
        // Clear form
        setName('')
        setPassword('')
        
        // Small delay to ensure token is set and state is updated
        setTimeout(() => {
          navigate('/files', { replace: true })
        }, 100)
      } else {
        // Map backend error messages to translated messages
        const errorMsg = result.error || ''
        if (errorMsg.includes('user not registered') || errorMsg.includes('用户未注册')) {
          setError(t('userNotRegistered'))
        } else if (errorMsg.includes('password incorrect') || errorMsg.includes('密码错误')) {
          setError(t('passwordIncorrect'))
        } else {
          setError(errorMsg || t('loginFailed'))
        }
        setLoading(false)
      }
    } catch (err) {
      console.error('Login error:', err)
      setError(t('loginFailed') + ': ' + (err.message || 'Unknown error'))
      setLoading(false)
    }
  }

  return (
    <div className="auth-container">
      <div className="auth-wrapper">
        <div className="auth-logo">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M10 4H4c-1.11 0-2 .89-2 2v12c0 1.11.89 2 2 2h16c1.11 0 2-.89 2-2V8c0-1.11-.89-2-2-2h-8l-2-2z" fill="#3385ff"/>
          </svg>
          <h1>CloudDisk</h1>
        </div>
        <div className="auth-card">
          <h2>{t('loginTitle')}</h2>
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label>{t('username')}</label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder={t('username')}
                required
              />
            </div>
            <div className="form-group">
              <label>{t('password')}</label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder={t('password')}
                required
              />
            </div>
            {error && <div className="error">{error}</div>}
            <button type="submit" disabled={loading} className="btn-submit">
              {loading ? t('loggingIn') : t('login')}
            </button>
            <p className="auth-link">
              <Link to="/forgot-password">{t('forgotPassword')}?</Link>
            </p>
            <p className="auth-link">
              {t('dontHaveAccount')} <Link to="/register">{t('register')}</Link>
            </p>
          </form>
        </div>
      </div>
    </div>
  )
}

export default Login

