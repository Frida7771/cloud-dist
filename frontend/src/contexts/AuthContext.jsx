import { createContext, useContext, useState, useEffect } from 'react'
import api from '../services/api'

const AuthContext = createContext()

export function AuthProvider({ children }) {
  // Load token from localStorage on initialization
  const [token, setToken] = useState(localStorage.getItem('token') || '')
  const [user, setUser] = useState(null)
  const [isInitialized, setIsInitialized] = useState(false)

  const fetchUserDetail = async (currentToken) => {
    if (!currentToken) {
      console.log('fetchUserDetail: No token provided')
      return
    }
    try {
      api.defaults.headers.common['Authorization'] = `Bearer ${currentToken}`
      console.log('fetchUserDetail: Calling /user/detail API')
      // Get user identity from token or use empty string (backend will get from context)
      const response = await api.post('/user/detail', {})
      console.log('User detail response:', response.data)
      console.log('User detail response fields:', {
        name: response.data?.name,
        email: response.data?.email,
        now_volume: response.data?.now_volume,
        total_volume: response.data?.total_volume
      })
      setUser(response.data)
    } catch (error) {
      console.error('Failed to fetch user details:', error)
      console.error('Error response:', error.response?.data)
      console.error('Error status:', error.response?.status)
      // Don't clear token on error - let the API interceptor handle 401
      // Just clear user data
      setUser(null)
    }
  }

  // Initialize: load token from localStorage and fetch user details
  useEffect(() => {
    const storedToken = localStorage.getItem('token')
    if (storedToken) {
      setToken(storedToken)
      api.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`
      // Fetch user details immediately
      fetchUserDetail(storedToken)
    }
    setIsInitialized(true)
  }, []) // Only run once on mount

  // Listen for custom event when token is cleared externally (e.g., by API interceptor)
  useEffect(() => {
    const handleTokenCleared = () => {
      const currentToken = localStorage.getItem('token')
      // Only clear if token was actually removed from localStorage
      if (!currentToken) {
        setToken('')
        setUser(null)
        delete api.defaults.headers.common['Authorization']
      }
    }

    window.addEventListener('tokenCleared', handleTokenCleared)
    return () => {
      window.removeEventListener('tokenCleared', handleTokenCleared)
    }
  }, [])

  // Handle token changes and fetch user details
  useEffect(() => {
    if (!isInitialized) return // Wait for initialization
    
    if (token) {
      api.defaults.headers.common['Authorization'] = `Bearer ${token}`
      fetchUserDetail(token)
    } else {
      delete api.defaults.headers.common['Authorization']
      setUser(null)
    }
  }, [token, isInitialized])

  const login = async (name, password) => {
    try {
      const response = await api.post('/user/login', { name, password })
      console.log('Login response:', response.data)
      
      // Check response format
      const newToken = response.data?.token || response.data?.Token
      if (!newToken) {
        console.error('No token in response:', response.data)
        return { success: false, error: 'Invalid response from server' }
      }
      
      // Set token immediately
      localStorage.setItem('token', newToken)
      api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
      
      // Update state synchronously - useEffect will handle fetching user details
      // Small delay to ensure token is persisted before navigation
      setToken(newToken)
      
      return { success: true }
    } catch (error) {
      console.error('Login error:', error)
      const errorMessage = error.response?.data?.error || 
                          error.response?.data?.msg || 
                          error.message || 
                          'Login failed'
      return { success: false, error: errorMessage }
    }
  }

  const register = async (email, name, password, code) => {
    try {
      await api.post('/user/register', { email, name, password, code })
      return { success: true }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Registration failed' }
    }
  }

  const sendVerificationCode = async (email) => {
    try {
      await api.post('/mail/code/send/register', { email })
      return { success: true }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to send code' }
    }
  }

  const logout = async () => {
    try {
      await api.post('/user/logout')
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      setToken('')
      setUser(null)
      localStorage.removeItem('token')
      delete api.defaults.headers.common['Authorization']
    }
  }

  return (
    <AuthContext.Provider value={{ token, user, isInitialized, login, register, sendVerificationCode, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}

