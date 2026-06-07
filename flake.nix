{
  description = "Development environment for job-hunting-saas";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
    in
    flake-utils.lib.eachSystem supportedSystems (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
        node = pkgs.nodejs_24;
        pnpm = pkgs.writeShellScriptBin "pnpm" ''
          exec ${node}/bin/corepack pnpm "$@"
        '';
        oapi-codegen = pkgs.buildGoModule rec {
          pname = "oapi-codegen";
          version = "2.6.0";
          src = pkgs.fetchFromGitHub {
            owner = "oapi-codegen";
            repo = "oapi-codegen";
            rev = "v${version}";
            hash = "sha256-VUSqwc6TsMhry4BEj9nMkSaKg9PNMYGktwc0CA3yx6c=";
          };
          vendorHash = "sha256-vgSMGi0mnGX/Hwxu/XalIXLCbm/L4CwQfIf7DEJVk1E=";
          subPackages = [ "cmd/oapi-codegen" ];
          ldflags = [ "-X main.noVCSVersionOverride=v${version}" ];
        };
        projectTools = [
          pkgs.go_1_26
          pkgs.sqlc
          oapi-codegen
          node
          pnpm
          pkgs.gnumake
        ];
        mkApp =
          name: target:
          {
            type = "app";
            program = "${
              pkgs.writeShellApplication {
                inherit name;
                runtimeInputs = projectTools;
                text = ''
                  export GOTOOLCHAIN=auto
                  exec make ${target} "$@"
                '';
              }
            }/bin/${name}";
          };
      in
      {
        apps = {
          test = mkApp "job-hunting-saas-test" "test";
          build = mkApp "job-hunting-saas-build" "build";
          gen = mkApp "job-hunting-saas-gen" "gen";
        };

        devShells.default = pkgs.mkShell {
          packages = projectTools ++ [
            # Keep runtime/tool versions aligned with project declarations:
            # Go follows backend/go.mod, and pnpm follows packageManager via Corepack.
            pkgs.gopls
            pkgs.gotools

            pkgs.docker-client
            pkgs.docker-compose
            pkgs.postgresql_16
            pkgs.curl
            pkgs.jq
            pkgs.lsof
          ];

          shellHook = ''
            export GOTOOLCHAIN=auto
            echo "job-hunting-saas dev shell"
            echo "Go:   $(go version)"
            echo "Node: $(node --version)"
            echo "pnpm: managed by Corepack from package.json"
          '';
        };
      }
    );
}
