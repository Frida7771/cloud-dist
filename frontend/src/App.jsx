import { Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './contexts/AuthContext'
import Login from './pages/Login'
import Register from './pages/Register'
import Files from './pages/Files'
import Share from './pages/Share'
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
    <AuthProvider>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
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
          <Route path="share" element={<Share />} />
          <Route path="profile" element={<Profile />} />
        </Route>
        {/* Default route: redirect to login if not authenticated */}
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    </AuthProvider>
  )
}

export default App

