{
  description = "Official Flixur plugins repository";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    {
      # default overlay exposing all packages under a flixurPlugins set
      overlays.default = final: prev: {
        flixurPlugins = import ./default.nix {
          inherit (final) lib;
          # pass final callPackage function instead of overridable one
          callPackage = final.lib.callPackageWith final;
        };
      };
    }
    # expose each plugin as a package too
    // flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ self.overlays.default ];
        };
      in
      {
        packages = pkgs.flixurPlugins;
      }
    );
}
