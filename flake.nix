{
  description = "Official Flixur plugins repository";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
    flake-utils.url = "github:numtide/flake-utils";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    flixur = {
      url = "git+https://forge.xela.codes/xela/flixur.git";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      flixur,
      treefmt-nix,
    }:
    {
      # default overlay exposing all packages under a flixurPlugins set
      overlays.default = final: prev: {
        flixurPlugins = import ./default.nix final;
      };
    }
    # expose each plugin as a package too
    // flake-utils.lib.eachDefaultSystem (
      system:
      let
        lib = nixpkgs.lib;
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ self.overlays.default ];
        };

        treefmtEval = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;

        pluginkitManifestSchema = "${flixur}/pluginkit/manifest.schema.json";
      in
      {
        packages = pkgs.flixurPlugins;

        # validate plugin manifests against the schema
        checks.manifests =
          let
            plugins = lib.filterAttrs (
              name: type: type == "directory" && builtins.pathExists "${./.}/${name}/manifest.json"
            ) (builtins.readDir ./.);
          in
          pkgs.runCommand "check-manifests"
            {
              nativeBuildInputs = [ pkgs.check-jsonschema ];
            }
            ''
              for manifest in ${
                lib.escapeShellArgs (lib.mapAttrsToList (name: _: "${./.}/${name}/manifest.json") plugins)
              }; do
                check-jsonschema --schemafile ${pluginkitManifestSchema} "$manifest"
              done
              touch $out
            '';

        formatter = treefmtEval.config.build.wrapper;
        checks.format = treefmtEval.config.build.check self;
      }
    );
}
