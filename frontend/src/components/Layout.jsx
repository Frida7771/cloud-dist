import { Outlet, Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import './Layout.css'

function Layout() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  return (
    <div className="layout">
      <header className="header">
        <div className="header-content">
          <h1>Cloud Dist</h1>
          <nav>
            <Link to="/files">Files</Link>
            <Link to="/share">Share</Link>
            <Link to="/profile">Profile</Link>
          </nav>
          <div className="user-info">
            <span>{user?.name || 'User'}</span>
            <button onClick={handleLogout}>Logout</button>
          </div>
        </div>
      </header>
      <main className="main">
        <Outlet />
      </main>
    </div>
  )
}

export default Layout

