# Upgrade a plugin's flixur-repo deps to master.
upgrade dir:
    # force upgrade, bypassing cache
    cd {{dir}} && env GOPROXY=direct GOPRIVATE='forge.xela.codes/*' \
        go get -u forge.xela.codes/xela/flixur/pluginkit@master
    cd {{dir}} && go mod tidy
    # add fake hash to vendorHash so it'll update
    sed -i 's|vendorHash = ".*"|vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="|' {{dir}}/default.nix
    # add in the new vendorHash
    hash=$(nix build .#{{dir}} 2>&1 | grep -oP 'got:\s+sha256-\K[A-Za-z0-9+/=]+' | head -n1) && \
        sed -i "s|vendorHash = .*|vendorHash = \"sha256-$hash\";|" {{dir}}/default.nix
    # then do a test build
    nix build .#{{dir}}
