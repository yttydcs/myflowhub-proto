package auth

import "encoding/json"

const (
	SubProtoAuth uint8 = 2

	ActionRegister                 = "register"
	ActionAssistRegister           = "assist_register"
	ActionRegisterResp             = "register_resp"
	ActionAssistRegisterResp       = "assist_register_resp"
	ActionLogin                    = "login"
	ActionAssistLogin              = "assist_login"
	ActionLoginResp                = "login_resp"
	ActionAssistLoginResp          = "assist_login_resp"
	ActionRevoke                   = "revoke"
	ActionRevokeResp               = "revoke_resp"
	ActionAssistQueryCred          = "assist_query_credential"
	ActionAssistQueryCredResp      = "assist_query_credential_resp"
	ActionOffline                  = "offline"
	ActionAssistOffline            = "assist_offline"
	ActionGetPerms                 = "get_perms"
	ActionGetPermsResp             = "get_perms_resp"
	ActionListRoles                = "list_roles"
	ActionListRolesResp            = "list_roles_resp"
	ActionPermsInvalidate          = "perms_invalidate"
	ActionPermsSnapshot            = "perms_snapshot"
	ActionAuthorityPolicySync      = "authority_policy_sync"
	ActionListPendingRegisters     = "list_pending_registers"
	ActionListPendingRegistersResp = "list_pending_registers_resp"
	ActionApproveRegister          = "approve_register"
	ActionApproveRegisterResp      = "approve_register_resp"
	ActionRejectRegister           = "reject_register"
	ActionRejectRegisterResp       = "reject_register_resp"
	ActionIssueRegisterPermit      = "issue_register_permit"
	ActionIssueRegisterPermitResp  = "issue_register_permit_resp"
	ActionRevokeRegisterPermit     = "revoke_register_permit"
	ActionRevokeRegisterPermitResp = "revoke_register_permit_resp"
	ActionUpLogin                  = "up_login"
	ActionUpLoginResp              = "up_login_resp"
)

type Message struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type RegisterData struct {
	DeviceID      string `json:"device_id"`
	NodeID        uint32 `json:"node_id,omitempty"`
	RequestedRole string `json:"requested_role,omitempty"`
	JoinPermit    string `json:"join_permit,omitempty"`
	PubKey        string `json:"pubkey,omitempty"`
	NodePub       string `json:"node_pub,omitempty"`
	DisplayName   string `json:"display_name,omitempty"`
	TS            int64  `json:"ts,omitempty"`
	Nonce         string `json:"nonce,omitempty"`
}

type LoginData struct {
	DeviceID    string `json:"device_id"`
	NodeID      uint32 `json:"node_id,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	TS          int64  `json:"ts,omitempty"`
	Nonce       string `json:"nonce,omitempty"`
	Sig         string `json:"sig,omitempty"`
	Alg         string `json:"alg,omitempty"`
}

type RevokeData struct {
	DeviceID string `json:"device_id"`
	NodeID   uint32 `json:"node_id,omitempty"`
}

type QueryCredData struct {
	DeviceID string `json:"device_id"`
	NodeID   uint32 `json:"node_id,omitempty"`
}

type OfflineData struct {
	DeviceID string `json:"device_id"`
	NodeID   uint32 `json:"node_id,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

type RespData struct {
	Code        int      `json:"code"`
	Msg         string   `json:"msg,omitempty"`
	DeviceID    string   `json:"device_id,omitempty"`
	NodeID      uint32   `json:"node_id,omitempty"`
	HubID       uint32   `json:"hub_id,omitempty"`
	Role        string   `json:"role,omitempty"`
	Perms       []string `json:"perms,omitempty"`
	PubKey      string   `json:"pubkey,omitempty"`
	NodePub     string   `json:"node_pub,omitempty"`
	DisplayName string   `json:"display_name,omitempty"`
	Status      string   `json:"status,omitempty"`
	RequestID   string   `json:"request_id,omitempty"`
	Reason      string   `json:"reason,omitempty"`
	TS          int64    `json:"ts,omitempty"`
	Nonce       string   `json:"nonce,omitempty"`
}

type PermsQueryData struct {
	NodeID uint32 `json:"node_id"`
}

type InvalidateData struct {
	NodeIDs []uint32 `json:"node_ids,omitempty"`
	Reason  string   `json:"reason,omitempty"`
	Refresh bool     `json:"refresh,omitempty"`
}

type AuthorityPolicySyncData struct {
	Mode                 string `json:"mode,omitempty"`
	EffectiveAuthorityID uint32 `json:"effective_authority_id,omitempty"`
	Epoch                uint64 `json:"epoch,omitempty"`
	TTLSec               uint32 `json:"ttl_sec,omitempty"`
}

type RolePermEntry struct {
	NodeID uint32   `json:"node_id,omitempty"`
	Role   string   `json:"role,omitempty"`
	Perms  []string `json:"perms,omitempty"`
}

type ListRolesReq struct {
	Offset  int      `json:"offset,omitempty"`
	Limit   int      `json:"limit,omitempty"`
	Role    string   `json:"role,omitempty"`
	NodeIDs []uint32 `json:"node_ids,omitempty"`
}

type PendingRegisterInfo struct {
	RequestID     string `json:"request_id,omitempty"`
	DeviceID      string `json:"device_id,omitempty"`
	RequestedRole string `json:"requested_role,omitempty"`
	DisplayName   string `json:"display_name,omitempty"`
	CreatedAt     int64  `json:"created_at,omitempty"`
	ExpiresAt     int64  `json:"expires_at,omitempty"`
}

type ListPendingRegistersReq struct {
	Offset   int    `json:"offset,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	DeviceID string `json:"device_id,omitempty"`
}

type ListPendingRegistersResp struct {
	Code  int                   `json:"code"`
	Msg   string                `json:"msg,omitempty"`
	Total int                   `json:"total"`
	Items []PendingRegisterInfo `json:"items,omitempty"`
}

type ApproveRegisterReq struct {
	RequestID string `json:"request_id"`
	Role      string `json:"role,omitempty"`
}

type ApproveRegisterResp struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	DeviceID  string `json:"device_id,omitempty"`
	NodeID    uint32 `json:"node_id,omitempty"`
	Role      string `json:"role,omitempty"`
	Status    string `json:"status,omitempty"`
}

type RejectRegisterReq struct {
	RequestID string `json:"request_id"`
	Reason    string `json:"reason,omitempty"`
}

type RejectRegisterResp struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	DeviceID  string `json:"device_id,omitempty"`
	Status    string `json:"status,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

type IssueRegisterPermitReq struct {
	DeviceID  string `json:"device_id"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
}

type IssueRegisterPermitResp struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg,omitempty"`
	Permit    string `json:"permit,omitempty"`
	DeviceID  string `json:"device_id,omitempty"`
	Role      string `json:"role,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
}

type RevokeRegisterPermitReq struct {
	Permit string `json:"permit"`
}

type RevokeRegisterPermitResp struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg,omitempty"`
	Permit   string `json:"permit,omitempty"`
	DeviceID string `json:"device_id,omitempty"`
	Role     string `json:"role,omitempty"`
}

type UpLoginData struct {
	NodeID      uint32 `json:"node_id"`
	DeviceID    string `json:"device_id,omitempty"`
	HubID       uint32 `json:"hub_id,omitempty"`
	PubKey      string `json:"pubkey,omitempty"`
	TS          int64  `json:"ts,omitempty"`
	Nonce       string `json:"nonce,omitempty"`
	DeviceTS    int64  `json:"device_ts,omitempty"`
	DeviceNonce string `json:"device_nonce,omitempty"`
	DeviceSig   string `json:"device_sig,omitempty"`
	DeviceAlg   string `json:"device_alg,omitempty"`
	SenderID    uint32 `json:"sender_id,omitempty"`
	SenderTS    int64  `json:"sender_ts,omitempty"`
	SenderNonce string `json:"sender_nonce,omitempty"`
	SenderSig   string `json:"sender_sig,omitempty"`
	SenderAlg   string `json:"sender_alg,omitempty"`
	SenderPub   string `json:"sender_pub,omitempty"`
	Sig         string `json:"sig,omitempty"`
	Alg         string `json:"alg,omitempty"`
}
