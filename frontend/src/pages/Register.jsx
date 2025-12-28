import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import './Auth.css'

function Register() {
  const [email, setEmail] = useState('')
  const [name, setName] = useState('')
  const [password, setPassword] = useState('')
  const [code, setCode] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [codeSent, setCodeSent] = useState(false)
  const { register, sendVerificationCode } = useAuth()
  const navigate = useNavigate()

  const handleSendCode = async () => {
    const result = await sendVerificationCode(email)
    if (result.success) {
      setCodeSent(true)
      alert('Verification code sent to your email')
    } else {
      setError(result.error)
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    const result = await register(email, name, password, code)
    setLoading(false)

    if (result.success) {
      navigate('/login')
    } else {
      setError(result.error)
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
          <h2>Register</h2>
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label>Email</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="Enter email"
                required
              />
              <button type="button" onClick={handleSendCode} className="code-btn" disabled={!email || codeSent}>
                {codeSent ? 'Code Sent' : 'Send Code'}
              </button>
            </div>
            <div className="form-group">
              <label>Verification Code</label>
              <input
                type="text"
                value={code}
                onChange={(e) => setCode(e.target.value)}
                placeholder="Enter verification code"
                required
                disabled={!codeSent}
                maxLength={6}
              />
            </div>
            <div className="form-group">
              <label>Username</label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Enter username"
                required
              />
            </div>
            <div className="form-group">
              <label>Password</label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Enter password"
                required
              />
            </div>
            {error && <div className="error">{error}</div>}
            <button type="submit" disabled={loading || !codeSent} className="btn-submit">
              {loading ? 'Registering...' : 'Register'}
            </button>
            <p className="auth-link">
              Already have an account? <Link to="/login">Login</Link>
            </p>
          </form>
        </div>
      </div>
    </div>
  )
}

export default Register

