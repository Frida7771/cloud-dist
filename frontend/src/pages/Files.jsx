import { useState, useEffect } from 'react'
import { fileService } from '../services/fileService'
import { useApp } from '../contexts/AppContext'
import './Files.css'

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

function Files() {
  const { t } = useApp()
  const [files, setFiles] = useState([])
  const [currentPath, setCurrentPath] = useState([{ id: 0, identity: '', name: 'Root' }])
  const [loading, setLoading] = useState(false)
  const [uploadProgress, setUploadProgress] = useState(0)
  const [showCreateFolder, setShowCreateFolder] = useState(false)
  const [showUploadModal, setShowUploadModal] = useState(false)
  const [selectedFile, setSelectedFile] = useState(null)
  const [selectedFolderId, setSelectedFolderId] = useState(null)
  const [folders, setFolders] = useState([])
  const [newFolderName, setNewFolderName] = useState('')
  const [editingFile, setEditingFile] = useState(null)
  const [editFileName, setEditFileName] = useState('')
  const [showMoveModal, setShowMoveModal] = useState(false)
  const [fileToMove, setFileToMove] = useState(null)
  const [moveTargetFolderId, setMoveTargetFolderId] = useState(null)
  const [showPreviewModal, setShowPreviewModal] = useState(false)
  const [previewFile, setPreviewFile] = useState(null)
  const [previewUrl, setPreviewUrl] = useState(null)
  const [previewContent, setPreviewContent] = useState(null)
  const [previewLoading, setPreviewLoading] = useState(false)

  useEffect(() => {
    loadFiles()
  }, [currentPath])

  useEffect(() => {
    if (showUploadModal) {
      loadAllFolders()
    }
  }, [showUploadModal])

  useEffect(() => {
    if (showMoveModal) {
      loadAllFolders()
    }
  }, [showMoveModal])

  const loadFiles = async () => {
    setLoading(true)
    try {
      // Use identity instead of id for API call
      // If in root (id === 0), pass empty string
      const currentFolder = currentPath[currentPath.length - 1]
      const folderIdentity = currentFolder.id === 0 ? '' : (currentFolder.identity || '')
      const response = await fileService.getFileList(folderIdentity)
      setFiles(response.data.list || [])
    } catch (error) {
      console.error('Failed to load files:', error)
    } finally {
      setLoading(false)
    }
  }

  const loadAllFolders = async () => {
    try {
      // Build folder list from current path (breadcrumb)
      // Each item in currentPath has id and name
      const folderOptions = currentPath.map((pathItem, index) => ({
        id: pathItem.id,
        name: pathItem.name,
        level: index,
      }))
      
      // Load all folders recursively to build a complete list
      const allFolders = [...folderOptions]
      
      // Recursively load all folders starting from root
      await loadFoldersRecursive('', allFolders, 0)
      
      setFolders(allFolders)
      // Set default to current folder (last item in currentPath)
      // But only if it's not root (id !== 0)
      const currentFolderId = currentPath[currentPath.length - 1].id
      setSelectedFolderId(currentFolderId !== 0 ? currentFolderId : null)
    } catch (error) {
      console.error('Failed to load folders:', error)
      // Fallback to just current path
      const folderOptions = currentPath.map((pathItem, index) => ({
        id: pathItem.id,
        name: pathItem.name,
        level: index,
      }))
      setFolders(folderOptions)
      const currentFolderId = currentPath[currentPath.length - 1].id
      setSelectedFolderId(currentFolderId !== 0 ? currentFolderId : null)
    }
  }

  const loadFoldersRecursive = async (parentIdentity, allFolders, level) => {
    try {
      const response = await fileService.getFolderList(parentIdentity)
      const folders = response.data.list || []
      
      for (const folder of folders) {
        // Try to find the folder in current files to get its database ID
        const fileItem = files.find(f => f.identity === folder.identity && f.ext === '')
        const folderId = fileItem ? fileItem.id : 0
        
        // Only add if not already in list (avoid duplicates from currentPath)
        if (!allFolders.find(f => f.id === folderId && f.id !== 0)) {
          allFolders.push({
            id: folderId,
            identity: folder.identity,
            name: folder.name,
            level: level + 1,
          })
        }
        
        // Recursively load subfolders
        await loadFoldersRecursive(folder.identity, allFolders, level + 1)
      }
    } catch (error) {
      console.error('Failed to load folders recursively:', error)
    }
  }

  const handleFileSelect = (e) => {
    const file = e.target.files[0]
    if (file) {
      setSelectedFile(file)
    }
  }

  const handleUpload = async () => {
    if (!selectedFile) {
      alert('Please select a file')
      return
    }

    if (!selectedFolderId || selectedFolderId === 0 || selectedFolderId === '') {
      alert('Please select a folder to upload the file to')
      return
    }

    try {
      setUploadProgress(0)
      
      // Upload file
      const uploadResponse = await fileService.uploadFile(selectedFile, (progress) => {
        setUploadProgress(progress)
      })
      
      // Save to user repository with selected folder
      try {
        await fileService.saveToRepository(
          selectedFolderId,
          uploadResponse.data.identity,
          uploadResponse.data.ext,
          uploadResponse.data.name
        )
        
        // Close modal and refresh
        setShowUploadModal(false)
        setSelectedFile(null)
        setSelectedFolderId(null)
        setUploadProgress(0)
        loadFiles()
      } catch (saveError) {
        // Check if file already exists
        const errorMessage = saveError.response?.data?.error || saveError.message
        if (errorMessage && errorMessage.toLowerCase().includes('already exists')) {
          alert('File already exists')
          // Still close modal and refresh to show existing file
          setShowUploadModal(false)
          setSelectedFile(null)
          setSelectedFolderId(null)
          setUploadProgress(0)
          loadFiles()
        } else {
          throw saveError // Re-throw other errors
        }
      }
    } catch (error) {
      console.error('Upload failed:', error)
      alert('Upload failed: ' + (error.response?.data?.error || error.message))
      setUploadProgress(0)
    }
  }

  const handleDelete = async (identity) => {
    if (!confirm('Are you sure you want to delete this file?')) return

    try {
      await fileService.deleteFile(identity)
      loadFiles()
    } catch (error) {
      console.error('Delete failed:', error)
      alert('Delete failed')
    }
  }

  const handleDownload = async (repositoryIdentity, fileName) => {
    try {
      const response = await fileService.downloadFile(repositoryIdentity)
      
      // Get the file name from Content-Disposition header or use provided fileName
      let downloadFileName = fileName || 'file'
      const contentDisposition = response.headers['content-disposition']
      if (contentDisposition) {
        // Try to extract filename from Content-Disposition header
        const fileNameMatch = contentDisposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/)
        if (fileNameMatch && fileNameMatch[1]) {
          downloadFileName = decodeURIComponent(fileNameMatch[1].replace(/['"]/g, ''))
        }
      }
      
      // Get content type from response headers
      const contentType = response.headers['content-type'] || 'application/octet-stream'
      
      // Create blob with correct MIME type
      const blob = new Blob([response.data], { type: contentType })
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', downloadFileName)
      document.body.appendChild(link)
      link.click()
      link.remove()
      // Clean up the object URL
      setTimeout(() => window.URL.revokeObjectURL(url), 100)
    } catch (error) {
      console.error('Download failed:', error)
      alert('Download failed: ' + (error.response?.data?.error || error.message))
    }
  }

  const handleFilePreview = async (file) => {
    if (file.ext === '') {
      // It's a folder, open it instead
      handleFolderClick(file)
      return
    }

    setPreviewFile(file)
    setShowPreviewModal(true)
    setPreviewLoading(true)
    setPreviewUrl(null)
    setPreviewContent(null)

    try {
      const response = await fileService.downloadFile(file.repository_identity)
      const contentType = response.headers['content-type'] || 'application/octet-stream'
      const extLower = file.ext.toLowerCase()
      
      // For text files, read as text
      if (['.txt', '.md', '.log', '.csv', '.js', '.jsx', '.ts', '.tsx', '.py', '.java', '.cpp', '.c', '.h', '.html', '.css', '.json', '.xml', '.yaml', '.yml'].includes(extLower) || contentType.startsWith('text/')) {
        // response.data is already a Blob, convert to text
        const text = await response.data.text()
        setPreviewContent(text)
      } else if (['.doc', '.docx'].includes(extLower)) {
        // Word documents cannot be previewed, show download option
        setPreviewUrl(null)
        setPreviewContent(null)
      } else {
        // For other files (images, PDF, video, audio), create blob URL
        // response.data is already a Blob
        const url = window.URL.createObjectURL(response.data)
        setPreviewUrl(url)
      }
    } catch (error) {
      console.error('Preview failed:', error)
      alert('Preview failed: ' + (error.response?.data?.error || error.message))
      setShowPreviewModal(false)
    } finally {
      setPreviewLoading(false)
    }
  }

  const closePreview = () => {
    if (previewUrl) {
      window.URL.revokeObjectURL(previewUrl)
    }
    setShowPreviewModal(false)
    setPreviewFile(null)
    setPreviewUrl(null)
    setPreviewContent(null)
    setPreviewLoading(false)
  }

  const canPreview = (ext) => {
    const extLower = ext.toLowerCase()
    const previewableTypes = [
      // Images
      '.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp', '.svg',
      // PDF
      '.pdf',
      // Text
      '.txt', '.md', '.log', '.csv',
      // Code
      '.js', '.jsx', '.ts', '.tsx', '.py', '.java', '.cpp', '.c', '.h',
      '.html', '.css', '.json', '.xml', '.yaml', '.yml',
      // Video
      '.mp4', '.webm', '.mov',
      // Audio
      '.mp3', '.wav', '.ogg',
      // Office documents (will show preview option, but may need download)
      '.doc', '.docx'
    ]
    return previewableTypes.includes(extLower)
  }

  const getPreviewComponent = (file, url, content, loading) => {
    if (loading) {
      return (
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <div>{t('loading')}</div>
        </div>
      )
    }

    if (!url && !content) {
      return (
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <div>{t('previewNotAvailable') || 'Preview not available'}</div>
        </div>
      )
    }

    const extLower = file.ext.toLowerCase()

    // Text files - display as text
    if (content !== null) {
      return (
        <pre style={{ 
          textAlign: 'left', 
          padding: '20px', 
          backgroundColor: '#f5f5f5', 
          borderRadius: '4px',
          overflow: 'auto',
          maxHeight: '80vh',
          margin: 0,
          fontSize: '14px',
          lineHeight: '1.5',
          fontFamily: 'monospace'
        }}>
          {content}
        </pre>
      )
    }

    // Images
    if (['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp', '.svg'].includes(extLower)) {
      return (
        <img 
          src={url} 
          alt={file.name}
          style={{ maxWidth: '100%', maxHeight: '80vh', objectFit: 'contain' }}
        />
      )
    }

    // PDF
    if (extLower === '.pdf') {
      return (
        <iframe
          src={url}
          style={{ width: '100%', height: '80vh', border: 'none' }}
          title={file.name}
        />
      )
    }

    // Video
    if (['.mp4', '.webm', '.mov'].includes(extLower)) {
      return (
        <video
          src={url}
          controls
          style={{ maxWidth: '100%', maxHeight: '80vh' }}
        >
          {t('videoNotSupported') || 'Your browser does not support video playback'}
        </video>
      )
    }

    // Audio
    if (['.mp3', '.wav', '.ogg'].includes(extLower)) {
      return (
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <audio src={url} controls style={{ width: '100%' }}>
            {t('audioNotSupported') || 'Your browser does not support audio playback'}
          </audio>
        </div>
      )
    }

    // Word documents - cannot preview directly in browser
    if (['.doc', '.docx'].includes(extLower)) {
      return (
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <div style={{ marginBottom: '20px', fontSize: '16px', color: '#666' }}>
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" style={{ width: '64px', height: '64px', margin: '0 auto 20px', opacity: 0.5 }}>
              <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z" fill="#2b579a"/>
            </svg>
            <p style={{ margin: '10px 0', fontWeight: '500' }}>
              {t('docPreviewNotSupported') || 'Word documents cannot be previewed in the browser'}
            </p>
            <p style={{ margin: '10px 0', fontSize: '14px', color: '#999', lineHeight: '1.6' }}>
              {t('docPreviewReason') || 'Browser does not natively support Word format. Please download the file to view it with Microsoft Word or other compatible software.'}
            </p>
            {extLower === '.doc' && (
              <p style={{ margin: '10px 0', fontSize: '13px', color: '#ff9800' }}>
                {t('docOldFormatNote') || 'Note: Old Word format (.doc) requires Microsoft Word or compatible software.'}
              </p>
            )}
          </div>
          <button onClick={() => handleDownload(file.repository_identity, file.name + file.ext)} className="btn-primary" style={{ marginTop: '10px' }}>
            {t('download')} {t('toView') || 'to View'}
          </button>
        </div>
      )
    }

    // Default: show download option
    return (
      <div style={{ textAlign: 'center', padding: '40px' }}>
        <p>{t('previewNotSupported') || 'Preview not supported for this file type'}</p>
        <button onClick={() => handleDownload(file.repository_identity, file.name + file.ext)} className="btn-primary">
          {t('download')}
        </button>
      </div>
    )
  }

  const handleFolderClick = (folder) => {
    setCurrentPath([...currentPath, { id: folder.id, identity: folder.identity, name: folder.name }])
  }

  const handleBreadcrumbClick = (index) => {
    setCurrentPath(currentPath.slice(0, index + 1))
  }

  const handleCreateFolder = async () => {
    if (!newFolderName.trim()) {
      alert('Please enter folder name')
      return
    }

    try {
      // Get parent ID: if in root (id === 0), use 0, otherwise use the folder's database ID
      const currentFolder = currentPath[currentPath.length - 1]
      const parentId = currentFolder.id === 0 ? 0 : currentFolder.id
      console.log('Creating folder with parentId:', parentId, 'currentPath:', currentPath)
      const response = await fileService.createFolder(parentId, newFolderName.trim())
      setNewFolderName('')
      setShowCreateFolder(false)
      // Refresh file list to show the new folder
      loadFiles()
    } catch (error) {
      console.error('Create folder failed:', error)
      alert('Failed to create folder: ' + (error.response?.data?.error || error.message))
    }
  }

  const handleRename = async (identity, currentName) => {
    if (!editFileName.trim()) {
      setEditingFile(null)
      return
    }

    try {
      await fileService.renameFile(identity, editFileName.trim())
      setEditingFile(null)
      setEditFileName('')
      loadFiles()
    } catch (error) {
      console.error('Rename failed:', error)
      alert('Failed to rename: ' + (error.response?.data?.error || error.message))
    }
  }

  const startRename = (file) => {
    setEditingFile(file.identity)
    setEditFileName(file.name)
  }

  const startMove = (file) => {
    setFileToMove(file)
    setMoveTargetFolderId(null)
    setShowMoveModal(true)
  }

  const handleMove = async () => {
    if (!fileToMove) {
      alert('No file selected')
      return
    }

    // Allow moving to root (moveTargetFolderId === 0 or null)
    // For root, we use empty string as parent_identity
    let targetParentIdentity = ''
    
    if (moveTargetFolderId !== null && moveTargetFolderId !== 0 && moveTargetFolderId !== '') {
      // Find the target folder identity
      const targetFolder = folders.find(f => f.id === moveTargetFolderId)
      if (!targetFolder) {
        alert('Invalid target folder')
        return
      }
      
      // If target folder has identity, use it; otherwise it's root
      if (targetFolder.identity) {
        targetParentIdentity = targetFolder.identity
      }
      
      // Prevent moving to the same folder
      const currentFolder = currentPath[currentPath.length - 1]
      if (targetFolder.id === currentFolder.id) {
        alert('File is already in this folder')
        return
      }
    } else {
      // Moving to root - check if already in root
      const currentFolder = currentPath[currentPath.length - 1]
      if (currentFolder.id === 0) {
        alert('File is already in root directory')
        return
      }
    }

    // Prevent moving a folder into itself
    if (fileToMove.ext === '' && moveTargetFolderId === fileToMove.id) {
      alert('Cannot move a folder into itself')
      return
    }

    try {
      await fileService.moveFile(fileToMove.identity, targetParentIdentity)
      setShowMoveModal(false)
      setFileToMove(null)
      setMoveTargetFolderId(null)
      loadFiles()
    } catch (error) {
      console.error('Move failed:', error)
      alert('Failed to move file: ' + (error.response?.data?.error || error.message))
    }
  }

  return (
    <div className="files">
      <div className="files-toolbar">
        <div className="breadcrumb">
          {currentPath.map((path, index) => (
            <span key={index} className="breadcrumb-item">
              <button onClick={() => handleBreadcrumbClick(index)}>
                {path.name}
              </button>
              {index < currentPath.length - 1 && <span className="separator">/</span>}
            </span>
          ))}
        </div>
        <div className="toolbar-actions">
          <button
            onClick={() => setShowUploadModal(true)}
            className="btn-primary"
          >
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" fill="currentColor"/>
            </svg>
            {t('upload')}
          </button>
          <button
            onClick={() => setShowCreateFolder(true)}
            className="btn-default"
          >
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M10 4H4c-1.11 0-2 .89-2 2v12c0 1.11.89 2 2 2h16c1.11 0 2-.89 2-2V8c0-1.11-.89-2-2-2h-8l-2-2z" fill="currentColor"/>
            </svg>
            {t('newFolder')}
          </button>
        </div>
      </div>

      {showUploadModal && (
        <div className="modal-overlay" onClick={() => {
          if (uploadProgress === 0) {
            setShowUploadModal(false)
            setSelectedFile(null)
            setSelectedFolderId(null)
          }
        }}>
          <div className="modal-content upload-modal" onClick={(e) => e.stopPropagation()}>
            <h3>{t('uploadFile')}</h3>
            
            <div className="form-group">
              <label>{t('selectFile')}</label>
              <input
                type="file"
                onChange={handleFileSelect}
                disabled={uploadProgress > 0}
              />
              {selectedFile && (
                <div className="file-info">
                  <span>File: {selectedFile.name}</span>
                  <span>Size: {(selectedFile.size / (1024 * 1024)).toFixed(2)} MB</span>
                </div>
              )}
            </div>

            <div className="form-group">
              <label>{t('uploadToFolder')}</label>
              <select
                value={selectedFolderId || ''}
                onChange={(e) => setSelectedFolderId(e.target.value ? Number(e.target.value) : null)}
                disabled={uploadProgress > 0}
                required
              >
                <option value="">-- {t('selectFolder')} --</option>
                {folders
                  .filter(folder => folder.id !== 0) // Exclude root directory
                  .map((folder, index) => (
                    <option key={folder.id || index} value={folder.id}>
                      {'  '.repeat(folder.level || 0)}{folder.name}
                    </option>
                  ))}
              </select>
              {folders.filter(f => f.id !== 0).length === 0 && (
                <p style={{ color: '#f5222d', fontSize: '0.9rem', marginTop: '0.5rem' }}>
                  No folders available. Please create a folder first.
                </p>
              )}
            </div>

            {uploadProgress > 0 && (
              <div className="upload-progress-bar">
                <div className="progress-fill" style={{ width: `${uploadProgress}%` }}></div>
                <span>Uploading: {uploadProgress}%</span>
              </div>
            )}

            <div className="modal-actions">
              <button
                onClick={handleUpload}
                disabled={!selectedFile || uploadProgress > 0}
                className="btn-primary"
              >
                {uploadProgress > 0 ? t('uploading') : t('upload')}
              </button>
              <button
                onClick={() => {
                  if (uploadProgress === 0) {
                    setShowUploadModal(false)
                    setSelectedFile(null)
                    setSelectedFolderId(0)
                  }
                }}
                disabled={uploadProgress > 0}
                className="btn-default"
              >
                {t('cancel')}
              </button>
            </div>
          </div>
        </div>
      )}

      {showCreateFolder && (
        <div className="modal-overlay" onClick={() => setShowCreateFolder(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h3>{t('newFolder')}</h3>
            <input
              type="text"
              placeholder={t('newFolder')}
              value={newFolderName}
              onChange={(e) => setNewFolderName(e.target.value)}
              onKeyPress={(e) => {
                if (e.key === 'Enter') handleCreateFolder()
              }}
              autoFocus
            />
            <div className="modal-actions">
              <button onClick={handleCreateFolder} className="btn-primary">{t('create')}</button>
              <button onClick={() => {
                setShowCreateFolder(false)
                setNewFolderName('')
              }} className="btn-default">{t('cancel')}</button>
            </div>
          </div>
        </div>
      )}

      {showPreviewModal && previewFile && (
        <div className="modal-overlay" onClick={closePreview}>
          <div className="modal-content preview-modal" onClick={(e) => e.stopPropagation()} style={{ maxWidth: '90vw', maxHeight: '90vh', overflow: 'auto' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px', paddingBottom: '10px', borderBottom: '1px solid #e0e0e0' }}>
              <h3 style={{ margin: 0 }}>{previewFile.name}</h3>
              <button onClick={closePreview} style={{ background: 'none', border: 'none', fontSize: '24px', cursor: 'pointer', padding: '0 10px' }}>×</button>
            </div>
            <div style={{ textAlign: 'center', minHeight: '400px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
              {getPreviewComponent(previewFile, previewUrl, previewContent, previewLoading)}
            </div>
            <div style={{ marginTop: '20px', paddingTop: '10px', borderTop: '1px solid #e0e0e0', display: 'flex', justifyContent: 'flex-end', gap: '10px' }}>
              <button onClick={() => handleDownload(previewFile.repository_identity, previewFile.name + previewFile.ext)} className="btn-primary">
                {t('download')}
              </button>
              <button onClick={closePreview} className="btn-default">
                {t('close')}
              </button>
            </div>
          </div>
        </div>
      )}

      {showMoveModal && (
        <div className="modal-overlay" onClick={() => setShowMoveModal(false)}>
          <div className="modal-content upload-modal" onClick={(e) => e.stopPropagation()}>
            <h3>{t('move')} {fileToMove?.name}</h3>
            
            <div className="form-group">
              <label>{t('selectTargetFolder')}</label>
              <select
                value={moveTargetFolderId === null ? '' : moveTargetFolderId}
                onChange={(e) => {
                  const value = e.target.value
                  if (value === '') {
                    setMoveTargetFolderId(null)
                  } else if (value === '0') {
                    setMoveTargetFolderId(0)
                  } else {
                    setMoveTargetFolderId(Number(value))
                  }
                }}
                required
              >
                <option value="">-- {t('selectFolder')} --</option>
                <option value="0" style={{ fontWeight: 'bold' }}>📁 {t('root')}</option>
                <option disabled>──────────</option>
                {folders
                  .filter(folder => {
                    // Include all folders except:
                    // 1. The file itself if it's a folder
                    if (fileToMove && fileToMove.ext === '' && folder.id === fileToMove.id) return false
                    // 2. Current folder (file is already there)
                    const currentFolder = currentPath[currentPath.length - 1]
                    if (folder.id === currentFolder.id) return false
                    // 3. Root (already shown as separate option, not a folder)
                    if (folder.id === 0) return false
                    return true
                  })
                  .map((folder, index) => (
                    <option key={folder.id || index} value={folder.id}>
                      {'  '.repeat(folder.level || 0)}{folder.name}
                    </option>
                  ))}
              </select>
            </div>

            <div className="modal-actions">
              <button
                onClick={handleMove}
                disabled={moveTargetFolderId === null}
                className="btn-primary"
              >
                {t('move')}
              </button>
              <button
                onClick={() => {
                  setShowMoveModal(false)
                  setFileToMove(null)
                  setMoveTargetFolderId(null)
                }}
                className="btn-default"
              >
                {t('cancel')}
              </button>
            </div>
          </div>
        </div>
      )}

      <div className="files-container">
        {loading ? (
          <div className="loading">
            <div>{t('loading')}</div>
          </div>
        ) : (
          <div className="file-table-wrapper">
            {(() => {
              // Root directory: only show folders
              // Inside folder: show both files and subfolders
              const isRoot = currentPath.length === 1
              const displayItems = isRoot
                ? files.filter(file => file.ext === '') // Only folders in root
                : files // Show both files and subfolders inside folder
              
              if (displayItems.length === 0) {
                return (
                  <div className="empty-state">
                    <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                      <path d="M10 4H4c-1.11 0-2 .89-2 2v12c0 1.11.89 2 2 2h16c1.11 0 2-.89 2-2V8c0-1.11-.89-2-2-2h-8l-2-2z" fill="#d9d9d9"/>
                    </svg>
                    <div className="empty-text">
                      {isRoot ? t('noFolders') : t('noFiles')}
                    </div>
                    <div className="empty-hint">
                      {isRoot ? t('clickNewFolder') : t('clickUpload')}
                    </div>
                  </div>
                )
              }
              
              return (
                <table className="file-table">
                  <thead>
                    <tr>
                      <th style={{ width: '50%' }}>{t('name')}</th>
                      <th style={{ width: '20%' }}>{t('size')}</th>
                      <th style={{ width: '30%' }}>{t('actions')}</th>
                    </tr>
                  </thead>
                  <tbody>
                    {displayItems.map((file) => (
                      <tr key={file.identity} className="file-row">
                        <td className="file-name-cell">
                          <div className="file-name-wrapper">
                            <span className="file-icon">
                              {file.ext === '' ? (
                                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                  <path d="M10 4H4c-1.11 0-2 .89-2 2v12c0 1.11.89 2 2 2h16c1.11 0 2-.89 2-2V8c0-1.11-.89-2-2-2h-8l-2-2z" fill="#ffa940"/>
                                </svg>
                              ) : (
                                <FileIcon ext={file.ext} />
                              )}
                            </span>
                            {editingFile === file.identity ? (
                              <input
                                type="text"
                                value={editFileName}
                                onChange={(e) => setEditFileName(e.target.value)}
                                onKeyPress={(e) => {
                                  if (e.key === 'Enter') {
                                    handleRename(file.identity, file.name)
                                  } else if (e.key === 'Escape') {
                                    setEditingFile(null)
                                    setEditFileName('')
                                  }
                                }}
                                onBlur={() => handleRename(file.identity, file.name)}
                                autoFocus
                                className="rename-input"
                              />
                            ) : (
                              <span 
                                className="file-name"
                                onClick={() => {
                                  if (file.ext === '') {
                                    handleFolderClick(file)
                                  } else if (canPreview(file.ext)) {
                                    handleFilePreview(file)
                                  }
                                }}
                                style={{ cursor: (file.ext === '' || canPreview(file.ext)) ? 'pointer' : 'default' }}
                              >
                                {file.name}
                              </span>
                            )}
                          </div>
                        </td>
                        <td className="file-size-cell">
                          {file.ext !== '' && file.size && file.size > 0 
                            ? `${(file.size / (1024 * 1024)).toFixed(2)} MB`
                            : '-'
                          }
                        </td>
                        <td className="file-actions-cell">
                          <div className="action-buttons">
                            {file.ext === '' ? (
                              <>
                                <button onClick={() => handleFolderClick(file)} className="btn-link">{t('open')}</button>
                                <button onClick={() => startMove(file)} className="btn-link">{t('move')}</button>
                                <button onClick={() => startRename(file)} className="btn-link">{t('rename')}</button>
                                <button onClick={() => handleDelete(file.identity)} className="btn-link danger">{t('delete')}</button>
                              </>
                            ) : (
                              <>
                                <button onClick={() => handleDownload(file.repository_identity, file.name + file.ext)} className="btn-link">{t('download')}</button>
                                <button onClick={() => startMove(file)} className="btn-link">{t('move')}</button>
                                <button onClick={() => startRename(file)} className="btn-link">{t('rename')}</button>
                                <button onClick={() => handleDelete(file.identity)} className="btn-link danger">{t('delete')}</button>
                              </>
                            )}
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )
            })()}
          </div>
        )}
      </div>
    </div>
  )
}

export default Files

