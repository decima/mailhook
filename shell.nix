# nix shell
{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShellNoCC {
    buildInputs = [
      pkgs.swaks
      pkgs.gnumake
    ];
  }