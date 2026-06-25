{
  description = "Go backend + React Router frontend monorepo, wired via ConnectRPC";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.go_1_24
            pkgs.bun
            pkgs.buf
            pkgs.just
            pkgs.air
          ];
          shellHook = ''
            echo "Dev shell ready. Run 'just' to see available tasks."
          '';
        };
      }
    );
}
