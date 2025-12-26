import { useState } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { userService } from '../services/userService'
import './ChangePassword.css'

function ChangePassword() {
  const { user } = useAuth()
  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [code, setCode] = useState('')
  const [codeSent, setCodeSent] = useState(false)
  const [loading, setLoading] = useState(false)
  const [sendingCode, setSendingCode] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const handleSendCode = async () => {
    if (!user?.email) {
      setError('User email not found')
      return
    }

    setSendingCode(true)
    setError('')
    try {
      await userService.sendPasswordUpdateCode(user.email)
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

    // Validation
    if (!oldPassword || !newPassword || !code) {
      setError('Please fill in all fields')
      return
    }

    if (newPassword.length < 6) {
      setError('New password must be at least 6 characters')
      return
    }

    if (newPassword !== confirmPassword) {
      setError('New passwords do not match')
      return
    }

    if (oldPassword === newPassword) {
      setError('New password must be different from old password')
      return
    }

    if (!codeSent) {
      setError('Please request a verification code first')
      return
    }

    setLoading(true)
    try {
      await userService.updatePassword(oldPassword, newPassword, code)
      setSuccess('Password updated successfully!')
      // Clear form
      setOldPassword('')
      setNewPassword('')
      setConfirmPassword('')
      setCode('')
      setCodeSent(false)
      setTimeout(() => {
        setSuccess('')
      }, 3000)
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to update password')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="change-password">
      <h2>Change Password</h2>
      <div className="password-form-container">
        <form onSubmit={handleSubmit} className="password-form">
          <div className="form-group">
            <label>Email</label>
            <input
              type="email"
              value={user?.email || ''}
              disabled
              className="disabled-input"
            />
            <button
              type="button"
              onClick={handleSendCode}
              disabled={sendingCode || codeSent}
              className="send-code-btn"
            >
              {sendingCode
                ? 'Sending...'
                : codeSent
                ? 'Code Sent'
                : 'Send Verification Code'}
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
            <label>Old Password</label>
            <input
              type="password"
              value={oldPassword}
              onChange={(e) => setOldPassword(e.target.value)}
              placeholder="Enter old password"
              required
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

          {error && <div className="error-message">{error}</div>}
          {success && <div className="success-message">{success}</div>}

          <button type="submit" disabled={loading || !codeSent} className="submit-btn">
            {loading ? 'Updating...' : 'Update Password'}
          </button>
        </form>
      </div>
    </div>
  )
}

export default ChangePassword

