import { Outlet, NavLink, useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { useApp } from '../contexts/AppContext'
import './Layout.css'

function Layout() {
  const { user, logout } = useAuth()
  const { t, language, theme, toggleLanguage, toggleTheme } = useApp()
  const navigate = useNavigate()
  const location = useLocation()

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  const handleMyFilesClick = (e) => {
    if (location.pathname === '/files') {
      // 如果已经在/files页面，触发重置到Root的事件
      e.preventDefault()
      window.dispatchEvent(new CustomEvent('resetFilesToRoot'))
    }
    // 如果不在/files页面，NavLink会正常导航
  }

  return (
    <div className="layout">
      <header className="header">
        <div className="header-content">
          <div className="logo">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M10 4H4c-1.11 0-2 .89-2 2v12c0 1.11.89 2 2 2h16c1.11 0 2-.89 2-2V8c0-1.11-.89-2-2-2h-8l-2-2z" fill="#3385ff"/>
            </svg>
            <span>CloudDist</span>
          </div>
          <div className="header-actions">
            <button onClick={toggleLanguage} className="icon-btn" title={t('language')}>
              {language === 'en' ? '中' : 'EN'}
            </button>
            <button onClick={toggleTheme} className="icon-btn" title={t('theme')}>
              {theme === 'light' ? (
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" style={{ width: '20px', height: '20px' }}>
                  <path d="M12 3c.132 0 .263 0 .393 0a7.5 7.5 0 0 0 7.92 12.446a9 9 0 1 1 -8.313-12.454z" fill="currentColor"/>
                </svg>
              ) : (
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" style={{ width: '20px', height: '20px' }}>
                  <circle cx="12" cy="12" r="5" fill="currentColor"/>
                </svg>
              )}
            </button>
            <div className="user-info">
              <span className="username">{user?.name || 'User'}</span>
              <button onClick={handleLogout} className="logout-btn">{t('logout')}</button>
            </div>
          </div>
        </div>
      </header>
      <div className="layout-body">
        <aside className="sidebar">
          <nav className="sidebar-nav">
            <NavLink 
              to="/files" 
              className={({ isActive }) => isActive ? 'nav-item active' : 'nav-item'}
              onClick={handleMyFilesClick}
            >
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M10 4H4c-1.11 0-2 .89-2 2v12c0 1.11.89 2 2 2h16c1.11 0 2-.89 2-2V8c0-1.11-.89-2-2-2h-8l-2-2z" fill="currentColor"/>
              </svg>
              <span>{t('myFiles')}</span>
            </NavLink>
            <NavLink to="/friends" className={({ isActive }) => isActive ? 'nav-item active' : 'nav-item'}>
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z" fill="currentColor"/>
              </svg>
              <span>{t('friends')}</span>
            </NavLink>
            <NavLink to="/profile" className={({ isActive }) => isActive ? 'nav-item active' : 'nav-item'}>
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z" fill="currentColor"/>
              </svg>
              <span>{t('profile')}</span>
            </NavLink>
          </nav>
        </aside>
        <main className="main-content">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

export default Layout

