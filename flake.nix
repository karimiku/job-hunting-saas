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
      in
      {
        devShells.default = pkgs.mkShell {
          packages = [
            # Keep runtime/tool versions aligned with project declarations:
            # Go follows backend/go.mod, and pnpm follows packageManager via Corepack.
            pkgs.go_1_25
            pkgs.gopls
            pkgs.gotools
            pkgs.sqlc
            pkgs.oapi-codegen

            node
            pnpm

            pkgs.docker-client
            pkgs.docker-compose
            pkgs.postgresql_16
            pkgs.gnumake
            pkgs.curl
            pkgs.jq
            pkgs.lsof
          ];

          shellHook = ''
            echo "job-hunting-saas dev shell"
            echo "Go:   $(go version)"
            echo "Node: $(node --version)"
            echo "pnpm: managed by Corepack from package.json"
          '';
        };
      }
    );
}
