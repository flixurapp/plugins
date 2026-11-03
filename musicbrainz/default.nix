{
  buildGoModule,
  defaultAttrs,
  lib,
}:
buildGoModule (
  lib.recursiveUpdate {
    version = "0.1.0";
    src = ./.;
    vendorHash = "sha256-Fc2bYpaeYQ1r0MSDAJ8m7aCxQKVciTLzG0OxxgTNSJ0=";

    meta = {
      description = "Integration with MusicBrainz.org";
      homepage = "https://musicbrainz.org";
    };
  } defaultAttrs
)
