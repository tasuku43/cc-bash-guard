package app

import (
	"github.com/tasuku43/cc-bash-guard/internal/adapter/claude"
	"github.com/tasuku43/cc-bash-guard/internal/app/doctoring"
	configrepo "github.com/tasuku43/cc-bash-guard/internal/infra/config"
)

func RunDoctor(env Env) DoctorResult {
	inputs := configrepo.ResolveEffectiveInputs(env.Cwd, env.Home, env.XDGConfigHome, claude.Tool)
	loaded := configrepo.LoadEffectiveForTool(env.Cwd, env.Home, env.XDGConfigHome, claude.Tool)
	report := doctoring.Run(loaded, claude.Tool, env.Cwd, env.Home)
	report.Tool = claude.Tool
	report.ConfigSources = inputs.ConfigFiles
	report.SettingsPaths = inputs.SettingsPaths
	report.EffectiveFingerprint = inputs.Fingerprint
	report = doctoring.AddVerifiedArtifactCheck(report, configrepo.VerifiedEffectiveArtifactStatus(env.Cwd, env.Home, env.XDGConfigHome, env.XDGCacheHome, claude.Tool))
	return DoctorResult{
		Report: report,
	}
}
