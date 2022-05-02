// Code generated by sqlc. DO NOT EDIT.

package ntpdb

import (
	"database/sql"
	"fmt"
	"time"
)

type AccountInvitesStatus string

const (
	AccountInvitesStatusPending  AccountInvitesStatus = "pending"
	AccountInvitesStatusAccepted AccountInvitesStatus = "accepted"
	AccountInvitesStatusExpired  AccountInvitesStatus = "expired"
)

func (e *AccountInvitesStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AccountInvitesStatus(s)
	case string:
		*e = AccountInvitesStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for AccountInvitesStatus: %T", src)
	}
	return nil
}

type AccountSubscriptionsStatus string

const (
	AccountSubscriptionsStatusIncomplete        AccountSubscriptionsStatus = "incomplete"
	AccountSubscriptionsStatusIncompleteExpired AccountSubscriptionsStatus = "incomplete_expired"
	AccountSubscriptionsStatusTrialing          AccountSubscriptionsStatus = "trialing"
	AccountSubscriptionsStatusActive            AccountSubscriptionsStatus = "active"
	AccountSubscriptionsStatusPastDue           AccountSubscriptionsStatus = "past_due"
	AccountSubscriptionsStatusCanceled          AccountSubscriptionsStatus = "canceled"
	AccountSubscriptionsStatusUnpaid            AccountSubscriptionsStatus = "unpaid"
	AccountSubscriptionsStatusEnded             AccountSubscriptionsStatus = "ended"
)

func (e *AccountSubscriptionsStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AccountSubscriptionsStatus(s)
	case string:
		*e = AccountSubscriptionsStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for AccountSubscriptionsStatus: %T", src)
	}
	return nil
}

type MonitorsIpVersion string

const (
	MonitorsIpVersionV4 MonitorsIpVersion = "v4"
	MonitorsIpVersionV6 MonitorsIpVersion = "v6"
)

func (e *MonitorsIpVersion) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MonitorsIpVersion(s)
	case string:
		*e = MonitorsIpVersion(s)
	default:
		return fmt.Errorf("unsupported scan type for MonitorsIpVersion: %T", src)
	}
	return nil
}

type MonitorsStatus string

const (
	MonitorsStatusPending MonitorsStatus = "pending"
	MonitorsStatusTesting MonitorsStatus = "testing"
	MonitorsStatusActive  MonitorsStatus = "active"
	MonitorsStatusPaused  MonitorsStatus = "paused"
	MonitorsStatusDeleted MonitorsStatus = "deleted"
)

func (e *MonitorsStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MonitorsStatus(s)
	case string:
		*e = MonitorsStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for MonitorsStatus: %T", src)
	}
	return nil
}

type ServerScoresStatus string

const (
	ServerScoresStatusInactive ServerScoresStatus = "inactive"
	ServerScoresStatusTesting  ServerScoresStatus = "testing"
	ServerScoresStatusActive   ServerScoresStatus = "active"
)

func (e *ServerScoresStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ServerScoresStatus(s)
	case string:
		*e = ServerScoresStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for ServerScoresStatus: %T", src)
	}
	return nil
}

type ServersIpVersion string

const (
	ServersIpVersionV4 ServersIpVersion = "v4"
	ServersIpVersionV6 ServersIpVersion = "v6"
)

func (e *ServersIpVersion) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ServersIpVersion(s)
	case string:
		*e = ServersIpVersion(s)
	default:
		return fmt.Errorf("unsupported scan type for ServersIpVersion: %T", src)
	}
	return nil
}

type UserEquipmentApplicationsStatus string

const (
	UserEquipmentApplicationsStatusNew      UserEquipmentApplicationsStatus = "New"
	UserEquipmentApplicationsStatusPending  UserEquipmentApplicationsStatus = "Pending"
	UserEquipmentApplicationsStatusMaybe    UserEquipmentApplicationsStatus = "Maybe"
	UserEquipmentApplicationsStatusNo       UserEquipmentApplicationsStatus = "No"
	UserEquipmentApplicationsStatusApproved UserEquipmentApplicationsStatus = "Approved"
)

func (e *UserEquipmentApplicationsStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserEquipmentApplicationsStatus(s)
	case string:
		*e = UserEquipmentApplicationsStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for UserEquipmentApplicationsStatus: %T", src)
	}
	return nil
}

type VendorZonesClientType string

const (
	VendorZonesClientTypeNtp  VendorZonesClientType = "ntp"
	VendorZonesClientTypeSntp VendorZonesClientType = "sntp"
	VendorZonesClientTypeAll  VendorZonesClientType = "all"
)

func (e *VendorZonesClientType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = VendorZonesClientType(s)
	case string:
		*e = VendorZonesClientType(s)
	default:
		return fmt.Errorf("unsupported scan type for VendorZonesClientType: %T", src)
	}
	return nil
}

type VendorZonesStatus string

const (
	VendorZonesStatusNew      VendorZonesStatus = "New"
	VendorZonesStatusPending  VendorZonesStatus = "Pending"
	VendorZonesStatusApproved VendorZonesStatus = "Approved"
	VendorZonesStatusRejected VendorZonesStatus = "Rejected"
)

func (e *VendorZonesStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = VendorZonesStatus(s)
	case string:
		*e = VendorZonesStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for VendorZonesStatus: %T", src)
	}
	return nil
}

type ZoneServerCountsIpVersion string

const (
	ZoneServerCountsIpVersionV4 ZoneServerCountsIpVersion = "v4"
	ZoneServerCountsIpVersionV6 ZoneServerCountsIpVersion = "v6"
)

func (e *ZoneServerCountsIpVersion) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ZoneServerCountsIpVersion(s)
	case string:
		*e = ZoneServerCountsIpVersion(s)
	default:
		return fmt.Errorf("unsupported scan type for ZoneServerCountsIpVersion: %T", src)
	}
	return nil
}

type Account struct {
	ID               int32          `json:"id"`
	Name             sql.NullString `json:"name"`
	OrganizationName sql.NullString `json:"organization_name"`
	OrganizationUrl  sql.NullString `json:"organization_url"`
	PublicProfile    bool           `json:"public_profile"`
	UrlSlug          sql.NullString `json:"url_slug"`
	CreatedOn        time.Time      `json:"created_on"`
	ModifiedOn       time.Time      `json:"modified_on"`
	StripeCustomerID sql.NullString `json:"stripe_customer_id"`
}

type AccountInvite struct {
	ID         int32                `json:"id"`
	AccountID  int32                `json:"account_id"`
	Email      string               `json:"email"`
	Status     AccountInvitesStatus `json:"status"`
	UserID     sql.NullInt32        `json:"user_id"`
	SentByID   int32                `json:"sent_by_id"`
	Code       string               `json:"code"`
	ExpiresOn  time.Time            `json:"expires_on"`
	CreatedOn  time.Time            `json:"created_on"`
	ModifiedOn time.Time            `json:"modified_on"`
}

type AccountSubscription struct {
	ID                   int32                      `json:"id"`
	AccountID            int32                      `json:"account_id"`
	StripeSubscriptionID sql.NullString             `json:"stripe_subscription_id"`
	Status               AccountSubscriptionsStatus `json:"status"`
	Name                 string                     `json:"name"`
	MaxZones             int32                      `json:"max_zones"`
	MaxDevices           int32                      `json:"max_devices"`
	CreatedOn            time.Time                  `json:"created_on"`
	EndedOn              sql.NullTime               `json:"ended_on"`
	ModifiedOn           time.Time                  `json:"modified_on"`
}

type AccountUser struct {
	AccountID int32 `json:"account_id"`
	UserID    int32 `json:"user_id"`
}

type ApiKey struct {
	ID         int32          `json:"id"`
	ApiKey     sql.NullString `json:"api_key"`
	Grants     sql.NullString `json:"grants"`
	CreatedOn  time.Time      `json:"created_on"`
	ModifiedOn time.Time      `json:"modified_on"`
}

type CombustCache struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Created    time.Time      `json:"created"`
	PurgeKey   sql.NullString `json:"purge_key"`
	Data       string         `json:"data"`
	Metadata   sql.NullString `json:"metadata"`
	Serialized bool           `json:"serialized"`
	Expire     time.Time      `json:"expire"`
}

type CombustSecret struct {
	SecretTs  int32          `json:"secret_ts"`
	ExpiresTs int32          `json:"expires_ts"`
	Type      string         `json:"type"`
	Secret    sql.NullString `json:"secret"`
}

type DnsRoot struct {
	ID              int32  `json:"id"`
	Origin          string `json:"origin"`
	VendorAvailable int32  `json:"vendor_available"`
	GeneralUse      int32  `json:"general_use"`
	NsList          string `json:"ns_list"`
}

type Log struct {
	ID           int32          `json:"id"`
	AccountID    sql.NullInt32  `json:"account_id"`
	ServerID     sql.NullInt32  `json:"server_id"`
	UserID       sql.NullInt32  `json:"user_id"`
	VendorZoneID sql.NullInt32  `json:"vendor_zone_id"`
	Type         sql.NullString `json:"type"`
	Message      sql.NullString `json:"message"`
	Changes      sql.NullString `json:"changes"`
	CreatedOn    time.Time      `json:"created_on"`
}

type LogScore struct {
	ID         int64           `json:"id"`
	MonitorID  sql.NullInt32   `json:"monitor_id"`
	ServerID   int32           `json:"server_id"`
	Ts         time.Time       `json:"ts"`
	Score      float64         `json:"score"`
	Step       float64         `json:"step"`
	Offset     sql.NullFloat64 `json:"offset"`
	Rtt        sql.NullInt32   `json:"rtt"`
	Attributes sql.NullString  `json:"attributes"`
}

type LogScoresArchiveStatus struct {
	ID         int32         `json:"id"`
	Archiver   string        `json:"archiver"`
	LogScoreID sql.NullInt64 `json:"log_score_id"`
	ModifiedOn time.Time     `json:"modified_on"`
}

type LogStatus struct {
	ServerID   int32     `json:"server_id"`
	LastCheck  time.Time `json:"last_check"`
	TsArchived time.Time `json:"ts_archived"`
}

type Monitor struct {
	ID            int32             `json:"id"`
	UserID        sql.NullInt32     `json:"user_id"`
	AccountID     sql.NullInt32     `json:"account_id"`
	Name          string            `json:"name"`
	Location      string            `json:"location"`
	Ip            string            `json:"ip"`
	IpVersion     MonitorsIpVersion `json:"ip_version"`
	TlsName       sql.NullString    `json:"tls_name"`
	ApiKey        sql.NullString    `json:"api_key"`
	Status        MonitorsStatus    `json:"status"`
	Config        string            `json:"config"`
	ClientVersion string            `json:"client_version"`
	LastSeen      sql.NullTime      `json:"last_seen"`
	LastSubmit    sql.NullTime      `json:"last_submit"`
	CreatedOn     time.Time         `json:"created_on"`
}

type SchemaRevision struct {
	Revision   int32  `json:"revision"`
	SchemaName string `json:"schema_name"`
}

type Server struct {
	ID           int32            `json:"id"`
	Ip           string           `json:"ip"`
	IpVersion    ServersIpVersion `json:"ip_version"`
	UserID       int32            `json:"user_id"`
	AccountID    sql.NullInt32    `json:"account_id"`
	Hostname     sql.NullString   `json:"hostname"`
	Stratum      sql.NullInt32    `json:"stratum"`
	InPool       int32            `json:"in_pool"`
	InServerList int32            `json:"in_server_list"`
	Netspeed     int32            `json:"netspeed"`
	CreatedOn    time.Time        `json:"created_on"`
	UpdatedOn    time.Time        `json:"updated_on"`
	ScoreTs      sql.NullTime     `json:"score_ts"`
	ScoreRaw     float64          `json:"score_raw"`
	DeletionOn   sql.NullTime     `json:"deletion_on"`
}

type ServerAlert struct {
	ServerID       int32        `json:"server_id"`
	LastScore      float64      `json:"last_score"`
	FirstEmailTime time.Time    `json:"first_email_time"`
	LastEmailTime  sql.NullTime `json:"last_email_time"`
}

type ServerNote struct {
	ID         int32     `json:"id"`
	ServerID   int32     `json:"server_id"`
	Name       string    `json:"name"`
	Note       string    `json:"note"`
	CreatedOn  time.Time `json:"created_on"`
	ModifiedOn time.Time `json:"modified_on"`
}

type ServerScore struct {
	ID         int64              `json:"id"`
	MonitorID  int32              `json:"monitor_id"`
	ServerID   int32              `json:"server_id"`
	ScoreTs    sql.NullTime       `json:"score_ts"`
	ScoreRaw   float64            `json:"score_raw"`
	Stratum    sql.NullInt32      `json:"stratum"`
	Status     ServerScoresStatus `json:"status"`
	CreatedOn  time.Time          `json:"created_on"`
	ModifiedOn time.Time          `json:"modified_on"`
}

type ServerUrl struct {
	ID       int32  `json:"id"`
	ServerID int32  `json:"server_id"`
	Url      string `json:"url"`
}

type ServerZone struct {
	ServerID int32 `json:"server_id"`
	ZoneID   int32 `json:"zone_id"`
}

type SystemSetting struct {
	ID         int32          `json:"id"`
	Key        sql.NullString `json:"key"`
	Value      sql.NullString `json:"value"`
	CreatedOn  time.Time      `json:"created_on"`
	ModifiedOn time.Time      `json:"modified_on"`
}

type User struct {
	ID            int32          `json:"id"`
	Email         string         `json:"email"`
	Name          sql.NullString `json:"name"`
	Username      sql.NullString `json:"username"`
	PublicProfile bool           `json:"public_profile"`
}

type UserEquipmentApplication struct {
	ID                 int32                           `json:"id"`
	UserID             int32                           `json:"user_id"`
	Application        sql.NullString                  `json:"application"`
	ContactInformation sql.NullString                  `json:"contact_information"`
	Status             UserEquipmentApplicationsStatus `json:"status"`
}

type UserIdentity struct {
	ID        int32          `json:"id"`
	ProfileID string         `json:"profile_id"`
	UserID    int32          `json:"user_id"`
	Provider  string         `json:"provider"`
	Data      sql.NullString `json:"data"`
	Email     sql.NullString `json:"email"`
}

type UserPrivilege struct {
	UserID             int32 `json:"user_id"`
	SeeAllServers      bool  `json:"see_all_servers"`
	SeeAllUserProfiles bool  `json:"see_all_user_profiles"`
	VendorAdmin        int32 `json:"vendor_admin"`
	EquipmentAdmin     int32 `json:"equipment_admin"`
	SupportStaff       int32 `json:"support_staff"`
}

type VendorZone struct {
	ID                 int32                 `json:"id"`
	ZoneName           string                `json:"zone_name"`
	Status             VendorZonesStatus     `json:"status"`
	UserID             sql.NullInt32         `json:"user_id"`
	OrganizationName   sql.NullString        `json:"organization_name"`
	ClientType         VendorZonesClientType `json:"client_type"`
	ContactInformation sql.NullString        `json:"contact_information"`
	RequestInformation sql.NullString        `json:"request_information"`
	DeviceCount        sql.NullInt32         `json:"device_count"`
	RtTicket           sql.NullInt32         `json:"rt_ticket"`
	ApprovedOn         sql.NullTime          `json:"approved_on"`
	CreatedOn          time.Time             `json:"created_on"`
	ModifiedOn         time.Time             `json:"modified_on"`
	DnsRootID          int32                 `json:"dns_root_id"`
	AccountID          sql.NullInt32         `json:"account_id"`
}

type Zone struct {
	ID          int32          `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	ParentID    sql.NullInt32  `json:"parent_id"`
	Dns         bool           `json:"dns"`
}

type ZoneServerCount struct {
	ID              int32                     `json:"id"`
	ZoneID          int32                     `json:"zone_id"`
	IpVersion       ZoneServerCountsIpVersion `json:"ip_version"`
	Date            time.Time                 `json:"date"`
	CountActive     int32                     `json:"count_active"`
	CountRegistered int32                     `json:"count_registered"`
	NetspeedActive  int32                     `json:"netspeed_active"`
}
