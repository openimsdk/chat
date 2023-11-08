package version

import (
	"fmt"
	"runtime"
)

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() Info {
	// These variables typically come from -ldflags settings and in
	// their absence fallback to the settings in ./base.go
	return Info{
		Major:      gitMajor,
		Minor:      gitMinor,
		GitVersion: gitVersion,
		GitCommit:  gitCommit,
		BuildDate:  buildDate,
		GoVersion:  runtime.Version(),
		Compiler:   runtime.Compiler,
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// GetSingleVersion returns single version of sealer
func GetSingleVersion() string {
	return gitVersion
}

type Output struct {
	OpenIMChatVersion     Info               `json:"OpenIMChatVersion,omitempty" yaml:"OpenIMChatVersion,omitempty"`
	OpenIMServerVersion     *OpenIMServerVersion           `json:"OpenIMServerVersion,omitempty" yaml:"OpenIMServerVersion,omitempty"`
}

type OpenIMServerVersion struct {
	ServerVersion string `json:"serverVersion,omitempty" yaml:"serverVersion,omitempty"`
	ClientVersion string `json:"clientVersion,omitempty" yaml:"clientVersion,omitempty"`	//sdk core version
}