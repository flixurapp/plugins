{
  description = "Official Flixur plugins repository";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
    flake-utils.url = "github:numtide/flake-utils";

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
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
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ self.overlays.default ];
        };

        treefmtEval = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;
      in
      {
        packages = pkgs.flixurPlugins;

        formatter = treefmtEval.config.build.wrapper;
        checks.format = treefmtEval.config.build.check self;
      }
    );
}
