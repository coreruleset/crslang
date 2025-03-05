package types

type Kind string

const (
	CommentKind       Kind = "comment"
	ConfigurationKind Kind = "configuration"
	DefaultActionKind Kind = "defaultaction"
	RuleKind          Kind = "rule"
	Remove            Kind = "remove"
	UpdateTarget      Kind = "update_target"
	UpdateAction      Kind = "update_action"
)
