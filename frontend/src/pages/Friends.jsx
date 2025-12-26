import { useState, useEffect } from 'react'
import { friendService } from '../services/friendService'
import './Friends.css'

function Friends() {
  const [friends, setFriends] = useState([])
  const [requests, setRequests] = useState([])
  const [sharedFiles, setSharedFiles] = useState([])
  const [activeTab, setActiveTab] = useState('friends')
  const [loading, setLoading] = useState(false)
  const [newFriendEmail, setNewFriendEmail] = useState('')

  useEffect(() => {
    if (activeTab === 'friends') loadFriends()
    if (activeTab === 'requests') loadRequests()
    if (activeTab === 'shared') loadSharedFiles()
  }, [activeTab])

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

  const loadRequests = async () => {
    setLoading(true)
    try {
      const response = await friendService.getFriendRequests('received')
      setRequests(response.data.list || [])
    } catch (error) {
      console.error('Failed to load requests:', error)
    } finally {
      setLoading(false)
    }
  }

  const loadSharedFiles = async () => {
    setLoading(true)
    try {
      const response = await friendService.getSharedFiles('received')
      setSharedFiles(response.data.list || [])
    } catch (error) {
      console.error('Failed to load shared files:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSendRequest = async () => {
    if (!newFriendEmail) return
    try {
      await friendService.sendFriendRequest(newFriendEmail)
      alert('Friend request sent')
      setNewFriendEmail('')
    } catch (error) {
      alert('Failed to send friend request')
    }
  }

  const handleRespondRequest = async (identity, action) => {
    try {
      await friendService.respondFriendRequest(identity, action)
      loadRequests()
      loadFriends()
    } catch (error) {
      alert('Failed to respond to request')
    }
  }

  return (
    <div className="friends">
      <h2>Friends</h2>

      <div className="tabs">
        <button
          className={activeTab === 'friends' ? 'active' : ''}
          onClick={() => setActiveTab('friends')}
        >
          Friends
        </button>
        <button
          className={activeTab === 'requests' ? 'active' : ''}
          onClick={() => setActiveTab('requests')}
        >
          Requests
        </button>
        <button
          className={activeTab === 'shared' ? 'active' : ''}
          onClick={() => setActiveTab('shared')}
        >
          Shared Files
        </button>
      </div>

      {activeTab === 'friends' && (
        <div className="tab-content">
          <div className="add-friend">
            <input
              type="text"
              placeholder="Enter email or user ID"
              value={newFriendEmail}
              onChange={(e) => setNewFriendEmail(e.target.value)}
            />
            <button onClick={handleSendRequest}>Send Request</button>
          </div>
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
      )}

      {activeTab === 'requests' && (
        <div className="tab-content">
          {loading ? (
            <div className="loading">Loading...</div>
          ) : (
            <div className="requests-list">
              {requests
                .filter((r) => r.status === 'pending')
                .map((request) => (
                  <div key={request.identity} className="request-item">
                    <div>
                      <strong>{request.from_user_name}</strong>
                      <p>{request.message || 'No message'}</p>
                    </div>
                    <div className="request-actions">
                      <button
                        onClick={() =>
                          handleRespondRequest(request.identity, 'accept')
                        }
                      >
                        Accept
                      </button>
                      <button
                        onClick={() =>
                          handleRespondRequest(request.identity, 'reject')
                        }
                      >
                        Reject
                      </button>
                    </div>
                  </div>
                ))}
              {requests.filter((r) => r.status === 'pending').length === 0 && (
                <div className="empty">No pending requests</div>
              )}
            </div>
          )}
        </div>
      )}

      {activeTab === 'shared' && (
        <div className="tab-content">
          {loading ? (
            <div className="loading">Loading...</div>
          ) : (
            <div className="shared-files-list">
              {sharedFiles.map((share) => (
                <div key={share.identity} className="shared-file-item">
                  <div>
                    <strong>{share.file_name}</strong>
                    <p>From: {share.from_user_name}</p>
                    <p>{share.message || 'No message'}</p>
                  </div>
                  <a
                    href={share.path}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="download-btn"
                  >
                    Download
                  </a>
                </div>
              ))}
              {sharedFiles.length === 0 && (
                <div className="empty">No shared files</div>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  )
}

export default Friends

