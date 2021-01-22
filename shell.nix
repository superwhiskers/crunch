{ sources ? import ./nix/sources.nix, pkgs ? import ./nix { inherit sources; } }:

pkgs.mkShell {
  name = "crunch-shell";

  buildInputs = with pkgs; [
    go
  ];
}
