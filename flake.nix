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
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go_1_25
            gopls
            gotools
            sqlc
            oapi-codegen

            nodejs_24
            pnpm

            docker-client
            docker-compose
            postgresql_16
            gnumake
            curl
            jq
            lsof
          ];

          shellHook = ''
            export GOTOOLCHAIN=local
            echo "job-hunting-saas dev shell"
            echo "Go:   $(go version)"
            echo "Node: $(node --version)"
            echo "pnpm: $(pnpm --version)"
          '';
        };
      }
    );
}
