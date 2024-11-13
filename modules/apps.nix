{
  perSystem = {config, ...}: {
    apps = {
      snekcheck = {
        program = "${config.packages.snekcheck}/bin/snekcheck";
        type = "app";
      };
    };
  };
}
