import { useAuth } from '../contexts/AuthContext'
import { Link } from 'react-router-dom'
import './Dashboard.css'

function Dashboard() {
  const { user } = useAuth()

  const formatBytes = (bytes) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
  }

  const usedPercent = user
    ? Math.round((user.now_volume / user.total_volume) * 100)
    : 0

  return (
    <div className="dashboard">
      <h2>Dashboard</h2>
      {user && (
        <div className="stats">
          <div className="stat-card">
            <h3>Storage Usage</h3>
            <div className="progress-bar">
              <div
                className="progress-fill"
                style={{ width: `${usedPercent}%` }}
              ></div>
            </div>
            <p>
              {formatBytes(user.now_volume)} / {formatBytes(user.total_volume)} ({usedPercent}%)
            </p>
          </div>
          <div className="stat-card">
            <h3>User Info</h3>
            <p>Name: {user.name}</p>
            <p>Email: {user.email}</p>
            <Link to="/change-password" className="change-password-link">
              Change Password
            </Link>
          </div>
        </div>
      )}
    </div>
  )
}

export default Dashboard

