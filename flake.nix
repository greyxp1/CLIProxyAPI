{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = {
    self,
    nixpkgs,
  }: let
    forAllSystems = nixpkgs.lib.genAttrs ["x86_64-linux" "aarch64-linux"];
  in {
    packages = forAllSystems (system: {
      default = nixpkgs.legacyPackages.${system}.buildGoModule {
        pname = "cliproxyapi";
        version = "pool-${self.shortRev or "dirty"}";
        src = self;
        vendorHash = "sha256-9x/JpCr/qC6Cn5omS4HmPkKKaQMQgd6v4JAn6Dn5ms8=";
        # The pool database uses SQLite through cgo.
        env.CGO_ENABLED = 1;
        subPackages = ["cmd/server" "cmd/pool-client"];
        ldflags = ["-s" "-w" "-X main.Version=pool-${self.shortRev or "dirty"}"];
        postInstall = ''
          mv $out/bin/server $out/bin/cliproxyapi
          mv $out/bin/pool-client $out/bin/cliproxy-client
        '';
        meta = {
          description = "CLIProxyAPI with a shared Claude and Codex subscription pool";
          license = nixpkgs.lib.licenses.mit;
          mainProgram = "cliproxyapi";
        };
      };
    });
  };
}
