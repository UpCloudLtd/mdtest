//nolint:gochecknoglobals // These are defined as var instead of const to be able to override them on compile time.
package globals

var (
	BuildDate = "unknown"
	Version   = "dev"
)
