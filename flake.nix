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
        pnpm = pkgs.pnpm_10;
        go_1_26_4 = pkgs.go_1_26.overrideAttrs (finalAttrs: _oldAttrs: {
          version = "1.26.4";
          src = pkgs.fetchurl {
            url = "https://go.dev/dl/go${finalAttrs.version}.src.tar.gz";
            hash = "sha256-T2aKMvv8ETLmqIH7lowvHa2mMUkqM5IRc1+7JVpCYC0=";
          };
        });
        buildGo1264Module = pkgs.buildGoModule.override {
          go = go_1_26_4;
        };
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
          doCheck = false;
          ldflags = [ "-X main.noVCSVersionOverride=v${version}" ];
        };
        projectTools = [
          go_1_26_4
          pkgs.golangci-lint
          pkgs.govulncheck
          pkgs.sqlc
          oapi-codegen
          node
          pnpm
          pkgs.gnumake
        ];
        backend = buildGo1264Module {
          pname = "job-hunting-saas-backend";
          version = "0.1.0";
          src = ./backend;
          vendorHash = "sha256-oPyyxYFvab6QjMRu+JHhNCwhUCESff4nXjWStWfUrwk=";
          subPackages = [ "cmd/server" ];
        };
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
        packages = {
          backend = backend;
          default = backend;
        };

        apps = {
          test = mkApp "job-hunting-saas-test" "test";
          build = mkApp "job-hunting-saas-build" "build";
          gen = mkApp "job-hunting-saas-gen" "gen";
        };

        devShells.default = pkgs.mkShell {
          packages = projectTools ++ [
            # Keep runtime/tool versions aligned with project declarations:
            # Go follows backend/go.mod; pnpm stays on v10 and resolves the
            # packageManager-pinned version inside frontend.
            pkgs.gopls
            pkgs.gotools

            pkgs.supabase-cli

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
            echo "pnpm: $(pnpm --version)"
          '';
        };
      }
    );
}
