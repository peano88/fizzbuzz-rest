{ pkgs ? import <nixpkgs> { } }:


pkgs.mkShell {
  buildInputs = with pkgs;[
    go
    gotools
    go-tools
    gopls
    go-outline
    gocode
    gopkgs
    gocode-gomod
    godef
    golint
  ];
}
