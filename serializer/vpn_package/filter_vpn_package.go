package vpnpackage

type FilterVPNPackagesSerializer struct {
	VPNPackages []*VPNPackageSerializer `json:"vpn_packages"`
}
