# Imports every non-dot directory as a plugin.
{
  lib,
  callPackage,
}:
let
  buildableDir =
    name: type:
    type == "directory" && !lib.hasPrefix "." name && builtins.pathExists "${./.}/${name}/default.nix";

  defaultAttrs = pname: {
    inherit pname;
    meta = {
      license = lib.licenses.agpl3Only;
      mainProgram = pname;
      platforms = lib.platforms.unix;
    };
  };
in
lib.mapAttrs (name: _: callPackage ./${name} { defaultAttrs = defaultAttrs name; }) (
  lib.filterAttrs buildableDir (builtins.readDir ./.)
)
