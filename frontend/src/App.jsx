import { Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './contexts/AuthContext'
import { AppProvider } from './contexts/AppContext'
import Login from './pages/Login'
import Register from './pages/Register'
import ForgotPassword from './pages/ForgotPassword'
import Files from './pages/Files'
import Friends from './pages/Friends'
import ShareDetail from './pages/ShareDetail'
import Profile from './pages/Profile'
import Layout from './components/Layout'

function PrivateRoute({ children }) {
  const { token, isInitialized } = useAuth()
  
  // Wait for initialization to complete before checking token
  // This prevents redirecting to login during initial token load
  if (!isInitialized) {
    return null // or a loading spinner
  }
  
  // Redirect to login if no token after initialization
  if (!token || token === '') {
    return <Navigate to="/login" replace />
  }
  
  // If token exists, show protected content
  return children
}

function App() {
  return (
    <AppProvider>
      <AuthProvider>
        <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route path="/share/:identity" element={<ShareDetail />} />
        <Route
          path="/"
          element={
            <PrivateRoute>
              <Layout />
            </PrivateRoute>
          }
        >
          <Route index element={<Navigate to="/files" replace />} />
          <Route path="files" element={<Files />} />
          <Route path="friends" element={<Friends />} />
          <Route path="profile" element={<Profile />} />
        </Route>
        {/* Default route: redirect to login if not authenticated */}
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
      </AuthProvider>
    </AppProvider>
  )
}

export default App

