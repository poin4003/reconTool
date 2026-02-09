package dto

type AsimBaseResponse[T any] struct {
	Data         *T     `json:"data"`
	IsSucceeded  bool   `json:"isSucceeded"`
	ErrorMessage string `json:"errorMessage"`
}

type AsimAutoExtend string

const (
	AutoExtendNo  AsimAutoExtend = "0"
	AutoExtendYes AsimAutoExtend = "1"
)

type GetCurrentPackageRequest struct {
	Isdn string `json:"isdn"`
}

type PackageComponent struct {
	Unit      string `json:"unit"`
	Total     string `json:"total"`
	Available string `json:"available"`
}

type PackageVoice struct {
	External *PackageComponent `json:"external"`
	Internal *PackageComponent `json:"internal"`
}

type PackageForList struct {
	GoodName   string            `json:"goodName"`
	AutoExtend AsimAutoExtend    `json:"autoExtend"`
	StartDate  string            `json:"startDate"`
	EndDate    string            `json:"endDate"`
	Data       *PackageComponent `json:"data"`
	Voice      *PackageVoice     `json:"voice"`
}

type WaitPackage struct {
	CreateDateTime    string `json:"createDateTime"`
	CreateUser        string `json:"createUser"`
	ModifyDateTime    string `json:"modifyDateTime"`
	ModifyUser        string `json:"modifyUser"`
	GoodGroupName     string `json:"goodGroupName"`
	GoodNameDisplay   string `json:"goodNameDisplay"`
	GoodPrice         int    `json:"goodPrice"`
	TypicalPrice      int    `json:"typicalPrice"`
	Discount          int    `json:"discount"`
	GoodPriceDiscount int    `json:"goodPriceDiscount"`
	WaitPck           string `json:"waitPck"`
	WaitTime          string `json:"waitTime"`
	OrgPck            string `json:"orgPck"`
	ExpiredDateWait   string `json:"expiredDateWait"`
}

type GetCurrentPackageResponse struct {
	Amount        string           `json:"amount"`
	ZoneCode      string           `json:"zoneCode"`
	MainPackages  []PackageForList `json:"mainPackages"`
	AddonPackages []PackageForList `json:"addonPackages"`
	IpPackages    []PackageForList `json:"ipPackages"`
	OtherPackages []PackageForList `json:"otherPackages"`
	WaitPackages  *WaitPackage     `json:"waitPackages"`
}
