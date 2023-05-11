package windows

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	//"github.com/gonutz/w32"
)

type USER_INFO_0 struct {
	Usri1_name *uint16
}

type LOCALGROUP_USERS_INFO_0 struct {
	lgrui0_name *uint16
}

type _USER_INFO_2 struct {
	usri2_name           *uint16
	usri2_password       *uint16
	usri2_password_age   uint32
	usri2_priv           uint32
	usri2_home_dir       *uint16
	usri2_comment        *uint16
	usri2_flags          uint32
	usri2_script_path    *uint16
	usri2_auth_flags     uint32
	usri2_full_name      *uint16
	usri2_usr_comment    *uint16
	usri2_parms          *uint16
	usri2_workstations   *uint16
	usri2_last_logon     uint32
	usri2_last_logoff    uint32
	usri2_acct_expires   uint32
	usri2_max_storage    uint32
	usri2_units_per_week uint32
	usri2_logon_hours    *uint16
	usri2_bad_pw_count   uint32
	usri2_num_logons     uint32
	usri2_logon_server   *uint16
	usri2_country_code   uint32
	usri2_code_page      uint32
}

type Userinfo struct {
	Groupware     string
	LastLoginTime uint32
	Rootlet       uint32
	Usertype      string
	Name          string
}

func GetUserInfo() {
	userinfo := &Userinfo{}

	//调用windows的netapi32库
	netapi32 := syscall.NewLazyDLL("netapi32.dll")

	//调用库中的函数
	NetUserEnum := netapi32.NewProc("NetUserEnum")
	NetUserGetInfo := netapi32.NewProc("NetUserGetInfo")
	NetUserGetLocalGroups := netapi32.NewProc("NetUserGetLocalGroups")
	NetApiBufferFree := netapi32.NewProc("NetApiBufferFree")

	var serverName [128]byte
	var puserdata uintptr
	var dwEntriesRead, dwTotalEntries uint32
	//调用windows api,获取用户
	bret, _, _ := NetUserEnum.Call(uintptr(unsafe.Pointer(&serverName)), uintptr(0), uintptr(0x2), uintptr(unsafe.Pointer(&puserdata)),
		uintptr(128), uintptr(unsafe.Pointer(&dwEntriesRead)), uintptr(unsafe.Pointer(&dwTotalEntries)), uintptr(0))

	if int(bret) != 0 {
		return
	}

	var iter = puserdata
	//循环获取用户相关信息
	for i := uint32(0); i < dwEntriesRead; i++ {

		//var userinfo USERINFO
		var data = (*USER_INFO_0)(unsafe.Pointer(iter))

		//获取用户组信息
		var pgroupinfo uintptr
		var group_entriesread, group_totalentries uint32
		bret, _, _ = NetUserGetLocalGroups.Call(uintptr(0), uintptr(unsafe.Pointer(data.Usri1_name)), uintptr(0), uintptr(0x1),
			uintptr(unsafe.Pointer(&pgroupinfo)), uintptr(0xFFFFFFFF), uintptr(unsafe.Pointer(&group_entriesread)),
			uintptr(unsafe.Pointer(&group_totalentries)))

		//获取用户组信息失败
		if int(bret) != 0 {
			iter = uintptr(iter + unsafe.Sizeof(USER_INFO_0{}))

			continue
		}

		var ppgroupinfo_itr = pgroupinfo
		for j := uint32(0); j < group_entriesread; j++ {
			groupinfo := (*LOCALGROUP_USERS_INFO_0)(unsafe.Pointer(ppgroupinfo_itr))
			//用户组
			userinfo.Groupware = windows.UTF16PtrToString(groupinfo.lgrui0_name)
			ppgroupinfo_itr = uintptr(ppgroupinfo_itr + unsafe.Sizeof(LOCALGROUP_USERS_INFO_0{}))
		}
		//释放资源
		NetApiBufferFree.Call(uintptr(unsafe.Pointer(pgroupinfo)))

		//获取用户相关信息
		var puserinfo uintptr
		bret, _, _ = NetUserGetInfo.Call(uintptr(0), uintptr(unsafe.Pointer(data.Usri1_name)), uintptr(2), uintptr(unsafe.Pointer(&puserinfo)))

		//获取用户信息失败
		if int(bret) != 0 {
			iter = uintptr(iter + unsafe.Sizeof(USER_INFO_0{}))

			continue
		}

		var userdata = (*_USER_INFO_2)(unsafe.Pointer(puserinfo))

		//最后登陆时间
		userinfo.LastLoginTime = userdata.usri2_last_logon
		//用户权限 0-来宾 1-普通用户 2-管理员.
		userinfo.Rootlet = userdata.usri2_priv
		switch userinfo.Rootlet {
		case 0:
			userinfo.Usertype = "GUEST"
		case 1:
			userinfo.Usertype = "USER"
		case 2:
			userinfo.Usertype = "ADMIN"
		default:
			userinfo.Usertype = "USER"
		}

		//用户名
		if userdata.usri2_name != nil {
			userinfo.Name = windows.UTF16PtrToString(userdata.usri2_name)
		}

		//释放资源
		NetApiBufferFree.Call(uintptr(unsafe.Pointer(puserinfo)))

		iter = uintptr(iter + unsafe.Sizeof(USER_INFO_0{}))

	}

	//释放资源
	NetApiBufferFree.Call(uintptr(unsafe.Pointer(puserdata)))

	fmt.Print(userinfo)

	return
}

func Get() {
	const path = `C:\\some file`

	//size := w32.GetFileVersionInfoSize(path)
	//if size <= 0 {
	//	panic("GetFileVersionInfoSize failed")
	//}
	//
	//info := make([]byte, size)
	//ok := w32.GetFileVersionInfo(path, info)
	//if !ok {
	//	panic("GetFileVersionInfo failed")
	//}

	//fixed, ok := w32.VerQueryValueRoot(info)
	//if !ok {
	//	panic("VerQueryValueRoot failed")
	//}
	//version := fixed.FileVersion()
	//fmt.Printf(
	//	"file version: %d.%d.%d.%d\
	//",
	//version&0xFFFF000000000000>>48,
	//	version&0x0000FFFF00000000>>32,
	//	version&0x00000000FFFF0000>>16,
	//	version&0x000000000000FFFF>>0,
	//)

	//translations, ok := w32.VerQueryValueTranslations(info)
	//if !ok {
	//	panic("VerQueryValueTranslations failed")
	//}
	//if len(translations) == 0 {
	//	panic("no translation found")
	//}
	//fmt.Println("translations:", translations)
	//
	//t := translations[0]
	//w32.CompanyName simply translates to"CompanyName"
	//company, ok := w32.VerQueryValueString(info, t, w32.CompanyName)
	//if !ok {
	//	panic("cannot get company name")
	//}
	//fmt.Println("company:", company)
}
