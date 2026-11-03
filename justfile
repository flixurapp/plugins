# Upgrade a plugin's flixur-repo deps to master.
upgrade dir:
    cd {{dir}} && env GOPROXY=direct go get -u \
        forge.xela.codes/xela/flixur/pluginkit@master \
        forge.xela.codes/xela/flixur/proto/go@master
