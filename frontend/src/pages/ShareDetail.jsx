import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useApp } from '../contexts/AppContext'
import { useAuth } from '../contexts/AuthContext'
import { fileService } from '../services/fileService'
import api from '../services/api'
import './ShareDetail.css'

// File icon component based on file extension
const FileIcon = ({ ext }) => {
  const getIcon = (extension) => {
    const extLower = extension.toLowerCase()
    
    // PDF
    if (extLower === '.pdf') {
      return (
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z" fill="#e74c3c"/>
        </svg>
      )
    }
    
    // Word documents
    if (['.doc', '.docx'].includes(extLower)) {
      return (
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z" fill="#2b579a"/>
        </svg>
      )
    }
    
    // Excel
    if (['.xls', '.xlsx'].includes(extLower)) {
      return (
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z" fill="#1d6f42"/>
        </svg>
      )
    }
    
    // PowerPoint
    if (['.ppt', '.pptx'].includes(extLower)) {
      return (
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z" fill="#d04423"/>
        </svg>
      )
    }
    
    // Images
    if (['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp', '.svg'].includes(extLower)) {
      return (
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M21 19V5c0-1.1-.9-2-2-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2zM8.5 13.5l2.5 3.01L14.5 12l4.5 6H5l3.5-4.5z" fill="#4caf50"/>
        </svg>
      )
    }
    
    // Videos
    if (['.mp4', '.avi', '.mov', '.wmv', '.flv', '.webm', '.mkv'].includes(extLower)) {
      return (
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M21 3H3c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h18c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-8 12.5v-9l6 4.5-6 4.5z" fill="#f44336"/>
        </svg>
      )
    }
    
    // Audio
    if (['.mp3', '.wav', '.flac', '.aac', '.ogg', '.m4a'].includes(extLower)) {
      return (
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z" fill="#9c27b0"/>
        </svg>
      )
    }
    
    // Text files
    if (['.txt', '.md', '.log', '.csv'].includes(extLower)) {
      return (
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z" fill="#757575"/>
        </svg>
      )
    }
    
    // Code files
    if (['.js', '.jsx', '.ts', '.tsx', '.py', '.java', '.cpp', '.c', '.html', '.css', '.json', '.xml', '.yaml', '.yml'].includes(extLower)) {
      return (
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M9.4 16.6L4.8 12l4.6-4.6L8 6l-6 6 6 6 1.4-1.4zm5.2 0L19.2 12l-4.6-4.6L16 6l6 6-6 6-1.4-1.4z" fill="#ff9800"/>
        </svg>
      )
    }
    
    // Archives
    if (['.zip', '.rar', '.7z', '.tar', '.gz'].includes(extLower)) {
      return (
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M20 6h-2.18c.11-.31.18-.65.18-1a2.996 2.996 0 0 0-5.5-1.65l-.5.67-.5-.68C10.96 2.54 10 2 9 2 7.34 2 6 3.34 6 5c0 .35.07.69.18 1H4c-1.11 0-1.99.89-1.99 2L2 19c0 1.11.89 2 2 2h16c1.11 0 2-.89 2-2V8c0-1.11-.89-2-2-2zm-5-2c.55 0 1 .45 1 1s-.45 1-1 1-1-.45-1-1 .45-1 1-1zM9 4c.55 0 1 .45 1 1s-.45 1-1 1-1-.45-1-1 .45-1 1-1zm11 15H4v-2h16v2zm0-5H4V8h5.08L7 10.83 8.62 12 11 8.76l1-1.36 1 1.36L15.38 12 17 10.83 14.92 8H20v6z" fill="#ff9800"/>
        </svg>
      )
    }
    
    // Default file icon
    return (
      <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z" fill="#3385ff"/>
      </svg>
    )
  }
  
  return getIcon(ext)
}

function ShareDetail() {
  const { t } = useApp()
  const { token } = useAuth()
  const { identity } = useParams()
  const navigate = useNavigate()
  const [shareInfo, setShareInfo] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [saving, setSaving] = useState(false)
  const [showSaveModal, setShowSaveModal] = useState(false)
  const [targetFolderId, setTargetFolderId] = useState(null)
  const [folders, setFolders] = useState([])
  const [showPreviewModal, setShowPreviewModal] = useState(false)

  useEffect(() => {
    loadShareDetail()
  }, [identity])

  useEffect(() => {
    if (showSaveModal && token) {
      loadFolders()
    }
  }, [showSaveModal, token])

  const loadShareDetail = async () => {
    try {
      setLoading(true)
      // This endpoint is public, no auth required
      const response = await api.get(`/share/basic/detail?identity=${identity}`)
      setShareInfo(response.data)
      setError(null)
    } catch (err) {
      console.error('Failed to load share detail:', err)
      setError(err.response?.data?.error || 'Share link is invalid or expired')
    } finally {
      setLoading(false)
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
        // Use folder.id if available, otherwise use 0
        const folderId = folder.id || 0
        // Only add folders with valid id (not 0) to avoid duplicates
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

  const handlePreview = () => {
    if (!shareInfo || !shareInfo.path) return
    
    const extLower = (shareInfo.ext || '').toLowerCase()
    const videoExts = ['.mp4', '.avi', '.mov', '.wmv', '.flv', '.webm', '.mkv']
    
    // For video files, show in modal preview
    if (videoExts.includes(extLower)) {
      setShowPreviewModal(true)
    } else {
      // For other files (images, PDF, etc.), open in new window
      window.open(shareInfo.path, '_blank')
    }
  }

  const handleDownload = () => {
    if (!shareInfo || !shareInfo.download_url) return
    
    // Use the download URL directly (with Content-Disposition: attachment)
    // This avoids CORS issues and ensures the browser downloads the file
    const link = document.createElement('a')
    link.href = shareInfo.download_url
    link.download = `${shareInfo.name}${shareInfo.ext || ''}`
    link.target = '_blank' // Open in new tab as fallback
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }

  const handleSaveToDisk = async () => {
    if (!token) {
      // Redirect to login
      navigate('/login', { state: { returnTo: `/share/${identity}` } })
      return
    }

    if (!shareInfo || !shareInfo.repository_identity) return

    setShowSaveModal(true)
  }

  const handleSave = async () => {
    if (!shareInfo || !shareInfo.repository_identity) return
    if (targetFolderId === null) {
      alert('Please select a folder')
      return
    }

    try {
      setSaving(true)
      await fileService.saveShareFile(shareInfo.repository_identity, targetFolderId)
      alert('File saved to your dist successfully!')
      setShowSaveModal(false)
      setTargetFolderId(null)
    } catch (error) {
      console.error('Save failed:', error)
      // Handle 401 unauthorized error
      if (error.response?.status === 401) {
        alert('Please login to save files to your dist')
        setShowSaveModal(false)
        navigate('/login', { state: { returnTo: `/share/${identity}` } })
      } else {
        alert('Failed to save file: ' + (error.response?.data?.error || error.message))
      }
    } finally {
      setSaving(false)
    }
  }

  const formatFileSize = (bytes) => {
    if (!bytes || bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
  }

  if (loading) {
    return (
      <div className="share-detail">
        <div className="loading">
          <div>{t('loading')}</div>
        </div>
      </div>
    )
  }

  if (error || !shareInfo) {
    return (
      <div className="share-detail">
        <div className="error-state">
          <h2>{t('shareNotFound')}</h2>
          <p>{error || t('shareLinkInvalid')}</p>
          <button onClick={() => navigate('/')} className="btn-primary">
            {t('goHome')}
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="share-detail">
      <div className="share-detail-container">
        <h1>{t('shareDetail')}</h1>
        
        <div className="share-file-card" onClick={handlePreview} style={{ cursor: 'pointer' }}>
          <div className="file-icon-large">
            <FileIcon ext={shareInfo.ext || ''} />
          </div>
          
          <div className="file-info">
            <h2>{shareInfo.name}{shareInfo.ext}</h2>
            <p className="file-size">{formatFileSize(shareInfo.size)}</p>
          </div>

          <div className="share-actions" onClick={(e) => e.stopPropagation()}>
            <button onClick={handleDownload} className="btn-primary">
              {t('download')}
            </button>
            <button onClick={handleSaveToDisk} className="btn-default">
              {t('saveToDisk')}
            </button>
          </div>
        </div>
      </div>

      {showPreviewModal && shareInfo && (
        <div className="modal-overlay" onClick={() => setShowPreviewModal(false)}>
          <div className="modal-content preview-modal" onClick={(e) => e.stopPropagation()}>
            <div className="preview-header">
              <h3>{shareInfo.name}{shareInfo.ext}</h3>
              <button onClick={() => setShowPreviewModal(false)} className="close-btn">×</button>
            </div>
            <div className="preview-body">
              <video
                src={shareInfo.path}
                controls
                autoPlay
                style={{ width: '100%', maxHeight: '80vh' }}
              >
                Your browser does not support video playback.
              </video>
            </div>
          </div>
        </div>
      )}

      {showSaveModal && (
        <div className="modal-overlay" onClick={() => {
          if (!saving) {
            setShowSaveModal(false)
            setTargetFolderId(null)
          }
        }}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h3>{t('saveToDisk')}</h3>
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

export default ShareDetail

