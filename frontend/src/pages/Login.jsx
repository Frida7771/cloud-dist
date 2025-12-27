import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import './Auth.css'

function Login() {
  const [name, setName] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { login } = useAuth()
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
        setError(result.error || 'Login failed')
        setLoading(false)
      }
    } catch (err) {
      console.error('Login error:', err)
      setError('Login failed: ' + (err.message || 'Unknown error'))
      setLoading(false)
    }
  }

  return (
    <div className="auth-container">
      <div className="login-wrapper">
        <div className="login-features">
          <h1>Cloud Dist</h1>
          <p className="tagline">Your Secure Cloud Storage Solution</p>
          <div className="features-list">
            <div className="feature-item">
              <span className="feature-icon">📁</span>
              <div>
                <h3>File Management</h3>
                <p>Upload, organize, and manage your files with folders. Support for large files with intelligent chunked upload.</p>
              </div>
            </div>
            <div className="feature-item">
              <span className="feature-icon">🔗</span>
              <div>
                <h3>File Sharing</h3>
                <p>Share files with friends through secure links with expiration control. Share directly with your friends network.</p>
              </div>
            </div>
            <div className="feature-item">
              <span className="feature-icon">👥</span>
              <div>
                <h3>Friend System</h3>
                <p>Connect with friends and share files directly. Send and receive friend requests easily.</p>
              </div>
            </div>
            <div className="feature-item">
              <span className="feature-icon">💾</span>
              <div>
                <h3>Expandable Storage</h3>
                <p>Start with 5GB free storage. Purchase additional capacity as you need with secure Stripe payment.</p>
              </div>
            </div>
            <div className="feature-item">
              <span className="feature-icon">🔒</span>
              <div>
                <h3>Secure & Fast</h3>
                <p>Files stored on AWS S3 with xxHash64 deduplication. JWT authentication and encrypted connections.</p>
              </div>
            </div>
          </div>
        </div>
        <div className="auth-card">
          <h2>Login</h2>
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label>Username</label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </div>
            <div className="form-group">
              <label>Password</label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
            {error && <div className="error">{error}</div>}
            <button type="submit" disabled={loading}>
              {loading ? 'Logging in...' : 'Login'}
            </button>
            <p className="auth-link">
              Don't have an account? <Link to="/register">Register</Link>
            </p>
          </form>
        </div>
      </div>
    </div>
  )
}

export default Login

