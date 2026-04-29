package command

import "testing"

func TestTerraformParserExtractsSemanticFields(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want TerraformSemantic
	}{
		{name: "validate", raw: "terraform validate", want: TerraformSemantic{Subcommand: "validate"}},
		{name: "plan out", raw: "terraform plan -out=tfplan", want: TerraformSemantic{Subcommand: "plan", Out: "tfplan"}},
		{name: "plan destroy", raw: "terraform plan -destroy", want: TerraformSemantic{Subcommand: "plan", Destroy: true}},
		{name: "apply", raw: "terraform apply", want: TerraformSemantic{Subcommand: "apply"}},
		{name: "apply auto approve", raw: "terraform apply -auto-approve", want: TerraformSemantic{Subcommand: "apply", AutoApprove: true}},
		{name: "apply plan file", raw: "terraform apply tfplan", want: TerraformSemantic{Subcommand: "apply", PlanFile: "tfplan"}},
		{name: "destroy auto approve", raw: "terraform destroy -auto-approve", want: TerraformSemantic{Subcommand: "destroy", Destroy: true, AutoApprove: true}},
		{name: "import", raw: "terraform import aws_s3_bucket.example bucket-name", want: TerraformSemantic{Subcommand: "import"}},
		{name: "state list", raw: "terraform state list", want: TerraformSemantic{Subcommand: "state", StateSubcommand: "list"}},
		{name: "state rm", raw: "terraform state rm aws_instance.bad", want: TerraformSemantic{Subcommand: "state", StateSubcommand: "rm"}},
		{name: "workspace select", raw: "terraform workspace select prod", want: TerraformSemantic{Subcommand: "workspace", WorkspaceSubcommand: "select"}},
		{name: "workspace delete", raw: "terraform workspace delete old", want: TerraformSemantic{Subcommand: "workspace", WorkspaceSubcommand: "delete"}},
		{name: "chdir target", raw: "terraform -chdir=infra plan -target=module.foo", want: TerraformSemantic{Subcommand: "plan", GlobalChdir: "infra", Target: true, Targets: []string{"module.foo"}}},
		{name: "init flags", raw: "terraform init -upgrade -reconfigure", want: TerraformSemantic{Subcommand: "init", Upgrade: true, Reconfigure: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := Parse(tt.raw)
			if len(plan.Commands) != 1 {
				t.Fatalf("commands = %d, want 1", len(plan.Commands))
			}
			cmd := plan.Commands[0]
			if cmd.Parser != "terraform" || cmd.SemanticParser != "terraform" || cmd.Terraform == nil {
				t.Fatalf("parser state = (%q, %q, %v), want terraform semantic", cmd.Parser, cmd.SemanticParser, cmd.Terraform)
			}
			got := cmd.Terraform
			if got.Subcommand != tt.want.Subcommand {
				t.Fatalf("Subcommand = %q, want %q", got.Subcommand, tt.want.Subcommand)
			}
			if got.GlobalChdir != tt.want.GlobalChdir {
				t.Fatalf("GlobalChdir = %q, want %q", got.GlobalChdir, tt.want.GlobalChdir)
			}
			if got.WorkspaceSubcommand != tt.want.WorkspaceSubcommand {
				t.Fatalf("WorkspaceSubcommand = %q, want %q", got.WorkspaceSubcommand, tt.want.WorkspaceSubcommand)
			}
			if got.StateSubcommand != tt.want.StateSubcommand {
				t.Fatalf("StateSubcommand = %q, want %q", got.StateSubcommand, tt.want.StateSubcommand)
			}
			if got.Destroy != tt.want.Destroy || got.AutoApprove != tt.want.AutoApprove ||
				got.Target != tt.want.Target || got.Upgrade != tt.want.Upgrade || got.Reconfigure != tt.want.Reconfigure {
				t.Fatalf("booleans = destroy:%v auto:%v target:%v upgrade:%v reconfigure:%v, want %+v", got.Destroy, got.AutoApprove, got.Target, got.Upgrade, got.Reconfigure, tt.want)
			}
			if got.Out != tt.want.Out || got.PlanFile != tt.want.PlanFile {
				t.Fatalf("values = out:%q plan_file:%q, want %+v", got.Out, got.PlanFile, tt.want)
			}
			if len(tt.want.Targets) > 0 && (len(got.Targets) != len(tt.want.Targets) || got.Targets[0] != tt.want.Targets[0]) {
				t.Fatalf("Targets = %v, want %v", got.Targets, tt.want.Targets)
			}
		})
	}
}
