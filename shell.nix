{ pkgs ? import (fetchTarball
  "https://github.com/NixOS/nixpkgs/archive/ea1cb3d8d4727a78ddf58d0ac53150142217a1d0.tar.gz")
  { } }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go_1_23
    pkgs.nodejs
  ];
}
