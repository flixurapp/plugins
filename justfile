# Upgrade a plugin's flixur-repo deps to master.
upgrade dir:
    #!/usr/bin/env bash
    set -euo pipefail
    # force upgrade, bypassing cache
    (cd {{dir}} && env GOPROXY=direct GOPRIVATE='forge.xela.codes/*' \
        go get -u forge.xela.codes/xela/flixur/pluginkit@master)
    (cd {{dir}} && go mod tidy)
    # add fake hash to vendorHash so it'll update
    sed -i 's|vendorHash = ".*"|vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="|' {{dir}}/default.nix
    # nix build will fail with a hash mismatch, grab the real hash from the error
    build_out=$(nix build .#{{dir}} 2>&1) || true
    hash=$(printf '%s\n' "$build_out" | grep -oP 'got:\s+sha256-\K[A-Za-z0-9+/=]+' | head -n1 || true)
    # if the build failed for another reason, show the failure and stop
    if [ -z "$hash" ]; then
        printf '%s\n' "$build_out"
        exit 1
    fi
    # add in the new vendorHash
    sed -i "s|vendorHash = .*|vendorHash = \"sha256-$hash\";|" {{dir}}/default.nix
    # then do a test build
    nix build .#{{dir}}
