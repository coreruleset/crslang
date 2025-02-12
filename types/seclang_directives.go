package types

type SeclangDirective interface {
	ToSeclang() string
}

type ChainableDirective interface {
	SeclangDirective
	ToSeclangWithIdent(string) string
	GetChainedDirective() ChainableDirective
	AppendChainedDirective(ChainableDirective)
	NonDisruptiveActionsCount() int
}

type ConfigurationDirective struct {
	Kind      string           `yaml:"kind"`
	Metadata  *CommentMetadata `yaml:",inline"`
	Name      string           `yaml:"name"`
	Parameter string           `yaml:"parameter"`
}

func NewConfigurationDirective() *ConfigurationDirective {
	c := new(ConfigurationDirective)
	c.Kind = "configuration"
	c.Metadata = new(CommentMetadata)
	return c
}

func (c ConfigurationDirective) GetMetadata() Metadata {
	return c.Metadata
}

func (c ConfigurationDirective) ToSeclang() string {
	return c.Metadata.Comment + c.Name + " " + c.Parameter + "\n"
}

type CommentDirective struct {
	Kind     string          `yaml:"kind"`
	Metadata CommentMetadata `yaml:",inline"`
}

func (d CommentDirective) ToSeclang() string {
	return d.Metadata.ToSeclang()
}
