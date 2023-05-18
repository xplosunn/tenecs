{ pkgs ? import (fetchTarball
  "https://github.com/NixOS/nixpkgs/archive/a14e2b4b1631efa899f0487231b0706937b02ad7.tar.gz")
  { } }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go_1_19
  ];
}
