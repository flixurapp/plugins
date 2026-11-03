{
  buildGoModule,
  defaultAttrs,
  lib,
}:
buildGoModule (
  lib.recursiveUpdate {
    version = "2026.11.03";
    src = ./.;
    vendorHash = "sha256-0ljfehTUqrn5rdxGX36beUaq2nMVpPJT3kwjNmwWBf4=";

    meta = {
      author = "xela.codes";
      description = "Integration with MusicBrainz.org";
      homepage = "https://musicbrainz.org";
    };
  } defaultAttrs
)
