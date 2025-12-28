import { useState, useEffect } from 'react'
import { friendService } from '../services/friendService'
import { fileService } from '../services/fileService'
import { useApp } from '../contexts/AppContext'
import './Friends.css'

function Friends() {
  const { t } = useApp()
  const [friends, setFriends] = useState([])
  const [requests, setRequests] = useState([])
  const [sharedFiles, setSharedFiles] = useState([])
  const [activeTab, setActiveTab] = useState('friends')
  const [loading, setLoading] = useState(false)
  
  // Friend request form
  const [newFriendEmail, setNewFriendEmail] = useState('')
  const [requestMessage, setRequestMessage] = useState('')
  
  // Share file form
  const [selectedFile, setSelectedFile] = useState('')
  const [selectedFriend, setSelectedFriend] = useState('')
  const [shareMessage, setShareMessage] = useState('')
  const [userFiles, setUserFiles] = useState([])
  const [shareType, setShareType] = useState('received') // 'received' or 'sent'
  
  // Save to disk modal
  const [showSaveModal, setShowSaveModal] = useState(false)
  const [fileToSave, setFileToSave] = useState(null)
  const [targetFolderId, setTargetFolderId] = useState(null)
  const [folders, setFolders] = useState([])
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (activeTab === 'friends') {
      loadFriends()
      loadRequests()
    } else if (activeTab === 'share') {
      loadFriends()
      loadUserFiles()
      loadSharedFiles('received')
    } else if (activeTab === 'shared') {
      loadSharedFiles(shareType)
    }
  }, [activeTab, shareType])

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
    try {
      const response = await friendService.getFriendRequests('received')
      setRequests(response.data.list || [])
    } catch (error) {
      console.error('Failed to load requests:', error)
    }
  }

  const loadUserFiles = async () => {
    try {
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
      const response = await fileService.getFileList(folderIdentity)
      const items = response.data.list || []
      
      const files = items.filter(item => item.ext !== '')
      allFiles.push(...files)
      
      const folders = items.filter(item => item.ext === '')
      for (const folder of folders) {
        await loadFilesRecursive(folder.identity, allFiles)
      }
    } catch (error) {
      console.error('Failed to load files recursively:', error)
    }
  }

  const loadFolders = async () => {
    try {
      const allFolders = [{ id: 0, identity: '', name: 'Root', level: 0 }]
      await loadFoldersRecursive('', allFolders, 0)
      setFolders(allFolders)
    } catch (error) {
      console.error('Failed to load folders:', error)
      setFolders([{ id: 0, identity: '', name: 'Root', level: 0 }])
    }
  }

  const loadFoldersRecursive = async (parentIdentity, allFolders, level) => {
    try {
      const response = await fileService.getFileList(parentIdentity)
      const items = response.data.list || []
      const folders = items.filter(item => item.ext === '')
      
      for (const folder of folders) {
        const folderId = folder.id || 0
        if (folderId !== 0 && !allFolders.find(f => f.id === folderId)) {
          allFolders.push({
            id: folderId,
            identity: folder.identity,
            name: folder.name,
            level: level + 1,
          })
          await loadFoldersRecursive(folder.identity, allFolders, level + 1)
        }
      }
    } catch (error) {
      console.error('Failed to load folders recursively:', error)
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

  const handleSendRequest = async () => {
    if (!newFriendEmail.trim()) {
      alert(t('pleaseEnterEmail'))
      return
    }
    try {
      await friendService.sendFriendRequest(newFriendEmail.trim(), requestMessage)
      alert(t('friendRequestSent'))
      setNewFriendEmail('')
      setRequestMessage('')
      loadRequests()
    } catch (error) {
      alert(t('failedToSendRequest') + ': ' + (error.response?.data?.error || error.message))
    }
  }

  const handleRespondRequest = async (identity, action) => {
    try {
      await friendService.respondFriendRequest(identity, action)
      loadRequests()
      loadFriends()
    } catch (error) {
      alert(t('failedToRespond') + ': ' + (error.response?.data?.error || error.message))
    }
  }

  const handleShareFile = async () => {
    if (!selectedFile || !selectedFriend) {
      alert(t('pleaseSelectFileAndFriend'))
      return
    }

    try {
      await friendService.shareFile(selectedFriend, selectedFile, shareMessage)
      alert(t('fileSharedSuccessfully'))
      setSelectedFile('')
      setSelectedFriend('')
      setShareMessage('')
      loadSharedFiles('sent')
    } catch (error) {
      alert(t('failedToShareFile') + ': ' + (error.response?.data?.error || error.message))
    }
  }

  const handleDownload = async (shareIdentity, fileName) => {
    try {
      // Use friend share download endpoint which verifies friendship
      const response = await friendService.downloadShareFile(shareIdentity)
      
      let downloadFileName = fileName || 'file'
      const contentDisposition = response.headers['content-disposition']
      if (contentDisposition) {
        const fileNameMatch = contentDisposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/)
        if (fileNameMatch && fileNameMatch[1]) {
          downloadFileName = decodeURIComponent(fileNameMatch[1].replace(/['"]/g, ''))
        }
      }
      
      const contentType = response.headers['content-type'] || 'application/octet-stream'
      const blob = new Blob([response.data], { type: contentType })
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', downloadFileName)
      document.body.appendChild(link)
      link.click()
      link.remove()
      setTimeout(() => window.URL.revokeObjectURL(url), 100)
    } catch (error) {
      console.error('Download failed:', error)
      alert(t('downloadFailed') + ': ' + (error.response?.data?.error || error.message))
    }
  }

  const handleSaveToDisk = (share) => {
    setFileToSave(share)
    setShowSaveModal(true)
    loadFolders()
  }

  const handleSave = async () => {
    if (!fileToSave || !fileToSave.identity) return
    if (targetFolderId === null) {
      alert(t('pleaseSelectFolder'))
      return
    }

    try {
      setSaving(true)
      await friendService.saveShareFile(fileToSave.identity, targetFolderId)
      alert(t('fileSavedSuccessfully'))
      setShowSaveModal(false)
      setFileToSave(null)
      setTargetFolderId(null)
    } catch (error) {
      console.error('Save failed:', error)
      alert(t('failedToSaveFile') + ': ' + (error.response?.data?.error || error.message))
    } finally {
      setSaving(false)
    }
  }

  const formatDate = (dateString) => {
    if (!dateString) return ''
    const date = new Date(dateString)
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }

  return (
    <div className="friends-page">
      <h2>{t('friends')}</h2>

      <div className="tabs">
        <button
          className={activeTab === 'friends' ? 'active' : ''}
          onClick={() => setActiveTab('friends')}
        >
          {t('friends')}
        </button>
        <button
          className={activeTab === 'share' ? 'active' : ''}
          onClick={() => setActiveTab('share')}
        >
          {t('shareFile')}
        </button>
        <button
          className={activeTab === 'shared' ? 'active' : ''}
          onClick={() => setActiveTab('shared')}
        >
          {t('sharedFiles')}
        </button>
      </div>

      {/* Friends Tab */}
      {activeTab === 'friends' && (
        <div className="tab-content">
          <div className="section">
            <h3>{t('addFriend')}</h3>
            <div className="add-friend-form">
              <div className="form-group">
                <label>{t('emailOrUserId')}</label>
                <input
                  type="text"
                  placeholder={t('enterEmailOrUserId')}
                  value={newFriendEmail}
                  onChange={(e) => setNewFriendEmail(e.target.value)}
                />
              </div>
              <div className="form-group">
                <label>{t('message')} ({t('optional')})</label>
                <textarea
                  value={requestMessage}
                  onChange={(e) => setRequestMessage(e.target.value)}
                  placeholder={t('addMessage')}
                  rows={2}
                />
              </div>
              <button onClick={handleSendRequest} className="btn-primary">
                {t('sendRequest')}
              </button>
            </div>
          </div>

          <div className="section">
            <h3>{t('friendRequests')}</h3>
            {loading ? (
              <div className="loading">{t('loading')}</div>
            ) : (
              <div className="requests-list">
                {requests
                  .filter((r) => r.status === 'pending')
                  .map((request) => (
                    <div key={request.identity} className="request-item">
                      <div className="request-info">
                        <strong>{request.from_user_name}</strong>
                        <p>{request.from_user_email}</p>
                        {request.message && <p className="message">{request.message}</p>}
                        <p className="date">{formatDate(request.created_at)}</p>
                      </div>
                      <div className="request-actions">
                        <button
                          onClick={() => handleRespondRequest(request.identity, 'accept')}
                          className="btn-primary"
                        >
                          {t('accept')}
                        </button>
                        <button
                          onClick={() => handleRespondRequest(request.identity, 'reject')}
                          className="btn-default"
                        >
                          {t('reject')}
                        </button>
                      </div>
                    </div>
                  ))}
                {requests.filter((r) => r.status === 'pending').length === 0 && (
                  <div className="empty">{t('noPendingRequests')}</div>
                )}
              </div>
            )}
          </div>

          <div className="section">
            <h3>{t('myFriends')}</h3>
            {loading ? (
              <div className="loading">{t('loading')}</div>
            ) : (
              <div className="friends-list">
                {friends.map((friend) => (
                  <div key={friend.identity} className="friend-item">
                    <div className="friend-info">
                      <strong>{friend.user_name}</strong>
                      <p>{friend.user_email}</p>
                    </div>
                  </div>
                ))}
                {friends.length === 0 && (
                  <div className="empty">{t('noFriendsYet')}</div>
                )}
              </div>
            )}
          </div>
        </div>
      )}

      {/* Share File Tab */}
      {activeTab === 'share' && (
        <div className="tab-content">
          <div className="share-form">
            <h3>{t('shareFileWithFriend')}</h3>
            <div className="form-group">
              <label>{t('selectFile')}</label>
              <select
                value={selectedFile}
                onChange={(e) => setSelectedFile(e.target.value)}
              >
                <option value="">-- {t('chooseFile')} --</option>
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
              <label>{t('selectFriend')}</label>
              <select
                value={selectedFriend}
                onChange={(e) => setSelectedFriend(e.target.value)}
              >
                <option value="">-- {t('chooseFriend')} --</option>
                {friends.map((friend) => (
                  <option key={friend.identity} value={friend.user_identity}>
                    {friend.user_name} ({friend.user_email})
                  </option>
                ))}
              </select>
            </div>
            <div className="form-group">
              <label>{t('message')} ({t('optional')})</label>
              <textarea
                value={shareMessage}
                onChange={(e) => setShareMessage(e.target.value)}
                placeholder={t('addMessage')}
                rows={3}
              />
            </div>
            <button onClick={handleShareFile} className="btn-primary">
              {t('shareFile')}
            </button>
          </div>
        </div>
      )}

      {/* Shared Files Tab */}
      {activeTab === 'shared' && (
        <div className="tab-content">
          <div className="shared-files-header">
            <div className="sub-tabs">
              <button
                className={shareType === 'received' ? 'active' : ''}
                onClick={() => setShareType('received')}
              >
                {t('received')}
              </button>
              <button
                className={shareType === 'sent' ? 'active' : ''}
                onClick={() => setShareType('sent')}
              >
                {t('sent')}
              </button>
            </div>
          </div>

          {loading ? (
            <div className="loading">{t('loading')}</div>
          ) : (
            <div className="shared-files-list">
              {sharedFiles.map((share) => (
                <div key={share.identity} className="shared-file-item">
                  <div className="file-info">
                    <strong>{share.file_name}{share.file_ext}</strong>
                    {shareType === 'received' ? (
                      <>
                        <p>{t('from')}: {share.from_user_name}</p>
                        {share.message && <p className="message">{share.message}</p>}
                      </>
                    ) : (
                      <>
                        <p>{t('to')}: {share.to_user_name}</p>
                        {share.message && <p className="message">{share.message}</p>}
                      </>
                    )}
                    <p className="date">{formatDate(share.created_at)}</p>
                  </div>
                  {shareType === 'received' && share.identity && (
                    <div className="file-actions">
                      <button
                        onClick={() => handleDownload(share.identity, share.file_name + share.file_ext)}
                        className="btn-primary"
                      >
                        {t('download')}
                      </button>
                      <button
                        onClick={() => handleSaveToDisk(share)}
                        className="btn-default"
                      >
                        {t('saveToDisk')}
                      </button>
                    </div>
                  )}
                </div>
              ))}
              {sharedFiles.length === 0 && (
                <div className="empty">
                  {shareType === 'received' ? t('noFilesSharedWithYou') : t('noFilesSharedYet')}
                </div>
              )}
            </div>
          )}
        </div>
      )}

      {/* Save to Disk Modal */}
      {showSaveModal && (
        <div className="modal-overlay" onClick={() => {
          if (!saving) {
            setShowSaveModal(false)
            setFileToSave(null)
            setTargetFolderId(null)
          }
        }}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h3>{t('saveToDisk')}</h3>
            {fileToSave && (
              <div className="share-file-info">
                <p><strong>{t('file')}:</strong> {fileToSave.file_name}{fileToSave.file_ext}</p>
              </div>
            )}
            <div className="form-group">
              <label>{t('selectTargetFolder')}</label>
              <select
                value={targetFolderId === null ? '' : targetFolderId}
                onChange={(e) => {
                  const value = e.target.value
                  if (value === '') {
                    setTargetFolderId(null)
                  } else {
                    setTargetFolderId(Number(value))
                  }
                }}
                disabled={saving}
              >
                <option value="">-- {t('selectFolder')} --</option>
                <option value="0" style={{ fontWeight: 'bold' }}>📁 {t('root')}</option>
                {folders
                  .filter(folder => folder.id !== 0)
                  .map((folder, index) => (
                    <option key={folder.id || index} value={folder.id}>
                      {'  '.repeat(folder.level || 0)}{folder.name}
                    </option>
                  ))}
              </select>
            </div>
            <div className="modal-actions">
              <button
                onClick={handleSave}
                disabled={targetFolderId === null || saving}
                className="btn-primary"
              >
                {saving ? t('saving') : t('save')}
              </button>
              <button
                onClick={() => {
                  setShowSaveModal(false)
                  setFileToSave(null)
                  setTargetFolderId(null)
                }}
                disabled={saving}
                className="btn-default"
              >
                {t('cancel')}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default Friends
