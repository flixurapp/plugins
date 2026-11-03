{
  projectRootFile = "flake.nix";

  settings.excludes = [
    # go build artifacts / lockfiles
    "vendor/"
    "*.lock"
    "*.sum"
  ];

  programs = {
    nixfmt.enable = true;
    gofmt.enable = true;
    prettier.enable = true;
    just.enable = true;
  };
}
