import { createContext, useContext, useState, useEffect } from 'react'

const AppContext = createContext()

const translations = {
  en: {
    // Common
    login: 'Login',
    register: 'Register',
    logout: 'Logout',
    username: 'Username',
    password: 'Password',
    email: 'Email',
    confirmPassword: 'Confirm Password',
    verificationCode: 'Verification Code',
    sendCode: 'Send Code',
    codeSent: 'Code Sent',
    submit: 'Submit',
    cancel: 'Cancel',
    save: 'Save',
    delete: 'Delete',
    edit: 'Edit',
    create: 'Create',
    upload: 'Upload',
    download: 'Download',
    rename: 'Rename',
    move: 'Move',
    open: 'Open',
    close: 'Close',
    loading: 'Loading...',
    error: 'Error',
    success: 'Success',
    
    // Auth
    loginTitle: 'Login',
    registerTitle: 'Register',
    forgotPassword: 'Forgot Password',
    forgotPasswordTitle: 'Forgot Password',
    dontHaveAccount: "Don't have an account?",
    alreadyHaveAccount: 'Already have an account?',
    rememberPassword: 'Remember your password?',
    loggingIn: 'Logging in...',
    registering: 'Registering...',
    userNotRegistered: 'User not registered',
    passwordIncorrect: 'Password incorrect',
    loginFailed: 'Login failed',
    
    // Files
    myFiles: 'My Files',
    share: 'Share',
    profile: 'Profile',
    newFolder: 'New Folder',
    uploadFile: 'Upload File',
    selectFile: 'Select File',
    uploadToFolder: 'Upload To Folder (Required)',
    selectFolder: 'Select a folder',
    selectTargetFolder: 'Select Target Folder',
    root: 'Root',
    uploading: 'Uploading...',
    noFolders: 'No folders',
    noFiles: 'No files',
    clickNewFolder: 'Click "New Folder" to get started',
    clickUpload: 'Click "Upload" to add files',
    name: 'Name',
    size: 'Size',
    actions: 'Actions',
    preview: 'Preview',
    previewNotAvailable: 'Preview not available',
    previewNotSupported: 'Preview not supported for this file type',
    videoNotSupported: 'Your browser does not support video playback',
    audioNotSupported: 'Your browser does not support audio playback',
    docPreviewNotSupported: 'Word documents cannot be previewed in the browser',
    docPreviewReason: 'Browser does not natively support Word format. Please download the file to view it with Microsoft Word or other compatible software.',
    docOldFormatNote: 'Note: Old Word format (.doc) requires Microsoft Word or compatible software.',
    pleaseDownloadToView: 'Please download the file to view it.',
    toView: 'to View',
    
    // Profile
    userInfo: 'User Info',
    changePassword: 'Change Password',
    friends: 'Friends',
    buyStorage: 'Buy Storage',
    storageUsage: 'Storage Usage',
    
    // Settings
    language: 'Language',
    theme: 'Theme',
    lightMode: 'Light Mode',
    darkMode: 'Dark Mode',
    english: 'English',
    chinese: '中文',
  },
  zh: {
    // Common
    login: '登录',
    register: '注册',
    logout: '退出',
    username: '用户名',
    password: '密码',
    email: '邮箱',
    confirmPassword: '确认密码',
    verificationCode: '验证码',
    sendCode: '发送验证码',
    codeSent: '已发送',
    submit: '提交',
    cancel: '取消',
    save: '保存',
    delete: '删除',
    edit: '编辑',
    create: '创建',
    upload: '上传',
    download: '下载',
    rename: '重命名',
    move: '移动',
    open: '打开',
    close: '关闭',
    loading: '加载中...',
    error: '错误',
    success: '成功',
    
    // Auth
    loginTitle: '登录',
    registerTitle: '注册',
    forgotPassword: '忘记密码',
    forgotPasswordTitle: '忘记密码',
    dontHaveAccount: '还没有账号？',
    alreadyHaveAccount: '已有账号？',
    rememberPassword: '想起密码了？',
    loggingIn: '登录中...',
    registering: '注册中...',
    userNotRegistered: '用户未注册',
    passwordIncorrect: '密码错误',
    loginFailed: '登录失败',
    
    // Files
    myFiles: '我的文件',
    share: '分享',
    profile: '个人中心',
    newFolder: '新建文件夹',
    uploadFile: '上传文件',
    selectFile: '选择文件',
    uploadToFolder: '上传到文件夹（必选）',
    selectFolder: '选择文件夹',
    selectTargetFolder: '选择目标文件夹',
    root: '根目录',
    uploading: '上传中...',
    noFolders: '暂无文件夹',
    noFiles: '暂无文件',
    clickNewFolder: '点击"新建文件夹"开始使用',
    clickUpload: '点击"上传"添加文件',
    name: '名称',
    size: '大小',
    actions: '操作',
    preview: '预览',
    previewNotAvailable: '预览不可用',
    previewNotSupported: '此文件类型不支持预览',
    videoNotSupported: '您的浏览器不支持视频播放',
    audioNotSupported: '您的浏览器不支持音频播放',
    docPreviewNotSupported: 'Word 文档无法在浏览器中预览',
    docPreviewReason: '浏览器原生不支持 Word 格式。请下载文件后使用 Microsoft Word 或其他兼容软件查看。',
    docOldFormatNote: '注意：旧版 Word 格式 (.doc) 需要 Microsoft Word 或兼容软件。',
    pleaseDownloadToView: '请下载文件后查看。',
    toView: '以查看',
    
    // Profile
    userInfo: '用户信息',
    changePassword: '修改密码',
    friends: '好友',
    buyStorage: '购买存储',
    storageUsage: '存储使用',
    
    // Settings
    language: '语言',
    theme: '主题',
    lightMode: '白天模式',
    darkMode: '黑夜模式',
    english: 'English',
    chinese: '中文',
  },
}

export function AppProvider({ children }) {
  const [language, setLanguage] = useState(() => {
    return localStorage.getItem('language') || 'en'
  })
  const [theme, setTheme] = useState(() => {
    return localStorage.getItem('theme') || 'light'
  })

  useEffect(() => {
    localStorage.setItem('language', language)
  }, [language])

  // Initialize theme on mount
  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme)
  }, [])

  useEffect(() => {
    localStorage.setItem('theme', theme)
    document.documentElement.setAttribute('data-theme', theme)
  }, [theme])

  const t = (key) => {
    return translations[language]?.[key] || key
  }

  const toggleLanguage = () => {
    setLanguage(prev => prev === 'en' ? 'zh' : 'en')
  }

  const toggleTheme = () => {
    setTheme(prev => prev === 'light' ? 'dark' : 'light')
  }

  return (
    <AppContext.Provider value={{ language, theme, t, toggleLanguage, toggleTheme }}>
      {children}
    </AppContext.Provider>
  )
}

export const useApp = () => {
  const context = useContext(AppContext)
  if (!context) {
    throw new Error('useApp must be used within AppProvider')
  }
  return context
}

