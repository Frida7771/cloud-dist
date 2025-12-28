import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { userService } from '../services/userService'
import './Auth.css'

function ForgotPassword() {
  const [email, setEmail] = useState('')
  const [code, setCode] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [loading, setLoading] = useState(false)
  const [codeSent, setCodeSent] = useState(false)
  const [sendingCode, setSendingCode] = useState(false)
  const navigate = useNavigate()

  const handleSendCode = async () => {
    if (!email) {
      setError('Please enter your email')
      return
    }

    setSendingCode(true)
    setError('')
    try {
      await userService.sendPasswordResetCode(email)
      setCodeSent(true)
      setSuccess('Verification code sent to your email')
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to send verification code')
    } finally {
      setSendingCode(false)
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    if (!email || !code || !newPassword || !confirmPassword) {
      setError('Please fill in all fields')
      return
    }

    if (newPassword.length < 6) {
      setError('Password must be at least 6 characters')
      return
    }

    if (newPassword !== confirmPassword) {
      setError('Passwords do not match')
      return
    }

    if (!codeSent) {
      setError('Please request a verification code first')
      return
    }

    setLoading(true)
    try {
      await userService.resetPassword(email, code, newPassword)
      setSuccess('Password reset successfully! Redirecting to login...')
      setTimeout(() => {
        navigate('/login')
      }, 2000)
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to reset password')
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
          <h2>Forgot Password</h2>
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label>Email</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="Enter your email"
                required
                disabled={codeSent}
              />
              {!codeSent && (
                <button type="button" onClick={handleSendCode} className="code-btn" disabled={sendingCode || !email}>
                  {sendingCode ? 'Sending...' : 'Send Code'}
                </button>
              )}
            </div>
            {codeSent && (
              <>
                <div className="form-group">
                  <label>Verification Code</label>
                  <input
                    type="text"
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                    placeholder="Enter verification code"
                    required
                    maxLength={6}
                  />
                </div>
                <div className="form-group">
                  <label>New Password</label>
                  <input
                    type="password"
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    placeholder="Enter new password (min 6 characters)"
                    required
                    minLength={6}
                  />
                </div>
                <div className="form-group">
                  <label>Confirm New Password</label>
                  <input
                    type="password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    placeholder="Confirm new password"
                    required
                    minLength={6}
                  />
                </div>
              </>
            )}
            {error && <div className="error">{error}</div>}
            {success && <div className="success-message">{success}</div>}
            {codeSent && (
              <button type="submit" disabled={loading} className="btn-submit">
                {loading ? 'Resetting...' : 'Reset Password'}
              </button>
            )}
            <p className="auth-link">
              Remember your password? <Link to="/login">Login</Link>
            </p>
          </form>
        </div>
      </div>
    </div>
  )
}

export default ForgotPassword

