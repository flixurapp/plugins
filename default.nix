# Imports every non-dot directory as a plugin.
{ lib, ... }@pkgs:
let
  buildableDir =
    name: type:
    type == "directory" && !lib.hasPrefix "." name && builtins.pathExists "${./.}/${name}/default.nix";

  manifestOf = name: builtins.fromJSON (builtins.readFile ./${name}/manifest.json);

  # Package metadata derived from the manifest.
  defaultAttrs = manifest: {
    pname = manifest.id;
    version = manifest.version;
    meta = {
      description = manifest.description;
      author = manifest.author;
      homepage = manifest.url;
      mainProgram = manifest.id;
      # from this repo
      license = lib.licenses.agpl3Only;
      platforms = lib.platforms.unix;
    };

    # install manifest.json in the package too
    postInstall = ''
      test -x "$out/bin/${manifest.id}"
      mkdir -p "$out/share/flixur-plugins/${manifest.id}"
      ${lib.getExe pkgs.jq} --arg bin "$out/bin/${manifest.id}" '.bin = $bin' \
        ${./${manifest.id}/manifest.json} > "$out/share/flixur-plugins/${manifest.id}/manifest.json"
    '';
  };
in
lib.mapAttrs (
  name: _:
  let
    manifest = manifestOf name;
    attrs = (
      assert manifest.id == name;
      defaultAttrs manifest
    );
    extraAttrs = import (./. + "/${name}") { inherit pkgs lib manifest; };
  in
  pkgs.buildGoModule (lib.recursiveUpdate attrs extraAttrs)
) (lib.filterAttrs buildableDir (builtins.readDir ./.))
