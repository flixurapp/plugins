{
  buildGoModule,
  defaultAttrs,
  lib,
}:
buildGoModule (
  lib.recursiveUpdate {
    version = "2026.11.03";
    src = ./.;
    vendorHash = "sha256-Fc2bYpaeYQ1r0MSDAJ8m7aCxQKVciTLzG0OxxgTNSJ0=";

    meta = {
      author = "xela.codes";
      description = "Integration with MusicBrainz.org";
      homepage = "https://musicbrainz.org";
    };
  } defaultAttrs
)
