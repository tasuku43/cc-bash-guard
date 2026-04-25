package app

import "github.com/tasuku43/cc-bash-guard/internal/infra/buildinfo"

func RunVersion() VersionResult {
	return VersionResult{Info: buildinfo.Read()}
}
