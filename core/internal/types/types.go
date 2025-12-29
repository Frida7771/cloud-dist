package types

type FileUploadChunkCompleteRequest struct {
	Md5      string       `json:"md5"`
	Name     string       `json:"name"`
	Ext      string       `json:"ext"`
	Size     int64        `json:"size"`
	Key      string       `json:"key"`
	UploadId string       `json:"upload_id"`
	Parts    []UploadPart `json:"cos_objects"`
}

type UploadPart struct {
	PartNumber int    `json:"part_number"`
	Etag       string `json:"etag"`
}

type FileUploadChunkCompleteReply struct {
	Identity string `json:"identity"` // Repository pool identity
}

type FileUploadChunkRequest struct {
}

type FileUploadChunkReply struct {
	Etag string `json:"etag"` // MD5
}

type FileUploadPrepareRequest struct {
	Md5  string `json:"md5"`
	Name string `json:"name"`
	Ext  string `json:"ext"`
}

type FileUploadPrepareReply struct {
	Identity string       `json:"identity"`
	UploadId string       `json:"upload_id"`
	Key      string       `json:"key"`
	Parts    []UploadPart `json:"parts,optional"` // Already uploaded parts for resume
}

type RefreshAuthorizationRequest struct {
}

type RefreshAuthorizationReply struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type UserPasswordUpdateRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	Code        string `json:"code"` // Email verification code
}

type UserPasswordUpdateReply struct {
}

type ShareBasicSaveRequest struct {
	RepositoryIdentity string `json:"repository_identity"`
	ParentId           int64  `json:"parent_id"`
}

type ShareBasicSaveReply struct {
	Identity string `json:"identity"`
}

type ShareBasicDetailRequest struct {
	Identity string `json:"identity,optional"`
}

type ShareBasicDetailReply struct {
	RepositoryIdentity string `json:"repository_identity"`
	Name               string `json:"name"`
	Ext                string `json:"ext"`
	Size               int64  `json:"size"`
	Path               string `json:"path"`         // Presigned URL for preview
	DownloadUrl        string `json:"download_url"` // Presigned URL for download (with Content-Disposition: attachment)
}

type ShareBasicCreateRequest struct {
	UserRepositoryIdentity string `json:"user_repository_identity"`
	ExpiredTime            int    `json:"expired_time"`
}

type ShareBasicCreateReply struct {
	Identity string `json:"identity"`
}

type UserFileMoveRequest struct {
	Idnetity       string `json:"identity"`
	ParentIdnetity string `json:"parent_identity"`
}

type UserFileMoveReply struct {
}

type UserFileDeleteRequest struct {
	Identity string `json:"identity"`
}

type UserFileDeleteReply struct {
}

type UserFolderCreateRequest struct {
	ParentId int64  `json:"parent_id"`
	Name     string `json:"name"`
}

type UserFolderCreateReply struct {
	Identity string `json:"identity"`
}

type UserFileNameUpdateRequest struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
}

type UserFileNameUpdateReply struct {
}

type UserFileListRequest struct {
	Identity string `json:"identity,optional"`
	Page     int    `json:"page,optional"`
	Size     int    `json:"size,optional"`
}

type UserFileListReply struct {
	List  []*UserFile `json:"list"`
	Count int64       `json:"count"`
}

type UserFile struct {
	Id                 int64  `json:"id"`
	Identity           string `json:"identity"`
	RepositoryIdentity string `json:"repository_identity"`
	Name               string `json:"name"`
	Ext                string `json:"ext"`
	Path               string `json:"path"`
	Size               int64  `json:"size"`
	CreatedAt          string `json:"created_at"`
}

type UserFileSearchRequest struct {
	Keyword  string `json:"keyword"`
	FileType string `json:"file_type,optional"` // File extension filter, e.g., ".pdf", ".jpg"
	Page     int    `json:"page,optional"`
	Size     int    `json:"size,optional"`
}

type UserFileSearchReply struct {
	List  []*UserFileSearchItem `json:"list"`
	Count int64                 `json:"count"`
}

type UserFileSearchItem struct {
	ID                 int64  `json:"id"`
	Identity           string `json:"identity"`
	RepositoryIdentity string `json:"repository_identity"`
	Ext                string `json:"ext"`
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	CreatedAt          string `json:"created_at"`
	Path               string `json:"path"`            // Download URL
	ParentPath         string `json:"parent_path"`     // Full folder path
	ParentId           int64  `json:"parent_id"`       // Parent folder ID
	ParentIdentity     string `json:"parent_identity"` // Parent folder identity
}

type UserFolderListRequest struct {
	Identity string `json:"identity,optional"`
}

type UserFolderListReply struct {
	List []*UserFolder `json:"list"`
}

type UserFolder struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
}

type UserRepositorySaveRequest struct {
	ParentId           int64  `json:"parentId"`
	RepositoryIdentity string `json:"repositoryIdentity"`
	Ext                string `json:"ext"`
	Name               string `json:"name"`
}

type UserRepositorySaveReply struct {
}

type FileUploadRequest struct {
	Hash string `json:"hash,optional"`
	Name string `json:"name,optional"`
	Ext  string `json:"ext,optional"`
	Size int64  `json:"size,optional"`
	Path string `json:"path,optional"`
}

type FileUploadReply struct {
	Identity string `json:"identity"`
	Ext      string `json:"ext"`
	Name     string `json:"name"`
}

type UserRegisterRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Code     string `json:"code"`
}

type UserRegisterReply struct {
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginReply struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type UserDetailRequest struct {
	Identity string `json:"identity"`
}

type UserDetailReply struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	NowVolume   int64  `json:"now_volume"`
	TotalVolume int64  `json:"total_volume"`
}

type MailCodeSendRequest struct {
	Email string `json:"email"`
}

type MailCodeSendReply struct {
}

type MailCodeSendPasswordUpdateRequest struct {
	Email string `json:"email"`
}

type MailCodeSendPasswordUpdateReply struct {
}

type MailCodeSendPasswordResetRequest struct {
	Email string `json:"email"`
}

type MailCodeSendPasswordResetReply struct {
}

type UserPasswordResetRequest struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

type UserPasswordResetReply struct {
}

type UserLogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserLogoutReply struct {
}

// Friend request types
type FriendRequestSendRequest struct {
	ToUserIdentity string `json:"to_user_identity"` // Email or user identity
	Message        string `json:"message,optional"`
}

type FriendRequestSendReply struct {
	Identity string `json:"identity"`
}

type FriendRequestListRequest struct {
	Type string `json:"type,optional"` // sent, received, all
}

type FriendRequestListReply struct {
	List []*FriendRequestItem `json:"list"`
}

type FriendRequestItem struct {
	Identity         string `json:"identity"`
	FromUserIdentity string `json:"from_user_identity"`
	ToUserIdentity   string `json:"to_user_identity"`
	FromUserName     string `json:"from_user_name"`
	ToUserName       string `json:"to_user_name"`
	Status           string `json:"status"`
	Message          string `json:"message"`
	CreatedAt        string `json:"created_at"`
}

type FriendRequestRespondRequest struct {
	Identity string `json:"identity"`
	Action   string `json:"action"` // accept, reject
}

type FriendRequestRespondReply struct {
}

// Friend list types
type FriendListRequest struct {
}

type FriendListReply struct {
	List []*FriendItem `json:"list"`
}

type FriendItem struct {
	Identity     string `json:"identity"`
	UserIdentity string `json:"user_identity"`
	UserName     string `json:"user_name"`
	UserEmail    string `json:"user_email"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}

// Friend share types
type FriendShareCreateRequest struct {
	ToUserIdentity         string `json:"to_user_identity"`         // Friend's user identity
	UserRepositoryIdentity string `json:"user_repository_identity"` // File to share
	Message                string `json:"message,optional"`
}

type FriendShareCreateReply struct {
	Identity string `json:"identity"`
}

type FriendShareListRequest struct {
	Type string `json:"type,optional"` // sent, received, all
}

type FriendShareListReply struct {
	List []*FriendShareItem `json:"list"`
}

type FriendShareItem struct {
	Identity               string `json:"identity"`
	FromUserIdentity       string `json:"from_user_identity"`
	FromUserName           string `json:"from_user_name"`
	ToUserIdentity         string `json:"to_user_identity"`
	ToUserName             string `json:"to_user_name"`
	RepositoryIdentity     string `json:"repository_identity"`
	UserRepositoryIdentity string `json:"user_repository_identity"`
	FileName               string `json:"file_name"`
	FileExt                string `json:"file_ext"`
	FileSize               int64  `json:"file_size"`
	Path                   string `json:"path"` // Download URL
	Message                string `json:"message"`
	IsRead                 bool   `json:"is_read"`
	CreatedAt              string `json:"created_at"`
}

type FriendShareMarkReadRequest struct {
	Identity string `json:"identity"`
}

type FriendShareMarkReadReply struct {
}

type FriendShareDownloadRequest struct {
	ShareIdentity string `json:"share_identity,optional"`
}

type FriendShareDownloadReply struct {
}

type FriendShareSaveRequest struct {
	ShareIdentity string `json:"share_identity"`
	ParentId      int64  `json:"parent_id"`
}

type FriendShareSaveReply struct {
	Identity string `json:"identity"`
}

// Storage Purchase Types
type StoragePurchaseCreateRequest struct {
	StorageAmount int64  `json:"storage_amount"`    // Storage capacity in bytes (e.g., 10737418240 for 10GB)
	Currency      string `json:"currency,optional"` // Currency code, default: usd
}

type StoragePurchaseCreateReply struct {
	SessionID string `json:"session_id"` // Stripe Checkout Session ID
	URL       string `json:"url"`        // Stripe Checkout URL
}

type StorageOrderListRequest struct {
	Status string `json:"status,optional"` // Filter by status: pending, paid, failed, refunded, or empty for all
}

type StorageOrderListReply struct {
	List []*StorageOrderItem `json:"list"`
}

type StorageOrderItem struct {
	Identity      string `json:"identity"`
	StorageAmount int64  `json:"storage_amount"` // Bytes
	PriceAmount   int64  `json:"price_amount"`   // Cents
	Currency      string `json:"currency"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// Sync order status from Stripe
type StoragePurchaseSyncRequest struct {
	SessionID string `json:"session_id"` // Stripe Checkout Session ID
}

type StoragePurchaseSyncReply struct {
	Status        string `json:"status"`         // Order status: paid, pending, failed
	StorageAmount int64  `json:"storage_amount"` // Storage capacity added
	Message       string `json:"message"`
}
