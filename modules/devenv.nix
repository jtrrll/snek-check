{
  inputs,
  self,
  ...
}: {
  imports = [
    inputs.devenv.flakeModule
  ];
  perSystem = {
    config,
    lib,
    pkgs,
    system,
    ...
  }: {
    devenv = {
      modules = [
        inputs.env-help.devenvModule
      ];
      shells.default = let
        buildInputs = config.packages.snekcheck.nativeBuildInputs;
        goPkg = lib.findFirst (pkg: builtins.match "go" pkg.pname != null) pkgs.go buildInputs;
      in {
        enterShell = ''
          printf "\n" | ${pkgs.lolcat}/bin/lolcat
          printf "\033[0;1;36mDEVSHELL ACTIVATED\033[0m\n"
        '';

        env-help.enable = true;

        languages = {
          go = {
            enable = true;
            package = goPkg;
          };
          nix.enable = true;
        };

        packages =
          [
            inputs.gomod2nix.legacyPackages.${system}.gomod2nix
            pkgs.commitizen
            pkgs.golangci-lint
          ]
          ++ buildInputs;

        pre-commit = {
          default_stages = ["pre-push"];
          hooks = {
            actionlint.enable = true;
            alejandra.enable = true;
            check-added-large-files = {
              enable = true;
              stages = ["pre-commit"];
            };
            check-yaml.enable = true;
            commitizen.enable = true;
            deadnix.enable = true;
            detect-private-keys = {
              enable = true;
              stages = ["pre-commit"];
            };
            end-of-file-fixer.enable = true;
            flake-checker.enable = true;
            gofmt.enable = true;
            golangci-lint.enable = true;
            govet.enable = true;
            markdownlint.enable = true;
            mixed-line-endings.enable = true;
            nil.enable = true;
            no-commit-to-branch = {
              enable = true;
              stages = ["pre-commit"];
            };
            ripsecrets = {
              enable = true;
              stages = ["pre-commit"];
            };
            shellcheck.enable = true;
            shfmt.enable = true;
            snekcheck = {
              enable = true;
              entry = "${self.packages.${system}.snekcheck}/bin/snekcheck";
              name = "snekcheck";
            };
            statix.enable = true;
          };
        };

        scripts = {
          bench = {
            description = "Runs all benchmark tests.";
            exec = ''
              ${pkgs.gum}/bin/gum spin --show-output --spinner line --title "go test -bench" -- \
                go test ./... -bench=.
            '';
          };
          build = {
            description = "Builds the project binary.";
            exec = ''
              ${pkgs.gum}/bin/gum spin --show-error --spinner line --title "gomod2nix" -- \
                gomod2nix
              nix build .#snekcheck
            '';
          };
          demo = {
            description = "Generates a demo GIF.";
            exec = ''
              ${pkgs.uutils-coreutils-noprefix}/bin/printf "TODO\n"
            '';
          };
          lint = {
            description = "Lints the project.";
            exec = ''
              ${pkgs.gum}/bin/gum spin --show-error --spinner line --title "alejandra ." -- \
                alejandra .
              ${pkgs.gum}/bin/gum spin --show-error --spinner line --title "go mod tidy" -- \
                go mod tidy
              ${pkgs.gum}/bin/gum spin --show-error --spinner line --title "go fmt" -- \
                go fmt ./...
              ${pkgs.gum}/bin/gum spin --show-error --spinner line --title "go vet" -- \
                go vet ./...
              ${pkgs.gum}/bin/gum spin --show-error --spinner line --title "golangci-lint" -- \
                golangci-lint run ./...
            '';
          };
          run = {
            description = "Runs the project.";
            exec = ''
              ${pkgs.gum}/bin/gum spin --show-error --spinner line --title "gomod2nix" -- \
                gomod2nix
              nix run .#snekcheck -- "$@"
            '';
          };
          unit = {
            description = "Runs all unit tests.";
            exec = ''
              ${pkgs.gum}/bin/gum spin --show-output --spinner line --title "go test" -- \
                go test ./...
            '';
          };
        };
      };
    };
  };
}
