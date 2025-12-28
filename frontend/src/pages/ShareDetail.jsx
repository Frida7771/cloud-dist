import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useApp } from '../contexts/AppContext'
import { useAuth } from '../contexts/AuthContext'
import { fileService } from '../services/fileService'
import api from '../services/api'
import './ShareDetail.css'

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

  const handleDownload = async () => {
    if (!shareInfo || !identity) return
    
    try {
      // Use share download endpoint which doesn't require authentication
      const response = await fileService.downloadShareFile(identity)
      
      let downloadFileName = shareInfo.name + shareInfo.ext
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
      alert('Download failed: ' + (error.response?.data?.error || error.message))
    }
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
      alert('File saved to your disk successfully!')
      setShowSaveModal(false)
      setTargetFolderId(null)
    } catch (error) {
      console.error('Save failed:', error)
      alert('Failed to save file: ' + (error.response?.data?.error || error.message))
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
        
        <div className="share-file-card">
          <div className="file-icon-large">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z" fill="#3385ff"/>
            </svg>
          </div>
          
          <div className="file-info">
            <h2>{shareInfo.name}{shareInfo.ext}</h2>
            <p className="file-size">{formatFileSize(shareInfo.size)}</p>
          </div>

          <div className="share-actions">
            <button onClick={handleDownload} className="btn-primary">
              {t('download')}
            </button>
            <button onClick={handleSaveToDisk} className="btn-default">
              {t('saveToDisk')}
            </button>
          </div>
        </div>
      </div>

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

