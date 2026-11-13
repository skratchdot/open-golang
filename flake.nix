{
  description = "open-golang development environment";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs =
    { nixpkgs, ... }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
    in
    {
      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              delve
              firefox
              go
              gopls
              golangci-lint
            ];
            shellHook = ''
              export PATH="${./vendor}:$PATH"
              export BROWSER=firefox
              # Keep the vendored xdg-open on its generic path so the host's
              # desktop session doesn't change test behavior.
              unset \
                  DBUS_SESSION_BUS_ADDRESS \
                  DESKTOP_SESSION \
                  DISPLAY \
                  GNOME_DESKTOP_SESSION_ID \
                  KDE_FULL_SESSION
            '';
          };
        }
      );

      checks = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};

          # The tests launch firefox directly. Provide an inert stand-in so no
          # browser (or display) is needed; xdg-open comes from ./vendor.
          firefoxStub = pkgs.writeShellScriptBin "firefox" "exit 0";

          mkGoCheck =
            name: nativeBuildInputs: command:
            pkgs.stdenvNoCC.mkDerivation {
              inherit name;
              src = ./.;
              inherit nativeBuildInputs;
              dontBuild = true;
              doCheck = true;
              checkPhase = ''
                export PATH="${./vendor}:${firefoxStub}/bin:$PATH"
                export BROWSER=firefox
                unset \
                  DBUS_SESSION_BUS_ADDRESS \
                  DESKTOP_SESSION \
                  DISPLAY \
                  GNOME_DESKTOP_SESSION_ID \
                  KDE_FULL_SESSION
                ${command}
              '';
              installPhase = "mkdir -p $out";
            };
        in
        {
          go-test = mkGoCheck "go-test" [
            pkgs.go
            pkgs.writableTmpDirAsHomeHook
          ] "go test -v ./...";

          go-lint = mkGoCheck "go-lint" [
            pkgs.go
            pkgs.golangci-lint
            pkgs.writableTmpDirAsHomeHook
          ] "golangci-lint run ./...";
        }
      );
    };
}
