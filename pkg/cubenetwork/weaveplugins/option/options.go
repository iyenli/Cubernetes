package option

const (
	Launch = "launch"
	Sudo   = "sudo"
)

const (
	DnsAdd        = "dns-add"
	DnsRemove     = "dns-remove"
	DefaultSuffix = ".weave.local"
)

const (
	WeaveName = "weave"
	Attach    = "attach"
	Detach    = "detach"
	Stop      = "stop"
)

const (
	InitScriptFileDir = "/etc/cubernetes/cubenetwork/"
	InitScriptFile    = InitScriptFileDir + "init.sh"
	InitScript        = `#!/bin/bash

eval $(weave env)
`
)
