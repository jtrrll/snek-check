{inputs, ...}: {
  perSystem = {system, ...}: {
    packages = {
      snekcheck = inputs.gomod2nix.legacyPackages.${system}.buildGoApplication {
        modules = ../gomod2nix.toml;
        pname = "snekcheck";
        src = ../.;
        version = "0.0";
      };
    };
  };
}
