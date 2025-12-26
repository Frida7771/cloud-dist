import { useState, useEffect } from 'react'
import { fileService } from '../services/fileService'
import './Files.css'

function Files() {
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

  useEffect(() => {
    loadFiles()
  }, [currentPath])

  useEffect(() => {
    if (showUploadModal) {
      loadAllFolders()
    }
  }, [showUploadModal])

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
    } catch (error) {
      console.error('Upload failed:', error)
      alert('Upload failed: ' + (error.response?.data?.error || error.message))
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
      const parentId = currentPath[currentPath.length - 1].id
      const response = await fileService.createFolder(parentId, newFolderName.trim())
      setNewFolderName('')
      setShowCreateFolder(false)
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

  return (
    <div className="files">
      <div className="files-header">
        <h2>Files</h2>
        <div className="action-buttons">
          <button
            onClick={() => setShowUploadModal(true)}
            className="action-btn upload-btn"
          >
            Upload File
          </button>
          <button
            onClick={() => setShowCreateFolder(true)}
            className="action-btn create-folder-btn"
          >
            Create Folder
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
            <h3>Upload File</h3>
            
            <div className="form-group">
              <label>Select File</label>
              <input
                type="file"
                onChange={handleFileSelect}
                disabled={uploadProgress > 0}
              />
              {selectedFile && (
                <div className="file-info">
                  <span>Selected: {selectedFile.name}</span>
                  <span>Size: {(selectedFile.size / (1024 * 1024)).toFixed(2)} MB</span>
                </div>
              )}
            </div>

            <div className="form-group">
              <label>Upload To Folder (Required)</label>
              <select
                value={selectedFolderId || ''}
                onChange={(e) => setSelectedFolderId(e.target.value ? Number(e.target.value) : null)}
                disabled={uploadProgress > 0}
                required
              >
                <option value="">-- Select a folder --</option>
                {folders
                  .filter(folder => folder.id !== 0) // Exclude root directory
                  .map((folder, index) => (
                    <option key={folder.id || index} value={folder.id}>
                      {'  '.repeat(folder.level || 0)}{folder.name}
                    </option>
                  ))}
              </select>
              {folders.filter(f => f.id !== 0).length === 0 && (
                <p style={{ color: '#e74c3c', fontSize: '0.9rem', marginTop: '0.5rem' }}>
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
              >
                {uploadProgress > 0 ? 'Uploading...' : 'Upload'}
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
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {showCreateFolder && (
        <div className="modal-overlay" onClick={() => setShowCreateFolder(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h3>Create Folder</h3>
            <input
              type="text"
              placeholder="Folder name"
              value={newFolderName}
              onChange={(e) => setNewFolderName(e.target.value)}
              onKeyPress={(e) => {
                if (e.key === 'Enter') handleCreateFolder()
              }}
              autoFocus
            />
            <div className="modal-actions">
              <button onClick={handleCreateFolder}>Create</button>
              <button onClick={() => {
                setShowCreateFolder(false)
                setNewFolderName('')
              }}>Cancel</button>
            </div>
          </div>
        </div>
      )}

      <div className="breadcrumb">
        {currentPath.map((path, index) => (
          <span key={index}>
            <button onClick={() => handleBreadcrumbClick(index)}>
              {path.name}
            </button>
            {index < currentPath.length - 1 && ' / '}
          </span>
        ))}
      </div>

      {loading ? (
        <div className="loading">Loading...</div>
      ) : (
        <div className="file-list">
          {(() => {
            // Root directory: only show folders
            // Inside folder: only show files (not subfolders)
            const isRoot = currentPath.length === 1
            const displayItems = isRoot
              ? files.filter(file => file.ext === '') // Only folders in root
              : files.filter(file => file.ext !== '') // Only files inside folder
            
            if (displayItems.length === 0) {
              return (
                <div className="empty">
                  {isRoot 
                    ? 'No folders in root directory' 
                    : 'No files in this folder'}
                </div>
              )
            }
            
            return displayItems.map((file) => (
              <div key={file.identity} className="file-item">
                <div className="file-info">
                  <span className="file-icon">
                    {file.ext === '' ? '📁' : '📄'}
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
                    <>
                      <span className="file-name">{file.name}</span>
                      {/* Only show size for files, not folders */}
                      {file.ext !== '' && file.size && file.size > 0 && (
                        <span className="file-size">
                          {(file.size / (1024 * 1024)).toFixed(2)} MB
                        </span>
                      )}
                    </>
                  )}
                </div>
                <div className="file-actions">
                  {file.ext === '' ? (
                    <>
                      <button onClick={() => handleFolderClick(file)}>Open</button>
                      <button onClick={() => startRename(file)}>Rename</button>
                      <button onClick={() => handleDelete(file.identity)}>Delete</button>
                    </>
                  ) : (
                    <>
                      <button onClick={() => handleDownload(file.repository_identity, file.name + file.ext)}>
                        Download
                      </button>
                      <button onClick={() => startRename(file)}>Rename</button>
                      <button onClick={() => handleDelete(file.identity)}>Delete</button>
                    </>
                  )}
                </div>
              </div>
            ))
          })()}
        </div>
      )}
    </div>
  )
}

export default Files

