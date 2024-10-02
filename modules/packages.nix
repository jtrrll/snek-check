{inputs, ...}: {
  perSystem = {system, ...}: {
    packages = {
      snek-check = inputs.gomod2nix.legacyPackages.${system}.buildGoApplication {
        modules = ../gomod2nix.toml;
        pname = "snek-check";
        src = ../.;
        version = "0.0";
      };
    };
  };
}
