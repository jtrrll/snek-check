{
  perSystem = {config, ...}: {
    apps = {
      snek = {
        program = "${config.packages.snek-check}/bin/snek";
        type = "app";
      };
    };
  };
}
