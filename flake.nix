{
  description = "An opinionated filename linter that loves snake case.";

  inputs = {
    devenv = {
      inputs.nixpkgs.follows = "nixpkgs";
      url = "github:cachix/devenv";
    };
    env-help = {
      inputs.nixpkgs.follows = "nixpkgs";
      url = "github:jtrrll/env-help";
    };
    flake-parts.url = "github:hercules-ci/flake-parts";
    gomod2nix = {
      inputs.nixpkgs.follows = "nixpkgs";
      # Uses a fork of "github:nix-community/gomod2nix" that generates vendor/modules.txt,
      # and therefore supports Go 1.23.
      url = "github:obreitwi/gomod2nix/fix/go_mod_vendor";
    };
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = {flake-parts, ...} @ inputs:
    flake-parts.lib.mkFlake {inherit inputs;} {
      imports = [./modules];
    };
}
