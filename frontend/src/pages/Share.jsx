import { useState, useEffect } from 'react'
import { friendService } from '../services/friendService'
import { fileService } from '../services/fileService'
import './Share.css'

function Share() {
  const [friends, setFriends] = useState([])
  const [sharedFiles, setSharedFiles] = useState([])
  const [requests, setRequests] = useState([])
  const [activeTab, setActiveTab] = useState('share')
  const [loading, setLoading] = useState(false)
  
  // Share file form
  const [selectedFile, setSelectedFile] = useState('')
  const [selectedFriend, setSelectedFriend] = useState('')
  const [shareMessage, setShareMessage] = useState('')
  const [userFiles, setUserFiles] = useState([])

  useEffect(() => {
    if (activeTab === 'share') {
      loadFriends()
      loadUserFiles()
    } else if (activeTab === 'received') {
      loadSharedFiles('received')
    } else if (activeTab === 'sent') {
      loadSharedFiles('sent')
    } else if (activeTab === 'requests') {
      loadRequests()
    }
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

  const loadUserFiles = async () => {
    try {
      // Recursively load all files from all folders
      const allFiles = []
      await loadFilesRecursive('', allFiles)
      setUserFiles(allFiles)
    } catch (error) {
      console.error('Failed to load files:', error)
      setUserFiles([])
    }
  }

  const loadFilesRecursive = async (folderIdentity, allFiles) => {
    try {
      // Load files from current folder
      const response = await fileService.getFileList(folderIdentity)
      const items = response.data.list || []
      
      // Add files (not folders) to the list
      const files = items.filter(item => item.ext !== '')
      allFiles.push(...files)
      
      // Get folders and recursively load their files
      const folders = items.filter(item => item.ext === '')
      for (const folder of folders) {
        await loadFilesRecursive(folder.identity, allFiles)
      }
    } catch (error) {
      console.error('Failed to load files recursively:', error)
    }
  }

  const loadSharedFiles = async (type) => {
    setLoading(true)
    try {
      const response = await friendService.getSharedFiles(type)
      setSharedFiles(response.data.list || [])
    } catch (error) {
      console.error('Failed to load shared files:', error)
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

  const handleShareFile = async () => {
    if (!selectedFile || !selectedFriend) {
      alert('Please select a file and a friend')
      return
    }

    try {
      await friendService.shareFile(selectedFriend, selectedFile, shareMessage)
      alert('File shared successfully!')
      setSelectedFile('')
      setSelectedFriend('')
      setShareMessage('')
      loadSharedFiles('sent')
    } catch (error) {
      alert('Failed to share file: ' + (error.response?.data?.error || error.message))
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

  const handleDownload = (path) => {
    window.open(path, '_blank')
  }

  return (
    <div className="share">
      <h2>Share</h2>

      <div className="tabs">
        <button
          className={activeTab === 'share' ? 'active' : ''}
          onClick={() => setActiveTab('share')}
        >
          Share File
        </button>
        <button
          className={activeTab === 'received' ? 'active' : ''}
          onClick={() => setActiveTab('received')}
        >
          Received
        </button>
        <button
          className={activeTab === 'sent' ? 'active' : ''}
          onClick={() => setActiveTab('sent')}
        >
          Sent
        </button>
        <button
          className={activeTab === 'requests' ? 'active' : ''}
          onClick={() => setActiveTab('requests')}
        >
          Friend Requests
        </button>
      </div>

      {activeTab === 'share' && (
        <div className="tab-content">
          <div className="share-form">
            <h3>Share File with Friend</h3>
            <div className="form-group">
              <label>Select File</label>
              <select
                value={selectedFile}
                onChange={(e) => setSelectedFile(e.target.value)}
              >
                <option value="">Choose a file...</option>
                {userFiles
                  .filter((file) => file.ext !== '')
                  .map((file) => (
                    <option key={file.identity} value={file.identity}>
                      {file.name}{file.ext}
                    </option>
                  ))}
              </select>
            </div>
            <div className="form-group">
              <label>Select Friend</label>
              <select
                value={selectedFriend}
                onChange={(e) => setSelectedFriend(e.target.value)}
              >
                <option value="">Choose a friend...</option>
                {friends.map((friend) => (
                  <option key={friend.identity} value={friend.user_identity}>
                    {friend.user_name} ({friend.user_email})
                  </option>
                ))}
              </select>
            </div>
            <div className="form-group">
              <label>Message (Optional)</label>
              <textarea
                value={shareMessage}
                onChange={(e) => setShareMessage(e.target.value)}
                placeholder="Add a message..."
                rows={3}
              />
            </div>
            <button onClick={handleShareFile} className="share-btn">
              Share File
            </button>
          </div>
        </div>
      )}

      {activeTab === 'received' && (
        <div className="tab-content">
          {loading ? (
            <div className="loading">Loading...</div>
          ) : (
            <div className="shared-files-list">
              {sharedFiles.map((share) => (
                <div key={share.identity} className="shared-file-item">
                  <div className="file-info">
                    <strong>{share.file_name}{share.file_ext}</strong>
                    <p>From: {share.from_user_name}</p>
                    {share.message && <p className="message">{share.message}</p>}
                    <p className="date">{share.created_at}</p>
                  </div>
                  <button
                    onClick={() => handleDownload(share.path)}
                    className="download-btn"
                  >
                    Download
                  </button>
                </div>
              ))}
              {sharedFiles.length === 0 && (
                <div className="empty">No files shared with you</div>
              )}
            </div>
          )}
        </div>
      )}

      {activeTab === 'sent' && (
        <div className="tab-content">
          {loading ? (
            <div className="loading">Loading...</div>
          ) : (
            <div className="shared-files-list">
              {sharedFiles.map((share) => (
                <div key={share.identity} className="shared-file-item">
                  <div className="file-info">
                    <strong>{share.file_name}{share.file_ext}</strong>
                    <p>To: {share.to_user_name}</p>
                    {share.message && <p className="message">{share.message}</p>}
                    <p className="date">{share.created_at}</p>
                  </div>
                </div>
              ))}
              {sharedFiles.length === 0 && (
                <div className="empty">No files shared yet</div>
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
    </div>
  )
}

export default Share

