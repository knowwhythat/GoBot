package forms

import (
	"github.com/google/uuid"
	"time"
)

type LoginReq struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

type RespData[T any] struct {
	Code int    `json:"code"`
	Data T      `json:"data"`
	Msg  string `json:"msg"`
}

type LoginResponse struct {
	User      SysUser `json:"user"`
	Token     string  `json:"token"`
	ExpiresAt int64   `json:"expiresAt"`
}

type Register struct {
	Username    string `json:"userName" example:"用户名"`
	Password    string `json:"passWord" example:"密码"`
	NickName    string `json:"nickName" example:"昵称"`
	AuthorityId uint   `json:"authorityId"`
}

type SysUser struct {
	ID          uint           `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt   time.Time      // 创建时间
	UpdatedAt   time.Time      // 更新时间
	UUID        uuid.UUID      `json:"uuid" gorm:"index;comment:用户UUID"`                                                     // 用户UUID
	Username    string         `json:"userName" gorm:"index;comment:用户登录名"`                                                  // 用户登录名
	Password    string         `json:"-"  gorm:"comment:用户登录密码"`                                                             // 用户登录密码
	NickName    string         `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                                            // 用户昵称
	SideMode    string         `json:"sideMode" gorm:"default:dark;comment:用户侧边主题"`                                          // 用户侧边主题
	HeaderImg   string         `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	BaseColor   string         `json:"baseColor" gorm:"default:#fff;comment:基础颜色"`                                           // 基础颜色
	AuthorityId uint           `json:"authorityId" gorm:"default:888;comment:用户角色ID"`                                        // 用户角色ID
	Authority   SysAuthority   `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
	Authorities []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
	Phone       string         `json:"phone"  gorm:"comment:用户手机号"`                     // 用户手机号
	Email       string         `json:"email"  gorm:"comment:用户邮箱"`                      // 用户邮箱
	Enable      int            `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"` //用户是否被冻结 1正常 2冻结
}

type SysAuthority struct {
	CreatedAt       time.Time       // 创建时间
	UpdatedAt       time.Time       // 更新时间
	DeletedAt       *time.Time      `sql:"index"`
	AuthorityId     uint            `json:"authorityId" gorm:"not null;unique;primary_key;comment:角色ID;size:90"` // 角色ID
	AuthorityName   string          `json:"authorityName" gorm:"comment:角色名"`                                    // 角色名
	ParentId        *uint           `json:"parentId" gorm:"comment:父角色ID"`                                       // 父角色ID
	DataAuthorityId []*SysAuthority `json:"dataAuthorityId" gorm:"many2many:sys_data_authority_id;"`
	Children        []SysAuthority  `json:"children" gorm:"-"`
	Users           []SysUser       `json:"-" gorm:"many2many:sys_user_authority;"`
	DefaultRouter   string          `json:"defaultRouter" gorm:"comment:默认菜单;default:dashboard"` // 默认菜单(默认dashboard)
}
