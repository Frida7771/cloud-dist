import { Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './contexts/AuthContext'
import Login from './pages/Login'
import Register from './pages/Register'
import Files from './pages/Files'
import Share from './pages/Share'
import Profile from './pages/Profile'
import Layout from './components/Layout'

function PrivateRoute({ children }) {
  const { token } = useAuth()
  
  // Check token from context (more reliable than localStorage check)
  // Only redirect if token is explicitly empty (not just undefined during initialization)
  if (token === '' || (token === null && !localStorage.getItem('token'))) {
    return <Navigate to="/login" replace />
  }
  
  // If token exists or is being initialized, show children
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
          <Route index element={<Files />} />
          <Route path="files" element={<Files />} />
          <Route path="share" element={<Share />} />
          <Route path="profile" element={<Profile />} />
        </Route>
      </Routes>
    </AuthProvider>
  )
}

export default App

