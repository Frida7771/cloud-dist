import { useState, useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { friendService } from '../services/friendService'
import { userService } from '../services/userService'
import { storageService, STORAGE_PLANS } from '../services/storageService'
import './Profile.css'

// Password change form component
function PasswordChangeForm() {
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
      setOldPassword('')
      setNewPassword('')
      setConfirmPassword('')
      setCode('')
      setCodeSent(false)
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to update password')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="password-form">
      <h3>Change Password</h3>
      <div className="form-group">
        <label>Email</label>
        <input type="email" value={user?.email || ''} disabled className="disabled-input" />
        <button type="button" onClick={handleSendCode} disabled={sendingCode || codeSent} className="code-btn">
          {sendingCode ? 'Sending...' : codeSent ? 'Code Sent' : 'Send Verification Code'}
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
  )
}

function Profile() {
  const { user, token } = useAuth()
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()
  const [activeTab, setActiveTab] = useState('info')
  const [friends, setFriends] = useState([])
  const [loading, setLoading] = useState(false)
  const [userLoading, setUserLoading] = useState(true)
  
  // Add friend form
  const [newFriendEmail, setNewFriendEmail] = useState('')
  const [friendMessage, setFriendMessage] = useState('')

  // Storage purchase
  const [orders, setOrders] = useState([])
  const [purchasing, setPurchasing] = useState(false)
  const [orderFilter, setOrderFilter] = useState('')

  // Debug: log user data
  useEffect(() => {
    console.log('Profile - User data:', user)
    console.log('Profile - Token:', token ? 'Present' : 'Missing')
    console.log('Profile - User fields:', {
      name: user?.name,
      email: user?.email,
      now_volume: user?.now_volume,
      total_volume: user?.total_volume
    })
    
    // If user data is not loaded but token exists, wait a bit for AuthContext to load
    if (!user && token) {
      setUserLoading(true)
      const timer = setTimeout(() => {
        setUserLoading(false)
      }, 1000)
      return () => clearTimeout(timer)
    } else {
      setUserLoading(false)
    }
  }, [user, token])

  useEffect(() => {
    // Check for payment result in URL
    const payment = searchParams.get('payment')
    
    if (payment === 'success') {
      // Payment successful, wait for webhook to update order status
      alert('Payment successful! Processing your order...')
      setSearchParams({}) // Clear URL params
      
      // Switch to storage tab to show orders
      setActiveTab('storage')
      
      // Start polling for order status update (webhook will update it)
      startOrderStatusPolling()
    } else if (payment === 'cancel') {
      alert('Payment was cancelled.')
      setSearchParams({}) // Clear URL params
    }
  }, [searchParams, setSearchParams])

  const startOrderStatusPolling = () => {
    // Poll for order status updates (webhook will update the status)
    let pollCount = 0
    const maxPolls = 30 // Poll for up to 30 seconds (30 * 1 second)
    
    const pollInterval = setInterval(async () => {
      pollCount++
      
      try {
        // Reload orders to check for status update
        const response = await storageService.getOrderList('')
        const updatedOrders = response.data.list || []
        setOrders(updatedOrders)
        
        // Check if any pending order became paid
        const hasPaidOrder = updatedOrders.some(order => order.status === 'paid')
        
        if (hasPaidOrder || pollCount >= maxPolls) {
          clearInterval(pollInterval)
          
          if (hasPaidOrder) {
            // Order updated by webhook, refresh user data
            alert('Order confirmed! Your storage capacity has been increased.')
            // Refresh user data to show updated storage
            window.location.reload()
          } else if (pollCount >= maxPolls) {
            // Timeout, show message
            alert('Payment successful! Order is being processed. Please refresh the page in a moment.')
          }
        }
      } catch (error) {
        console.error('Failed to poll order status:', error)
        if (pollCount >= maxPolls) {
          clearInterval(pollInterval)
        }
      }
    }, 1000) // Poll every 1 second
    
    // Cleanup on unmount
    return () => clearInterval(pollInterval)
  }

  useEffect(() => {
    if (activeTab === 'friends') {
      loadFriends()
    } else if (activeTab === 'storage') {
      loadOrders()
    }
  }, [activeTab, orderFilter])

  const loadFriends = async () => {
    setLoading(true)
    try {
      const response = await friendService.getFriends()
      setFriends(response.data.list || [])
    } catch (error) {
      console.error('Failed to load friends:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSendFriendRequest = async () => {
    if (!newFriendEmail) {
      alert('Please enter email or user ID')
      return
    }

    try {
      await friendService.sendFriendRequest(newFriendEmail, friendMessage)
      alert('Friend request sent!')
      setNewFriendEmail('')
      setFriendMessage('')
    } catch (error) {
      alert('Failed to send friend request: ' + (error.response?.data?.error || error.message))
    }
  }

  const handlePurchaseStorage = async (storageBytes) => {
    if (!confirm(`Are you sure you want to purchase this storage plan?`)) {
      return
    }

    setPurchasing(true)
    try {
      const response = await storageService.createPurchaseSession(storageBytes)
      // Redirect to Stripe Checkout
      window.location.href = response.data.url
    } catch (error) {
      alert('Failed to create payment session: ' + (error.response?.data?.error || error.message))
      setPurchasing(false)
    }
  }

  const loadOrders = async () => {
    setLoading(true)
    try {
      const response = await storageService.getOrderList(orderFilter)
      setOrders(response.data.list || [])
    } catch (error) {
      console.error('Failed to load orders:', error)
      setOrders([])
    } finally {
      setLoading(false)
    }
  }

  const formatBytes = (bytes) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
  }

  // Calculate storage usage percentage (precise value for display and progress bar)
  const usedPercent = (() => {
    if (!user) return 0
    const nowVolume = user.now_volume || user?.NowVolume || 0
    const totalVolume = user.total_volume || user?.TotalVolume || 0
    if (totalVolume === 0) return 0
    const percent = (nowVolume / totalVolume) * 100
    // Ensure percentage doesn't exceed 100%
    return Math.min(percent, 100)
  })()

  // Format percentage for display (show at least 2 decimal places if less than 1%)
  const formatPercent = (percent) => {
    if (percent === 0) return '0%'
    if (percent < 1) {
      return percent.toFixed(2) + '%'
    }
    return Math.round(percent) + '%'
  }

  return (
    <div className="profile">
      <h2>My Profile</h2>

      <div className="tabs">
        <button
          className={activeTab === 'info' ? 'active' : ''}
          onClick={() => setActiveTab('info')}
        >
          User Info
        </button>
        <button
          className={activeTab === 'password' ? 'active' : ''}
          onClick={() => setActiveTab('password')}
        >
          Change Password
        </button>
        <button
          className={activeTab === 'friends' ? 'active' : ''}
          onClick={() => setActiveTab('friends')}
        >
          Friends
        </button>
        <button
          className={activeTab === 'storage' ? 'active' : ''}
          onClick={() => setActiveTab('storage')}
        >
          Buy Storage
        </button>
      </div>

      {activeTab === 'info' && (
        <div className="tab-content">
          {userLoading ? (
            <div className="loading">Loading user information...</div>
          ) : !user ? (
            <div className="empty">No user information available. Please login again.</div>
          ) : (
            <div className="info-card">
              <h3>User Information</h3>
              <div className="info-item">
                <label>Name:</label>
                <span>{user.name || user.Name || '-'}</span>
              </div>
              <div className="info-item">
                <label>Email:</label>
                <span>{user.email || user.Email || '-'}</span>
              </div>
              <div className="info-item">
                <label>Storage Usage:</label>
                <div className="storage-info">
                  <div className="progress-bar">
                    <div
                      className="progress-fill"
                      style={{ width: `${usedPercent}%` }}
                    ></div>
                  </div>
                  <span>
                    {formatBytes(user.now_volume || user.NowVolume || 0)} / {formatBytes(user.total_volume || user.TotalVolume || 0)} ({formatPercent(usedPercent)})
                  </span>
                </div>
              </div>
              {/* Debug info - remove in production */}
              {process.env.NODE_ENV === 'development' && (
                <div className="info-item" style={{ fontSize: '0.8rem', color: '#999', marginTop: '1rem', padding: '1rem', background: '#f5f5f5', borderRadius: '4px' }}>
                  <label>Debug Info:</label>
                  <pre style={{ fontSize: '0.7rem', overflow: 'auto', marginTop: '0.5rem' }}>
                    {JSON.stringify(user, null, 2)}
                  </pre>
                </div>
              )}
            </div>
          )}
        </div>
      )}

      {activeTab === 'password' && (
        <div className="tab-content">
          <PasswordChangeForm />
        </div>
      )}

      {activeTab === 'friends' && (
        <div className="tab-content">
          <div className="add-friend-section">
            <h3>Add Friend</h3>
            <div className="add-friend-form">
              <input
                type="text"
                placeholder="Enter email or user ID"
                value={newFriendEmail}
                onChange={(e) => setNewFriendEmail(e.target.value)}
              />
              <textarea
                placeholder="Message (optional)"
                value={friendMessage}
                onChange={(e) => setFriendMessage(e.target.value)}
                rows={2}
              />
              <button onClick={handleSendFriendRequest} className="add-friend-btn">
                Send Friend Request
              </button>
            </div>
          </div>

          <div className="friends-section">
            <h3>My Friends</h3>
            {loading ? (
              <div className="loading">Loading...</div>
            ) : (
              <div className="friends-list">
                {friends.map((friend) => (
                  <div key={friend.identity} className="friend-item">
                    <div>
                      <strong>{friend.user_name}</strong>
                      <p>{friend.user_email}</p>
                    </div>
                  </div>
                ))}
                {friends.length === 0 && (
                  <div className="empty">No friends yet</div>
                )}
              </div>
            )}
          </div>
        </div>
      )}

      {activeTab === 'storage' && (
        <div className="tab-content">
          <div className="storage-purchase-section">
            <h3>Purchase Storage</h3>
            <p className="storage-description">
              Choose a storage plan to increase your cloud storage capacity.
            </p>
            <div className="storage-plans">
              {STORAGE_PLANS.map((plan) => (
                <div key={plan.bytes} className="storage-plan-card">
                  <div className="plan-header">
                    <h4>{plan.name}</h4>
                    <div className="plan-price">${plan.price}</div>
                  </div>
                  <div className="plan-details">
                    <p>{formatBytes(plan.bytes)} of additional storage</p>
                  </div>
                  <button
                    className="purchase-btn"
                    onClick={() => handlePurchaseStorage(plan.bytes)}
                    disabled={purchasing}
                  >
                    {purchasing ? 'Processing...' : 'Purchase'}
                  </button>
                </div>
              ))}
            </div>
          </div>

          <div className="orders-section">
            <h3>Order History</h3>
            <div className="order-filters">
              <button
                className={orderFilter === '' ? 'active' : ''}
                onClick={() => setOrderFilter('')}
              >
                All
              </button>
              <button
                className={orderFilter === 'pending' ? 'active' : ''}
                onClick={() => setOrderFilter('pending')}
              >
                Pending
              </button>
              <button
                className={orderFilter === 'paid' ? 'active' : ''}
                onClick={() => setOrderFilter('paid')}
              >
                Paid
              </button>
              <button
                className={orderFilter === 'failed' ? 'active' : ''}
                onClick={() => setOrderFilter('failed')}
              >
                Failed
              </button>
            </div>
            {loading ? (
              <div className="loading">Loading orders...</div>
            ) : (
              <div className="orders-list">
                {orders.map((order) => (
                  <div key={order.identity} className="order-item">
                    <div className="order-info">
                      <div className="order-header">
                        <span className="order-storage">{formatBytes(order.storage_amount)}</span>
                        <span className={`order-status ${order.status}`}>{order.status}</span>
                      </div>
                      <div className="order-details">
                        <span>Price: ${(order.price_amount / 100).toFixed(2)} {order.currency.toUpperCase()}</span>
                        <span>Date: {new Date(order.created_at).toLocaleString()}</span>
                      </div>
                    </div>
                  </div>
                ))}
                {orders.length === 0 && (
                  <div className="empty">No orders found</div>
                )}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

export default Profile

